# SignUpWithEmailAndPassword200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Token** | Pointer to **string** | Authentication token for the session | [optional] 
**User** | [**SignUpWithEmailAndPassword200ResponseUser**](SignUpWithEmailAndPassword200ResponseUser.md) |  | 

## Methods

### NewSignUpWithEmailAndPassword200Response

`func NewSignUpWithEmailAndPassword200Response(user SignUpWithEmailAndPassword200ResponseUser, ) *SignUpWithEmailAndPassword200Response`

NewSignUpWithEmailAndPassword200Response instantiates a new SignUpWithEmailAndPassword200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSignUpWithEmailAndPassword200ResponseWithDefaults

`func NewSignUpWithEmailAndPassword200ResponseWithDefaults() *SignUpWithEmailAndPassword200Response`

NewSignUpWithEmailAndPassword200ResponseWithDefaults instantiates a new SignUpWithEmailAndPassword200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetToken

`func (o *SignUpWithEmailAndPassword200Response) GetToken() string`

GetToken returns the Token field if non-nil, zero value otherwise.

### GetTokenOk

`func (o *SignUpWithEmailAndPassword200Response) GetTokenOk() (*string, bool)`

GetTokenOk returns a tuple with the Token field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetToken

`func (o *SignUpWithEmailAndPassword200Response) SetToken(v string)`

SetToken sets Token field to given value.

### HasToken

`func (o *SignUpWithEmailAndPassword200Response) HasToken() bool`

HasToken returns a boolean if a field has been set.

### GetUser

`func (o *SignUpWithEmailAndPassword200Response) GetUser() SignUpWithEmailAndPassword200ResponseUser`

GetUser returns the User field if non-nil, zero value otherwise.

### GetUserOk

`func (o *SignUpWithEmailAndPassword200Response) GetUserOk() (*SignUpWithEmailAndPassword200ResponseUser, bool)`

GetUserOk returns a tuple with the User field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUser

`func (o *SignUpWithEmailAndPassword200Response) SetUser(v SignUpWithEmailAndPassword200ResponseUser)`

SetUser sets User field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


