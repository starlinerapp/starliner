import { z } from "@hono/zod-openapi";

export const UserProfileSchema = z
  .object({
    id: z.string().openapi({ example: "user_abc123" }),
    name: z.string().openapi({ example: "Jane Doe" }),
    email: z.string().openapi({ example: "jane@example.com" }),
  })
  .openapi("UserProfile");

export const BulkUserLookupResponseSchema = z
  .object({
    users: z.array(UserProfileSchema),
  })
  .openapi("BulkUserLookupResponse");

export type BulkUserLookupResponse = z.infer<typeof BulkUserLookupResponseSchema>;
