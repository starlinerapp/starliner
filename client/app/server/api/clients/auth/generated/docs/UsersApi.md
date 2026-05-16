# UsersApi

All URIs are relative to _http://localhost_

| Method                                | HTTP request    | Description         |
| ------------------------------------- | --------------- | ------------------- |
| [**bulkUserLookup**](#bulkuserlookup) | **POST** /users | Look up users by ID |

# **bulkUserLookup**

> BulkUserLookupResponse bulkUserLookup(bulkUserLookupRequest)

Returns user profiles for the given IDs. Unknown IDs are omitted. At most 200 unique IDs per request.

### Example

```typescript
import { UsersApi, Configuration, BulkUserLookupRequest } from "./api";

const configuration = new Configuration();
const apiInstance = new UsersApi(configuration);

let bulkUserLookupRequest: BulkUserLookupRequest; //

const { status, data } = await apiInstance.bulkUserLookup(
  bulkUserLookupRequest,
);
```

### Parameters

| Name                      | Type                      | Description | Notes |
| ------------------------- | ------------------------- | ----------- | ----- |
| **bulkUserLookupRequest** | **BulkUserLookupRequest** |             |       |

### Return type

**BulkUserLookupResponse**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

### HTTP response details

| Status code | Description                       | Response headers |
| ----------- | --------------------------------- | ---------------- |
| **200**     | Matching user profiles            | -                |
| **400**     | Invalid JSON body or too many IDs | -                |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)
