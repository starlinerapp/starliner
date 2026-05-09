# LinkSocialAccountRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CallbackURL** | Pointer to **NullableString** | The URL to redirect to after the user has signed in | [optional] 
**Provider** | **string** |  | 
**IdToken** | Pointer to [**LinkSocialAccountRequestIdToken**](LinkSocialAccountRequestIdToken.md) |  | [optional] 
**RequestSignUp** | Pointer to **NullableBool** |  | [optional] 
**Scopes** | Pointer to **[]interface{}** | Additional scopes to request from the provider | [optional] 
**ErrorCallbackURL** | Pointer to **NullableString** | The URL to redirect to if there is an error during the link process | [optional] 
**DisableRedirect** | Pointer to **NullableBool** | Disable automatic redirection to the provider. Useful for handling the redirection yourself | [optional] 
**AdditionalData** | Pointer to **NullableString** |  | [optional] 

## Methods

### NewLinkSocialAccountRequest

`func NewLinkSocialAccountRequest(provider string, ) *LinkSocialAccountRequest`

NewLinkSocialAccountRequest instantiates a new LinkSocialAccountRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLinkSocialAccountRequestWithDefaults

`func NewLinkSocialAccountRequestWithDefaults() *LinkSocialAccountRequest`

NewLinkSocialAccountRequestWithDefaults instantiates a new LinkSocialAccountRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCallbackURL

`func (o *LinkSocialAccountRequest) GetCallbackURL() string`

GetCallbackURL returns the CallbackURL field if non-nil, zero value otherwise.

### GetCallbackURLOk

`func (o *LinkSocialAccountRequest) GetCallbackURLOk() (*string, bool)`

GetCallbackURLOk returns a tuple with the CallbackURL field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallbackURL

`func (o *LinkSocialAccountRequest) SetCallbackURL(v string)`

SetCallbackURL sets CallbackURL field to given value.

### HasCallbackURL

`func (o *LinkSocialAccountRequest) HasCallbackURL() bool`

HasCallbackURL returns a boolean if a field has been set.

### SetCallbackURLNil

`func (o *LinkSocialAccountRequest) SetCallbackURLNil(b bool)`

 SetCallbackURLNil sets the value for CallbackURL to be an explicit nil

### UnsetCallbackURL
`func (o *LinkSocialAccountRequest) UnsetCallbackURL()`

UnsetCallbackURL ensures that no value is present for CallbackURL, not even an explicit nil
### GetProvider

`func (o *LinkSocialAccountRequest) GetProvider() string`

GetProvider returns the Provider field if non-nil, zero value otherwise.

### GetProviderOk

`func (o *LinkSocialAccountRequest) GetProviderOk() (*string, bool)`

GetProviderOk returns a tuple with the Provider field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProvider

`func (o *LinkSocialAccountRequest) SetProvider(v string)`

SetProvider sets Provider field to given value.


### GetIdToken

`func (o *LinkSocialAccountRequest) GetIdToken() LinkSocialAccountRequestIdToken`

GetIdToken returns the IdToken field if non-nil, zero value otherwise.

### GetIdTokenOk

`func (o *LinkSocialAccountRequest) GetIdTokenOk() (*LinkSocialAccountRequestIdToken, bool)`

GetIdTokenOk returns a tuple with the IdToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIdToken

`func (o *LinkSocialAccountRequest) SetIdToken(v LinkSocialAccountRequestIdToken)`

SetIdToken sets IdToken field to given value.

### HasIdToken

`func (o *LinkSocialAccountRequest) HasIdToken() bool`

HasIdToken returns a boolean if a field has been set.

### GetRequestSignUp

`func (o *LinkSocialAccountRequest) GetRequestSignUp() bool`

GetRequestSignUp returns the RequestSignUp field if non-nil, zero value otherwise.

### GetRequestSignUpOk

`func (o *LinkSocialAccountRequest) GetRequestSignUpOk() (*bool, bool)`

GetRequestSignUpOk returns a tuple with the RequestSignUp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestSignUp

`func (o *LinkSocialAccountRequest) SetRequestSignUp(v bool)`

SetRequestSignUp sets RequestSignUp field to given value.

### HasRequestSignUp

`func (o *LinkSocialAccountRequest) HasRequestSignUp() bool`

HasRequestSignUp returns a boolean if a field has been set.

### SetRequestSignUpNil

`func (o *LinkSocialAccountRequest) SetRequestSignUpNil(b bool)`

 SetRequestSignUpNil sets the value for RequestSignUp to be an explicit nil

### UnsetRequestSignUp
`func (o *LinkSocialAccountRequest) UnsetRequestSignUp()`

UnsetRequestSignUp ensures that no value is present for RequestSignUp, not even an explicit nil
### GetScopes

`func (o *LinkSocialAccountRequest) GetScopes() []interface{}`

GetScopes returns the Scopes field if non-nil, zero value otherwise.

### GetScopesOk

`func (o *LinkSocialAccountRequest) GetScopesOk() (*[]interface{}, bool)`

GetScopesOk returns a tuple with the Scopes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScopes

`func (o *LinkSocialAccountRequest) SetScopes(v []interface{})`

SetScopes sets Scopes field to given value.

### HasScopes

`func (o *LinkSocialAccountRequest) HasScopes() bool`

HasScopes returns a boolean if a field has been set.

### SetScopesNil

`func (o *LinkSocialAccountRequest) SetScopesNil(b bool)`

 SetScopesNil sets the value for Scopes to be an explicit nil

### UnsetScopes
`func (o *LinkSocialAccountRequest) UnsetScopes()`

UnsetScopes ensures that no value is present for Scopes, not even an explicit nil
### GetErrorCallbackURL

`func (o *LinkSocialAccountRequest) GetErrorCallbackURL() string`

GetErrorCallbackURL returns the ErrorCallbackURL field if non-nil, zero value otherwise.

### GetErrorCallbackURLOk

`func (o *LinkSocialAccountRequest) GetErrorCallbackURLOk() (*string, bool)`

GetErrorCallbackURLOk returns a tuple with the ErrorCallbackURL field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorCallbackURL

`func (o *LinkSocialAccountRequest) SetErrorCallbackURL(v string)`

SetErrorCallbackURL sets ErrorCallbackURL field to given value.

### HasErrorCallbackURL

`func (o *LinkSocialAccountRequest) HasErrorCallbackURL() bool`

HasErrorCallbackURL returns a boolean if a field has been set.

### SetErrorCallbackURLNil

`func (o *LinkSocialAccountRequest) SetErrorCallbackURLNil(b bool)`

 SetErrorCallbackURLNil sets the value for ErrorCallbackURL to be an explicit nil

### UnsetErrorCallbackURL
`func (o *LinkSocialAccountRequest) UnsetErrorCallbackURL()`

UnsetErrorCallbackURL ensures that no value is present for ErrorCallbackURL, not even an explicit nil
### GetDisableRedirect

`func (o *LinkSocialAccountRequest) GetDisableRedirect() bool`

GetDisableRedirect returns the DisableRedirect field if non-nil, zero value otherwise.

### GetDisableRedirectOk

`func (o *LinkSocialAccountRequest) GetDisableRedirectOk() (*bool, bool)`

GetDisableRedirectOk returns a tuple with the DisableRedirect field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisableRedirect

`func (o *LinkSocialAccountRequest) SetDisableRedirect(v bool)`

SetDisableRedirect sets DisableRedirect field to given value.

### HasDisableRedirect

`func (o *LinkSocialAccountRequest) HasDisableRedirect() bool`

HasDisableRedirect returns a boolean if a field has been set.

### SetDisableRedirectNil

`func (o *LinkSocialAccountRequest) SetDisableRedirectNil(b bool)`

 SetDisableRedirectNil sets the value for DisableRedirect to be an explicit nil

### UnsetDisableRedirect
`func (o *LinkSocialAccountRequest) UnsetDisableRedirect()`

UnsetDisableRedirect ensures that no value is present for DisableRedirect, not even an explicit nil
### GetAdditionalData

`func (o *LinkSocialAccountRequest) GetAdditionalData() string`

GetAdditionalData returns the AdditionalData field if non-nil, zero value otherwise.

### GetAdditionalDataOk

`func (o *LinkSocialAccountRequest) GetAdditionalDataOk() (*string, bool)`

GetAdditionalDataOk returns a tuple with the AdditionalData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAdditionalData

`func (o *LinkSocialAccountRequest) SetAdditionalData(v string)`

SetAdditionalData sets AdditionalData field to given value.

### HasAdditionalData

`func (o *LinkSocialAccountRequest) HasAdditionalData() bool`

HasAdditionalData returns a boolean if a field has been set.

### SetAdditionalDataNil

`func (o *LinkSocialAccountRequest) SetAdditionalDataNil(b bool)`

 SetAdditionalDataNil sets the value for AdditionalData to be an explicit nil

### UnsetAdditionalData
`func (o *LinkSocialAccountRequest) UnsetAdditionalData()`

UnsetAdditionalData ensures that no value is present for AdditionalData, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


