# SocialSignInRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CallbackURL** | Pointer to **NullableString** | Callback URL to redirect to after the user has signed in | [optional] 
**NewUserCallbackURL** | Pointer to **NullableString** |  | [optional] 
**ErrorCallbackURL** | Pointer to **NullableString** | Callback URL to redirect to if an error happens | [optional] 
**Provider** | **string** |  | 
**DisableRedirect** | Pointer to **NullableBool** | Disable automatic redirection to the provider. Useful for handling the redirection yourself | [optional] 
**IdToken** | Pointer to [**SocialSignInRequestIdToken**](SocialSignInRequestIdToken.md) |  | [optional] 
**Scopes** | Pointer to **[]interface{}** | Array of scopes to request from the provider. This will override the default scopes passed. | [optional] 
**RequestSignUp** | Pointer to **NullableBool** | Explicitly request sign-up. Useful when disableImplicitSignUp is true for this provider | [optional] 
**LoginHint** | Pointer to **NullableString** | The login hint to use for the authorization code request | [optional] 
**AdditionalData** | Pointer to **NullableString** |  | [optional] 

## Methods

### NewSocialSignInRequest

`func NewSocialSignInRequest(provider string, ) *SocialSignInRequest`

NewSocialSignInRequest instantiates a new SocialSignInRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSocialSignInRequestWithDefaults

`func NewSocialSignInRequestWithDefaults() *SocialSignInRequest`

NewSocialSignInRequestWithDefaults instantiates a new SocialSignInRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCallbackURL

`func (o *SocialSignInRequest) GetCallbackURL() string`

GetCallbackURL returns the CallbackURL field if non-nil, zero value otherwise.

### GetCallbackURLOk

`func (o *SocialSignInRequest) GetCallbackURLOk() (*string, bool)`

GetCallbackURLOk returns a tuple with the CallbackURL field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallbackURL

`func (o *SocialSignInRequest) SetCallbackURL(v string)`

SetCallbackURL sets CallbackURL field to given value.

### HasCallbackURL

`func (o *SocialSignInRequest) HasCallbackURL() bool`

HasCallbackURL returns a boolean if a field has been set.

### SetCallbackURLNil

`func (o *SocialSignInRequest) SetCallbackURLNil(b bool)`

 SetCallbackURLNil sets the value for CallbackURL to be an explicit nil

### UnsetCallbackURL
`func (o *SocialSignInRequest) UnsetCallbackURL()`

UnsetCallbackURL ensures that no value is present for CallbackURL, not even an explicit nil
### GetNewUserCallbackURL

`func (o *SocialSignInRequest) GetNewUserCallbackURL() string`

GetNewUserCallbackURL returns the NewUserCallbackURL field if non-nil, zero value otherwise.

### GetNewUserCallbackURLOk

`func (o *SocialSignInRequest) GetNewUserCallbackURLOk() (*string, bool)`

GetNewUserCallbackURLOk returns a tuple with the NewUserCallbackURL field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNewUserCallbackURL

`func (o *SocialSignInRequest) SetNewUserCallbackURL(v string)`

SetNewUserCallbackURL sets NewUserCallbackURL field to given value.

### HasNewUserCallbackURL

`func (o *SocialSignInRequest) HasNewUserCallbackURL() bool`

HasNewUserCallbackURL returns a boolean if a field has been set.

### SetNewUserCallbackURLNil

`func (o *SocialSignInRequest) SetNewUserCallbackURLNil(b bool)`

 SetNewUserCallbackURLNil sets the value for NewUserCallbackURL to be an explicit nil

### UnsetNewUserCallbackURL
`func (o *SocialSignInRequest) UnsetNewUserCallbackURL()`

UnsetNewUserCallbackURL ensures that no value is present for NewUserCallbackURL, not even an explicit nil
### GetErrorCallbackURL

`func (o *SocialSignInRequest) GetErrorCallbackURL() string`

GetErrorCallbackURL returns the ErrorCallbackURL field if non-nil, zero value otherwise.

### GetErrorCallbackURLOk

`func (o *SocialSignInRequest) GetErrorCallbackURLOk() (*string, bool)`

GetErrorCallbackURLOk returns a tuple with the ErrorCallbackURL field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorCallbackURL

`func (o *SocialSignInRequest) SetErrorCallbackURL(v string)`

SetErrorCallbackURL sets ErrorCallbackURL field to given value.

### HasErrorCallbackURL

`func (o *SocialSignInRequest) HasErrorCallbackURL() bool`

HasErrorCallbackURL returns a boolean if a field has been set.

### SetErrorCallbackURLNil

`func (o *SocialSignInRequest) SetErrorCallbackURLNil(b bool)`

 SetErrorCallbackURLNil sets the value for ErrorCallbackURL to be an explicit nil

### UnsetErrorCallbackURL
`func (o *SocialSignInRequest) UnsetErrorCallbackURL()`

UnsetErrorCallbackURL ensures that no value is present for ErrorCallbackURL, not even an explicit nil
### GetProvider

`func (o *SocialSignInRequest) GetProvider() string`

GetProvider returns the Provider field if non-nil, zero value otherwise.

### GetProviderOk

`func (o *SocialSignInRequest) GetProviderOk() (*string, bool)`

GetProviderOk returns a tuple with the Provider field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProvider

`func (o *SocialSignInRequest) SetProvider(v string)`

SetProvider sets Provider field to given value.


### GetDisableRedirect

`func (o *SocialSignInRequest) GetDisableRedirect() bool`

GetDisableRedirect returns the DisableRedirect field if non-nil, zero value otherwise.

### GetDisableRedirectOk

`func (o *SocialSignInRequest) GetDisableRedirectOk() (*bool, bool)`

GetDisableRedirectOk returns a tuple with the DisableRedirect field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisableRedirect

`func (o *SocialSignInRequest) SetDisableRedirect(v bool)`

SetDisableRedirect sets DisableRedirect field to given value.

### HasDisableRedirect

`func (o *SocialSignInRequest) HasDisableRedirect() bool`

HasDisableRedirect returns a boolean if a field has been set.

### SetDisableRedirectNil

`func (o *SocialSignInRequest) SetDisableRedirectNil(b bool)`

 SetDisableRedirectNil sets the value for DisableRedirect to be an explicit nil

### UnsetDisableRedirect
`func (o *SocialSignInRequest) UnsetDisableRedirect()`

UnsetDisableRedirect ensures that no value is present for DisableRedirect, not even an explicit nil
### GetIdToken

`func (o *SocialSignInRequest) GetIdToken() SocialSignInRequestIdToken`

GetIdToken returns the IdToken field if non-nil, zero value otherwise.

### GetIdTokenOk

`func (o *SocialSignInRequest) GetIdTokenOk() (*SocialSignInRequestIdToken, bool)`

GetIdTokenOk returns a tuple with the IdToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIdToken

`func (o *SocialSignInRequest) SetIdToken(v SocialSignInRequestIdToken)`

SetIdToken sets IdToken field to given value.

### HasIdToken

`func (o *SocialSignInRequest) HasIdToken() bool`

HasIdToken returns a boolean if a field has been set.

### GetScopes

`func (o *SocialSignInRequest) GetScopes() []interface{}`

GetScopes returns the Scopes field if non-nil, zero value otherwise.

### GetScopesOk

`func (o *SocialSignInRequest) GetScopesOk() (*[]interface{}, bool)`

GetScopesOk returns a tuple with the Scopes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScopes

`func (o *SocialSignInRequest) SetScopes(v []interface{})`

SetScopes sets Scopes field to given value.

### HasScopes

`func (o *SocialSignInRequest) HasScopes() bool`

HasScopes returns a boolean if a field has been set.

### SetScopesNil

`func (o *SocialSignInRequest) SetScopesNil(b bool)`

 SetScopesNil sets the value for Scopes to be an explicit nil

### UnsetScopes
`func (o *SocialSignInRequest) UnsetScopes()`

UnsetScopes ensures that no value is present for Scopes, not even an explicit nil
### GetRequestSignUp

`func (o *SocialSignInRequest) GetRequestSignUp() bool`

GetRequestSignUp returns the RequestSignUp field if non-nil, zero value otherwise.

### GetRequestSignUpOk

`func (o *SocialSignInRequest) GetRequestSignUpOk() (*bool, bool)`

GetRequestSignUpOk returns a tuple with the RequestSignUp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestSignUp

`func (o *SocialSignInRequest) SetRequestSignUp(v bool)`

SetRequestSignUp sets RequestSignUp field to given value.

### HasRequestSignUp

`func (o *SocialSignInRequest) HasRequestSignUp() bool`

HasRequestSignUp returns a boolean if a field has been set.

### SetRequestSignUpNil

`func (o *SocialSignInRequest) SetRequestSignUpNil(b bool)`

 SetRequestSignUpNil sets the value for RequestSignUp to be an explicit nil

### UnsetRequestSignUp
`func (o *SocialSignInRequest) UnsetRequestSignUp()`

UnsetRequestSignUp ensures that no value is present for RequestSignUp, not even an explicit nil
### GetLoginHint

`func (o *SocialSignInRequest) GetLoginHint() string`

GetLoginHint returns the LoginHint field if non-nil, zero value otherwise.

### GetLoginHintOk

`func (o *SocialSignInRequest) GetLoginHintOk() (*string, bool)`

GetLoginHintOk returns a tuple with the LoginHint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLoginHint

`func (o *SocialSignInRequest) SetLoginHint(v string)`

SetLoginHint sets LoginHint field to given value.

### HasLoginHint

`func (o *SocialSignInRequest) HasLoginHint() bool`

HasLoginHint returns a boolean if a field has been set.

### SetLoginHintNil

`func (o *SocialSignInRequest) SetLoginHintNil(b bool)`

 SetLoginHintNil sets the value for LoginHint to be an explicit nil

### UnsetLoginHint
`func (o *SocialSignInRequest) UnsetLoginHint()`

UnsetLoginHint ensures that no value is present for LoginHint, not even an explicit nil
### GetAdditionalData

`func (o *SocialSignInRequest) GetAdditionalData() string`

GetAdditionalData returns the AdditionalData field if non-nil, zero value otherwise.

### GetAdditionalDataOk

`func (o *SocialSignInRequest) GetAdditionalDataOk() (*string, bool)`

GetAdditionalDataOk returns a tuple with the AdditionalData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAdditionalData

`func (o *SocialSignInRequest) SetAdditionalData(v string)`

SetAdditionalData sets AdditionalData field to given value.

### HasAdditionalData

`func (o *SocialSignInRequest) HasAdditionalData() bool`

HasAdditionalData returns a boolean if a field has been set.

### SetAdditionalDataNil

`func (o *SocialSignInRequest) SetAdditionalDataNil(b bool)`

 SetAdditionalDataNil sets the value for AdditionalData to be an explicit nil

### UnsetAdditionalData
`func (o *SocialSignInRequest) UnsetAdditionalData()`

UnsetAdditionalData ensures that no value is present for AdditionalData, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


