# VerifyEmailGet200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**User** | [**User**](User.md) |  | 
**Status** | **bool** | Indicates if the email was verified successfully | 

## Methods

### NewVerifyEmailGet200Response

`func NewVerifyEmailGet200Response(user User, status bool, ) *VerifyEmailGet200Response`

NewVerifyEmailGet200Response instantiates a new VerifyEmailGet200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVerifyEmailGet200ResponseWithDefaults

`func NewVerifyEmailGet200ResponseWithDefaults() *VerifyEmailGet200Response`

NewVerifyEmailGet200ResponseWithDefaults instantiates a new VerifyEmailGet200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUser

`func (o *VerifyEmailGet200Response) GetUser() User`

GetUser returns the User field if non-nil, zero value otherwise.

### GetUserOk

`func (o *VerifyEmailGet200Response) GetUserOk() (*User, bool)`

GetUserOk returns a tuple with the User field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUser

`func (o *VerifyEmailGet200Response) SetUser(v User)`

SetUser sets User field to given value.


### GetStatus

`func (o *VerifyEmailGet200Response) GetStatus() bool`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *VerifyEmailGet200Response) GetStatusOk() (*bool, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *VerifyEmailGet200Response) SetStatus(v bool)`

SetStatus sets Status field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


