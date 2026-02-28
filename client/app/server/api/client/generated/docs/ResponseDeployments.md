# ResponseDeployments

## Properties

| Name          | Type                                                                         | Description | Notes                  |
| ------------- | ---------------------------------------------------------------------------- | ----------- | ---------------------- |
| **databases** | [**Array&lt;ResponseDatabaseDeployment&gt;**](ResponseDatabaseDeployment.md) |             | [default to undefined] |
| **images**    | [**Array&lt;ResponseImageDeployment&gt;**](ResponseImageDeployment.md)       |             | [default to undefined] |
| **ingresses** | [**Array&lt;ResponseIngressDeployment&gt;**](ResponseIngressDeployment.md)   |             | [default to undefined] |

## Example

```typescript
import { ResponseDeployments } from "./api";

const instance: ResponseDeployments = {
  databases,
  images,
  ingresses,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
