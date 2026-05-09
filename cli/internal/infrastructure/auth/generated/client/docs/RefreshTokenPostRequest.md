# RefreshTokenPostRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ProviderId** | **string** | The provider ID for the OAuth provider | 
**AccountId** | Pointer to **NullableString** | The account ID associated with the refresh token | [optional] 
**UserId** | Pointer to **NullableString** | The user ID associated with the account | [optional] 

## Methods

### NewRefreshTokenPostRequest

`func NewRefreshTokenPostRequest(providerId string, ) *RefreshTokenPostRequest`

NewRefreshTokenPostRequest instantiates a new RefreshTokenPostRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRefreshTokenPostRequestWithDefaults

`func NewRefreshTokenPostRequestWithDefaults() *RefreshTokenPostRequest`

NewRefreshTokenPostRequestWithDefaults instantiates a new RefreshTokenPostRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetProviderId

`func (o *RefreshTokenPostRequest) GetProviderId() string`

GetProviderId returns the ProviderId field if non-nil, zero value otherwise.

### GetProviderIdOk

`func (o *RefreshTokenPostRequest) GetProviderIdOk() (*string, bool)`

GetProviderIdOk returns a tuple with the ProviderId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProviderId

`func (o *RefreshTokenPostRequest) SetProviderId(v string)`

SetProviderId sets ProviderId field to given value.


### GetAccountId

`func (o *RefreshTokenPostRequest) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *RefreshTokenPostRequest) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *RefreshTokenPostRequest) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.

### HasAccountId

`func (o *RefreshTokenPostRequest) HasAccountId() bool`

HasAccountId returns a boolean if a field has been set.

### SetAccountIdNil

`func (o *RefreshTokenPostRequest) SetAccountIdNil(b bool)`

 SetAccountIdNil sets the value for AccountId to be an explicit nil

### UnsetAccountId
`func (o *RefreshTokenPostRequest) UnsetAccountId()`

UnsetAccountId ensures that no value is present for AccountId, not even an explicit nil
### GetUserId

`func (o *RefreshTokenPostRequest) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *RefreshTokenPostRequest) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *RefreshTokenPostRequest) SetUserId(v string)`

SetUserId sets UserId field to given value.

### HasUserId

`func (o *RefreshTokenPostRequest) HasUserId() bool`

HasUserId returns a boolean if a field has been set.

### SetUserIdNil

`func (o *RefreshTokenPostRequest) SetUserIdNil(b bool)`

 SetUserIdNil sets the value for UserId to be an explicit nil

### UnsetUserId
`func (o *RefreshTokenPostRequest) UnsetUserId()`

UnsetUserId ensures that no value is present for UserId, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


