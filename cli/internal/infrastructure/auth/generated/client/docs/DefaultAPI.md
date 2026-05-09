# \DefaultAPI

All URIs are relative to *https://dev.starliner.app/api/auth*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AccountInfoGet**](DefaultAPI.md#AccountInfoGet) | **Get** /account-info | 
[**CallbackIdGet**](DefaultAPI.md#CallbackIdGet) | **Get** /callback/{id} | 
[**CallbackIdPost**](DefaultAPI.md#CallbackIdPost) | **Post** /callback/{id} | 
[**ChangeEmail**](DefaultAPI.md#ChangeEmail) | **Post** /change-email | 
[**ChangePassword**](DefaultAPI.md#ChangePassword) | **Post** /change-password | 
[**DeleteUser**](DefaultAPI.md#DeleteUser) | **Post** /delete-user | 
[**DeleteUserCallbackGet**](DefaultAPI.md#DeleteUserCallbackGet) | **Get** /delete-user/callback | 
[**ErrorGet**](DefaultAPI.md#ErrorGet) | **Get** /error | 
[**GetAccessTokenPost**](DefaultAPI.md#GetAccessTokenPost) | **Post** /get-access-token | 
[**GetSession**](DefaultAPI.md#GetSession) | **Get** /get-session | 
[**GetSession_0**](DefaultAPI.md#GetSession_0) | **Post** /get-session | 
[**LinkSocialAccount**](DefaultAPI.md#LinkSocialAccount) | **Post** /link-social | 
[**ListUserAccounts**](DefaultAPI.md#ListUserAccounts) | **Get** /list-accounts | 
[**ListUserSessions**](DefaultAPI.md#ListUserSessions) | **Get** /list-sessions | 
[**OkGet**](DefaultAPI.md#OkGet) | **Get** /ok | 
[**RefreshTokenPost**](DefaultAPI.md#RefreshTokenPost) | **Post** /refresh-token | 
[**RequestPasswordReset**](DefaultAPI.md#RequestPasswordReset) | **Post** /request-password-reset | 
[**ResetPassword**](DefaultAPI.md#ResetPassword) | **Post** /reset-password | 
[**ResetPasswordCallback**](DefaultAPI.md#ResetPasswordCallback) | **Get** /reset-password/{token} | 
[**RevokeOtherSessionsPost**](DefaultAPI.md#RevokeOtherSessionsPost) | **Post** /revoke-other-sessions | 
[**RevokeSessionPost**](DefaultAPI.md#RevokeSessionPost) | **Post** /revoke-session | 
[**RevokeSessionsPost**](DefaultAPI.md#RevokeSessionsPost) | **Post** /revoke-sessions | 
[**SendVerificationEmail**](DefaultAPI.md#SendVerificationEmail) | **Post** /send-verification-email | 
[**SignInEmail**](DefaultAPI.md#SignInEmail) | **Post** /sign-in/email | 
[**SignOut**](DefaultAPI.md#SignOut) | **Post** /sign-out | 
[**SignUpWithEmailAndPassword**](DefaultAPI.md#SignUpWithEmailAndPassword) | **Post** /sign-up/email | 
[**SocialSignIn**](DefaultAPI.md#SocialSignIn) | **Post** /sign-in/social | 
[**UnlinkAccountPost**](DefaultAPI.md#UnlinkAccountPost) | **Post** /unlink-account | 
[**UpdateSession**](DefaultAPI.md#UpdateSession) | **Post** /update-session | 
[**UpdateUser**](DefaultAPI.md#UpdateUser) | **Post** /update-user | 
[**VerifyEmailGet**](DefaultAPI.md#VerifyEmailGet) | **Get** /verify-email | 
[**VerifyPassword**](DefaultAPI.md#VerifyPassword) | **Post** /verify-password | 



## AccountInfoGet

> AccountInfoGet200Response AccountInfoGet(ctx).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.AccountInfoGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.AccountInfoGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AccountInfoGet`: AccountInfoGet200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.AccountInfoGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiAccountInfoGetRequest struct via the builder pattern


### Return type

[**AccountInfoGet200Response**](AccountInfoGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CallbackIdGet

> CallbackIdGet(ctx).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DefaultAPI.CallbackIdGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.CallbackIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiCallbackIdGetRequest struct via the builder pattern


### Return type

 (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CallbackIdPost

> CallbackIdPost(ctx).Body(body).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	body := map[string]interface{}{ ... } // map[string]interface{} |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DefaultAPI.CallbackIdPost(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.CallbackIdPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCallbackIdPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | **map[string]interface{}** |  | 

### Return type

 (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ChangeEmail

> ChangeEmail200Response ChangeEmail(ctx).ChangeEmailRequest(changeEmailRequest).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	changeEmailRequest := *openapiclient.NewChangeEmailRequest("NewEmail_example") // ChangeEmailRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ChangeEmail(context.Background()).ChangeEmailRequest(changeEmailRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ChangeEmail``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ChangeEmail`: ChangeEmail200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ChangeEmail`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiChangeEmailRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changeEmailRequest** | [**ChangeEmailRequest**](ChangeEmailRequest.md) |  | 

### Return type

[**ChangeEmail200Response**](ChangeEmail200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ChangePassword

> ChangePassword200Response ChangePassword(ctx).ChangePasswordRequest(changePasswordRequest).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	changePasswordRequest := *openapiclient.NewChangePasswordRequest("NewPassword_example", "CurrentPassword_example") // ChangePasswordRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ChangePassword(context.Background()).ChangePasswordRequest(changePasswordRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ChangePassword``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ChangePassword`: ChangePassword200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ChangePassword`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiChangePasswordRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changePasswordRequest** | [**ChangePasswordRequest**](ChangePasswordRequest.md) |  | 

### Return type

[**ChangePassword200Response**](ChangePassword200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteUser

> DeleteUser200Response DeleteUser(ctx).DeleteUserRequest(deleteUserRequest).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	deleteUserRequest := *openapiclient.NewDeleteUserRequest() // DeleteUserRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.DeleteUser(context.Background()).DeleteUserRequest(deleteUserRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.DeleteUser``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteUser`: DeleteUser200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.DeleteUser`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDeleteUserRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deleteUserRequest** | [**DeleteUserRequest**](DeleteUserRequest.md) |  | 

### Return type

[**DeleteUser200Response**](DeleteUser200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteUserCallbackGet

> DeleteUserCallbackGet200Response DeleteUserCallbackGet(ctx).Token(token).CallbackURL(callbackURL).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	token := "token_example" // string |  (optional)
	callbackURL := "callbackURL_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.DeleteUserCallbackGet(context.Background()).Token(token).CallbackURL(callbackURL).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.DeleteUserCallbackGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteUserCallbackGet`: DeleteUserCallbackGet200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.DeleteUserCallbackGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDeleteUserCallbackGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **token** | **string** |  | 
 **callbackURL** | **string** |  | 

### Return type

[**DeleteUserCallbackGet200Response**](DeleteUserCallbackGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ErrorGet

> string ErrorGet(ctx).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ErrorGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ErrorGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ErrorGet`: string
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ErrorGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiErrorGetRequest struct via the builder pattern


### Return type

**string**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/html, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAccessTokenPost

> GetAccessTokenPost200Response GetAccessTokenPost(ctx).RefreshTokenPostRequest(refreshTokenPostRequest).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	refreshTokenPostRequest := *openapiclient.NewRefreshTokenPostRequest("ProviderId_example") // RefreshTokenPostRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.GetAccessTokenPost(context.Background()).RefreshTokenPostRequest(refreshTokenPostRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetAccessTokenPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAccessTokenPost`: GetAccessTokenPost200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetAccessTokenPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetAccessTokenPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **refreshTokenPostRequest** | [**RefreshTokenPostRequest**](RefreshTokenPostRequest.md) |  | 

### Return type

[**GetAccessTokenPost200Response**](GetAccessTokenPost200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetSession

> GetSession200Response GetSession(ctx).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.GetSession(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetSession``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetSession`: GetSession200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetSession`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetSessionRequest struct via the builder pattern


### Return type

[**GetSession200Response**](GetSession200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetSession_0

> GetSession_0(ctx).Body(body).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	body := map[string]interface{}{ ... } // map[string]interface{} |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DefaultAPI.GetSession_0(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetSession_0``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetSession_1Request struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | **map[string]interface{}** |  | 

### Return type

 (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## LinkSocialAccount

> LinkSocialAccount200Response LinkSocialAccount(ctx).LinkSocialAccountRequest(linkSocialAccountRequest).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	linkSocialAccountRequest := *openapiclient.NewLinkSocialAccountRequest("Provider_example") // LinkSocialAccountRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.LinkSocialAccount(context.Background()).LinkSocialAccountRequest(linkSocialAccountRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.LinkSocialAccount``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `LinkSocialAccount`: LinkSocialAccount200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.LinkSocialAccount`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiLinkSocialAccountRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **linkSocialAccountRequest** | [**LinkSocialAccountRequest**](LinkSocialAccountRequest.md) |  | 

### Return type

[**LinkSocialAccount200Response**](LinkSocialAccount200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListUserAccounts

> []ListUserAccounts200ResponseInner ListUserAccounts(ctx).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ListUserAccounts(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ListUserAccounts``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListUserAccounts`: []ListUserAccounts200ResponseInner
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ListUserAccounts`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListUserAccountsRequest struct via the builder pattern


### Return type

[**[]ListUserAccounts200ResponseInner**](ListUserAccounts200ResponseInner.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListUserSessions

> []Session ListUserSessions(ctx).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ListUserSessions(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ListUserSessions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListUserSessions`: []Session
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ListUserSessions`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListUserSessionsRequest struct via the builder pattern


### Return type

[**[]Session**](Session.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OkGet

> OkGet200Response OkGet(ctx).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.OkGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.OkGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OkGet`: OkGet200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.OkGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiOkGetRequest struct via the builder pattern


### Return type

[**OkGet200Response**](OkGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RefreshTokenPost

> RefreshTokenPost200Response RefreshTokenPost(ctx).RefreshTokenPostRequest(refreshTokenPostRequest).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	refreshTokenPostRequest := *openapiclient.NewRefreshTokenPostRequest("ProviderId_example") // RefreshTokenPostRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.RefreshTokenPost(context.Background()).RefreshTokenPostRequest(refreshTokenPostRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.RefreshTokenPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RefreshTokenPost`: RefreshTokenPost200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.RefreshTokenPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRefreshTokenPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **refreshTokenPostRequest** | [**RefreshTokenPostRequest**](RefreshTokenPostRequest.md) |  | 

### Return type

[**RefreshTokenPost200Response**](RefreshTokenPost200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RequestPasswordReset

> RequestPasswordReset200Response RequestPasswordReset(ctx).RequestPasswordResetRequest(requestPasswordResetRequest).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	requestPasswordResetRequest := *openapiclient.NewRequestPasswordResetRequest("Email_example") // RequestPasswordResetRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.RequestPasswordReset(context.Background()).RequestPasswordResetRequest(requestPasswordResetRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.RequestPasswordReset``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RequestPasswordReset`: RequestPasswordReset200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.RequestPasswordReset`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRequestPasswordResetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestPasswordResetRequest** | [**RequestPasswordResetRequest**](RequestPasswordResetRequest.md) |  | 

### Return type

[**RequestPasswordReset200Response**](RequestPasswordReset200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ResetPassword

> ResetPassword200Response ResetPassword(ctx).ResetPasswordRequest(resetPasswordRequest).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	resetPasswordRequest := *openapiclient.NewResetPasswordRequest("NewPassword_example") // ResetPasswordRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ResetPassword(context.Background()).ResetPasswordRequest(resetPasswordRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ResetPassword``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ResetPassword`: ResetPassword200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ResetPassword`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiResetPasswordRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **resetPasswordRequest** | [**ResetPasswordRequest**](ResetPasswordRequest.md) |  | 

### Return type

[**ResetPassword200Response**](ResetPassword200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ResetPasswordCallback

> ResetPasswordCallback200Response ResetPasswordCallback(ctx, token).CallbackURL(callbackURL).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	token := "token_example" // string | The token to reset the password
	callbackURL := "callbackURL_example" // string | The URL to redirect the user to reset their password

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ResetPasswordCallback(context.Background(), token).CallbackURL(callbackURL).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ResetPasswordCallback``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ResetPasswordCallback`: ResetPasswordCallback200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ResetPasswordCallback`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**token** | **string** | The token to reset the password | 

### Other Parameters

Other parameters are passed through a pointer to a apiResetPasswordCallbackRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **callbackURL** | **string** | The URL to redirect the user to reset their password | 

### Return type

[**ResetPasswordCallback200Response**](ResetPasswordCallback200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RevokeOtherSessionsPost

> RevokeOtherSessionsPost200Response RevokeOtherSessionsPost(ctx).Body(body).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	body := map[string]interface{}{ ... } // map[string]interface{} |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.RevokeOtherSessionsPost(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.RevokeOtherSessionsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RevokeOtherSessionsPost`: RevokeOtherSessionsPost200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.RevokeOtherSessionsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRevokeOtherSessionsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | **map[string]interface{}** |  | 

### Return type

[**RevokeOtherSessionsPost200Response**](RevokeOtherSessionsPost200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RevokeSessionPost

> RevokeSessionPost200Response RevokeSessionPost(ctx).RevokeSessionPostRequest(revokeSessionPostRequest).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	revokeSessionPostRequest := *openapiclient.NewRevokeSessionPostRequest("Token_example") // RevokeSessionPostRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.RevokeSessionPost(context.Background()).RevokeSessionPostRequest(revokeSessionPostRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.RevokeSessionPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RevokeSessionPost`: RevokeSessionPost200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.RevokeSessionPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRevokeSessionPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **revokeSessionPostRequest** | [**RevokeSessionPostRequest**](RevokeSessionPostRequest.md) |  | 

### Return type

[**RevokeSessionPost200Response**](RevokeSessionPost200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RevokeSessionsPost

> RevokeSessionsPost200Response RevokeSessionsPost(ctx).Body(body).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	body := map[string]interface{}{ ... } // map[string]interface{} |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.RevokeSessionsPost(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.RevokeSessionsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RevokeSessionsPost`: RevokeSessionsPost200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.RevokeSessionsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRevokeSessionsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | **map[string]interface{}** |  | 

### Return type

[**RevokeSessionsPost200Response**](RevokeSessionsPost200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SendVerificationEmail

> SendVerificationEmail200Response SendVerificationEmail(ctx).SendVerificationEmailRequest(sendVerificationEmailRequest).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	sendVerificationEmailRequest := *openapiclient.NewSendVerificationEmailRequest("user@example.com") // SendVerificationEmailRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.SendVerificationEmail(context.Background()).SendVerificationEmailRequest(sendVerificationEmailRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.SendVerificationEmail``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SendVerificationEmail`: SendVerificationEmail200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.SendVerificationEmail`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSendVerificationEmailRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **sendVerificationEmailRequest** | [**SendVerificationEmailRequest**](SendVerificationEmailRequest.md) |  | 

### Return type

[**SendVerificationEmail200Response**](SendVerificationEmail200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SignInEmail

> SignInEmail200Response SignInEmail(ctx).SignInEmailRequest(signInEmailRequest).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	signInEmailRequest := *openapiclient.NewSignInEmailRequest("Email_example", "Password_example") // SignInEmailRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.SignInEmail(context.Background()).SignInEmailRequest(signInEmailRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.SignInEmail``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SignInEmail`: SignInEmail200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.SignInEmail`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSignInEmailRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **signInEmailRequest** | [**SignInEmailRequest**](SignInEmailRequest.md) |  | 

### Return type

[**SignInEmail200Response**](SignInEmail200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SignOut

> SignOut200Response SignOut(ctx).Body(body).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	body := map[string]interface{}{ ... } // map[string]interface{} |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.SignOut(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.SignOut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SignOut`: SignOut200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.SignOut`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSignOutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | **map[string]interface{}** |  | 

### Return type

[**SignOut200Response**](SignOut200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SignUpWithEmailAndPassword

> SignUpWithEmailAndPassword200Response SignUpWithEmailAndPassword(ctx).SignUpWithEmailAndPasswordRequest(signUpWithEmailAndPasswordRequest).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	signUpWithEmailAndPasswordRequest := *openapiclient.NewSignUpWithEmailAndPasswordRequest("Name_example", "Email_example", "Password_example") // SignUpWithEmailAndPasswordRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.SignUpWithEmailAndPassword(context.Background()).SignUpWithEmailAndPasswordRequest(signUpWithEmailAndPasswordRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.SignUpWithEmailAndPassword``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SignUpWithEmailAndPassword`: SignUpWithEmailAndPassword200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.SignUpWithEmailAndPassword`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSignUpWithEmailAndPasswordRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **signUpWithEmailAndPasswordRequest** | [**SignUpWithEmailAndPasswordRequest**](SignUpWithEmailAndPasswordRequest.md) |  | 

### Return type

[**SignUpWithEmailAndPassword200Response**](SignUpWithEmailAndPassword200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SocialSignIn

> SocialSignIn200Response SocialSignIn(ctx).SocialSignInRequest(socialSignInRequest).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	socialSignInRequest := *openapiclient.NewSocialSignInRequest("Provider_example") // SocialSignInRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.SocialSignIn(context.Background()).SocialSignInRequest(socialSignInRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.SocialSignIn``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SocialSignIn`: SocialSignIn200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.SocialSignIn`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSocialSignInRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **socialSignInRequest** | [**SocialSignInRequest**](SocialSignInRequest.md) |  | 

### Return type

[**SocialSignIn200Response**](SocialSignIn200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UnlinkAccountPost

> ResetPassword200Response UnlinkAccountPost(ctx).UnlinkAccountPostRequest(unlinkAccountPostRequest).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	unlinkAccountPostRequest := *openapiclient.NewUnlinkAccountPostRequest("ProviderId_example") // UnlinkAccountPostRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.UnlinkAccountPost(context.Background()).UnlinkAccountPostRequest(unlinkAccountPostRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.UnlinkAccountPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UnlinkAccountPost`: ResetPassword200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.UnlinkAccountPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUnlinkAccountPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **unlinkAccountPostRequest** | [**UnlinkAccountPostRequest**](UnlinkAccountPostRequest.md) |  | 

### Return type

[**ResetPassword200Response**](ResetPassword200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateSession

> UpdateSession200Response UpdateSession(ctx).Body(body).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	body := map[string]interface{}{ ... } // map[string]interface{} |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.UpdateSession(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.UpdateSession``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateSession`: UpdateSession200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.UpdateSession`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateSessionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | **map[string]interface{}** |  | 

### Return type

[**UpdateSession200Response**](UpdateSession200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateUser

> UpdateUser200Response UpdateUser(ctx).UpdateUserRequest(updateUserRequest).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	updateUserRequest := *openapiclient.NewUpdateUserRequest() // UpdateUserRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.UpdateUser(context.Background()).UpdateUserRequest(updateUserRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.UpdateUser``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateUser`: UpdateUser200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.UpdateUser`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateUserRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **updateUserRequest** | [**UpdateUserRequest**](UpdateUserRequest.md) |  | 

### Return type

[**UpdateUser200Response**](UpdateUser200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## VerifyEmailGet

> VerifyEmailGet200Response VerifyEmailGet(ctx).Token(token).CallbackURL(callbackURL).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	token := "token_example" // string | The token to verify the email
	callbackURL := "callbackURL_example" // string | The URL to redirect to after email verification (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.VerifyEmailGet(context.Background()).Token(token).CallbackURL(callbackURL).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.VerifyEmailGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `VerifyEmailGet`: VerifyEmailGet200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.VerifyEmailGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiVerifyEmailGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **token** | **string** | The token to verify the email | 
 **callbackURL** | **string** | The URL to redirect to after email verification | 

### Return type

[**VerifyEmailGet200Response**](VerifyEmailGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## VerifyPassword

> ResetPassword200Response VerifyPassword(ctx).VerifyPasswordRequest(verifyPasswordRequest).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/authClient"
)

func main() {
	verifyPasswordRequest := *openapiclient.NewVerifyPasswordRequest("Password_example") // VerifyPasswordRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.VerifyPassword(context.Background()).VerifyPasswordRequest(verifyPasswordRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.VerifyPassword``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `VerifyPassword`: ResetPassword200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.VerifyPassword`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiVerifyPasswordRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **verifyPasswordRequest** | [**VerifyPasswordRequest**](VerifyPasswordRequest.md) |  | 

### Return type

[**ResetPassword200Response**](ResetPassword200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

