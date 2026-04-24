# BuildApi

All URIs are relative to _http://localhost_

| Method                                  | HTTP request                     | Description       |
| --------------------------------------- | -------------------------------- | ----------------- |
| [**getBuildLogs**](#getbuildlogs)       | **GET** /builds/{id}/logs        | Get Build Logs    |
| [**streamBuildLogs**](#streambuildlogs) | **GET** /builds/{id}/logs/stream | Stream build logs |

# **getBuildLogs**

> ResponseBuildLogs getBuildLogs()

### Example

```typescript
import { BuildApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new BuildApi(configuration);

let xUserID: string; //User ID (default to undefined)
let id: number; //Build ID (default to undefined)

const { status, data } = await apiInstance.getBuildLogs(xUserID, id);
```

### Parameters

| Name        | Type         | Description | Notes                 |
| ----------- | ------------ | ----------- | --------------------- |
| **xUserID** | [**string**] | User ID     | defaults to undefined |
| **id**      | [**number**] | Build ID    | defaults to undefined |

### Return type

**ResponseBuildLogs**

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

# **streamBuildLogs**

> streamBuildLogs()

### Example

```typescript
import { BuildApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new BuildApi(configuration);

let xUserID: string; //User ID (default to undefined)
let id: number; //Build ID (default to undefined)

const { status, data } = await apiInstance.streamBuildLogs(xUserID, id);
```

### Parameters

| Name        | Type         | Description | Notes                 |
| ----------- | ------------ | ----------- | --------------------- |
| **xUserID** | [**string**] | User ID     | defaults to undefined |
| **id**      | [**number**] | Build ID    | defaults to undefined |

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

### HTTP response details

| Status code | Description | Response headers                                                                                        |
| ----------- | ----------- | ------------------------------------------------------------------------------------------------------- |
| **200**     | OK          | _ Cache-Control - no-cache <br> _ Connection - keep-alive <br> \* Content-Type - text/event-stream <br> |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)
