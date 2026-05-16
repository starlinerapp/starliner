import { createRoute } from "@hono/zod-openapi";
import { BulkUserLookupRequestSchema } from "~/presentation/http/dto/request/user";
import { BulkUserLookupResponseSchema } from "~/presentation/http/dto/response/user";
import { ErrorResponseSchema } from "~/presentation/http/dto/response/common";

export const bulkUserLookupRoute = createRoute({
  method: "post",
  path: "/users",
  tags: ["users"],
  operationId: "bulkUserLookup",
  summary: "Look up users by ID",
  description:
    "Returns user profiles for the given IDs. Unknown IDs are omitted. At most 200 unique IDs per request.",
  request: {
    body: {
      content: {
        "application/json": {
          schema: BulkUserLookupRequestSchema,
        },
      },
      required: true,
    },
  },
  responses: {
    200: {
      content: {
        "application/json": {
          schema: BulkUserLookupResponseSchema,
        },
      },
      description: "Matching user profiles",
    },
    400: {
      content: {
        "application/json": {
          schema: ErrorResponseSchema,
        },
      },
      description: "Invalid JSON body or too many IDs",
    },
  },
});
