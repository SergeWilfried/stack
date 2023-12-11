// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/formancehq/formance-sdk-go/pkg/utils"
	"math/big"
)

type BulkElementBulkElementDeleteMetadataData struct {
	Key        string     `json:"key"`
	TargetID   TargetID   `json:"targetId"`
	TargetType TargetType `json:"targetType"`
}

func (o *BulkElementBulkElementDeleteMetadataData) GetKey() string {
	if o == nil {
		return ""
	}
	return o.Key
}

func (o *BulkElementBulkElementDeleteMetadataData) GetTargetID() TargetID {
	if o == nil {
		return TargetID{}
	}
	return o.TargetID
}

func (o *BulkElementBulkElementDeleteMetadataData) GetTargetType() TargetType {
	if o == nil {
		return TargetType("")
	}
	return o.TargetType
}

type BulkElementBulkElementDeleteMetadata struct {
	Action string                                    `json:"action"`
	Data   *BulkElementBulkElementDeleteMetadataData `json:"data,omitempty"`
	Ik     *string                                   `json:"ik,omitempty"`
}

func (o *BulkElementBulkElementDeleteMetadata) GetAction() string {
	if o == nil {
		return ""
	}
	return o.Action
}

func (o *BulkElementBulkElementDeleteMetadata) GetData() *BulkElementBulkElementDeleteMetadataData {
	if o == nil {
		return nil
	}
	return o.Data
}

func (o *BulkElementBulkElementDeleteMetadata) GetIk() *string {
	if o == nil {
		return nil
	}
	return o.Ik
}

type BulkElementBulkElementRevertTransactionData struct {
	Force *bool    `json:"force,omitempty"`
	ID    *big.Int `json:"id"`
}

func (b BulkElementBulkElementRevertTransactionData) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(b, "", false)
}

func (b *BulkElementBulkElementRevertTransactionData) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &b, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *BulkElementBulkElementRevertTransactionData) GetForce() *bool {
	if o == nil {
		return nil
	}
	return o.Force
}

func (o *BulkElementBulkElementRevertTransactionData) GetID() *big.Int {
	if o == nil {
		return big.NewInt(0)
	}
	return o.ID
}

type BulkElementBulkElementRevertTransaction struct {
	Action string                                       `json:"action"`
	Data   *BulkElementBulkElementRevertTransactionData `json:"data,omitempty"`
	Ik     *string                                      `json:"ik,omitempty"`
}

func (o *BulkElementBulkElementRevertTransaction) GetAction() string {
	if o == nil {
		return ""
	}
	return o.Action
}

func (o *BulkElementBulkElementRevertTransaction) GetData() *BulkElementBulkElementRevertTransactionData {
	if o == nil {
		return nil
	}
	return o.Data
}

func (o *BulkElementBulkElementRevertTransaction) GetIk() *string {
	if o == nil {
		return nil
	}
	return o.Ik
}

type BulkElementBulkElementAddMetadataData struct {
	Metadata   map[string]string `json:"metadata"`
	TargetID   TargetID          `json:"targetId"`
	TargetType TargetType        `json:"targetType"`
}

func (o *BulkElementBulkElementAddMetadataData) GetMetadata() map[string]string {
	if o == nil {
		return map[string]string{}
	}
	return o.Metadata
}

func (o *BulkElementBulkElementAddMetadataData) GetTargetID() TargetID {
	if o == nil {
		return TargetID{}
	}
	return o.TargetID
}

func (o *BulkElementBulkElementAddMetadataData) GetTargetType() TargetType {
	if o == nil {
		return TargetType("")
	}
	return o.TargetType
}

type BulkElementBulkElementAddMetadata struct {
	Action string                                 `json:"action"`
	Data   *BulkElementBulkElementAddMetadataData `json:"data,omitempty"`
	Ik     *string                                `json:"ik,omitempty"`
}

func (o *BulkElementBulkElementAddMetadata) GetAction() string {
	if o == nil {
		return ""
	}
	return o.Action
}

func (o *BulkElementBulkElementAddMetadata) GetData() *BulkElementBulkElementAddMetadataData {
	if o == nil {
		return nil
	}
	return o.Data
}

func (o *BulkElementBulkElementAddMetadata) GetIk() *string {
	if o == nil {
		return nil
	}
	return o.Ik
}

type BulkElementBulkElementCreateTransaction struct {
	Action string           `json:"action"`
	Data   *PostTransaction `json:"data,omitempty"`
	Ik     *string          `json:"ik,omitempty"`
}

func (o *BulkElementBulkElementCreateTransaction) GetAction() string {
	if o == nil {
		return ""
	}
	return o.Action
}

func (o *BulkElementBulkElementCreateTransaction) GetData() *PostTransaction {
	if o == nil {
		return nil
	}
	return o.Data
}

func (o *BulkElementBulkElementCreateTransaction) GetIk() *string {
	if o == nil {
		return nil
	}
	return o.Ik
}

type BulkElementType string

const (
	BulkElementTypeAddMetadata       BulkElementType = "ADD_METADATA"
	BulkElementTypeCreateTransaction BulkElementType = "CREATE_TRANSACTION"
	BulkElementTypeDeleteMetadata    BulkElementType = "DELETE_METADATA"
	BulkElementTypeRevertTransaction BulkElementType = "REVERT_TRANSACTION"
)

type BulkElement struct {
	BulkElementBulkElementCreateTransaction *BulkElementBulkElementCreateTransaction
	BulkElementBulkElementAddMetadata       *BulkElementBulkElementAddMetadata
	BulkElementBulkElementRevertTransaction *BulkElementBulkElementRevertTransaction
	BulkElementBulkElementDeleteMetadata    *BulkElementBulkElementDeleteMetadata

	Type BulkElementType
}

func CreateBulkElementAddMetadata(addMetadata BulkElementBulkElementAddMetadata) BulkElement {
	typ := BulkElementTypeAddMetadata
	typStr := string(typ)
	addMetadata.Action = typStr

	return BulkElement{
		BulkElementBulkElementAddMetadata: &addMetadata,
		Type:                              typ,
	}
}

func CreateBulkElementCreateTransaction(createTransaction BulkElementBulkElementCreateTransaction) BulkElement {
	typ := BulkElementTypeCreateTransaction
	typStr := string(typ)
	createTransaction.Action = typStr

	return BulkElement{
		BulkElementBulkElementCreateTransaction: &createTransaction,
		Type:                                    typ,
	}
}

func CreateBulkElementDeleteMetadata(deleteMetadata BulkElementBulkElementDeleteMetadata) BulkElement {
	typ := BulkElementTypeDeleteMetadata
	typStr := string(typ)
	deleteMetadata.Action = typStr

	return BulkElement{
		BulkElementBulkElementDeleteMetadata: &deleteMetadata,
		Type:                                 typ,
	}
}

func CreateBulkElementRevertTransaction(revertTransaction BulkElementBulkElementRevertTransaction) BulkElement {
	typ := BulkElementTypeRevertTransaction
	typStr := string(typ)
	revertTransaction.Action = typStr

	return BulkElement{
		BulkElementBulkElementRevertTransaction: &revertTransaction,
		Type:                                    typ,
	}
}

func (u *BulkElement) UnmarshalJSON(data []byte) error {

	type discriminator struct {
		Action string
	}

	dis := new(discriminator)
	if err := json.Unmarshal(data, &dis); err != nil {
		return fmt.Errorf("could not unmarshal discriminator: %w", err)
	}

	switch dis.Action {
	case "ADD_METADATA":
		bulkElementBulkElementAddMetadata := new(BulkElementBulkElementAddMetadata)
		if err := utils.UnmarshalJSON(data, &bulkElementBulkElementAddMetadata, "", true, true); err != nil {
			return fmt.Errorf("could not unmarshal expected type: %w", err)
		}

		u.BulkElementBulkElementAddMetadata = bulkElementBulkElementAddMetadata
		u.Type = BulkElementTypeAddMetadata
		return nil
	case "CREATE_TRANSACTION":
		bulkElementBulkElementCreateTransaction := new(BulkElementBulkElementCreateTransaction)
		if err := utils.UnmarshalJSON(data, &bulkElementBulkElementCreateTransaction, "", true, true); err != nil {
			return fmt.Errorf("could not unmarshal expected type: %w", err)
		}

		u.BulkElementBulkElementCreateTransaction = bulkElementBulkElementCreateTransaction
		u.Type = BulkElementTypeCreateTransaction
		return nil
	case "DELETE_METADATA":
		bulkElementBulkElementDeleteMetadata := new(BulkElementBulkElementDeleteMetadata)
		if err := utils.UnmarshalJSON(data, &bulkElementBulkElementDeleteMetadata, "", true, true); err != nil {
			return fmt.Errorf("could not unmarshal expected type: %w", err)
		}

		u.BulkElementBulkElementDeleteMetadata = bulkElementBulkElementDeleteMetadata
		u.Type = BulkElementTypeDeleteMetadata
		return nil
	case "REVERT_TRANSACTION":
		bulkElementBulkElementRevertTransaction := new(BulkElementBulkElementRevertTransaction)
		if err := utils.UnmarshalJSON(data, &bulkElementBulkElementRevertTransaction, "", true, true); err != nil {
			return fmt.Errorf("could not unmarshal expected type: %w", err)
		}

		u.BulkElementBulkElementRevertTransaction = bulkElementBulkElementRevertTransaction
		u.Type = BulkElementTypeRevertTransaction
		return nil
	}

	return errors.New("could not unmarshal into supported union types")
}

func (u BulkElement) MarshalJSON() ([]byte, error) {
	if u.BulkElementBulkElementCreateTransaction != nil {
		return utils.MarshalJSON(u.BulkElementBulkElementCreateTransaction, "", true)
	}

	if u.BulkElementBulkElementAddMetadata != nil {
		return utils.MarshalJSON(u.BulkElementBulkElementAddMetadata, "", true)
	}

	if u.BulkElementBulkElementRevertTransaction != nil {
		return utils.MarshalJSON(u.BulkElementBulkElementRevertTransaction, "", true)
	}

	if u.BulkElementBulkElementDeleteMetadata != nil {
		return utils.MarshalJSON(u.BulkElementBulkElementDeleteMetadata, "", true)
	}

	return nil, errors.New("could not marshal union type: all fields are null")
}