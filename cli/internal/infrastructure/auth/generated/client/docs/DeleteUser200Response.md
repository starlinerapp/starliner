# DeleteUser200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Success** | **bool** | Indicates if the operation was successful | 
**Message** | **string** | Status message of the deletion process | 

## Methods

### NewDeleteUser200Response

`func NewDeleteUser200Response(success bool, message string, ) *DeleteUser200Response`

NewDeleteUser200Response instantiates a new DeleteUser200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeleteUser200ResponseWithDefaults

`func NewDeleteUser200ResponseWithDefaults() *DeleteUser200Response`

NewDeleteUser200ResponseWithDefaults instantiates a new DeleteUser200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSuccess

`func (o *DeleteUser200Response) GetSuccess() bool`

GetSuccess returns the Success field if non-nil, zero value otherwise.

### GetSuccessOk

`func (o *DeleteUser200Response) GetSuccessOk() (*bool, bool)`

GetSuccessOk returns a tuple with the Success field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuccess

`func (o *DeleteUser200Response) SetSuccess(v bool)`

SetSuccess sets Success field to given value.


### GetMessage

`func (o *DeleteUser200Response) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *DeleteUser200Response) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *DeleteUser200Response) SetMessage(v string)`

SetMessage sets Message field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


