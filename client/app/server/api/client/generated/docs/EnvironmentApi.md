# EnvironmentApi

All URIs are relative to _http://localhost_

| Method                                                      | HTTP request                           | Description                 |
| ----------------------------------------------------------- | -------------------------------------- | --------------------------- |
| [**createEnvironment**](#createenvironment)                 | **POST** /environments                 | Create Environment          |
| [**getEnvironmentDeployments**](#getenvironmentdeployments) | **GET** /environments/{id}/deployments | Get Environment Deployments |

# **createEnvironment**

> createEnvironment(data)

### Example

```typescript
import { EnvironmentApi, Configuration, RequestCreateEnvironment } from "./api";

const configuration = new Configuration();
const apiInstance = new EnvironmentApi(configuration);

let xUserID: string; //User ID (default to undefined)
let data: RequestCreateEnvironment; //Create Environment

const { status, data } = await apiInstance.createEnvironment(xUserID, data);
```

### Parameters

| Name        | Type                         | Description        | Notes                 |
| ----------- | ---------------------------- | ------------------ | --------------------- |
| **data**    | **RequestCreateEnvironment** | Create Environment |                       |
| **xUserID** | [**string**]                 | User ID            | defaults to undefined |

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
| **201**     | Created     | -                |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getEnvironmentDeployments**

> ResponseDeployments getEnvironmentDeployments()

### Example

```typescript
import { EnvironmentApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new EnvironmentApi(configuration);

let xUserID: string; //User ID (default to undefined)
let id: number; //Environment ID (default to undefined)

const { status, data } = await apiInstance.getEnvironmentDeployments(
  xUserID,
  id,
);
```

### Parameters

| Name        | Type         | Description    | Notes                 |
| ----------- | ------------ | -------------- | --------------------- |
| **xUserID** | [**string**] | User ID        | defaults to undefined |
| **id**      | [**number**] | Environment ID | defaults to undefined |

### Return type

**ResponseDeployments**

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
