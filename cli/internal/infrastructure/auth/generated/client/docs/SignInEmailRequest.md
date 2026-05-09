# SignInEmailRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | **string** | Email of the user | 
**Password** | **string** | Password of the user | 
**CallbackURL** | Pointer to **NullableString** | Callback URL to use as a redirect for email verification | [optional] 
**RememberMe** | Pointer to **NullableBool** | If this is false, the session will not be remembered. Default is &#x60;true&#x60;. | [optional] [default to true]

## Methods

### NewSignInEmailRequest

`func NewSignInEmailRequest(email string, password string, ) *SignInEmailRequest`

NewSignInEmailRequest instantiates a new SignInEmailRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSignInEmailRequestWithDefaults

`func NewSignInEmailRequestWithDefaults() *SignInEmailRequest`

NewSignInEmailRequestWithDefaults instantiates a new SignInEmailRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *SignInEmailRequest) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *SignInEmailRequest) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *SignInEmailRequest) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetPassword

`func (o *SignInEmailRequest) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *SignInEmailRequest) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *SignInEmailRequest) SetPassword(v string)`

SetPassword sets Password field to given value.


### GetCallbackURL

`func (o *SignInEmailRequest) GetCallbackURL() string`

GetCallbackURL returns the CallbackURL field if non-nil, zero value otherwise.

### GetCallbackURLOk

`func (o *SignInEmailRequest) GetCallbackURLOk() (*string, bool)`

GetCallbackURLOk returns a tuple with the CallbackURL field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallbackURL

`func (o *SignInEmailRequest) SetCallbackURL(v string)`

SetCallbackURL sets CallbackURL field to given value.

### HasCallbackURL

`func (o *SignInEmailRequest) HasCallbackURL() bool`

HasCallbackURL returns a boolean if a field has been set.

### SetCallbackURLNil

`func (o *SignInEmailRequest) SetCallbackURLNil(b bool)`

 SetCallbackURLNil sets the value for CallbackURL to be an explicit nil

### UnsetCallbackURL
`func (o *SignInEmailRequest) UnsetCallbackURL()`

UnsetCallbackURL ensures that no value is present for CallbackURL, not even an explicit nil
### GetRememberMe

`func (o *SignInEmailRequest) GetRememberMe() bool`

GetRememberMe returns the RememberMe field if non-nil, zero value otherwise.

### GetRememberMeOk

`func (o *SignInEmailRequest) GetRememberMeOk() (*bool, bool)`

GetRememberMeOk returns a tuple with the RememberMe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRememberMe

`func (o *SignInEmailRequest) SetRememberMe(v bool)`

SetRememberMe sets RememberMe field to given value.

### HasRememberMe

`func (o *SignInEmailRequest) HasRememberMe() bool`

HasRememberMe returns a boolean if a field has been set.

### SetRememberMeNil

`func (o *SignInEmailRequest) SetRememberMeNil(b bool)`

 SetRememberMeNil sets the value for RememberMe to be an explicit nil

### UnsetRememberMe
`func (o *SignInEmailRequest) UnsetRememberMe()`

UnsetRememberMe ensures that no value is present for RememberMe, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


