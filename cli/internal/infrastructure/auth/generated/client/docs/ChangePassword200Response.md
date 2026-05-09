# ChangePassword200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Token** | Pointer to **string** | New session token if other sessions were revoked | [optional] 
**User** | [**SignUpWithEmailAndPassword200ResponseUser**](SignUpWithEmailAndPassword200ResponseUser.md) |  | 

## Methods

### NewChangePassword200Response

`func NewChangePassword200Response(user SignUpWithEmailAndPassword200ResponseUser, ) *ChangePassword200Response`

NewChangePassword200Response instantiates a new ChangePassword200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewChangePassword200ResponseWithDefaults

`func NewChangePassword200ResponseWithDefaults() *ChangePassword200Response`

NewChangePassword200ResponseWithDefaults instantiates a new ChangePassword200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetToken

`func (o *ChangePassword200Response) GetToken() string`

GetToken returns the Token field if non-nil, zero value otherwise.

### GetTokenOk

`func (o *ChangePassword200Response) GetTokenOk() (*string, bool)`

GetTokenOk returns a tuple with the Token field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetToken

`func (o *ChangePassword200Response) SetToken(v string)`

SetToken sets Token field to given value.

### HasToken

`func (o *ChangePassword200Response) HasToken() bool`

HasToken returns a boolean if a field has been set.

### GetUser

`func (o *ChangePassword200Response) GetUser() SignUpWithEmailAndPassword200ResponseUser`

GetUser returns the User field if non-nil, zero value otherwise.

### GetUserOk

`func (o *ChangePassword200Response) GetUserOk() (*SignUpWithEmailAndPassword200ResponseUser, bool)`

GetUserOk returns a tuple with the User field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUser

`func (o *ChangePassword200Response) SetUser(v SignUpWithEmailAndPassword200ResponseUser)`

SetUser sets User field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


