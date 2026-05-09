# ResetPasswordRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NewPassword** | **string** | The new password to set | 
**Token** | Pointer to **NullableString** | The token to reset the password | [optional] 

## Methods

### NewResetPasswordRequest

`func NewResetPasswordRequest(newPassword string, ) *ResetPasswordRequest`

NewResetPasswordRequest instantiates a new ResetPasswordRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewResetPasswordRequestWithDefaults

`func NewResetPasswordRequestWithDefaults() *ResetPasswordRequest`

NewResetPasswordRequestWithDefaults instantiates a new ResetPasswordRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNewPassword

`func (o *ResetPasswordRequest) GetNewPassword() string`

GetNewPassword returns the NewPassword field if non-nil, zero value otherwise.

### GetNewPasswordOk

`func (o *ResetPasswordRequest) GetNewPasswordOk() (*string, bool)`

GetNewPasswordOk returns a tuple with the NewPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNewPassword

`func (o *ResetPasswordRequest) SetNewPassword(v string)`

SetNewPassword sets NewPassword field to given value.


### GetToken

`func (o *ResetPasswordRequest) GetToken() string`

GetToken returns the Token field if non-nil, zero value otherwise.

### GetTokenOk

`func (o *ResetPasswordRequest) GetTokenOk() (*string, bool)`

GetTokenOk returns a tuple with the Token field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetToken

`func (o *ResetPasswordRequest) SetToken(v string)`

SetToken sets Token field to given value.

### HasToken

`func (o *ResetPasswordRequest) HasToken() bool`

HasToken returns a boolean if a field has been set.

### SetTokenNil

`func (o *ResetPasswordRequest) SetTokenNil(b bool)`

 SetTokenNil sets the value for Token to be an explicit nil

### UnsetToken
`func (o *ResetPasswordRequest) UnsetToken()`

UnsetToken ensures that no value is present for Token, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


