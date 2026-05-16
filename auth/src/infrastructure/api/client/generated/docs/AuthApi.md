# AuthApi

All URIs are relative to _http://localhost_

| Method                                              | HTTP request                           | Description               |
| --------------------------------------------------- | -------------------------------------- | ------------------------- |
| [**sendResetPassword**](#sendresetpassword)         | **POST** /auth/send-reset-password     | Send password reset email |
| [**sendVerificationEmail**](#sendverificationemail) | **POST** /auth/send-verification-email | Send email verification   |

# **sendResetPassword**

> sendResetPassword(data)

### Example

```typescript
import { AuthApi, Configuration, RequestSendResetPasswordRequest } from "./api";

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let data: RequestSendResetPasswordRequest; //Password reset

const { status, data } = await apiInstance.sendResetPassword(data);
```

### Parameters

| Name     | Type                                | Description    | Notes |
| -------- | ----------------------------------- | -------------- | ----- |
| **data** | **RequestSendResetPasswordRequest** | Password reset |       |

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
  AuthApi,
  Configuration,
  RequestSendVerificationEmailRequest,
} from "./api";

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let data: RequestSendVerificationEmailRequest; //Verification

const { status, data } = await apiInstance.sendVerificationEmail(data);
```

### Parameters

| Name     | Type                                    | Description  | Notes |
| -------- | --------------------------------------- | ------------ | ----- |
| **data** | **RequestSendVerificationEmailRequest** | Verification |       |

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
