# ChangeEmailRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NewEmail** | **string** | The new email address to set must be a valid email address | 
**CallbackURL** | Pointer to **NullableString** | The URL to redirect to after email verification | [optional] 

## Methods

### NewChangeEmailRequest

`func NewChangeEmailRequest(newEmail string, ) *ChangeEmailRequest`

NewChangeEmailRequest instantiates a new ChangeEmailRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewChangeEmailRequestWithDefaults

`func NewChangeEmailRequestWithDefaults() *ChangeEmailRequest`

NewChangeEmailRequestWithDefaults instantiates a new ChangeEmailRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNewEmail

`func (o *ChangeEmailRequest) GetNewEmail() string`

GetNewEmail returns the NewEmail field if non-nil, zero value otherwise.

### GetNewEmailOk

`func (o *ChangeEmailRequest) GetNewEmailOk() (*string, bool)`

GetNewEmailOk returns a tuple with the NewEmail field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNewEmail

`func (o *ChangeEmailRequest) SetNewEmail(v string)`

SetNewEmail sets NewEmail field to given value.


### GetCallbackURL

`func (o *ChangeEmailRequest) GetCallbackURL() string`

GetCallbackURL returns the CallbackURL field if non-nil, zero value otherwise.

### GetCallbackURLOk

`func (o *ChangeEmailRequest) GetCallbackURLOk() (*string, bool)`

GetCallbackURLOk returns a tuple with the CallbackURL field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallbackURL

`func (o *ChangeEmailRequest) SetCallbackURL(v string)`

SetCallbackURL sets CallbackURL field to given value.

### HasCallbackURL

`func (o *ChangeEmailRequest) HasCallbackURL() bool`

HasCallbackURL returns a boolean if a field has been set.

### SetCallbackURLNil

`func (o *ChangeEmailRequest) SetCallbackURLNil(b bool)`

 SetCallbackURLNil sets the value for CallbackURL to be an explicit nil

### UnsetCallbackURL
`func (o *ChangeEmailRequest) UnsetCallbackURL()`

UnsetCallbackURL ensures that no value is present for CallbackURL, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


