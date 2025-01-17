package atlar

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/atlar/client"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/contextutil"
	"github.com/get-momo/atlar-v1-go-client/client/transactions"
	atlar_models "github.com/get-momo/atlar-v1-go-client/models"
	"go.opentelemetry.io/otel/attribute"
)

func FetchTransactionsTask(config Config, client *client.Client) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		resolver task.StateResolver,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"atlar.taskFetchTransactions",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
		)
		defer span.End()

		// Pagination works by cursor token.
		for token := ""; ; {
			requestCtx, cancel := contextutil.DetachedWithTimeout(ctx, 30*time.Second)
			defer cancel()
			pagedTransactions, err := client.GetV1Transactions(requestCtx, token, int64(config.PageSize))
			if err != nil {
				otel.RecordError(span, err)
				return err
			}

			token = pagedTransactions.Payload.NextToken

			if err := ingestPaymentsBatch(ctx, connectorID, ingester, pagedTransactions); err != nil {
				otel.RecordError(span, err)
				return err
			}

			if token == "" {
				break
			}
		}

		return nil
	}
}

func ingestPaymentsBatch(
	ctx context.Context,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	pagedTransactions *transactions.GetV1TransactionsOK,
) error {
	batch := ingestion.PaymentBatch{}

	for _, item := range pagedTransactions.Payload.Items {
		raw, err := json.Marshal(item)
		if err != nil {
			return err
		}

		paymentType := determinePaymentType(item)

		itemAmount := item.Amount
		precision := supportedCurrenciesWithDecimal[*itemAmount.Currency]

		amount, err := atlarTransactionAmountToPaymentAbsoluteAmount(*itemAmount.StringValue, precision)
		if err != nil {
			return err
		}

		createdAt, err := ParseAtlarTimestamp(item.Created)
		if err != nil {
			return err
		}

		paymentId := models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: item.ID,
				Type:      paymentType,
			},
			ConnectorID: connectorID,
		}

		batchElement := ingestion.PaymentBatchElement{
			Payment: &models.Payment{
				ID:            paymentId,
				Reference:     item.ID,
				Type:          paymentType,
				ConnectorID:   connectorID,
				CreatedAt:     createdAt,
				Status:        determinePaymentStatus(item),
				Scheme:        determinePaymentScheme(item),
				Amount:        amount,
				InitialAmount: amount,
				Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, *item.Amount.Currency),
				Metadata:      ExtractPaymentMetadata(paymentId, item),
				RawData:       raw,
			},
		}

		if *itemAmount.Value >= 0 {
			// DEBIT
			batchElement.Payment.DestinationAccountID = &models.AccountID{
				Reference:   *item.Account.ID,
				ConnectorID: connectorID,
			}
		} else {
			// CREDIT
			batchElement.Payment.SourceAccountID = &models.AccountID{
				Reference:   *item.Account.ID,
				ConnectorID: connectorID,
			}
		}

		batch = append(batch, batchElement)
	}

	if err := ingester.IngestPayments(ctx, batch); err != nil {
		return err
	}

	return nil
}

func determinePaymentType(item *atlar_models.Transaction) models.PaymentType {
	if *item.Amount.Value >= 0 {
		return models.PaymentTypePayIn
	} else {
		return models.PaymentTypePayOut
	}
}

func determinePaymentStatus(item *atlar_models.Transaction) models.PaymentStatus {
	if item.Reconciliation.Status == atlar_models.ReconciliationDetailsStatusEXPECTED {
		// A payment initiated by the owner of the accunt through the Atlar API,
		// which was not yet reconciled with a payment from the statement
		return models.PaymentStatusPending
	}
	if item.Reconciliation.Status == atlar_models.ReconciliationDetailsStatusBOOKED {
		// A payment comissioned with the bank, which was not yet reconciled with a
		// payment from the statement
		return models.PaymentStatusSucceeded
	}
	if item.Reconciliation.Status == atlar_models.ReconciliationDetailsStatusRECONCILED {
		return models.PaymentStatusSucceeded
	}
	return models.PaymentStatusOther
}

func determinePaymentScheme(item *atlar_models.Transaction) models.PaymentScheme {
	// item.Characteristics.BankTransactionCode.Domain
	// item.Characteristics.BankTransactionCode.Family
	// TODO: fees and interest -> models.PaymentSchemeOther with additional info on metadata. Will need example transactions for that

	if *item.Amount.Value > 0 {
		return models.PaymentSchemeSepaDebit
	} else if *item.Amount.Value < 0 {
		return models.PaymentSchemeSepaCredit
	}
	return models.PaymentSchemeSepa
}

func ExtractPaymentMetadata(paymentId models.PaymentID, transaction *atlar_models.Transaction) []*models.PaymentMetadata {
	result := []*models.PaymentMetadata{}
	if transaction.Date != "" {
		result = append(result, ComputePaymentMetadata(paymentId, "date", transaction.Date))
	}
	if transaction.ValueDate != "" {
		result = append(result, ComputePaymentMetadata(paymentId, "valueDate", transaction.ValueDate))
	}
	result = append(result, ComputePaymentMetadata(paymentId, "remittanceInformation/type", *transaction.RemittanceInformation.Type))
	result = append(result, ComputePaymentMetadata(paymentId, "remittanceInformation/value", *transaction.RemittanceInformation.Value))
	result = append(result, ComputePaymentMetadata(paymentId, "btc/domain", transaction.Characteristics.BankTransactionCode.Domain))
	result = append(result, ComputePaymentMetadata(paymentId, "btc/familiy", transaction.Characteristics.BankTransactionCode.Family))
	result = append(result, ComputePaymentMetadata(paymentId, "btc/subfamiliy", transaction.Characteristics.BankTransactionCode.Subfamily))
	result = append(result, ComputePaymentMetadata(paymentId, "btc/description", transaction.Characteristics.BankTransactionCode.Description))
	result = append(result, ComputePaymentMetadataBool(paymentId, "returned", transaction.Characteristics.Returned))
	if transaction.CounterpartyDetails != nil && transaction.CounterpartyDetails.Name != "" {
		result = append(result, ComputePaymentMetadata(paymentId, "counterparty/name", transaction.CounterpartyDetails.Name))
		if transaction.CounterpartyDetails.ExternalAccount != nil && transaction.CounterpartyDetails.ExternalAccount.Identifier != nil {
			result = append(result, ComputePaymentMetadata(paymentId, "counterparty/bank/bic", transaction.CounterpartyDetails.ExternalAccount.Bank.Bic))
			result = append(result, ComputePaymentMetadata(paymentId, "counterparty/bank/name", transaction.CounterpartyDetails.ExternalAccount.Bank.Name))
			result = append(result, ComputePaymentMetadata(paymentId,
				fmt.Sprintf("counterparty/identifier/%s", transaction.CounterpartyDetails.ExternalAccount.Identifier.Type),
				transaction.CounterpartyDetails.ExternalAccount.Identifier.Number))
		}
	}
	if transaction.Characteristics.Returned {
		result = append(result, ComputePaymentMetadata(paymentId, "returnReason/code", transaction.Characteristics.ReturnReason.Code))
		result = append(result, ComputePaymentMetadata(paymentId, "returnReason/description", transaction.Characteristics.ReturnReason.Description))
		result = append(result, ComputePaymentMetadata(paymentId, "returnReason/btc/domain", transaction.Characteristics.ReturnReason.OriginalBankTransactionCode.Domain))
		result = append(result, ComputePaymentMetadata(paymentId, "returnReason/btc/family", transaction.Characteristics.ReturnReason.OriginalBankTransactionCode.Family))
		result = append(result, ComputePaymentMetadata(paymentId, "returnReason/btc/subfamily", transaction.Characteristics.ReturnReason.OriginalBankTransactionCode.Subfamily))
		result = append(result, ComputePaymentMetadata(paymentId, "returnReason/btc/description", transaction.Characteristics.ReturnReason.OriginalBankTransactionCode.Description))
	}
	if transaction.Characteristics.VirtualAccount != nil {
		result = append(result, ComputePaymentMetadata(paymentId, "virtualAccount/market", transaction.Characteristics.VirtualAccount.Market))
		result = append(result, ComputePaymentMetadata(paymentId, "virtualAccount/rawIdentifier", transaction.Characteristics.VirtualAccount.RawIdentifier))
		result = append(result, ComputePaymentMetadata(paymentId, "virtualAccount/bank/id", transaction.Characteristics.VirtualAccount.Bank.ID))
		result = append(result, ComputePaymentMetadata(paymentId, "virtualAccount/bank/name", transaction.Characteristics.VirtualAccount.Bank.Name))
		result = append(result, ComputePaymentMetadata(paymentId, "virtualAccount/bank/bic", transaction.Characteristics.VirtualAccount.Bank.Bic))
		result = append(result, ComputePaymentMetadata(paymentId, "virtualAccount/identifier/holderName", *transaction.Characteristics.VirtualAccount.Identifier.HolderName))
		result = append(result, ComputePaymentMetadata(paymentId, "virtualAccount/identifier/market", transaction.Characteristics.VirtualAccount.Identifier.Market))
		result = append(result, ComputePaymentMetadata(paymentId, "virtualAccount/identifier/type", transaction.Characteristics.VirtualAccount.Identifier.Type))
		result = append(result, ComputePaymentMetadata(paymentId, "virtualAccount/identifier/number", transaction.Characteristics.VirtualAccount.Identifier.Number))
	}
	result = append(result, ComputePaymentMetadata(paymentId, "reconciliation/status", transaction.Reconciliation.Status))
	result = append(result, ComputePaymentMetadata(paymentId, "reconciliation/transactableId", transaction.Reconciliation.TransactableID))
	result = append(result, ComputePaymentMetadata(paymentId, "reconciliation/transactableType", transaction.Reconciliation.TransactableType))
	if transaction.Characteristics.CurrencyExchange != nil {
		result = append(result, ComputePaymentMetadata(paymentId, "currencyExchange/sourceCurrency", transaction.Characteristics.CurrencyExchange.SourceCurrency))
		result = append(result, ComputePaymentMetadata(paymentId, "currencyExchange/targetCurrency", transaction.Characteristics.CurrencyExchange.TargetCurrency))
		result = append(result, ComputePaymentMetadata(paymentId, "currencyExchange/exchangeRate", transaction.Characteristics.CurrencyExchange.ExchangeRate))
		result = append(result, ComputePaymentMetadata(paymentId, "currencyExchange/unitCurrency", transaction.Characteristics.CurrencyExchange.UnitCurrency))
	}
	if transaction.CounterpartyDetails.MandateReference != "" {
		result = append(result, ComputePaymentMetadata(paymentId, "mandateReference", transaction.CounterpartyDetails.MandateReference))
	}

	return result
}

func atlarTransactionAmountToPaymentAbsoluteAmount(atlarAmount string, precision int) (*big.Int, error) {
	var amount big.Float
	_, ok := amount.SetString(atlarAmount)
	amount.Abs(&amount)
	if !ok {
		return nil, fmt.Errorf("failed to parse amount %s", atlarAmount)
	}

	var amountInt big.Int
	amount.Mul(&amount, big.NewFloat(math.Pow(10, float64(precision)))).Int(&amountInt)

	return &amountInt, nil
}
