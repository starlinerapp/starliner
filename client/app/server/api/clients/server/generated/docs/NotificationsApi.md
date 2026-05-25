# NotificationsApi

All URIs are relative to _http://localhost_

| Method                                                      | HTTP request           | Description                 |
| ----------------------------------------------------------- | ---------------------- | --------------------------- |
| [**streamGlobalNotifications**](#streamglobalnotifications) | **GET** /notifications | Stream global notifications |

# **streamGlobalNotifications**

> streamGlobalNotifications()

### Example

```typescript
import { NotificationsApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new NotificationsApi(configuration);

let xUserID: string; //User ID (default to undefined)
let organizationId: number; //Organization ID (default to undefined)

const { status, data } = await apiInstance.streamGlobalNotifications(
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
