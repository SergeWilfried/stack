package ledgerstore

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"

	"github.com/formancehq/ledger/internal/storage/sqlutils"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/uptrace/bun"
)

func fetch[T any](s *Store, addModel bool, ctx context.Context, builders ...func(query *bun.SelectQuery) *bun.SelectQuery) (T, error) {

	var ret T
	ret = reflect.New(reflect.TypeOf(ret).Elem()).Interface().(T)

	query := s.bucket.db.NewSelect()

	if addModel {
		query = query.Model(ret)
	}

	for _, builder := range builders {
		query = query.Apply(builder)
	}

	if err := query.Scan(ctx, ret); err != nil {
		return ret, sqlutils.PostgresError(err)
	}

	return ret, nil
}

func paginateWithOffset[FILTERS any, RETURN any](s *Store, ctx context.Context,
	q *bunpaginate.OffsetPaginatedQuery[FILTERS], builders ...func(query *bun.SelectQuery) *bun.SelectQuery) (*api.Cursor[RETURN], error) {

	//var ret RETURN
	query := s.bucket.db.NewSelect()
	for _, builder := range builders {
		query = query.Apply(builder)
	}
	//if query.GetModel() == nil && query.GetTableName() == "" {
	//	query = query.Model(ret)
	//}

	return bunpaginate.UsingOffset[FILTERS, RETURN](ctx, query, *q)
}

func paginateWithColumn[FILTERS any, RETURN any](s *Store, ctx context.Context, q *bunpaginate.ColumnPaginatedQuery[FILTERS], builders ...func(query *bun.SelectQuery) *bun.SelectQuery) (*api.Cursor[RETURN], error) {
	query := s.bucket.db.NewSelect()
	for _, builder := range builders {
		query = query.Apply(builder)
	}

	ret, err := bunpaginate.UsingColumn[FILTERS, RETURN](ctx, query, *q)
	if err != nil {
		return nil, sqlutils.PostgresError(err)
	}

	return ret, nil
}

func count[T any](s *Store, addModel bool, ctx context.Context, builders ...func(query *bun.SelectQuery) *bun.SelectQuery) (int, error) {
	query := s.bucket.db.NewSelect()
	if addModel {
		query = query.Model((*T)(nil))
	}
	for _, builder := range builders {
		query = query.Apply(builder)
	}
	return s.bucket.db.NewSelect().
		TableExpr("(" + query.String() + ") data").
		Count(ctx)
}

func filterAccountAddress(address, key string) string {
	parts := make([]string, 0)
	src := strings.Split(address, ":")

	needSegmentCheck := false
	for _, segment := range src {
		needSegmentCheck = segment == ""
		if needSegmentCheck {
			break
		}
	}

	if needSegmentCheck {
		parts = append(parts, fmt.Sprintf("jsonb_array_length(%s_array) = %d", key, len(src)))

		for i, segment := range src {
			if len(segment) == 0 {
				continue
			}
			parts = append(parts, fmt.Sprintf("%s_array @@ ('$[%d] == \"%s\"')::jsonpath", key, i, segment))
		}
	} else {
		parts = append(parts, fmt.Sprintf("%s = '%s'", key, address))
	}

	return strings.Join(parts, " and ")
}

func filterAccountAddressOnTransactions(address string, source, destination bool) string {
	src := strings.Split(address, ":")

	needSegmentCheck := false
	for _, segment := range src {
		needSegmentCheck = segment == ""
		if needSegmentCheck {
			break
		}
	}

	if needSegmentCheck {
		m := map[string]any{
			fmt.Sprint(len(src)): nil,
		}
		parts := make([]string, 0)

		for i, segment := range src {
			if len(segment) == 0 {
				continue
			}
			m[fmt.Sprint(i)] = segment
		}

		data, err := json.Marshal([]any{m})
		if err != nil {
			panic(err)
		}

		if source {
			parts = append(parts, fmt.Sprintf("sources_arrays @> '%s'", string(data)))
		}
		if destination {
			parts = append(parts, fmt.Sprintf("destinations_arrays @> '%s'", string(data)))
		}
		return strings.Join(parts, " or ")
	} else {
		data, err := json.Marshal([]string{address})
		if err != nil {
			panic(err)
		}

		parts := make([]string, 0)
		if source {
			parts = append(parts, fmt.Sprintf("sources @> '%s'", string(data)))
		}
		if destination {
			parts = append(parts, fmt.Sprintf("destinations @> '%s'", string(data)))
		}
		return strings.Join(parts, " or ")
	}
}

func filterPIT(pit *ledger.Time, column string) func(query *bun.SelectQuery) *bun.SelectQuery {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		if pit == nil || pit.IsZero() {
			return query
		}
		return query.Where(fmt.Sprintf("%s <= ?", column), pit)
	}
}

type PaginatedQueryOptions[T any] struct {
	QueryBuilder query.Builder `json:"qb"`
	PageSize     uint64        `json:"pageSize"`
	Options      T             `json:"options"`
}

func (opts PaginatedQueryOptions[T]) WithQueryBuilder(qb query.Builder) PaginatedQueryOptions[T] {
	opts.QueryBuilder = qb

	return opts
}

func (opts PaginatedQueryOptions[T]) WithPageSize(pageSize uint64) PaginatedQueryOptions[T] {
	opts.PageSize = pageSize

	return opts
}

func NewPaginatedQueryOptions[T any](options T) PaginatedQueryOptions[T] {
	return PaginatedQueryOptions[T]{
		Options:  options,
		PageSize: bunpaginate.QueryDefaultPageSize,
	}
}

type PITFilter struct {
	PIT *ledger.Time `json:"pit"`
}

type PITFilterWithVolumes struct {
	PITFilter
	ExpandVolumes          bool `json:"volumes"`
	ExpandEffectiveVolumes bool `json:"effectiveVolumes"`
}
