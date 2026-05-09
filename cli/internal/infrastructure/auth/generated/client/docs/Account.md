# Account

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**AccountId** | **string** |  | 
**ProviderId** | **string** |  | 
**UserId** | **string** |  | 
**AccessToken** | Pointer to **string** |  | [optional] 
**RefreshToken** | Pointer to **string** |  | [optional] 
**IdToken** | Pointer to **string** |  | [optional] 
**AccessTokenExpiresAt** | Pointer to **time.Time** |  | [optional] 
**RefreshTokenExpiresAt** | Pointer to **time.Time** |  | [optional] 
**Scope** | Pointer to **string** |  | [optional] 
**Password** | Pointer to **string** |  | [optional] 
**CreatedAt** | **time.Time** |  | [default to "Generated at runtime"]
**UpdatedAt** | **time.Time** |  | 

## Methods

### NewAccount

`func NewAccount(accountId string, providerId string, userId string, createdAt time.Time, updatedAt time.Time, ) *Account`

NewAccount instantiates a new Account object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAccountWithDefaults

`func NewAccountWithDefaults() *Account`

NewAccountWithDefaults instantiates a new Account object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Account) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Account) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Account) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *Account) HasId() bool`

HasId returns a boolean if a field has been set.

### GetAccountId

`func (o *Account) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *Account) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *Account) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.


### GetProviderId

`func (o *Account) GetProviderId() string`

GetProviderId returns the ProviderId field if non-nil, zero value otherwise.

### GetProviderIdOk

`func (o *Account) GetProviderIdOk() (*string, bool)`

GetProviderIdOk returns a tuple with the ProviderId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProviderId

`func (o *Account) SetProviderId(v string)`

SetProviderId sets ProviderId field to given value.


### GetUserId

`func (o *Account) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *Account) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *Account) SetUserId(v string)`

SetUserId sets UserId field to given value.


### GetAccessToken

`func (o *Account) GetAccessToken() string`

GetAccessToken returns the AccessToken field if non-nil, zero value otherwise.

### GetAccessTokenOk

`func (o *Account) GetAccessTokenOk() (*string, bool)`

GetAccessTokenOk returns a tuple with the AccessToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessToken

`func (o *Account) SetAccessToken(v string)`

SetAccessToken sets AccessToken field to given value.

### HasAccessToken

`func (o *Account) HasAccessToken() bool`

HasAccessToken returns a boolean if a field has been set.

### GetRefreshToken

`func (o *Account) GetRefreshToken() string`

GetRefreshToken returns the RefreshToken field if non-nil, zero value otherwise.

### GetRefreshTokenOk

`func (o *Account) GetRefreshTokenOk() (*string, bool)`

GetRefreshTokenOk returns a tuple with the RefreshToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshToken

`func (o *Account) SetRefreshToken(v string)`

SetRefreshToken sets RefreshToken field to given value.

### HasRefreshToken

`func (o *Account) HasRefreshToken() bool`

HasRefreshToken returns a boolean if a field has been set.

### GetIdToken

`func (o *Account) GetIdToken() string`

GetIdToken returns the IdToken field if non-nil, zero value otherwise.

### GetIdTokenOk

`func (o *Account) GetIdTokenOk() (*string, bool)`

GetIdTokenOk returns a tuple with the IdToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIdToken

`func (o *Account) SetIdToken(v string)`

SetIdToken sets IdToken field to given value.

### HasIdToken

`func (o *Account) HasIdToken() bool`

HasIdToken returns a boolean if a field has been set.

### GetAccessTokenExpiresAt

`func (o *Account) GetAccessTokenExpiresAt() time.Time`

GetAccessTokenExpiresAt returns the AccessTokenExpiresAt field if non-nil, zero value otherwise.

### GetAccessTokenExpiresAtOk

`func (o *Account) GetAccessTokenExpiresAtOk() (*time.Time, bool)`

GetAccessTokenExpiresAtOk returns a tuple with the AccessTokenExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessTokenExpiresAt

`func (o *Account) SetAccessTokenExpiresAt(v time.Time)`

SetAccessTokenExpiresAt sets AccessTokenExpiresAt field to given value.

### HasAccessTokenExpiresAt

`func (o *Account) HasAccessTokenExpiresAt() bool`

HasAccessTokenExpiresAt returns a boolean if a field has been set.

### GetRefreshTokenExpiresAt

`func (o *Account) GetRefreshTokenExpiresAt() time.Time`

GetRefreshTokenExpiresAt returns the RefreshTokenExpiresAt field if non-nil, zero value otherwise.

### GetRefreshTokenExpiresAtOk

`func (o *Account) GetRefreshTokenExpiresAtOk() (*time.Time, bool)`

GetRefreshTokenExpiresAtOk returns a tuple with the RefreshTokenExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshTokenExpiresAt

`func (o *Account) SetRefreshTokenExpiresAt(v time.Time)`

SetRefreshTokenExpiresAt sets RefreshTokenExpiresAt field to given value.

### HasRefreshTokenExpiresAt

`func (o *Account) HasRefreshTokenExpiresAt() bool`

HasRefreshTokenExpiresAt returns a boolean if a field has been set.

### GetScope

`func (o *Account) GetScope() string`

GetScope returns the Scope field if non-nil, zero value otherwise.

### GetScopeOk

`func (o *Account) GetScopeOk() (*string, bool)`

GetScopeOk returns a tuple with the Scope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScope

`func (o *Account) SetScope(v string)`

SetScope sets Scope field to given value.

### HasScope

`func (o *Account) HasScope() bool`

HasScope returns a boolean if a field has been set.

### GetPassword

`func (o *Account) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *Account) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *Account) SetPassword(v string)`

SetPassword sets Password field to given value.

### HasPassword

`func (o *Account) HasPassword() bool`

HasPassword returns a boolean if a field has been set.

### GetCreatedAt

`func (o *Account) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *Account) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *Account) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetUpdatedAt

`func (o *Account) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *Account) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *Account) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


