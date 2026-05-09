# SignUpWithEmailAndPasswordRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | The name of the user | 
**Email** | **string** | The email of the user | 
**Password** | **string** | The password of the user | 
**Image** | Pointer to **string** | The profile image URL of the user | [optional] 
**CallbackURL** | Pointer to **string** | The URL to use for email verification callback | [optional] 
**RememberMe** | Pointer to **bool** | If this is false, the session will not be remembered. Default is &#x60;true&#x60;. | [optional] 

## Methods

### NewSignUpWithEmailAndPasswordRequest

`func NewSignUpWithEmailAndPasswordRequest(name string, email string, password string, ) *SignUpWithEmailAndPasswordRequest`

NewSignUpWithEmailAndPasswordRequest instantiates a new SignUpWithEmailAndPasswordRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSignUpWithEmailAndPasswordRequestWithDefaults

`func NewSignUpWithEmailAndPasswordRequestWithDefaults() *SignUpWithEmailAndPasswordRequest`

NewSignUpWithEmailAndPasswordRequestWithDefaults instantiates a new SignUpWithEmailAndPasswordRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *SignUpWithEmailAndPasswordRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SignUpWithEmailAndPasswordRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SignUpWithEmailAndPasswordRequest) SetName(v string)`

SetName sets Name field to given value.


### GetEmail

`func (o *SignUpWithEmailAndPasswordRequest) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *SignUpWithEmailAndPasswordRequest) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *SignUpWithEmailAndPasswordRequest) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetPassword

`func (o *SignUpWithEmailAndPasswordRequest) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *SignUpWithEmailAndPasswordRequest) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *SignUpWithEmailAndPasswordRequest) SetPassword(v string)`

SetPassword sets Password field to given value.


### GetImage

`func (o *SignUpWithEmailAndPasswordRequest) GetImage() string`

GetImage returns the Image field if non-nil, zero value otherwise.

### GetImageOk

`func (o *SignUpWithEmailAndPasswordRequest) GetImageOk() (*string, bool)`

GetImageOk returns a tuple with the Image field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImage

`func (o *SignUpWithEmailAndPasswordRequest) SetImage(v string)`

SetImage sets Image field to given value.

### HasImage

`func (o *SignUpWithEmailAndPasswordRequest) HasImage() bool`

HasImage returns a boolean if a field has been set.

### GetCallbackURL

`func (o *SignUpWithEmailAndPasswordRequest) GetCallbackURL() string`

GetCallbackURL returns the CallbackURL field if non-nil, zero value otherwise.

### GetCallbackURLOk

`func (o *SignUpWithEmailAndPasswordRequest) GetCallbackURLOk() (*string, bool)`

GetCallbackURLOk returns a tuple with the CallbackURL field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallbackURL

`func (o *SignUpWithEmailAndPasswordRequest) SetCallbackURL(v string)`

SetCallbackURL sets CallbackURL field to given value.

### HasCallbackURL

`func (o *SignUpWithEmailAndPasswordRequest) HasCallbackURL() bool`

HasCallbackURL returns a boolean if a field has been set.

### GetRememberMe

`func (o *SignUpWithEmailAndPasswordRequest) GetRememberMe() bool`

GetRememberMe returns the RememberMe field if non-nil, zero value otherwise.

### GetRememberMeOk

`func (o *SignUpWithEmailAndPasswordRequest) GetRememberMeOk() (*bool, bool)`

GetRememberMeOk returns a tuple with the RememberMe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRememberMe

`func (o *SignUpWithEmailAndPasswordRequest) SetRememberMe(v bool)`

SetRememberMe sets RememberMe field to given value.

### HasRememberMe

`func (o *SignUpWithEmailAndPasswordRequest) HasRememberMe() bool`

HasRememberMe returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


