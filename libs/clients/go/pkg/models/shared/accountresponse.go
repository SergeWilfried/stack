// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type AccountResponse struct {
	Data Account `json:"data"`
}

func (o *AccountResponse) GetData() Account {
	if o == nil {
		return Account{}
	}
	return o.Data
}