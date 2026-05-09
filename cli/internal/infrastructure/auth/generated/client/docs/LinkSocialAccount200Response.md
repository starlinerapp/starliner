# LinkSocialAccount200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Url** | Pointer to **string** | The authorization URL to redirect the user to | [optional] 
**Redirect** | **bool** | Indicates if the user should be redirected to the authorization URL | 
**Status** | Pointer to **bool** |  | [optional] 

## Methods

### NewLinkSocialAccount200Response

`func NewLinkSocialAccount200Response(redirect bool, ) *LinkSocialAccount200Response`

NewLinkSocialAccount200Response instantiates a new LinkSocialAccount200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLinkSocialAccount200ResponseWithDefaults

`func NewLinkSocialAccount200ResponseWithDefaults() *LinkSocialAccount200Response`

NewLinkSocialAccount200ResponseWithDefaults instantiates a new LinkSocialAccount200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUrl

`func (o *LinkSocialAccount200Response) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *LinkSocialAccount200Response) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *LinkSocialAccount200Response) SetUrl(v string)`

SetUrl sets Url field to given value.

### HasUrl

`func (o *LinkSocialAccount200Response) HasUrl() bool`

HasUrl returns a boolean if a field has been set.

### GetRedirect

`func (o *LinkSocialAccount200Response) GetRedirect() bool`

GetRedirect returns the Redirect field if non-nil, zero value otherwise.

### GetRedirectOk

`func (o *LinkSocialAccount200Response) GetRedirectOk() (*bool, bool)`

GetRedirectOk returns a tuple with the Redirect field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirect

`func (o *LinkSocialAccount200Response) SetRedirect(v bool)`

SetRedirect sets Redirect field to given value.


### GetStatus

`func (o *LinkSocialAccount200Response) GetStatus() bool`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *LinkSocialAccount200Response) GetStatusOk() (*bool, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *LinkSocialAccount200Response) SetStatus(v bool)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *LinkSocialAccount200Response) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


