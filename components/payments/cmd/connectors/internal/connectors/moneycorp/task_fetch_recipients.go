package moneycorp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/moneycorp/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

type fetchRecipientsState struct {
	LastPage int `json:"last_page"`
	// Moneycorp does not allow us to sort by , but we can still
	// sort by ID created (which is incremental when creating accounts).
	LastCreatedAt time.Time `json:"last_created_at"`
}

func taskFetchRecipients(client *client.Client, accountID string) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		resolver task.StateResolver,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"moneycorp.taskFetchRecipients",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
			attribute.String("accountID", accountID),
		)
		defer span.End()

		state := task.MustResolveTo(ctx, resolver, fetchRecipientsState{})

		newState, err := fetchRecipients(ctx, client, accountID, connectorID, ingester, scheduler, state)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		if err := ingester.UpdateTaskState(ctx, newState); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func fetchRecipients(
	ctx context.Context,
	client *client.Client,
	accountID string,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	scheduler task.Scheduler,
	state fetchRecipientsState,
) (fetchRecipientsState, error) {
	newState := fetchRecipientsState{
		LastPage:      state.LastPage,
		LastCreatedAt: state.LastCreatedAt,
	}

	for page := 0; ; page++ {
		newState.LastPage = page

		pagedRecipients, err := client.GetRecipients(ctx, accountID, page, pageSize)
		if err != nil {
			return fetchRecipientsState{}, err
		}

		if len(pagedRecipients) == 0 {
			break
		}

		batch := ingestion.AccountBatch{}
		for _, recipient := range pagedRecipients {
			createdAt, err := time.Parse("2006-01-02T15:04:05.999999999", recipient.Attributes.CreatedAt)
			if err != nil {
				return fetchRecipientsState{}, fmt.Errorf("failed to parse transaction date: %w", err)
			}

			switch createdAt.Compare(state.LastCreatedAt) {
			case -1, 0:
				continue
			default:
			}

			raw, err := json.Marshal(recipient)
			if err != nil {
				return fetchRecipientsState{}, err
			}

			batch = append(batch, &models.Account{
				ID: models.AccountID{
					Reference:   recipient.ID,
					ConnectorID: connectorID,
				},
				// Moneycorp does not send the opening date of the account
				CreatedAt:    createdAt,
				Reference:    recipient.ID,
				ConnectorID:  connectorID,
				DefaultAsset: currency.FormatAsset(supportedCurrenciesWithDecimal, recipient.Attributes.BankAccountCurrency),
				AccountName:  recipient.Attributes.BankAccountName,
				Type:         models.AccountTypeExternal,
				RawData:      raw,
			})

			newState.LastCreatedAt = createdAt
		}

		if err := ingester.IngestAccounts(ctx, batch); err != nil {
			return fetchRecipientsState{}, err
		}

		if len(pagedRecipients) < pageSize {
			break
		}
	}

	return newState, nil
}
