# TeamApi

All URIs are relative to _http://localhost_

| Method                                                  | HTTP request                                    | Description                              |
| ------------------------------------------------------- | ----------------------------------------------- | ---------------------------------------- |
| [**addTeamMember**](#addteammember)                     | **POST** /teams/{teamId}/members                | Add organization member to team          |
| [**assignClusterToTeam**](#assignclustertoteam)         | **POST** /teams/{teamId}/clusters/{clusterId}   | Assign a cluster to a team               |
| [**assignRepoToTeam**](#assignrepototeam)               | **POST** /teams/{teamId}/repos                  | Assign a GitHub repository to a team     |
| [**createTeam**](#createteam)                           | **POST** /organizations/{id}/teams              | Create team                              |
| [**getTeamClusters**](#getteamclusters)                 | **GET** /teams/{teamId}/clusters                | Get clusters assigned to a team          |
| [**getTeamMembers**](#getteammembers)                   | **GET** /teams/{teamId}/members                 | Get Team Members                         |
| [**getTeamRepositories**](#getteamrepositories)         | **GET** /teams/{teamId}/repos                   | Get repositories assigned to a team      |
| [**getUserTeams**](#getuserteams)                       | **GET** /organizations/{id}/teams               | Get User Teams                           |
| [**joinTeam**](#jointeam)                               | **POST** /organizations/{id}/teams/join         | Join a team by slug                      |
| [**removeTeamMember**](#removeteammember)               | **DELETE** /teams/{teamId}/members              | Remove organization member from team     |
| [**unassignClusterFromTeam**](#unassignclusterfromteam) | **DELETE** /teams/{teamId}/clusters/{clusterId} | Unassign a cluster from a team           |
| [**unassignRepoFromTeam**](#unassignrepofromteam)       | **DELETE** /teams/{teamId}/repos/{repoId}       | Unassign a GitHub repository from a team |

# **addTeamMember**

> addTeamMember(data)

### Example

```typescript
import { TeamApi, Configuration, RequestAddTeamMember } from "./api";

const configuration = new Configuration();
const apiInstance = new TeamApi(configuration);

let xUserID: string; //User ID (default to undefined)
let teamId: number; //Team ID (default to undefined)
let data: RequestAddTeamMember; //ID of member to add

const { status, data } = await apiInstance.addTeamMember(xUserID, teamId, data);
```

### Parameters

| Name        | Type                     | Description         | Notes                 |
| ----------- | ------------------------ | ------------------- | --------------------- |
| **data**    | **RequestAddTeamMember** | ID of member to add |                       |
| **xUserID** | [**string**]             | User ID             | defaults to undefined |
| **teamId**  | [**number**]             | Team ID             | defaults to undefined |

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

# **assignClusterToTeam**

> assignClusterToTeam()

### Example

```typescript
import { TeamApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new TeamApi(configuration);

let xUserID: string; //User ID (default to undefined)
let teamId: number; //Team ID (default to undefined)
let clusterId: number; //Cluster ID (default to undefined)

const { status, data } = await apiInstance.assignClusterToTeam(
  xUserID,
  teamId,
  clusterId,
);
```

### Parameters

| Name          | Type         | Description | Notes                 |
| ------------- | ------------ | ----------- | --------------------- |
| **xUserID**   | [**string**] | User ID     | defaults to undefined |
| **teamId**    | [**number**] | Team ID     | defaults to undefined |
| **clusterId** | [**number**] | Cluster ID  | defaults to undefined |

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

# **assignRepoToTeam**

> assignRepoToTeam(data)

### Example

```typescript
import { TeamApi, Configuration, RequestAssignRepoToTeam } from "./api";

const configuration = new Configuration();
const apiInstance = new TeamApi(configuration);

let xUserID: string; //User ID (default to undefined)
let teamId: number; //Team ID (default to undefined)
let data: RequestAssignRepoToTeam; //Assign Repo

const { status, data } = await apiInstance.assignRepoToTeam(
  xUserID,
  teamId,
  data,
);
```

### Parameters

| Name        | Type                        | Description | Notes                 |
| ----------- | --------------------------- | ----------- | --------------------- |
| **data**    | **RequestAssignRepoToTeam** | Assign Repo |                       |
| **xUserID** | [**string**]                | User ID     | defaults to undefined |
| **teamId**  | [**number**]                | Team ID     | defaults to undefined |

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

# **createTeam**

> ResponseTeam createTeam(data)

### Example

```typescript
import { TeamApi, Configuration, RequestCreateTeam } from "./api";

const configuration = new Configuration();
const apiInstance = new TeamApi(configuration);

let xUserID: string; //User ID (default to undefined)
let id: number; //Organization ID (default to undefined)
let data: RequestCreateTeam; //Team slug (lowercase, alphanumeric, hyphens only)

const { status, data } = await apiInstance.createTeam(xUserID, id, data);
```

### Parameters

| Name        | Type                  | Description                                       | Notes                 |
| ----------- | --------------------- | ------------------------------------------------- | --------------------- |
| **data**    | **RequestCreateTeam** | Team slug (lowercase, alphanumeric, hyphens only) |                       |
| **xUserID** | [**string**]          | User ID                                           | defaults to undefined |
| **id**      | [**number**]          | Organization ID                                   | defaults to undefined |

### Return type

**ResponseTeam**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: _/_

### HTTP response details

| Status code | Description | Response headers |
| ----------- | ----------- | ---------------- |
| **201**     | Created     | -                |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getTeamClusters**

> Array<ResponseTeamCluster> getTeamClusters()

### Example

```typescript
import { TeamApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new TeamApi(configuration);

let xUserID: string; //User ID (default to undefined)
let teamId: number; //Team ID (default to undefined)

const { status, data } = await apiInstance.getTeamClusters(xUserID, teamId);
```

### Parameters

| Name        | Type         | Description | Notes                 |
| ----------- | ------------ | ----------- | --------------------- |
| **xUserID** | [**string**] | User ID     | defaults to undefined |
| **teamId**  | [**number**] | Team ID     | defaults to undefined |

### Return type

**Array<ResponseTeamCluster>**

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

# **getTeamMembers**

> Array<ResponseUser> getTeamMembers()

### Example

```typescript
import { TeamApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new TeamApi(configuration);

let xUserID: string; //User ID (default to undefined)
let teamId: number; //Team ID (default to undefined)

const { status, data } = await apiInstance.getTeamMembers(xUserID, teamId);
```

### Parameters

| Name        | Type         | Description | Notes                 |
| ----------- | ------------ | ----------- | --------------------- |
| **xUserID** | [**string**] | User ID     | defaults to undefined |
| **teamId**  | [**number**] | Team ID     | defaults to undefined |

### Return type

**Array<ResponseUser>**

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

# **getTeamRepositories**

> Array<ResponseTeamRepo> getTeamRepositories()

### Example

```typescript
import { TeamApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new TeamApi(configuration);

let xUserID: string; //User ID (default to undefined)
let teamId: number; //Team ID (default to undefined)

const { status, data } = await apiInstance.getTeamRepositories(xUserID, teamId);
```

### Parameters

| Name        | Type         | Description | Notes                 |
| ----------- | ------------ | ----------- | --------------------- |
| **xUserID** | [**string**] | User ID     | defaults to undefined |
| **teamId**  | [**number**] | Team ID     | defaults to undefined |

### Return type

**Array<ResponseTeamRepo>**

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

# **getUserTeams**

> Array<ResponseTeam> getUserTeams()

### Example

```typescript
import { TeamApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new TeamApi(configuration);

let xUserID: string; //User ID (default to undefined)
let id: number; //Organization ID (default to undefined)

const { status, data } = await apiInstance.getUserTeams(xUserID, id);
```

### Parameters

| Name        | Type         | Description     | Notes                 |
| ----------- | ------------ | --------------- | --------------------- |
| **xUserID** | [**string**] | User ID         | defaults to undefined |
| **id**      | [**number**] | Organization ID | defaults to undefined |

### Return type

**Array<ResponseTeam>**

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

# **joinTeam**

> joinTeam(data)

### Example

```typescript
import { TeamApi, Configuration, RequestJoinTeam } from "./api";

const configuration = new Configuration();
const apiInstance = new TeamApi(configuration);

let xUserID: string; //User ID (default to undefined)
let id: number; //Organization ID (default to undefined)
let data: RequestJoinTeam; //Join Team

const { status, data } = await apiInstance.joinTeam(xUserID, id, data);
```

### Parameters

| Name        | Type                | Description     | Notes                 |
| ----------- | ------------------- | --------------- | --------------------- |
| **data**    | **RequestJoinTeam** | Join Team       |                       |
| **xUserID** | [**string**]        | User ID         | defaults to undefined |
| **id**      | [**number**]        | Organization ID | defaults to undefined |

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

# **removeTeamMember**

> removeTeamMember(data)

### Example

```typescript
import { TeamApi, Configuration, RequestRemoveTeamMember } from "./api";

const configuration = new Configuration();
const apiInstance = new TeamApi(configuration);

let xUserID: string; //User ID (default to undefined)
let teamId: number; //Team ID (default to undefined)
let data: RequestRemoveTeamMember; //ID of the organization member to remove from the team

const { status, data } = await apiInstance.removeTeamMember(
  xUserID,
  teamId,
  data,
);
```

### Parameters

| Name        | Type                        | Description                                           | Notes                 |
| ----------- | --------------------------- | ----------------------------------------------------- | --------------------- |
| **data**    | **RequestRemoveTeamMember** | ID of the organization member to remove from the team |                       |
| **xUserID** | [**string**]                | User ID                                               | defaults to undefined |
| **teamId**  | [**number**]                | Team ID                                               | defaults to undefined |

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
| **204**     | No Content  | -                |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **unassignClusterFromTeam**

> unassignClusterFromTeam()

### Example

```typescript
import { TeamApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new TeamApi(configuration);

let xUserID: string; //User ID (default to undefined)
let teamId: number; //Team ID (default to undefined)
let clusterId: number; //Cluster ID (default to undefined)

const { status, data } = await apiInstance.unassignClusterFromTeam(
  xUserID,
  teamId,
  clusterId,
);
```

### Parameters

| Name          | Type         | Description | Notes                 |
| ------------- | ------------ | ----------- | --------------------- |
| **xUserID**   | [**string**] | User ID     | defaults to undefined |
| **teamId**    | [**number**] | Team ID     | defaults to undefined |
| **clusterId** | [**number**] | Cluster ID  | defaults to undefined |

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
| **204**     | No Content  | -                |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **unassignRepoFromTeam**

> unassignRepoFromTeam()

### Example

```typescript
import { TeamApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new TeamApi(configuration);

let xUserID: string; //User ID (default to undefined)
let teamId: number; //Team ID (default to undefined)
let repoId: number; //GitHub Repo ID (default to undefined)

const { status, data } = await apiInstance.unassignRepoFromTeam(
  xUserID,
  teamId,
  repoId,
);
```

### Parameters

| Name        | Type         | Description    | Notes                 |
| ----------- | ------------ | -------------- | --------------------- |
| **xUserID** | [**string**] | User ID        | defaults to undefined |
| **teamId**  | [**number**] | Team ID        | defaults to undefined |
| **repoId**  | [**number**] | GitHub Repo ID | defaults to undefined |

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
| **204**     | No Content  | -                |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)
