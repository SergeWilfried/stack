// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

type PaymentType string

const (
	PaymentTypePayIn    PaymentType = "PAY-IN"
	PaymentTypePayout   PaymentType = "PAYOUT"
	PaymentTypeTransfer PaymentType = "TRANSFER"
	PaymentTypeOther    PaymentType = "OTHER"
)

func (e PaymentType) ToPointer() *PaymentType {
	return &e
}

func (e *PaymentType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "PAY-IN":
		fallthrough
	case "PAYOUT":
		fallthrough
	case "TRANSFER":
		fallthrough
	case "OTHER":
		*e = PaymentType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for PaymentType: %v", v)
	}
}