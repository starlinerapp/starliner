# SignUpWithEmailAndPassword200ResponseUser

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | The unique identifier of the user | 
**Email** | **string** | The email address of the user | 
**Name** | **string** | The name of the user | 
**Image** | Pointer to **string** | The profile image URL of the user | [optional] 
**EmailVerified** | **bool** | Whether the email has been verified | 
**CreatedAt** | **time.Time** | When the user was created | 
**UpdatedAt** | **time.Time** | When the user was last updated | 

## Methods

### NewSignUpWithEmailAndPassword200ResponseUser

`func NewSignUpWithEmailAndPassword200ResponseUser(id string, email string, name string, emailVerified bool, createdAt time.Time, updatedAt time.Time, ) *SignUpWithEmailAndPassword200ResponseUser`

NewSignUpWithEmailAndPassword200ResponseUser instantiates a new SignUpWithEmailAndPassword200ResponseUser object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSignUpWithEmailAndPassword200ResponseUserWithDefaults

`func NewSignUpWithEmailAndPassword200ResponseUserWithDefaults() *SignUpWithEmailAndPassword200ResponseUser`

NewSignUpWithEmailAndPassword200ResponseUserWithDefaults instantiates a new SignUpWithEmailAndPassword200ResponseUser object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *SignUpWithEmailAndPassword200ResponseUser) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *SignUpWithEmailAndPassword200ResponseUser) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *SignUpWithEmailAndPassword200ResponseUser) SetId(v string)`

SetId sets Id field to given value.


### GetEmail

`func (o *SignUpWithEmailAndPassword200ResponseUser) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *SignUpWithEmailAndPassword200ResponseUser) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *SignUpWithEmailAndPassword200ResponseUser) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetName

`func (o *SignUpWithEmailAndPassword200ResponseUser) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SignUpWithEmailAndPassword200ResponseUser) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SignUpWithEmailAndPassword200ResponseUser) SetName(v string)`

SetName sets Name field to given value.


### GetImage

`func (o *SignUpWithEmailAndPassword200ResponseUser) GetImage() string`

GetImage returns the Image field if non-nil, zero value otherwise.

### GetImageOk

`func (o *SignUpWithEmailAndPassword200ResponseUser) GetImageOk() (*string, bool)`

GetImageOk returns a tuple with the Image field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImage

`func (o *SignUpWithEmailAndPassword200ResponseUser) SetImage(v string)`

SetImage sets Image field to given value.

### HasImage

`func (o *SignUpWithEmailAndPassword200ResponseUser) HasImage() bool`

HasImage returns a boolean if a field has been set.

### GetEmailVerified

`func (o *SignUpWithEmailAndPassword200ResponseUser) GetEmailVerified() bool`

GetEmailVerified returns the EmailVerified field if non-nil, zero value otherwise.

### GetEmailVerifiedOk

`func (o *SignUpWithEmailAndPassword200ResponseUser) GetEmailVerifiedOk() (*bool, bool)`

GetEmailVerifiedOk returns a tuple with the EmailVerified field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmailVerified

`func (o *SignUpWithEmailAndPassword200ResponseUser) SetEmailVerified(v bool)`

SetEmailVerified sets EmailVerified field to given value.


### GetCreatedAt

`func (o *SignUpWithEmailAndPassword200ResponseUser) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *SignUpWithEmailAndPassword200ResponseUser) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *SignUpWithEmailAndPassword200ResponseUser) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetUpdatedAt

`func (o *SignUpWithEmailAndPassword200ResponseUser) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *SignUpWithEmailAndPassword200ResponseUser) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *SignUpWithEmailAndPassword200ResponseUser) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


