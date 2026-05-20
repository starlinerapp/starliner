# RootApi

All URIs are relative to _http://localhost_

| Method                  | HTTP request | Description |
| ----------------------- | ------------ | ----------- |
| [**getRoot**](#getroot) | **GET** /    | Get root    |

# **getRoot**

> ResponseRoot getRoot()

### Example

```typescript
import { RootApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new RootApi(configuration);

let xUserID: string; //User ID (default to undefined)

const { status, data } = await apiInstance.getRoot(xUserID);
```

### Parameters

| Name        | Type         | Description | Notes                 |
| ----------- | ------------ | ----------- | --------------------- |
| **xUserID** | [**string**] | User ID     | defaults to undefined |

### Return type

**ResponseRoot**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: _/_

### HTTP response details

| Status code | Description | Response headers |
| ----------- | ----------- | ---------------- |
| **200**     | OK          | -                |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)
