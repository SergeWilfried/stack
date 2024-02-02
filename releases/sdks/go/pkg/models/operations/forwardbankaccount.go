// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"net/http"
)

type ForwardBankAccountRequest struct {
	ForwardBankAccountRequest shared.ForwardBankAccountRequest `request:"mediaType=application/json"`
	// The bank account ID.
	BankAccountID string `pathParam:"style=simple,explode=false,name=bankAccountId"`
}

func (o *ForwardBankAccountRequest) GetForwardBankAccountRequest() shared.ForwardBankAccountRequest {
	if o == nil {
		return shared.ForwardBankAccountRequest{}
	}
	return o.ForwardBankAccountRequest
}

func (o *ForwardBankAccountRequest) GetBankAccountID() string {
	if o == nil {
		return ""
	}
	return o.BankAccountID
}

type ForwardBankAccountResponse struct {
	// OK
	BankAccountResponse *shared.BankAccountResponse
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
}

func (o *ForwardBankAccountResponse) GetBankAccountResponse() *shared.BankAccountResponse {
	if o == nil {
		return nil
	}
	return o.BankAccountResponse
}

func (o *ForwardBankAccountResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *ForwardBankAccountResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *ForwardBankAccountResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}