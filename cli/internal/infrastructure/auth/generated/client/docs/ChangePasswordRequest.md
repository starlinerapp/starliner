# ChangePasswordRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NewPassword** | **string** | The new password to set | 
**CurrentPassword** | **string** | The current password is required | 
**RevokeOtherSessions** | Pointer to **NullableBool** | Must be a boolean value | [optional] 

## Methods

### NewChangePasswordRequest

`func NewChangePasswordRequest(newPassword string, currentPassword string, ) *ChangePasswordRequest`

NewChangePasswordRequest instantiates a new ChangePasswordRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewChangePasswordRequestWithDefaults

`func NewChangePasswordRequestWithDefaults() *ChangePasswordRequest`

NewChangePasswordRequestWithDefaults instantiates a new ChangePasswordRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNewPassword

`func (o *ChangePasswordRequest) GetNewPassword() string`

GetNewPassword returns the NewPassword field if non-nil, zero value otherwise.

### GetNewPasswordOk

`func (o *ChangePasswordRequest) GetNewPasswordOk() (*string, bool)`

GetNewPasswordOk returns a tuple with the NewPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNewPassword

`func (o *ChangePasswordRequest) SetNewPassword(v string)`

SetNewPassword sets NewPassword field to given value.


### GetCurrentPassword

`func (o *ChangePasswordRequest) GetCurrentPassword() string`

GetCurrentPassword returns the CurrentPassword field if non-nil, zero value otherwise.

### GetCurrentPasswordOk

`func (o *ChangePasswordRequest) GetCurrentPasswordOk() (*string, bool)`

GetCurrentPasswordOk returns a tuple with the CurrentPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrentPassword

`func (o *ChangePasswordRequest) SetCurrentPassword(v string)`

SetCurrentPassword sets CurrentPassword field to given value.


### GetRevokeOtherSessions

`func (o *ChangePasswordRequest) GetRevokeOtherSessions() bool`

GetRevokeOtherSessions returns the RevokeOtherSessions field if non-nil, zero value otherwise.

### GetRevokeOtherSessionsOk

`func (o *ChangePasswordRequest) GetRevokeOtherSessionsOk() (*bool, bool)`

GetRevokeOtherSessionsOk returns a tuple with the RevokeOtherSessions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRevokeOtherSessions

`func (o *ChangePasswordRequest) SetRevokeOtherSessions(v bool)`

SetRevokeOtherSessions sets RevokeOtherSessions field to given value.

### HasRevokeOtherSessions

`func (o *ChangePasswordRequest) HasRevokeOtherSessions() bool`

HasRevokeOtherSessions returns a boolean if a field has been set.

### SetRevokeOtherSessionsNil

`func (o *ChangePasswordRequest) SetRevokeOtherSessionsNil(b bool)`

 SetRevokeOtherSessionsNil sets the value for RevokeOtherSessions to be an explicit nil

### UnsetRevokeOtherSessions
`func (o *ChangePasswordRequest) UnsetRevokeOtherSessions()`

UnsetRevokeOtherSessions ensures that no value is present for RevokeOtherSessions, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


