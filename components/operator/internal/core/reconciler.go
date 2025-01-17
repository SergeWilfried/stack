package core

import (
	"context"
	"reflect"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"k8s.io/client-go/util/workqueue"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func MapObjectToReconcileRequests[T client.Object](items ...T) []reconcile.Request {
	return Map(items, func(object T) reconcile.Request {
		return reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      object.GetName(),
				Namespace: object.GetNamespace(),
			},
		}
	})
}

type Initializer func(mgr Manager) error

var initializers = make([]Initializer, 0)

func Init(i ...Initializer) {
	initializers = append(initializers, i...)
}

type ReconcilerOptionsWatch struct {
	Handler func(mgr Manager, builder *builder.Builder, target client.Object) handler.EventHandler
}

type Finalizer[T client.Object] func(ctx Context, t T) error

type ReconcilerOptions[T client.Object] struct {
	Owns       map[client.Object][]builder.OwnsOption
	Watchers   map[client.Object]ReconcilerOptionsWatch
	Finalizers map[string]Finalizer[T]
	Raws       []func(Context, *builder.Builder) error
}

type ReconcilerOption[T client.Object] func(*ReconcilerOptions[T])

func WithOwn[T client.Object](v client.Object, opts ...builder.OwnsOption) ReconcilerOption[T] {
	return func(options *ReconcilerOptions[T]) {
		options.Owns[v] = opts
	}
}

func WithRaw[T client.Object](fn func(Context, *builder.Builder) error) ReconcilerOption[T] {
	return func(options *ReconcilerOptions[T]) {
		options.Raws = append(options.Raws, fn)
	}
}

func BuildReconcileRequests(ctx context.Context, client client.Client, scheme *runtime.Scheme, target client.Object, opts ...client.ListOption) []reconcile.Request {
	kinds, _, err := scheme.ObjectKinds(target)
	if err != nil {
		return []reconcile.Request{}
	}

	us := &unstructured.UnstructuredList{}
	us.SetGroupVersionKind(kinds[0])
	if err := client.List(ctx, us, opts...); err != nil {
		return []reconcile.Request{}
	}

	return MapObjectToReconcileRequests(
		Map(us.Items, ToPointer[unstructured.Unstructured])...,
	)
}

func WithFinalizer[T client.Object](name string, callback Finalizer[T]) ReconcilerOption[T] {
	return func(r *ReconcilerOptions[T]) {
		r.Finalizers[name] = callback
	}
}

func WithWatchSettings[T client.Object]() ReconcilerOption[T] {
	return func(options *ReconcilerOptions[T]) {
		options.Watchers[&v1beta1.Settings{}] = ReconcilerOptionsWatch{
			Handler: func(mgr Manager, builder *builder.Builder, target client.Object) handler.EventHandler {
				return handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, object client.Object) []reconcile.Request {
					settings := object.(*v1beta1.Settings)

					ret := make([]reconcile.Request, 0)
					if !settings.IsWildcard() {
						for _, stack := range settings.GetStacks() {
							ret = append(ret, BuildReconcileRequests(ctx, mgr.GetClient(), mgr.GetScheme(), target, client.MatchingFields{
								"stack": stack,
							})...)
						}
					} else {
						ret = append(ret, BuildReconcileRequests(ctx, mgr.GetClient(), mgr.GetScheme(), target)...)
					}

					return ret
				})
			},
		}
	}
}

func WithWatchDependency[T client.Object](t v1beta1.Dependent) ReconcilerOption[T] {
	return func(options *ReconcilerOptions[T]) {
		options.Watchers[t] = ReconcilerOptionsWatch{
			Handler: func(mgr Manager, builder *builder.Builder, target client.Object) handler.EventHandler {
				return handler.EnqueueRequestsFromMapFunc(WatchDependents(mgr, target))
			},
		}
	}
}

func WithWatchStack[T client.Object]() ReconcilerOption[T] {
	return func(options *ReconcilerOptions[T]) {
		options.Watchers[&v1beta1.Stack{}] = ReconcilerOptionsWatch{
			Handler: func(mgr Manager, builder *builder.Builder, target client.Object) handler.EventHandler {
				return handler.EnqueueRequestsFromMapFunc(Watch(mgr, target))
			},
		}
	}
}

func WithWatch[T client.Object, WATCHED client.Object](fn func(ctx Context, object WATCHED) []reconcile.Request) ReconcilerOption[T] {
	var watched WATCHED
	watched = reflect.New(reflect.TypeOf(watched).Elem()).Interface().(WATCHED)
	return func(options *ReconcilerOptions[T]) {
		options.Watchers[watched] = ReconcilerOptionsWatch{
			Handler: func(mgr Manager, builder *builder.Builder, target client.Object) handler.EventHandler {
				return handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, object client.Object) []reconcile.Request {
					return fn(NewContext(mgr.GetClient(), mgr.GetScheme(), mgr.GetPlatform(), ctx), object.(WATCHED))
				})
			},
		}
	}
}

func WithReconciler[T client.Object](controller ObjectController[T], opts ...ReconcilerOption[T]) Initializer {
	return func(mgr Manager) error {

		options := ReconcilerOptions[T]{
			Owns:       map[client.Object][]builder.OwnsOption{},
			Watchers:   map[client.Object]ReconcilerOptionsWatch{},
			Finalizers: map[string]Finalizer[T]{},
		}
		for _, opt := range opts {
			opt(&options)
		}

		var t T
		t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
		b := ctrl.NewControllerManagedBy(mgr).
			For(t, builder.WithPredicates(predicate.Or(
				predicate.GenerationChangedPredicate{},
				predicate.Funcs{
					CreateFunc: func(event event.CreateEvent) bool {
						return true
					},
					DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
						return true
					},
					UpdateFunc: func(updateEvent event.UpdateEvent) bool {
					l:
						for _, referenceFromNew := range updateEvent.ObjectNew.GetOwnerReferences() {
							for _, referenceFromOld := range updateEvent.ObjectOld.GetOwnerReferences() {
								if referenceFromNew.UID == referenceFromOld.UID {
									continue l
								}
							}
							return true
						}

						return len(updateEvent.ObjectOld.GetOwnerReferences()) != len(updateEvent.ObjectNew.GetOwnerReferences())
					},
					GenericFunc: func(genericEvent event.GenericEvent) bool {
						return true
					},
				},
			)))

		for object, ownsOptions := range options.Owns {
			b = b.Owns(object, ownsOptions...)
		}
		for object, watch := range options.Watchers {
			b = b.Watches(object, watch.Handler(mgr, b, t))
		}
		for _, raw := range options.Raws {
			if err := raw(NewContext(mgr.GetClient(), mgr.GetScheme(), mgr.GetPlatform(), context.Background()), b); err != nil {
				return err
			}
		}

		return b.Complete(reconcile.Func(reconcileObject(mgr, controller, options.Finalizers)))
	}
}

func reconcileObject[T client.Object](mgr Manager, controller ObjectController[T], finalizers map[string]Finalizer[T]) func(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	return func(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {

		var object T
		object = reflect.New(reflect.TypeOf(object).Elem()).Interface().(T)
		if err := mgr.GetClient().Get(ctx, types.NamespacedName{
			Name: request.Name,
		}, object); err != nil {
			if apierrors.IsNotFound(err) {
				return ctrl.Result{}, nil
			}
			return ctrl.Result{}, err
		}

		reconcileContext := NewContext(mgr.GetClient(), mgr.GetScheme(), mgr.GetPlatform(), ctx)
		if !object.GetDeletionTimestamp().IsZero() {
			log.FromContext(ctx).Info("Resource " + request.Name + " deleted, calling finalizers...")
			for name, f := range finalizers {

				if !Contains(object.GetFinalizers(), name) {
					continue
				}

				if err := f(reconcileContext, object); err != nil {
					if IsApplicationError(err) {
						return reconcile.Result{}, nil
					}
					return reconcile.Result{}, errors.Wrapf(err, "executing finalizer '%s'", name)
				}

				patch := client.MergeFrom(object.DeepCopyObject().(T))
				if controllerutil.RemoveFinalizer(object, name) {
					if err := mgr.GetClient().Patch(ctx, object, patch); err != nil {
						if apierrors.IsConflict(err) {
							return reconcile.Result{Requeue: true}, nil
						}
						return reconcile.Result{}, errors.Wrapf(err, "patching resource to remove finalizer '%s'", name)
					}
				}
			}

			return reconcile.Result{}, nil
		}

		log.FromContext(ctx).Info("Reconcile " + request.Name)
		missingFinalizers := make([]string, 0)
		for name := range finalizers {
			if !Contains(object.GetFinalizers(), name) {
				missingFinalizers = append(missingFinalizers, name)
			}
		}
		if len(missingFinalizers) > 0 {
			patch := client.MergeFrom(object.DeepCopyObject().(T))
			finalizers := object.GetFinalizers()
			finalizers = append(finalizers, missingFinalizers...)
			object.SetFinalizers(finalizers)

			if err := mgr.GetClient().Patch(ctx, object, patch); err != nil {
				return reconcile.Result{}, errors.Wrap(err, "patching missing finalizers")
			}
		}

		cp := object.DeepCopyObject().(T)
		patch := client.MergeFrom(cp)

		var reconcilerError error
		err := controller(reconcileContext, object)
		if err != nil {
			if !IsApplicationError(err) {
				reconcilerError = errors.Wrap(err, "reconciling resource")
			}
		}

		if err := mgr.GetClient().Status().Patch(ctx, object, patch); err != nil {
			if apierrors.IsNotFound(err) {
				// Ignore resource deleted
				return ctrl.Result{}, nil
			}
			return ctrl.Result{}, errors.Wrap(err, "patching resource to update status")
		}

		if apierrors.IsConflict(reconcilerError) {
			return ctrl.Result{
				Requeue: true,
			}, nil
		}

		return ctrl.Result{}, reconcilerError
	}
}

func WithStdReconciler[T v1beta1.Object](ctrl ObjectController[T], opts ...ReconcilerOption[T]) Initializer {
	return WithReconciler(ForObjectController(ctrl), opts...)
}

func WithStackDependencyReconciler[T v1beta1.Dependent](fn StackDependentObjectController[T], opts ...ReconcilerOption[T]) Initializer {
	opts = append(opts, WithWatchStack[T]())
	return WithStdReconciler(ForStackDependency(fn), opts...)
}

func WithResourceReconciler[T v1beta1.Dependent](fn StackDependentObjectController[T], opts ...ReconcilerOption[T]) Initializer {
	return WithStackDependencyReconciler(ForResource(fn), opts...)
}

func WithModuleReconciler[T v1beta1.Module](fn ModuleController[T], opts ...ReconcilerOption[T]) Initializer {
	opts = append(opts, WithWatchVersions[T])
	return WithStackDependencyReconciler(ForModule(fn), opts...)
}

func WithWatchVersions[T client.Object](options *ReconcilerOptions[T]) {

	reconcileModule := func(ctx context.Context, mgr Manager, target client.Object, versionFileName string, limitingInterface workqueue.RateLimitingInterface) {
		stackList := &v1beta1.StackList{}
		if err := mgr.GetClient().List(ctx, stackList, client.MatchingFields{
			".spec.versionsFromFile": versionFileName,
		}); err != nil {
			panic(err)
		}

		kinds, _, err := mgr.GetScheme().ObjectKinds(target)
		if err != nil {
			panic(err)
		}

		for _, stack := range stackList.Items {
			list := &unstructured.UnstructuredList{}
			list.SetGroupVersionKind(kinds[0])
			if err := mgr.GetClient().List(ctx, list, client.MatchingFields{
				"stack": stack.Name,
			}); err != nil {
				panic(err)
			}

			for _, item := range list.Items {
				limitingInterface.Add(reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name: item.GetName(),
					},
				})
			}
		}
	}

	options.Watchers[&v1beta1.Versions{}] = ReconcilerOptionsWatch{
		Handler: func(mgr Manager, builder *builder.Builder, target client.Object) handler.EventHandler {
			return handler.Funcs{
				CreateFunc: func(ctx context.Context, createEvent event.CreateEvent, limitingInterface workqueue.RateLimitingInterface) {
					reconcileModule(ctx, mgr, target, createEvent.Object.GetName(), limitingInterface)
				},
				UpdateFunc: func(ctx context.Context, updateEvent event.UpdateEvent, limitingInterface workqueue.RateLimitingInterface) {
					oldObject := updateEvent.ObjectOld.(*v1beta1.Versions)
					newObject := updateEvent.ObjectNew.(*v1beta1.Versions)

					kinds, _, err := mgr.GetScheme().ObjectKinds(target)
					if err != nil {
						panic(err)
					}
					kind := strings.ToLower(kinds[0].Kind)
					if oldObject.Spec[kind] == newObject.Spec[kind] {
						return
					}

					reconcileModule(ctx, mgr, target, updateEvent.ObjectNew.GetName(), limitingInterface)
				},
				DeleteFunc: func(ctx context.Context, deleteEvent event.DeleteEvent, limitingInterface workqueue.RateLimitingInterface) {
					reconcileModule(ctx, mgr, target, deleteEvent.Object.GetName(), limitingInterface)
				},
			}
		},
	}
}

func WithIndex[T client.Object](name string, eval func(t T) []string) Initializer {
	return func(mgr Manager) error {
		var t T
		t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
		return mgr.GetFieldIndexer().
			IndexField(context.Background(), t, name, func(rawObj client.Object) []string {
				return eval(rawObj.(T))
			})
	}
}

func WithSimpleIndex[T client.Object](name string, eval func(t T) string) Initializer {
	return WithIndex(name, func(t T) []string {
		return []string{eval(t)}
	})
}
