# ClusterApi

All URIs are relative to _http://localhost_

| Method                              | HTTP request              | Description    |
| ----------------------------------- | ------------------------- | -------------- |
| [**createCluster**](#createcluster) | **POST** /clusters        | Create Cluster |
| [**deleteCluster**](#deletecluster) | **DELETE** /clusters/{id} | Delete Cluster |
| [**getCluster**](#getcluster)       | **GET** /clusters/{id}    | Get Cluster    |

# **createCluster**

> createCluster(data)

### Example

```typescript
import { ClusterApi, Configuration, RequestCreateCluster } from "./api";

const configuration = new Configuration();
const apiInstance = new ClusterApi(configuration);

let xUserID: string; //User ID (default to undefined)
let data: RequestCreateCluster; //Create Cluster

const { status, data } = await apiInstance.createCluster(xUserID, data);
```

### Parameters

| Name        | Type                     | Description    | Notes                 |
| ----------- | ------------------------ | -------------- | --------------------- |
| **data**    | **RequestCreateCluster** | Create Cluster |                       |
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

# **deleteCluster**

> deleteCluster()

### Example

```typescript
import { ClusterApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new ClusterApi(configuration);

let xUserID: string; //User ID (default to undefined)
let id: number; //Cluster ID (default to undefined)

const { status, data } = await apiInstance.deleteCluster(xUserID, id);
```

### Parameters

| Name        | Type         | Description | Notes                 |
| ----------- | ------------ | ----------- | --------------------- |
| **xUserID** | [**string**] | User ID     | defaults to undefined |
| **id**      | [**number**] | Cluster ID  | defaults to undefined |

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

# **getCluster**

> ResponseCluster getCluster()

### Example

```typescript
import { ClusterApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new ClusterApi(configuration);

let xUserID: string; //User ID (default to undefined)
let id: number; //Cluster ID (default to undefined)

const { status, data } = await apiInstance.getCluster(xUserID, id);
```

### Parameters

| Name        | Type         | Description | Notes                 |
| ----------- | ------------ | ----------- | --------------------- |
| **xUserID** | [**string**] | User ID     | defaults to undefined |
| **id**      | [**number**] | Cluster ID  | defaults to undefined |

### Return type

**ResponseCluster**

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
