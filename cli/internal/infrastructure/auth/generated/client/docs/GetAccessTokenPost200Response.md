# GetAccessTokenPost200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TokenType** | Pointer to **string** |  | [optional] 
**IdToken** | Pointer to **string** |  | [optional] 
**AccessToken** | Pointer to **string** |  | [optional] 
**AccessTokenExpiresAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewGetAccessTokenPost200Response

`func NewGetAccessTokenPost200Response() *GetAccessTokenPost200Response`

NewGetAccessTokenPost200Response instantiates a new GetAccessTokenPost200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGetAccessTokenPost200ResponseWithDefaults

`func NewGetAccessTokenPost200ResponseWithDefaults() *GetAccessTokenPost200Response`

NewGetAccessTokenPost200ResponseWithDefaults instantiates a new GetAccessTokenPost200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTokenType

`func (o *GetAccessTokenPost200Response) GetTokenType() string`

GetTokenType returns the TokenType field if non-nil, zero value otherwise.

### GetTokenTypeOk

`func (o *GetAccessTokenPost200Response) GetTokenTypeOk() (*string, bool)`

GetTokenTypeOk returns a tuple with the TokenType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTokenType

`func (o *GetAccessTokenPost200Response) SetTokenType(v string)`

SetTokenType sets TokenType field to given value.

### HasTokenType

`func (o *GetAccessTokenPost200Response) HasTokenType() bool`

HasTokenType returns a boolean if a field has been set.

### GetIdToken

`func (o *GetAccessTokenPost200Response) GetIdToken() string`

GetIdToken returns the IdToken field if non-nil, zero value otherwise.

### GetIdTokenOk

`func (o *GetAccessTokenPost200Response) GetIdTokenOk() (*string, bool)`

GetIdTokenOk returns a tuple with the IdToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIdToken

`func (o *GetAccessTokenPost200Response) SetIdToken(v string)`

SetIdToken sets IdToken field to given value.

### HasIdToken

`func (o *GetAccessTokenPost200Response) HasIdToken() bool`

HasIdToken returns a boolean if a field has been set.

### GetAccessToken

`func (o *GetAccessTokenPost200Response) GetAccessToken() string`

GetAccessToken returns the AccessToken field if non-nil, zero value otherwise.

### GetAccessTokenOk

`func (o *GetAccessTokenPost200Response) GetAccessTokenOk() (*string, bool)`

GetAccessTokenOk returns a tuple with the AccessToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessToken

`func (o *GetAccessTokenPost200Response) SetAccessToken(v string)`

SetAccessToken sets AccessToken field to given value.

### HasAccessToken

`func (o *GetAccessTokenPost200Response) HasAccessToken() bool`

HasAccessToken returns a boolean if a field has been set.

### GetAccessTokenExpiresAt

`func (o *GetAccessTokenPost200Response) GetAccessTokenExpiresAt() time.Time`

GetAccessTokenExpiresAt returns the AccessTokenExpiresAt field if non-nil, zero value otherwise.

### GetAccessTokenExpiresAtOk

`func (o *GetAccessTokenPost200Response) GetAccessTokenExpiresAtOk() (*time.Time, bool)`

GetAccessTokenExpiresAtOk returns a tuple with the AccessTokenExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessTokenExpiresAt

`func (o *GetAccessTokenPost200Response) SetAccessTokenExpiresAt(v time.Time)`

SetAccessTokenExpiresAt sets AccessTokenExpiresAt field to given value.

### HasAccessTokenExpiresAt

`func (o *GetAccessTokenPost200Response) HasAccessTokenExpiresAt() bool`

HasAccessTokenExpiresAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


