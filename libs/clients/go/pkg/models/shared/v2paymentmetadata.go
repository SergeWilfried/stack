// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type V2PaymentMetadata struct {
	Key *string `json:"key,omitempty"`
}

func (o *V2PaymentMetadata) GetKey() *string {
	if o == nil {
		return nil
	}
	return o.Key
}