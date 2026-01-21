# DeploymentApi

All URIs are relative to _http://localhost_

| Method                                | HTTP request                           | Description     |
| ------------------------------------- | -------------------------------------- | --------------- |
| [**deleteDatabase**](#deletedatabase) | **DELETE** /deployments/databases/{id} | Delete database |
| [**deployDatabase**](#deploydatabase) | **POST** /deployments/databases        | Deploy database |

# **deleteDatabase**

> deleteDatabase()

### Example

```typescript
import { DeploymentApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new DeploymentApi(configuration);

let xUserID: string; //User ID (default to undefined)
let id: number; //Deployment ID (default to undefined)

const { status, data } = await apiInstance.deleteDatabase(xUserID, id);
```

### Parameters

| Name        | Type         | Description   | Notes                 |
| ----------- | ------------ | ------------- | --------------------- |
| **xUserID** | [**string**] | User ID       | defaults to undefined |
| **id**      | [**number**] | Deployment ID | defaults to undefined |

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

# **deployDatabase**

> deployDatabase(data)

### Example

```typescript
import { DeploymentApi, Configuration, RequestDeployDatabase } from "./api";

const configuration = new Configuration();
const apiInstance = new DeploymentApi(configuration);

let xUserID: string; //User ID (default to undefined)
let data: RequestDeployDatabase; //Deploy Database

const { status, data } = await apiInstance.deployDatabase(xUserID, data);
```

### Parameters

| Name        | Type                      | Description     | Notes                 |
| ----------- | ------------------------- | --------------- | --------------------- |
| **data**    | **RequestDeployDatabase** | Deploy Database |                       |
| **xUserID** | [**string**]              | User ID         | defaults to undefined |

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
