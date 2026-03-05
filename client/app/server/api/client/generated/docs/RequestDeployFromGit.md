# RequestDeployFromGit

## Properties

| Name                      | Type                                               | Description | Notes                  |
| ------------------------- | -------------------------------------------------- | ----------- | ---------------------- |
| **dockerfilePath**        | **string**                                         |             | [default to undefined] |
| **environmentId**         | **number**                                         |             | [default to undefined] |
| **envs**                  | [**Array&lt;RequestEnvVar&gt;**](RequestEnvVar.md) |             | [default to undefined] |
| **gitUrl**                | **string**                                         |             | [default to undefined] |
| **port**                  | **number**                                         |             | [default to undefined] |
| **projectRepositoryPath** | **string**                                         |             | [default to undefined] |
| **serviceName**           | **string**                                         |             | [default to undefined] |

## Example

```typescript
import { RequestDeployFromGit } from "./api";

const instance: RequestDeployFromGit = {
  dockerfilePath,
  environmentId,
  envs,
  gitUrl,
  port,
  projectRepositoryPath,
  serviceName,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
