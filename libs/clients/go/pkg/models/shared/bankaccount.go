// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"github.com/formancehq/formance-sdk-go/pkg/utils"
	"time"
)

type BankAccount struct {
	AccountNumber *string           `json:"accountNumber,omitempty"`
	ConnectorID   string            `json:"connectorID"`
	Country       string            `json:"country"`
	CreatedAt     time.Time         `json:"createdAt"`
	Iban          *string           `json:"iban,omitempty"`
	ID            string            `json:"id"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	Provider      *string           `json:"provider,omitempty"`
	SwiftBicCode  *string           `json:"swiftBicCode,omitempty"`
}

func (b BankAccount) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(b, "", false)
}

func (b *BankAccount) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &b, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *BankAccount) GetAccountNumber() *string {
	if o == nil {
		return nil
	}
	return o.AccountNumber
}

func (o *BankAccount) GetConnectorID() string {
	if o == nil {
		return ""
	}
	return o.ConnectorID
}

func (o *BankAccount) GetCountry() string {
	if o == nil {
		return ""
	}
	return o.Country
}

func (o *BankAccount) GetCreatedAt() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.CreatedAt
}

func (o *BankAccount) GetIban() *string {
	if o == nil {
		return nil
	}
	return o.Iban
}

func (o *BankAccount) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *BankAccount) GetMetadata() map[string]string {
	if o == nil {
		return nil
	}
	return o.Metadata
}

func (o *BankAccount) GetProvider() *string {
	if o == nil {
		return nil
	}
	return o.Provider
}

func (o *BankAccount) GetSwiftBicCode() *string {
	if o == nil {
		return nil
	}
	return o.SwiftBicCode
}
