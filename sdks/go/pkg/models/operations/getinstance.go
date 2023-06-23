// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"net/http"
)

type GetInstanceRequest struct {
	// The instance id
	InstanceID string `pathParam:"style=simple,explode=false,name=instanceID"`
}

type GetInstanceResponse struct {
	ContentType string
	// General error
	Error *shared.Error
	// The workflow instance
	GetWorkflowInstanceResponse *shared.GetWorkflowInstanceResponse
	StatusCode                  int
	RawResponse                 *http.Response
}