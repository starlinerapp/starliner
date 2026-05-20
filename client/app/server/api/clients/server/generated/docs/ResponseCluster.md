# ResponseCluster

## Properties

| Name               | Type                                                  | Description | Notes                  |
| ------------------ | ----------------------------------------------------- | ----------- | ---------------------- |
| **createdAt**      | **string**                                            |             | [default to undefined] |
| **id**             | **number**                                            |             | [default to undefined] |
| **ipv4Address**    | **string**                                            |             | [default to undefined] |
| **name**           | **string**                                            |             | [default to undefined] |
| **organizationId** | **number**                                            |             | [default to undefined] |
| **serverType**     | **string**                                            |             | [default to undefined] |
| **status**         | [**ResponseClusterStatus**](ResponseClusterStatus.md) |             | [default to undefined] |
| **teamSlugs**      | **Array&lt;string&gt;**                               |             | [default to undefined] |
| **user**           | **string**                                            |             | [default to undefined] |

## Example

```typescript
import { ResponseCluster } from "./api";

const instance: ResponseCluster = {
  createdAt,
  id,
  ipv4Address,
  name,
  organizationId,
  serverType,
  status,
  teamSlugs,
  user,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
