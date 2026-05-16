# ResponseGitDeployment

## Properties

| Name                      | Type                                                 | Description | Notes                  |
| ------------------------- | ---------------------------------------------------- | ----------- | ---------------------- |
| **args**                  | [**Array&lt;ResponseArg&gt;**](ResponseArg.md)       |             | [default to undefined] |
| **dockerfilePath**        | **string**                                           |             | [default to undefined] |
| **envVars**               | [**Array&lt;ResponseEnvVar&gt;**](ResponseEnvVar.md) |             | [default to undefined] |
| **gitUrl**                | **string**                                           |             | [default to undefined] |
| **id**                    | **number**                                           |             | [default to undefined] |
| **internalEndpoint**      | **string**                                           |             | [default to undefined] |
| **port**                  | **string**                                           |             | [default to undefined] |
| **projectRepositoryPath** | **string**                                           |             | [default to undefined] |
| **serviceName**           | **string**                                           |             | [default to undefined] |
| **status**                | **string**                                           |             | [default to undefined] |

## Example

```typescript
import { ResponseGitDeployment } from "./api";

const instance: ResponseGitDeployment = {
  args,
  dockerfilePath,
  envVars,
  gitUrl,
  id,
  internalEndpoint,
  port,
  projectRepositoryPath,
  serviceName,
  status,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
