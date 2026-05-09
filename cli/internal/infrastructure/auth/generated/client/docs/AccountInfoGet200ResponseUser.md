# AccountInfoGet200ResponseUser

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**Name** | Pointer to **string** |  | [optional] 
**Email** | Pointer to **string** |  | [optional] 
**Image** | Pointer to **string** |  | [optional] 
**EmailVerified** | **bool** |  | 

## Methods

### NewAccountInfoGet200ResponseUser

`func NewAccountInfoGet200ResponseUser(id string, emailVerified bool, ) *AccountInfoGet200ResponseUser`

NewAccountInfoGet200ResponseUser instantiates a new AccountInfoGet200ResponseUser object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAccountInfoGet200ResponseUserWithDefaults

`func NewAccountInfoGet200ResponseUserWithDefaults() *AccountInfoGet200ResponseUser`

NewAccountInfoGet200ResponseUserWithDefaults instantiates a new AccountInfoGet200ResponseUser object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *AccountInfoGet200ResponseUser) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *AccountInfoGet200ResponseUser) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *AccountInfoGet200ResponseUser) SetId(v string)`

SetId sets Id field to given value.


### GetName

`func (o *AccountInfoGet200ResponseUser) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *AccountInfoGet200ResponseUser) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *AccountInfoGet200ResponseUser) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *AccountInfoGet200ResponseUser) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEmail

`func (o *AccountInfoGet200ResponseUser) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *AccountInfoGet200ResponseUser) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *AccountInfoGet200ResponseUser) SetEmail(v string)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *AccountInfoGet200ResponseUser) HasEmail() bool`

HasEmail returns a boolean if a field has been set.

### GetImage

`func (o *AccountInfoGet200ResponseUser) GetImage() string`

GetImage returns the Image field if non-nil, zero value otherwise.

### GetImageOk

`func (o *AccountInfoGet200ResponseUser) GetImageOk() (*string, bool)`

GetImageOk returns a tuple with the Image field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImage

`func (o *AccountInfoGet200ResponseUser) SetImage(v string)`

SetImage sets Image field to given value.

### HasImage

`func (o *AccountInfoGet200ResponseUser) HasImage() bool`

HasImage returns a boolean if a field has been set.

### GetEmailVerified

`func (o *AccountInfoGet200ResponseUser) GetEmailVerified() bool`

GetEmailVerified returns the EmailVerified field if non-nil, zero value otherwise.

### GetEmailVerifiedOk

`func (o *AccountInfoGet200ResponseUser) GetEmailVerifiedOk() (*bool, bool)`

GetEmailVerifiedOk returns a tuple with the EmailVerified field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmailVerified

`func (o *AccountInfoGet200ResponseUser) SetEmailVerified(v bool)`

SetEmailVerified sets EmailVerified field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


