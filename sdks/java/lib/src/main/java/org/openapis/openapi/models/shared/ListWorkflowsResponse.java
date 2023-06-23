/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package org.openapis.openapi.models.shared;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * ListWorkflowsResponse - List of workflows
 */
public class ListWorkflowsResponse {
    @JsonProperty("data")
    public Workflow[] data;

    public ListWorkflowsResponse withData(Workflow[] data) {
        this.data = data;
        return this;
    }
    
    public ListWorkflowsResponse(@JsonProperty("data") Workflow[] data) {
        this.data = data;
  }
}