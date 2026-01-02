# ProjectApi

All URIs are relative to _http://localhost_

| Method                              | HTTP request           | Description    |
| ----------------------------------- | ---------------------- | -------------- |
| [**createProject**](#createproject) | **POST** /projects     | Create Project |
| [**getProject**](#getproject)       | **GET** /projects/{id} | Get Project    |

# **createProject**

> createProject(data)

### Example

```typescript
import { ProjectApi, Configuration, RequestCreateProject } from "./api";

const configuration = new Configuration();
const apiInstance = new ProjectApi(configuration);

let xUserID: string; //User ID (default to undefined)
let data: RequestCreateProject; //Create Project

const { status, data } = await apiInstance.createProject(xUserID, data);
```

### Parameters

| Name        | Type                     | Description    | Notes                 |
| ----------- | ------------------------ | -------------- | --------------------- |
| **data**    | **RequestCreateProject** | Create Project |                       |
| **xUserID** | [**string**]             | User ID        | defaults to undefined |

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

# **getProject**

> ResponseProject getProject()

### Example

```typescript
import { ProjectApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new ProjectApi(configuration);

let xUserID: string; //User ID (default to undefined)
let id: number; //Project ID (default to undefined)

const { status, data } = await apiInstance.getProject(xUserID, id);
```

### Parameters

| Name        | Type         | Description | Notes                 |
| ----------- | ------------ | ----------- | --------------------- |
| **xUserID** | [**string**] | User ID     | defaults to undefined |
| **id**      | [**number**] | Project ID  | defaults to undefined |

### Return type

**ResponseProject**

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
