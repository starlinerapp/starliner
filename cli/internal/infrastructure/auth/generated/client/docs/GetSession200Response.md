# GetSession200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Session** | [**Session**](Session.md) |  | 
**User** | [**User**](User.md) |  | 

## Methods

### NewGetSession200Response

`func NewGetSession200Response(session Session, user User, ) *GetSession200Response`

NewGetSession200Response instantiates a new GetSession200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGetSession200ResponseWithDefaults

`func NewGetSession200ResponseWithDefaults() *GetSession200Response`

NewGetSession200ResponseWithDefaults instantiates a new GetSession200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSession

`func (o *GetSession200Response) GetSession() Session`

GetSession returns the Session field if non-nil, zero value otherwise.

### GetSessionOk

`func (o *GetSession200Response) GetSessionOk() (*Session, bool)`

GetSessionOk returns a tuple with the Session field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSession

`func (o *GetSession200Response) SetSession(v Session)`

SetSession sets Session field to given value.


### GetUser

`func (o *GetSession200Response) GetUser() User`

GetUser returns the User field if non-nil, zero value otherwise.

### GetUserOk

`func (o *GetSession200Response) GetUserOk() (*User, bool)`

GetUserOk returns a tuple with the User field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUser

`func (o *GetSession200Response) SetUser(v User)`

SetUser sets User field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


