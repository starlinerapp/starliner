# RequestDeployImage

## Properties

| Name                | Type                                               | Description | Notes                             |
| ------------------- | -------------------------------------------------- | ----------- | --------------------------------- |
| **environmentId**   | **number**                                         |             | [default to undefined]            |
| **envs**            | [**Array&lt;RequestEnvVar&gt;**](RequestEnvVar.md) |             | [default to undefined]            |
| **imageName**       | **string**                                         |             | [default to undefined]            |
| **port**            | **number**                                         |             | [default to undefined]            |
| **serviceName**     | **string**                                         |             | [default to undefined]            |
| **tag**             | **string**                                         |             | [default to undefined]            |
| **volumeMountPath** | **string**                                         |             | [optional] [default to undefined] |
| **volumeSizeMiB**   | **number**                                         |             | [optional] [default to undefined] |

## Example

```typescript
import { RequestDeployImage } from "./api";

const instance: RequestDeployImage = {
  environmentId,
  envs,
  imageName,
  port,
  serviceName,
  tag,
  volumeMountPath,
  volumeSizeMiB,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
