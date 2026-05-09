# ChangeEmail200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**User** | Pointer to [**User**](User.md) |  | [optional] 
**Status** | **bool** | Indicates if the request was successful | 
**Message** | Pointer to **string** | Status message of the email change process | [optional] 

## Methods

### NewChangeEmail200Response

`func NewChangeEmail200Response(status bool, ) *ChangeEmail200Response`

NewChangeEmail200Response instantiates a new ChangeEmail200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewChangeEmail200ResponseWithDefaults

`func NewChangeEmail200ResponseWithDefaults() *ChangeEmail200Response`

NewChangeEmail200ResponseWithDefaults instantiates a new ChangeEmail200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUser

`func (o *ChangeEmail200Response) GetUser() User`

GetUser returns the User field if non-nil, zero value otherwise.

### GetUserOk

`func (o *ChangeEmail200Response) GetUserOk() (*User, bool)`

GetUserOk returns a tuple with the User field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUser

`func (o *ChangeEmail200Response) SetUser(v User)`

SetUser sets User field to given value.

### HasUser

`func (o *ChangeEmail200Response) HasUser() bool`

HasUser returns a boolean if a field has been set.

### GetStatus

`func (o *ChangeEmail200Response) GetStatus() bool`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ChangeEmail200Response) GetStatusOk() (*bool, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ChangeEmail200Response) SetStatus(v bool)`

SetStatus sets Status field to given value.


### GetMessage

`func (o *ChangeEmail200Response) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ChangeEmail200Response) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *ChangeEmail200Response) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *ChangeEmail200Response) HasMessage() bool`

HasMessage returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


