# UnlinkAccountPostRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ProviderId** | **string** |  | 
**AccountId** | Pointer to **NullableString** |  | [optional] 

## Methods

### NewUnlinkAccountPostRequest

`func NewUnlinkAccountPostRequest(providerId string, ) *UnlinkAccountPostRequest`

NewUnlinkAccountPostRequest instantiates a new UnlinkAccountPostRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUnlinkAccountPostRequestWithDefaults

`func NewUnlinkAccountPostRequestWithDefaults() *UnlinkAccountPostRequest`

NewUnlinkAccountPostRequestWithDefaults instantiates a new UnlinkAccountPostRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetProviderId

`func (o *UnlinkAccountPostRequest) GetProviderId() string`

GetProviderId returns the ProviderId field if non-nil, zero value otherwise.

### GetProviderIdOk

`func (o *UnlinkAccountPostRequest) GetProviderIdOk() (*string, bool)`

GetProviderIdOk returns a tuple with the ProviderId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProviderId

`func (o *UnlinkAccountPostRequest) SetProviderId(v string)`

SetProviderId sets ProviderId field to given value.


### GetAccountId

`func (o *UnlinkAccountPostRequest) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *UnlinkAccountPostRequest) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *UnlinkAccountPostRequest) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.

### HasAccountId

`func (o *UnlinkAccountPostRequest) HasAccountId() bool`

HasAccountId returns a boolean if a field has been set.

### SetAccountIdNil

`func (o *UnlinkAccountPostRequest) SetAccountIdNil(b bool)`

 SetAccountIdNil sets the value for AccountId to be an explicit nil

### UnsetAccountId
`func (o *UnlinkAccountPostRequest) UnsetAccountId()`

UnsetAccountId ensures that no value is present for AccountId, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


