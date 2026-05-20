# GithubApi

All URIs are relative to _http://localhost_

| Method                                              | HTTP request                                                                | Description                                   |
| --------------------------------------------------- | --------------------------------------------------------------------------- | --------------------------------------------- |
| [**getAllRepositories**](#getallrepositories)       | **GET** /github/all-repositories/{organizationId}                           | Get All Repositories (owner only, unfiltered) |
| [**getFileContent**](#getfilecontent)               | **GET** /github/repositories/{organizationId}/{owner}/{repository}/file     | Get File Content                              |
| [**getRepositories**](#getrepositories)             | **GET** /github/repositories/{organizationId}                               | Get Repositories                              |
| [**getRepositoryContents**](#getrepositorycontents) | **GET** /github/repositories/{organizationId}/{owner}/{repository}/contents | Get Repository Content                        |

# **getAllRepositories**

> Array<ResponseRepository> getAllRepositories()

### Example

```typescript
import { GithubApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new GithubApi(configuration);

let xUserID: string; //User ID (default to undefined)
let organizationId: number; //Organization ID (default to undefined)

const { status, data } = await apiInstance.getAllRepositories(
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

# **getFileContent**

> ResponseFileContent getFileContent()

### Example

```typescript
import { GithubApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new GithubApi(configuration);

let xUserID: string; //User ID (default to undefined)
let organizationId: number; //Organization ID (default to undefined)
let owner: string; //Repository owner (user or org) (default to undefined)
let repository: string; //Repository name (default to undefined)
let path: string; //Path to the file within the repository (default to undefined)

const { status, data } = await apiInstance.getFileContent(
  xUserID,
  organizationId,
  owner,
  repository,
  path,
);
```

### Parameters

| Name               | Type         | Description                            | Notes                 |
| ------------------ | ------------ | -------------------------------------- | --------------------- |
| **xUserID**        | [**string**] | User ID                                | defaults to undefined |
| **organizationId** | [**number**] | Organization ID                        | defaults to undefined |
| **owner**          | [**string**] | Repository owner (user or org)         | defaults to undefined |
| **repository**     | [**string**] | Repository name                        | defaults to undefined |
| **path**           | [**string**] | Path to the file within the repository | defaults to undefined |

### Return type

**ResponseFileContent**

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

# **getRepositoryContents**

> Array<ResponseRepositoryFile> getRepositoryContents()

### Example

```typescript
import { GithubApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new GithubApi(configuration);

let xUserID: string; //User ID (default to undefined)
let organizationId: number; //Organization ID (default to undefined)
let owner: string; //Repository owner (user or org) (default to undefined)
let repository: string; //Repository name (default to undefined)
let path: string; //Path within the repository (e.g., src or src/main.go) (optional) (default to undefined)

const { status, data } = await apiInstance.getRepositoryContents(
  xUserID,
  organizationId,
  owner,
  repository,
  path,
);
```

### Parameters

| Name               | Type         | Description                                           | Notes                            |
| ------------------ | ------------ | ----------------------------------------------------- | -------------------------------- |
| **xUserID**        | [**string**] | User ID                                               | defaults to undefined            |
| **organizationId** | [**number**] | Organization ID                                       | defaults to undefined            |
| **owner**          | [**string**] | Repository owner (user or org)                        | defaults to undefined            |
| **repository**     | [**string**] | Repository name                                       | defaults to undefined            |
| **path**           | [**string**] | Path within the repository (e.g., src or src/main.go) | (optional) defaults to undefined |

### Return type

**Array<ResponseRepositoryFile>**

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
