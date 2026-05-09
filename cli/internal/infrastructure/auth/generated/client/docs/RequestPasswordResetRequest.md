# RequestPasswordResetRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | **string** | The email address of the user to send a password reset email to | 
**RedirectTo** | Pointer to **NullableString** | The URL to redirect the user to reset their password. If the token isn&#39;t valid or expired, it&#39;ll be redirected with a query parameter &#x60;?error&#x3D;INVALID_TOKEN&#x60;. If the token is valid, it&#39;ll be redirected with a query parameter &#x60;?token&#x3D;VALID_TOKEN | [optional] 

## Methods

### NewRequestPasswordResetRequest

`func NewRequestPasswordResetRequest(email string, ) *RequestPasswordResetRequest`

NewRequestPasswordResetRequest instantiates a new RequestPasswordResetRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRequestPasswordResetRequestWithDefaults

`func NewRequestPasswordResetRequestWithDefaults() *RequestPasswordResetRequest`

NewRequestPasswordResetRequestWithDefaults instantiates a new RequestPasswordResetRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *RequestPasswordResetRequest) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *RequestPasswordResetRequest) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *RequestPasswordResetRequest) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetRedirectTo

`func (o *RequestPasswordResetRequest) GetRedirectTo() string`

GetRedirectTo returns the RedirectTo field if non-nil, zero value otherwise.

### GetRedirectToOk

`func (o *RequestPasswordResetRequest) GetRedirectToOk() (*string, bool)`

GetRedirectToOk returns a tuple with the RedirectTo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirectTo

`func (o *RequestPasswordResetRequest) SetRedirectTo(v string)`

SetRedirectTo sets RedirectTo field to given value.

### HasRedirectTo

`func (o *RequestPasswordResetRequest) HasRedirectTo() bool`

HasRedirectTo returns a boolean if a field has been set.

### SetRedirectToNil

`func (o *RequestPasswordResetRequest) SetRedirectToNil(b bool)`

 SetRedirectToNil sets the value for RedirectTo to be an explicit nil

### UnsetRedirectTo
`func (o *RequestPasswordResetRequest) UnsetRedirectTo()`

UnsetRedirectTo ensures that no value is present for RedirectTo, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


