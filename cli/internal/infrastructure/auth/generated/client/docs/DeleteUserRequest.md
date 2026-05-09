# DeleteUserRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CallbackURL** | Pointer to **string** | The callback URL to redirect to after the user is deleted | [optional] 
**Password** | Pointer to **string** | The user&#39;s password. Required if session is not fresh | [optional] 
**Token** | Pointer to **string** | The deletion verification token | [optional] 

## Methods

### NewDeleteUserRequest

`func NewDeleteUserRequest() *DeleteUserRequest`

NewDeleteUserRequest instantiates a new DeleteUserRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeleteUserRequestWithDefaults

`func NewDeleteUserRequestWithDefaults() *DeleteUserRequest`

NewDeleteUserRequestWithDefaults instantiates a new DeleteUserRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCallbackURL

`func (o *DeleteUserRequest) GetCallbackURL() string`

GetCallbackURL returns the CallbackURL field if non-nil, zero value otherwise.

### GetCallbackURLOk

`func (o *DeleteUserRequest) GetCallbackURLOk() (*string, bool)`

GetCallbackURLOk returns a tuple with the CallbackURL field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallbackURL

`func (o *DeleteUserRequest) SetCallbackURL(v string)`

SetCallbackURL sets CallbackURL field to given value.

### HasCallbackURL

`func (o *DeleteUserRequest) HasCallbackURL() bool`

HasCallbackURL returns a boolean if a field has been set.

### GetPassword

`func (o *DeleteUserRequest) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *DeleteUserRequest) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *DeleteUserRequest) SetPassword(v string)`

SetPassword sets Password field to given value.

### HasPassword

`func (o *DeleteUserRequest) HasPassword() bool`

HasPassword returns a boolean if a field has been set.

### GetToken

`func (o *DeleteUserRequest) GetToken() string`

GetToken returns the Token field if non-nil, zero value otherwise.

### GetTokenOk

`func (o *DeleteUserRequest) GetTokenOk() (*string, bool)`

GetTokenOk returns a tuple with the Token field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetToken

`func (o *DeleteUserRequest) SetToken(v string)`

SetToken sets Token field to given value.

### HasToken

`func (o *DeleteUserRequest) HasToken() bool`

HasToken returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


