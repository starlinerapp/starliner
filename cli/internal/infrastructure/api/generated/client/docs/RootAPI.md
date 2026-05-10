# \RootAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetRoot**](RootAPI.md#GetRoot) | **Get** / | Get root



## GetRoot

> ResponseRoot GetRoot(ctx).XUserID(xUserID).Execute()

Get root

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/apiClient"
)

func main() {
	xUserID := "xUserID_example" // string | User ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.RootAPI.GetRoot(context.Background()).XUserID(xUserID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RootAPI.GetRoot``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetRoot`: ResponseRoot
	fmt.Fprintf(os.Stdout, "Response from `RootAPI.GetRoot`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetRootRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xUserID** | **string** | User ID | 

### Return type

[**ResponseRoot**](ResponseRoot.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

