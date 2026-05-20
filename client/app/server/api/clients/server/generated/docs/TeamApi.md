# TeamApi

All URIs are relative to _http://localhost_

| Method                                          | HTTP request                            | Description                          |
| ----------------------------------------------- | --------------------------------------- | ------------------------------------ |
| [**addTeamMember**](#addteammember)             | **POST** /teams/{teamId}/members        | Add organization member to team      |
| [**createTeam**](#createteam)                   | **POST** /organizations/{id}/teams      | Create team                          |
| [**getTeamClusters**](#getteamclusters)         | **GET** /teams/{teamId}/clusters        | Get clusters assigned to a team      |
| [**getTeamMembers**](#getteammembers)           | **GET** /teams/{teamId}/members         | Get Team Members                     |
| [**getTeamRepositories**](#getteamrepositories) | **GET** /teams/{teamId}/repos           | Get repositories assigned to a team  |
| [**getUserTeams**](#getuserteams)               | **GET** /organizations/{id}/teams       | Get User Teams                       |
| [**joinTeam**](#jointeam)                       | **POST** /organizations/{id}/teams/join | Join a team by slug                  |
| [**removeTeamMember**](#removeteammember)       | **DELETE** /teams/{teamId}/members      | Remove organization member from team |
| [**setTeamClusters**](#setteamclusters)         | **PUT** /teams/{teamId}/clusters        | Set clusters assigned to a team      |
| [**setTeamRepositories**](#setteamrepositories) | **PUT** /teams/{teamId}/repos           | Set repositories assigned to a team  |

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

# **setTeamClusters**

> setTeamClusters(data)

### Example

```typescript
import { TeamApi, Configuration, RequestSetTeamClusters } from "./api";

const configuration = new Configuration();
const apiInstance = new TeamApi(configuration);

let xUserID: string; //User ID (default to undefined)
let teamId: number; //Team ID (default to undefined)
let data: RequestSetTeamClusters; //Team Clusters

const { status, data } = await apiInstance.setTeamClusters(
  xUserID,
  teamId,
  data,
);
```

### Parameters

| Name        | Type                       | Description   | Notes                 |
| ----------- | -------------------------- | ------------- | --------------------- |
| **data**    | **RequestSetTeamClusters** | Team Clusters |                       |
| **xUserID** | [**string**]               | User ID       | defaults to undefined |
| **teamId**  | [**number**]               | Team ID       | defaults to undefined |

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

# **setTeamRepositories**

> setTeamRepositories(data)

### Example

```typescript
import { TeamApi, Configuration, RequestSetTeamRepositories } from "./api";

const configuration = new Configuration();
const apiInstance = new TeamApi(configuration);

let xUserID: string; //User ID (default to undefined)
let teamId: number; //Team ID (default to undefined)
let data: RequestSetTeamRepositories; //Team Repositories

const { status, data } = await apiInstance.setTeamRepositories(
  xUserID,
  teamId,
  data,
);
```

### Parameters

| Name        | Type                           | Description       | Notes                 |
| ----------- | ------------------------------ | ----------------- | --------------------- |
| **data**    | **RequestSetTeamRepositories** | Team Repositories |                       |
| **xUserID** | [**string**]                   | User ID           | defaults to undefined |
| **teamId**  | [**number**]                   | Team ID           | defaults to undefined |

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
