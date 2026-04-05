# GithubappApi

All URIs are relative to _http://localhost_

| Method                                  | HTTP request                         | Description       |
| --------------------------------------- | ------------------------------------ | ----------------- |
| [**createGithubApp**](#creategithubapp) | **POST** /githubapps                 | Create GitHub App |
| [**getGithubApp**](#getgithubapp)       | **GET** /githubapps/{organizationId} | Get GitHub App    |

# **createGithubApp**

> createGithubApp(data)

### Example

```typescript
import { GithubappApi, Configuration, RequestCreateGithubApp } from "./api";

const configuration = new Configuration();
const apiInstance = new GithubappApi(configuration);

let xUserID: string; //User ID (default to undefined)
let data: RequestCreateGithubApp; //Create GitHub App

const { status, data } = await apiInstance.createGithubApp(xUserID, data);
```

### Parameters

| Name        | Type                       | Description       | Notes                 |
| ----------- | -------------------------- | ----------------- | --------------------- |
| **data**    | **RequestCreateGithubApp** | Create GitHub App |                       |
| **xUserID** | [**string**]               | User ID           | defaults to undefined |

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

# **getGithubApp**

> ResponseGithubApp getGithubApp()

### Example

```typescript
import { GithubappApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new GithubappApi(configuration);

let xUserID: string; //User ID (default to undefined)
let organizationId: number; //Organization ID (default to undefined)

const { status, data } = await apiInstance.getGithubApp(
  xUserID,
  organizationId,
);
```

### Parameters

| Name               | Type         | Description     | Notes                 |
| ------------------ | ------------ | --------------- | --------------------- |
| **xUserID**        | [**string**] | User ID         | defaults to undefined |
| **organizationId** | [**number**] | Organization ID | defaults to undefined |

### Return type

**ResponseGithubApp**

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
