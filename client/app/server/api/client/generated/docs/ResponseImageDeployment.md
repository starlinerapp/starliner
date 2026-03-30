# ResponseImageDeployment

## Properties

| Name                 | Type                                                 | Description | Notes                             |
| -------------------- | ---------------------------------------------------- | ----------- | --------------------------------- |
| **envVars**          | [**Array&lt;ResponseEnvVar&gt;**](ResponseEnvVar.md) |             | [default to undefined]            |
| **id**               | **number**                                           |             | [default to undefined]            |
| **imageName**        | **string**                                           |             | [default to undefined]            |
| **internalEndpoint** | **string**                                           |             | [default to undefined]            |
| **port**             | **string**                                           |             | [default to undefined]            |
| **serviceName**      | **string**                                           |             | [default to undefined]            |
| **status**           | **string**                                           |             | [default to undefined]            |
| **tag**              | **string**                                           |             | [default to undefined]            |
| **volumeSizeMB**     | **number**                                           |             | [optional] [default to undefined] |

## Example

```typescript
import { ResponseImageDeployment } from "./api";

const instance: ResponseImageDeployment = {
  envVars,
  id,
  imageName,
  internalEndpoint,
  port,
  serviceName,
  status,
  tag,
  volumeSizeMB,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
