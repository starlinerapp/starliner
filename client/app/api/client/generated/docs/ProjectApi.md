# ProjectApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createProject**](#createproject) | **POST** /projects | Create Project|
|[**getProject**](#getproject) | **GET** /projects | Get Project|

# **createProject**
> createProject(data)


### Example

```typescript
import {
    ProjectApi,
    Configuration,
    RequestCreateProject
} from './api';

const configuration = new Configuration();
const apiInstance = new ProjectApi(configuration);

let data: RequestCreateProject; //Create Project

const { status, data } = await apiInstance.createProject(
    data
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **data** | **RequestCreateProject**| Create Project | |


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getProject**
> ResponseProject getProject(data)


### Example

```typescript
import {
    ProjectApi,
    Configuration,
    RequestGetProject
} from './api';

const configuration = new Configuration();
const apiInstance = new ProjectApi(configuration);

let data: RequestGetProject; //Get Project

const { status, data } = await apiInstance.getProject(
    data
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **data** | **RequestGetProject**| Get Project | |


### Return type

**ResponseProject**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

