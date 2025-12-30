# BuildApi

All URIs are relative to _http://localhost_

| Method                            | HTTP request     | Description   |
| --------------------------------- | ---------------- | ------------- |
| [**triggerBuild**](#triggerbuild) | **POST** /builds | Trigger Build |

# **triggerBuild**

> triggerBuild()

### Example

```typescript
import { BuildApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new BuildApi(configuration);

let xUserID: string; //User ID (default to undefined)

const { status, data } = await apiInstance.triggerBuild(xUserID);
```

### Parameters

| Name        | Type         | Description | Notes                 |
| ----------- | ------------ | ----------- | --------------------- |
| **xUserID** | [**string**] | User ID     | defaults to undefined |

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

### HTTP response details

| Status code | Description | Response headers |
| ----------- | ----------- | ---------------- |
| **200**     | OK          | -                |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)
