# DeploymentApi

All URIs are relative to _http://localhost_

| Method                                                              | HTTP request                                  | Description                |
| ------------------------------------------------------------------- | --------------------------------------------- | -------------------------- |
| [**deleteDeployment**](#deletedeployment)                           | **DELETE** /deployments/{id}                  | Delete deployment          |
| [**deployDatabase**](#deploydatabase)                               | **POST** /deployments/databases               | Deploy database            |
| [**deployFromGitRepository**](#deployfromgitrepository)             | **POST** /deployments/git                     | Deploy from Git Repository |
| [**deployImage**](#deployimage)                                     | **POST** /deployments/images                  | Deploy image               |
| [**deployIngress**](#deployingress)                                 | **POST** /deployments/ingresses               | Deploy ingress             |
| [**streamDeploymentLogs**](#streamdeploymentlogs)                   | **GET** /deployments/{id}/logs                | Stream deployment logs     |
| [**updateDeployFromGitRepository**](#updatedeployfromgitrepository) | **PUT** /deployments/git/{deploymentId}       | Update Deploy from Git     |
| [**updateImageDeployment**](#updateimagedeployment)                 | **PUT** /deployments/images/{deploymentId}    | Update image deployment    |
| [**updateIngressDeployment**](#updateingressdeployment)             | **PUT** /deployments/ingresses/{deploymentId} | Update ingress deployment  |

# **deleteDeployment**

> deleteDeployment()

### Example

```typescript
import { DeploymentApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new DeploymentApi(configuration);

let xUserID: string; //User ID (default to undefined)
let id: number; //Deployment ID (default to undefined)

const { status, data } = await apiInstance.deleteDeployment(xUserID, id);
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

# **deployFromGitRepository**

> deployFromGitRepository(data)

### Example

```typescript
import { DeploymentApi, Configuration, RequestDeployFromGit } from "./api";

const configuration = new Configuration();
const apiInstance = new DeploymentApi(configuration);

let xUserID: string; //User ID (default to undefined)
let data: RequestDeployFromGit; //Deploy from Git

const { status, data } = await apiInstance.deployFromGitRepository(
  xUserID,
  data,
);
```

### Parameters

| Name        | Type                     | Description     | Notes                 |
| ----------- | ------------------------ | --------------- | --------------------- |
| **data**    | **RequestDeployFromGit** | Deploy from Git |                       |
| **xUserID** | [**string**]             | User ID         | defaults to undefined |

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

# **deployImage**

> deployImage(data)

### Example

```typescript
import { DeploymentApi, Configuration, RequestDeployImage } from "./api";

const configuration = new Configuration();
const apiInstance = new DeploymentApi(configuration);

let xUserID: string; //User ID (default to undefined)
let data: RequestDeployImage; //Deploy Image

const { status, data } = await apiInstance.deployImage(xUserID, data);
```

### Parameters

| Name        | Type                   | Description  | Notes                 |
| ----------- | ---------------------- | ------------ | --------------------- |
| **data**    | **RequestDeployImage** | Deploy Image |                       |
| **xUserID** | [**string**]           | User ID      | defaults to undefined |

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

# **deployIngress**

> deployIngress(data)

### Example

```typescript
import { DeploymentApi, Configuration, RequestDeployIngress } from "./api";

const configuration = new Configuration();
const apiInstance = new DeploymentApi(configuration);

let xUserID: string; //User ID (default to undefined)
let data: RequestDeployIngress; //Deploy Ingress

const { status, data } = await apiInstance.deployIngress(xUserID, data);
```

### Parameters

| Name        | Type                     | Description    | Notes                 |
| ----------- | ------------------------ | -------------- | --------------------- |
| **data**    | **RequestDeployIngress** | Deploy Ingress |                       |
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

# **streamDeploymentLogs**

> streamDeploymentLogs()

### Example

```typescript
import { DeploymentApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new DeploymentApi(configuration);

let xUserID: string; //User ID (default to undefined)
let id: number; //Deployment ID (default to undefined)

const { status, data } = await apiInstance.streamDeploymentLogs(xUserID, id);
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

| Status code | Description | Response headers                                                                                        |
| ----------- | ----------- | ------------------------------------------------------------------------------------------------------- |
| **200**     | OK          | _ Cache-Control - no-cache <br> _ Connection - keep-alive <br> \* Content-Type - text/event-stream <br> |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateDeployFromGitRepository**

> updateDeployFromGitRepository(data)

### Example

```typescript
import {
  DeploymentApi,
  Configuration,
  RequestUpdateDeployFromGit,
} from "./api";

const configuration = new Configuration();
const apiInstance = new DeploymentApi(configuration);

let xUserID: string; //User ID (default to undefined)
let deploymentId: number; //Deployment ID (default to undefined)
let data: RequestUpdateDeployFromGit; //Update Deploy from Git

const { status, data } = await apiInstance.updateDeployFromGitRepository(
  xUserID,
  deploymentId,
  data,
);
```

### Parameters

| Name             | Type                           | Description            | Notes                 |
| ---------------- | ------------------------------ | ---------------------- | --------------------- |
| **data**         | **RequestUpdateDeployFromGit** | Update Deploy from Git |                       |
| **xUserID**      | [**string**]                   | User ID                | defaults to undefined |
| **deploymentId** | [**number**]                   | Deployment ID          | defaults to undefined |

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

# **updateImageDeployment**

> updateImageDeployment(data)

### Example

```typescript
import { DeploymentApi, Configuration, RequestUpdateImage } from "./api";

const configuration = new Configuration();
const apiInstance = new DeploymentApi(configuration);

let xUserID: string; //User ID (default to undefined)
let deploymentId: number; //Deployment ID (default to undefined)
let data: RequestUpdateImage; //Update Image

const { status, data } = await apiInstance.updateImageDeployment(
  xUserID,
  deploymentId,
  data,
);
```

### Parameters

| Name             | Type                   | Description   | Notes                 |
| ---------------- | ---------------------- | ------------- | --------------------- |
| **data**         | **RequestUpdateImage** | Update Image  |                       |
| **xUserID**      | [**string**]           | User ID       | defaults to undefined |
| **deploymentId** | [**number**]           | Deployment ID | defaults to undefined |

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

# **updateIngressDeployment**

> updateIngressDeployment(data)

### Example

```typescript
import { DeploymentApi, Configuration, RequestUpdateIngress } from "./api";

const configuration = new Configuration();
const apiInstance = new DeploymentApi(configuration);

let xUserID: string; //User ID (default to undefined)
let deploymentId: number; //Deployment ID (default to undefined)
let data: RequestUpdateIngress; //Update Ingress

const { status, data } = await apiInstance.updateIngressDeployment(
  xUserID,
  deploymentId,
  data,
);
```

### Parameters

| Name             | Type                     | Description    | Notes                 |
| ---------------- | ------------------------ | -------------- | --------------------- |
| **data**         | **RequestUpdateIngress** | Update Ingress |                       |
| **xUserID**      | [**string**]             | User ID        | defaults to undefined |
| **deploymentId** | [**number**]             | Deployment ID  | defaults to undefined |

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
