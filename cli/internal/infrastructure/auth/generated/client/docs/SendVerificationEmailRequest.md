# SendVerificationEmailRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | **string** | The email to send the verification email to | 
**CallbackURL** | Pointer to **string** | The URL to use for email verification callback | [optional] 

## Methods

### NewSendVerificationEmailRequest

`func NewSendVerificationEmailRequest(email string, ) *SendVerificationEmailRequest`

NewSendVerificationEmailRequest instantiates a new SendVerificationEmailRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSendVerificationEmailRequestWithDefaults

`func NewSendVerificationEmailRequestWithDefaults() *SendVerificationEmailRequest`

NewSendVerificationEmailRequestWithDefaults instantiates a new SendVerificationEmailRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *SendVerificationEmailRequest) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *SendVerificationEmailRequest) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *SendVerificationEmailRequest) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetCallbackURL

`func (o *SendVerificationEmailRequest) GetCallbackURL() string`

GetCallbackURL returns the CallbackURL field if non-nil, zero value otherwise.

### GetCallbackURLOk

`func (o *SendVerificationEmailRequest) GetCallbackURLOk() (*string, bool)`

GetCallbackURLOk returns a tuple with the CallbackURL field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallbackURL

`func (o *SendVerificationEmailRequest) SetCallbackURL(v string)`

SetCallbackURL sets CallbackURL field to given value.

### HasCallbackURL

`func (o *SendVerificationEmailRequest) HasCallbackURL() bool`

HasCallbackURL returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


