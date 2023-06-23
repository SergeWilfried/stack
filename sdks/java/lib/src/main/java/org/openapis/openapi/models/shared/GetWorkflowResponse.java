/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package org.openapis.openapi.models.shared;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * GetWorkflowResponse - The workflow
 */
public class GetWorkflowResponse {
    @JsonProperty("data")
    public Workflow data;

    public GetWorkflowResponse withData(Workflow data) {
        this.data = data;
        return this;
    }
    
    public GetWorkflowResponse(@JsonProperty("data") Workflow data) {
        this.data = data;
  }
}