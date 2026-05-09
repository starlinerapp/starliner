# AccountInfoGet200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**User** | [**AccountInfoGet200ResponseUser**](AccountInfoGet200ResponseUser.md) |  | 
**Data** | **map[string]interface{}** |  | 

## Methods

### NewAccountInfoGet200Response

`func NewAccountInfoGet200Response(user AccountInfoGet200ResponseUser, data map[string]interface{}, ) *AccountInfoGet200Response`

NewAccountInfoGet200Response instantiates a new AccountInfoGet200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAccountInfoGet200ResponseWithDefaults

`func NewAccountInfoGet200ResponseWithDefaults() *AccountInfoGet200Response`

NewAccountInfoGet200ResponseWithDefaults instantiates a new AccountInfoGet200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUser

`func (o *AccountInfoGet200Response) GetUser() AccountInfoGet200ResponseUser`

GetUser returns the User field if non-nil, zero value otherwise.

### GetUserOk

`func (o *AccountInfoGet200Response) GetUserOk() (*AccountInfoGet200ResponseUser, bool)`

GetUserOk returns a tuple with the User field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUser

`func (o *AccountInfoGet200Response) SetUser(v AccountInfoGet200ResponseUser)`

SetUser sets User field to given value.


### GetData

`func (o *AccountInfoGet200Response) GetData() map[string]interface{}`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *AccountInfoGet200Response) GetDataOk() (*map[string]interface{}, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *AccountInfoGet200Response) SetData(v map[string]interface{})`

SetData sets Data field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


