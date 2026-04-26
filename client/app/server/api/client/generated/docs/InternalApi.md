# InternalApi

All URIs are relative to _http://localhost_

| Method                                              | HTTP request                               | Description               |
| --------------------------------------------------- | ------------------------------------------ | ------------------------- |
| [**sendResetPassword**](#sendresetpassword)         | **POST** /internal/send-reset-password     | Send password reset email |
| [**sendVerificationEmail**](#sendverificationemail) | **POST** /internal/send-verification-email | Send email verification   |

# **sendResetPassword**

> sendResetPassword(data)

### Example

```typescript
import {
  InternalApi,
  Configuration,
  RequestSendResetPasswordRequest,
} from "./api";

const configuration = new Configuration();
const apiInstance = new InternalApi(configuration);

let xUserID: string; //User ID (default to undefined)
let data: RequestSendResetPasswordRequest; //Password reset

const { status, data } = await apiInstance.sendResetPassword(xUserID, data);
```

### Parameters

| Name        | Type                                | Description    | Notes                 |
| ----------- | ----------------------------------- | -------------- | --------------------- |
| **data**    | **RequestSendResetPasswordRequest** | Password reset |                       |
| **xUserID** | [**string**]                        | User ID        | defaults to undefined |

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

# **sendVerificationEmail**

> sendVerificationEmail(data)

### Example

```typescript
import {
  InternalApi,
  Configuration,
  RequestSendVerificationEmailRequest,
} from "./api";

const configuration = new Configuration();
const apiInstance = new InternalApi(configuration);

let xUserID: string; //User ID (default to undefined)
let data: RequestSendVerificationEmailRequest; //Verification

const { status, data } = await apiInstance.sendVerificationEmail(xUserID, data);
```

### Parameters

| Name        | Type                                    | Description  | Notes                 |
| ----------- | --------------------------------------- | ------------ | --------------------- |
| **data**    | **RequestSendVerificationEmailRequest** | Verification |                       |
| **xUserID** | [**string**]                            | User ID      | defaults to undefined |

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
