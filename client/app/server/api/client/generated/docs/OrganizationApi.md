# OrganizationApi

All URIs are relative to _http://localhost_

| Method                                                  | HTTP request                         | Description               |
| ------------------------------------------------------- | ------------------------------------ | ------------------------- |
| [**createOrganization**](#createorganization)           | **POST** /organizations              | Create organization       |
| [**getOrganizationClusters**](#getorganizationclusters) | **GET** /organizations/{id}/clusters | Get Organization Clusters |
| [**getOrganizationProjects**](#getorganizationprojects) | **GET** /organizations/{id}/projects | Get Organization Projects |
| [**getUserOrganizations**](#getuserorganizations)       | **GET** /organizations               | Get user organizations    |

# **createOrganization**

> createOrganization(data)

### Example

```typescript
import {
  OrganizationApi,
  Configuration,
  RequestCreateOrganization,
} from "./api";

const configuration = new Configuration();
const apiInstance = new OrganizationApi(configuration);

let data: RequestCreateOrganization; //Create Organization

const { status, data } = await apiInstance.createOrganization(data);
```

### Parameters

| Name     | Type                          | Description         | Notes |
| -------- | ----------------------------- | ------------------- | ----- |
| **data** | **RequestCreateOrganization** | Create Organization |       |

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

# **getOrganizationClusters**

> Array<ResponseCluster> getOrganizationClusters()

### Example

```typescript
import { OrganizationApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new OrganizationApi(configuration);

let xUserID: string; //User ID (default to undefined)
let id: number; //Organization ID (default to undefined)

const { status, data } = await apiInstance.getOrganizationClusters(xUserID, id);
```

### Parameters

| Name        | Type         | Description     | Notes                 |
| ----------- | ------------ | --------------- | --------------------- |
| **xUserID** | [**string**] | User ID         | defaults to undefined |
| **id**      | [**number**] | Organization ID | defaults to undefined |

### Return type

**Array<ResponseCluster>**

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

# **getOrganizationProjects**

> Array<ResponseProject> getOrganizationProjects()

### Example

```typescript
import { OrganizationApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new OrganizationApi(configuration);

let id: number; //Organization ID (default to undefined)

const { status, data } = await apiInstance.getOrganizationProjects(id);
```

### Parameters

| Name   | Type         | Description     | Notes                 |
| ------ | ------------ | --------------- | --------------------- |
| **id** | [**number**] | Organization ID | defaults to undefined |

### Return type

**Array<ResponseProject>**

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

# **getUserOrganizations**

> Array<ResponseOrganization> getUserOrganizations()

### Example

```typescript
import { OrganizationApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new OrganizationApi(configuration);

const { status, data } = await apiInstance.getUserOrganizations();
```

### Parameters

This endpoint does not have any parameters.

### Return type

**Array<ResponseOrganization>**

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
