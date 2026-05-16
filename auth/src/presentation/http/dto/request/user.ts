import { z } from "@hono/zod-openapi";

export const BulkUserLookupRequestSchema = z
  .object({
    ids: z
      .array(z.string())
      .openapi({ example: ["user_abc123", "user_def456"] }),
  })
  .openapi("BulkUserLookupRequest");

export type BulkUserLookupRequest = z.infer<typeof BulkUserLookupRequestSchema>;
