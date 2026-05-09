# SocialSignInRequestIdToken

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Token** | **string** | ID token from the provider | 
**Nonce** | Pointer to **NullableString** | Nonce used to generate the token | [optional] 
**AccessToken** | Pointer to **NullableString** | Access token from the provider | [optional] 
**RefreshToken** | Pointer to **NullableString** | Refresh token from the provider | [optional] 
**ExpiresAt** | Pointer to **NullableFloat32** | Expiry date of the token | [optional] 
**User** | Pointer to [**SocialSignInRequestIdTokenUser**](SocialSignInRequestIdTokenUser.md) |  | [optional] 

## Methods

### NewSocialSignInRequestIdToken

`func NewSocialSignInRequestIdToken(token string, ) *SocialSignInRequestIdToken`

NewSocialSignInRequestIdToken instantiates a new SocialSignInRequestIdToken object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSocialSignInRequestIdTokenWithDefaults

`func NewSocialSignInRequestIdTokenWithDefaults() *SocialSignInRequestIdToken`

NewSocialSignInRequestIdTokenWithDefaults instantiates a new SocialSignInRequestIdToken object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetToken

`func (o *SocialSignInRequestIdToken) GetToken() string`

GetToken returns the Token field if non-nil, zero value otherwise.

### GetTokenOk

`func (o *SocialSignInRequestIdToken) GetTokenOk() (*string, bool)`

GetTokenOk returns a tuple with the Token field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetToken

`func (o *SocialSignInRequestIdToken) SetToken(v string)`

SetToken sets Token field to given value.


### GetNonce

`func (o *SocialSignInRequestIdToken) GetNonce() string`

GetNonce returns the Nonce field if non-nil, zero value otherwise.

### GetNonceOk

`func (o *SocialSignInRequestIdToken) GetNonceOk() (*string, bool)`

GetNonceOk returns a tuple with the Nonce field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNonce

`func (o *SocialSignInRequestIdToken) SetNonce(v string)`

SetNonce sets Nonce field to given value.

### HasNonce

`func (o *SocialSignInRequestIdToken) HasNonce() bool`

HasNonce returns a boolean if a field has been set.

### SetNonceNil

`func (o *SocialSignInRequestIdToken) SetNonceNil(b bool)`

 SetNonceNil sets the value for Nonce to be an explicit nil

### UnsetNonce
`func (o *SocialSignInRequestIdToken) UnsetNonce()`

UnsetNonce ensures that no value is present for Nonce, not even an explicit nil
### GetAccessToken

`func (o *SocialSignInRequestIdToken) GetAccessToken() string`

GetAccessToken returns the AccessToken field if non-nil, zero value otherwise.

### GetAccessTokenOk

`func (o *SocialSignInRequestIdToken) GetAccessTokenOk() (*string, bool)`

GetAccessTokenOk returns a tuple with the AccessToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessToken

`func (o *SocialSignInRequestIdToken) SetAccessToken(v string)`

SetAccessToken sets AccessToken field to given value.

### HasAccessToken

`func (o *SocialSignInRequestIdToken) HasAccessToken() bool`

HasAccessToken returns a boolean if a field has been set.

### SetAccessTokenNil

`func (o *SocialSignInRequestIdToken) SetAccessTokenNil(b bool)`

 SetAccessTokenNil sets the value for AccessToken to be an explicit nil

### UnsetAccessToken
`func (o *SocialSignInRequestIdToken) UnsetAccessToken()`

UnsetAccessToken ensures that no value is present for AccessToken, not even an explicit nil
### GetRefreshToken

`func (o *SocialSignInRequestIdToken) GetRefreshToken() string`

GetRefreshToken returns the RefreshToken field if non-nil, zero value otherwise.

### GetRefreshTokenOk

`func (o *SocialSignInRequestIdToken) GetRefreshTokenOk() (*string, bool)`

GetRefreshTokenOk returns a tuple with the RefreshToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshToken

`func (o *SocialSignInRequestIdToken) SetRefreshToken(v string)`

SetRefreshToken sets RefreshToken field to given value.

### HasRefreshToken

`func (o *SocialSignInRequestIdToken) HasRefreshToken() bool`

HasRefreshToken returns a boolean if a field has been set.

### SetRefreshTokenNil

`func (o *SocialSignInRequestIdToken) SetRefreshTokenNil(b bool)`

 SetRefreshTokenNil sets the value for RefreshToken to be an explicit nil

### UnsetRefreshToken
`func (o *SocialSignInRequestIdToken) UnsetRefreshToken()`

UnsetRefreshToken ensures that no value is present for RefreshToken, not even an explicit nil
### GetExpiresAt

`func (o *SocialSignInRequestIdToken) GetExpiresAt() float32`

GetExpiresAt returns the ExpiresAt field if non-nil, zero value otherwise.

### GetExpiresAtOk

`func (o *SocialSignInRequestIdToken) GetExpiresAtOk() (*float32, bool)`

GetExpiresAtOk returns a tuple with the ExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresAt

`func (o *SocialSignInRequestIdToken) SetExpiresAt(v float32)`

SetExpiresAt sets ExpiresAt field to given value.

### HasExpiresAt

`func (o *SocialSignInRequestIdToken) HasExpiresAt() bool`

HasExpiresAt returns a boolean if a field has been set.

### SetExpiresAtNil

`func (o *SocialSignInRequestIdToken) SetExpiresAtNil(b bool)`

 SetExpiresAtNil sets the value for ExpiresAt to be an explicit nil

### UnsetExpiresAt
`func (o *SocialSignInRequestIdToken) UnsetExpiresAt()`

UnsetExpiresAt ensures that no value is present for ExpiresAt, not even an explicit nil
### GetUser

`func (o *SocialSignInRequestIdToken) GetUser() SocialSignInRequestIdTokenUser`

GetUser returns the User field if non-nil, zero value otherwise.

### GetUserOk

`func (o *SocialSignInRequestIdToken) GetUserOk() (*SocialSignInRequestIdTokenUser, bool)`

GetUserOk returns a tuple with the User field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUser

`func (o *SocialSignInRequestIdToken) SetUser(v SocialSignInRequestIdTokenUser)`

SetUser sets User field to given value.

### HasUser

`func (o *SocialSignInRequestIdToken) HasUser() bool`

HasUser returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


