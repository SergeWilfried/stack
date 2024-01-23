/*
Membership API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package membershipclient

import (
	"encoding/json"
)

// checks if the UpdateStackUserRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UpdateStackUserRequest{}

// UpdateStackUserRequest struct for UpdateStackUserRequest
type UpdateStackUserRequest struct {
	Role Role `json:"role"`
}

// NewUpdateStackUserRequest instantiates a new UpdateStackUserRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUpdateStackUserRequest(role Role) *UpdateStackUserRequest {
	this := UpdateStackUserRequest{}
	this.Role = role
	return &this
}

// NewUpdateStackUserRequestWithDefaults instantiates a new UpdateStackUserRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUpdateStackUserRequestWithDefaults() *UpdateStackUserRequest {
	this := UpdateStackUserRequest{}
	var role Role = EMPTY
	this.Role = role
	return &this
}

// GetRole returns the Role field value
func (o *UpdateStackUserRequest) GetRole() Role {
	if o == nil {
		var ret Role
		return ret
	}

	return o.Role
}

// GetRoleOk returns a tuple with the Role field value
// and a boolean to check if the value has been set.
func (o *UpdateStackUserRequest) GetRoleOk() (*Role, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Role, true
}

// SetRole sets field value
func (o *UpdateStackUserRequest) SetRole(v Role) {
	o.Role = v
}

func (o UpdateStackUserRequest) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UpdateStackUserRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["role"] = o.Role
	return toSerialize, nil
}

type NullableUpdateStackUserRequest struct {
	value *UpdateStackUserRequest
	isSet bool
}

func (v NullableUpdateStackUserRequest) Get() *UpdateStackUserRequest {
	return v.value
}

func (v *NullableUpdateStackUserRequest) Set(val *UpdateStackUserRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableUpdateStackUserRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableUpdateStackUserRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUpdateStackUserRequest(val *UpdateStackUserRequest) *NullableUpdateStackUserRequest {
	return &NullableUpdateStackUserRequest{value: val, isSet: true}
}

func (v NullableUpdateStackUserRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUpdateStackUserRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


