# TeamApi

All URIs are relative to _http://localhost_

| Method                                    | HTTP request                                        | Description              |
| ----------------------------------------- | --------------------------------------------------- | ------------------------ |
| [**addTeamMember**](#addteammember)       | **POST** /teams/{organizationId}/{teamId}/members   | Add current user to team |
| [**createTeam**](#createteam)             | **POST** /teams/{organizationId}                    | Create team              |
| [**getTeamMembers**](#getteammembers)     | **GET** /teams/{organizationId}/{teamId}/members    | Get Team Members         |
| [**getUserTeams**](#getuserteams)         | **GET** /teams/{organizationId}                     | Get User Teams           |
| [**joinTeam**](#jointeam)                 | **POST** /teams/{organizationId}/join               | Join a team by slug      |
| [**removeTeamMember**](#removeteammember) | **DELETE** /teams/{organizationId}/{teamId}/members | Remove Team Member       |

# **addTeamMember**

> addTeamMember()

### Example

```typescript
import { TeamApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new TeamApi(configuration);

let xUserID: string; //User ID (default to undefined)
let organizationId: number; //Organization ID (default to undefined)
let teamId: number; //Team ID (default to undefined)

const { status, data } = await apiInstance.addTeamMember(
  xUserID,
  organizationId,
  teamId,
);
```

### Parameters

| Name               | Type         | Description     | Notes                 |
| ------------------ | ------------ | --------------- | --------------------- |
| **xUserID**        | [**string**] | User ID         | defaults to undefined |
| **organizationId** | [**number**] | Organization ID | defaults to undefined |
| **teamId**         | [**number**] | Team ID         | defaults to undefined |

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
let organizationId: number; //Organization ID (default to undefined)
let data: RequestCreateTeam; //Create Team

const { status, data } = await apiInstance.createTeam(
  xUserID,
  organizationId,
  data,
);
```

### Parameters

| Name               | Type                  | Description     | Notes                 |
| ------------------ | --------------------- | --------------- | --------------------- |
| **data**           | **RequestCreateTeam** | Create Team     |                       |
| **xUserID**        | [**string**]          | User ID         | defaults to undefined |
| **organizationId** | [**number**]          | Organization ID | defaults to undefined |

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

# **getTeamMembers**

> Array<ResponseUser> getTeamMembers()

### Example

```typescript
import { TeamApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new TeamApi(configuration);

let xUserID: string; //User ID (default to undefined)
let organizationId: number; //Organization ID (default to undefined)
let teamId: number; //Team ID (default to undefined)

const { status, data } = await apiInstance.getTeamMembers(
  xUserID,
  organizationId,
  teamId,
);
```

### Parameters

| Name               | Type         | Description     | Notes                 |
| ------------------ | ------------ | --------------- | --------------------- |
| **xUserID**        | [**string**] | User ID         | defaults to undefined |
| **organizationId** | [**number**] | Organization ID | defaults to undefined |
| **teamId**         | [**number**] | Team ID         | defaults to undefined |

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

# **getUserTeams**

> Array<ResponseTeam> getUserTeams()

### Example

```typescript
import { TeamApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new TeamApi(configuration);

let xUserID: string; //User ID (default to undefined)
let organizationId: number; //Organization ID (default to undefined)

const { status, data } = await apiInstance.getUserTeams(
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
let organizationId: number; //Organization ID (default to undefined)
let data: RequestJoinTeam; //Join Team

const { status, data } = await apiInstance.joinTeam(
  xUserID,
  organizationId,
  data,
);
```

### Parameters

| Name               | Type                | Description     | Notes                 |
| ------------------ | ------------------- | --------------- | --------------------- |
| **data**           | **RequestJoinTeam** | Join Team       |                       |
| **xUserID**        | [**string**]        | User ID         | defaults to undefined |
| **organizationId** | [**number**]        | Organization ID | defaults to undefined |

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

> removeTeamMember()

### Example

```typescript
import { TeamApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new TeamApi(configuration);

let xUserID: string; //User ID (default to undefined)
let organizationId: number; //Organization ID (default to undefined)
let teamId: number; //Team ID (default to undefined)

const { status, data } = await apiInstance.removeTeamMember(
  xUserID,
  organizationId,
  teamId,
);
```

### Parameters

| Name               | Type         | Description     | Notes                 |
| ------------------ | ------------ | --------------- | --------------------- |
| **xUserID**        | [**string**] | User ID         | defaults to undefined |
| **organizationId** | [**number**] | Organization ID | defaults to undefined |
| **teamId**         | [**number**] | Team ID         | defaults to undefined |

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
