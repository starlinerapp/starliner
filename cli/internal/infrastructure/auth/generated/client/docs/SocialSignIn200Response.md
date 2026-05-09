# SocialSignIn200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Token** | **string** |  | 
**User** | [**User**](User.md) |  | 
**Url** | Pointer to **string** |  | [optional] 
**Redirect** | **bool** |  | 

## Methods

### NewSocialSignIn200Response

`func NewSocialSignIn200Response(token string, user User, redirect bool, ) *SocialSignIn200Response`

NewSocialSignIn200Response instantiates a new SocialSignIn200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSocialSignIn200ResponseWithDefaults

`func NewSocialSignIn200ResponseWithDefaults() *SocialSignIn200Response`

NewSocialSignIn200ResponseWithDefaults instantiates a new SocialSignIn200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetToken

`func (o *SocialSignIn200Response) GetToken() string`

GetToken returns the Token field if non-nil, zero value otherwise.

### GetTokenOk

`func (o *SocialSignIn200Response) GetTokenOk() (*string, bool)`

GetTokenOk returns a tuple with the Token field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetToken

`func (o *SocialSignIn200Response) SetToken(v string)`

SetToken sets Token field to given value.


### GetUser

`func (o *SocialSignIn200Response) GetUser() User`

GetUser returns the User field if non-nil, zero value otherwise.

### GetUserOk

`func (o *SocialSignIn200Response) GetUserOk() (*User, bool)`

GetUserOk returns a tuple with the User field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUser

`func (o *SocialSignIn200Response) SetUser(v User)`

SetUser sets User field to given value.


### GetUrl

`func (o *SocialSignIn200Response) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *SocialSignIn200Response) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *SocialSignIn200Response) SetUrl(v string)`

SetUrl sets Url field to given value.

### HasUrl

`func (o *SocialSignIn200Response) HasUrl() bool`

HasUrl returns a boolean if a field has been set.

### GetRedirect

`func (o *SocialSignIn200Response) GetRedirect() bool`

GetRedirect returns the Redirect field if non-nil, zero value otherwise.

### GetRedirectOk

`func (o *SocialSignIn200Response) GetRedirectOk() (*bool, bool)`

GetRedirectOk returns a tuple with the Redirect field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirect

`func (o *SocialSignIn200Response) SetRedirect(v bool)`

SetRedirect sets Redirect field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


