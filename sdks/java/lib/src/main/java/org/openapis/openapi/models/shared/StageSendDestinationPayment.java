/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package org.openapis.openapi.models.shared;

import com.fasterxml.jackson.annotation.JsonProperty;

public class StageSendDestinationPayment {
    @JsonProperty("psp")
    public String psp;

    public StageSendDestinationPayment withPsp(String psp) {
        this.psp = psp;
        return this;
    }
    
    public StageSendDestinationPayment(@JsonProperty("psp") String psp) {
        this.psp = psp;
  }
}