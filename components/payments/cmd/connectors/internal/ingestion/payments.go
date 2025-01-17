package ingestion

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/publish"
)

type PaymentBatchElement struct {
	Payment    *models.Payment
	Adjustment *models.PaymentAdjustment
	Metadata   []*models.PaymentMetadata
}

type PaymentBatch []PaymentBatchElement

type IngesterFn func(ctx context.Context, batch PaymentBatch, commitState any) error

func (fn IngesterFn) IngestPayments(ctx context.Context, batch PaymentBatch, commitState any) error {
	return fn(ctx, batch, commitState)
}

func (i *DefaultIngester) IngestPayments(
	ctx context.Context,
	batch PaymentBatch,
) error {
	startingAt := time.Now()

	logging.FromContext(ctx).WithFields(map[string]interface{}{
		"size":       len(batch),
		"startingAt": startingAt,
	}).Debugf("Ingest batch")

	var allPayments []*models.Payment //nolint:prealloc // length is unknown
	var allMetadata []*models.PaymentMetadata
	var allAdjustments []*models.PaymentAdjustment

	for batchIdx := range batch {
		payment := batch[batchIdx].Payment
		metadata := batch[batchIdx].Metadata
		adjustment := batch[batchIdx].Adjustment

		if metadata != nil {
			for _, data := range metadata {
				data.Changelog = append(data.Changelog,
					models.MetadataChangelog{
						CreatedAt: time.Now(),
						Value:     data.Value,
					})

				allMetadata = append(allMetadata, data)
			}
		}

		if payment != nil {
			allPayments = append(allPayments, payment)
		}

		if adjustment != nil && adjustment.Reference != "" {
			allAdjustments = append(allAdjustments, adjustment)
		}
	}

	// Insert first all payments
	idsInserted, err := i.store.UpsertPayments(ctx, allPayments)
	if err != nil {
		return fmt.Errorf("error upserting payments: %w", err)
	}

	// Then insert all metadata
	if err := i.store.UpsertPaymentsMetadata(ctx, allMetadata); err != nil {
		return fmt.Errorf("error upserting payments metadata: %w", err)
	}

	// Then insert all adjustments
	if err := i.store.UpsertPaymentsAdjustments(ctx, allAdjustments); err != nil {
		return fmt.Errorf("error upserting payments adjustments: %w", err)
	}

	idsInsertedMap := make(map[string]struct{}, len(idsInserted))
	for idx := range idsInserted {
		idsInsertedMap[idsInserted[idx].String()] = struct{}{}
	}

	for paymentIdx := range allPayments {
		_, ok := idsInsertedMap[allPayments[paymentIdx].ID.String()]
		if !ok {
			// No need to publish an event for an already existing payment
			continue
		}
		err = i.publisher.Publish(events.TopicPayments,
			publish.NewMessage(ctx, i.messages.NewEventSavedPayments(i.provider, allPayments[paymentIdx])))
		if err != nil {
			logging.FromContext(ctx).Errorf("Publishing message: %w", err)

			continue
		}
	}

	endedAt := time.Now()

	logging.FromContext(ctx).WithFields(map[string]interface{}{
		"size":    len(batch),
		"endedAt": endedAt,
		"latency": endedAt.Sub(startingAt).String(),
	}).Debugf("Batch ingested")

	return nil
}
