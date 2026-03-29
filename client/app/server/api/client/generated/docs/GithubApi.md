# GithubApi

All URIs are relative to _http://localhost_

| Method                                  | HTTP request                                  | Description      |
| --------------------------------------- | --------------------------------------------- | ---------------- |
| [**getRepositories**](#getrepositories) | **GET** /github/repositories/{organizationId} | Get Repositories |

# **getRepositories**

> Array<ResponseRepository> getRepositories()

### Example

```typescript
import { GithubApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new GithubApi(configuration);

let xUserID: string; //User ID (default to undefined)
let organizationId: number; //Organization ID (default to undefined)

const { status, data } = await apiInstance.getRepositories(
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

**Array<ResponseRepository>**

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
