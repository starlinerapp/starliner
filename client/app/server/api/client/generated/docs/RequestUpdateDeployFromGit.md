# RequestUpdateDeployFromGit

## Properties

| Name                      | Type                                               | Description | Notes                  |
| ------------------------- | -------------------------------------------------- | ----------- | ---------------------- |
| **dockerfilePath**        | **string**                                         |             | [default to undefined] |
| **environmentId**         | **number**                                         |             | [default to undefined] |
| **envs**                  | [**Array&lt;RequestEnvVar&gt;**](RequestEnvVar.md) |             | [default to undefined] |
| **port**                  | **number**                                         |             | [default to undefined] |
| **projectRepositoryPath** | **string**                                         |             | [default to undefined] |

## Example

```typescript
import { RequestUpdateDeployFromGit } from "./api";

const instance: RequestUpdateDeployFromGit = {
  dockerfilePath,
  environmentId,
  envs,
  port,
  projectRepositoryPath,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
