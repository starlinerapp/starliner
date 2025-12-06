# CreateEnvironmentApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createEnvironment**](#createenvironment) | **POST** /environments | Create Environment|

# **createEnvironment**
> createEnvironment(data)


### Example

```typescript
import {
    CreateEnvironmentApi,
    Configuration,
    RequestCreateEnvironment
} from './api';

const configuration = new Configuration();
const apiInstance = new CreateEnvironmentApi(configuration);

let data: RequestCreateEnvironment; //Create Environment

const { status, data } = await apiInstance.createEnvironment(
    data
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **data** | **RequestCreateEnvironment**| Create Environment | |


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

