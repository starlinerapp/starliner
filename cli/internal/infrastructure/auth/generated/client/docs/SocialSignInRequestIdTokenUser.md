# SocialSignInRequestIdTokenUser

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to [**SocialSignInRequestIdTokenUserName**](SocialSignInRequestIdTokenUserName.md) |  | [optional] 
**Email** | Pointer to **NullableString** |  | [optional] 

## Methods

### NewSocialSignInRequestIdTokenUser

`func NewSocialSignInRequestIdTokenUser() *SocialSignInRequestIdTokenUser`

NewSocialSignInRequestIdTokenUser instantiates a new SocialSignInRequestIdTokenUser object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSocialSignInRequestIdTokenUserWithDefaults

`func NewSocialSignInRequestIdTokenUserWithDefaults() *SocialSignInRequestIdTokenUser`

NewSocialSignInRequestIdTokenUserWithDefaults instantiates a new SocialSignInRequestIdTokenUser object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *SocialSignInRequestIdTokenUser) GetName() SocialSignInRequestIdTokenUserName`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SocialSignInRequestIdTokenUser) GetNameOk() (*SocialSignInRequestIdTokenUserName, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SocialSignInRequestIdTokenUser) SetName(v SocialSignInRequestIdTokenUserName)`

SetName sets Name field to given value.

### HasName

`func (o *SocialSignInRequestIdTokenUser) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEmail

`func (o *SocialSignInRequestIdTokenUser) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *SocialSignInRequestIdTokenUser) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *SocialSignInRequestIdTokenUser) SetEmail(v string)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *SocialSignInRequestIdTokenUser) HasEmail() bool`

HasEmail returns a boolean if a field has been set.

### SetEmailNil

`func (o *SocialSignInRequestIdTokenUser) SetEmailNil(b bool)`

 SetEmailNil sets the value for Email to be an explicit nil

### UnsetEmail
`func (o *SocialSignInRequestIdTokenUser) UnsetEmail()`

UnsetEmail ensures that no value is present for Email, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


