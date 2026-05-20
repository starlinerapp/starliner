import type { OpenAPIHono, RouteHandler } from "@hono/zod-openapi";
import { bulkUserLookupRoute } from "../routes/user";
import {
  type UserApplication,
  UserLookupError,
} from "../../../application/user";
import { toBulkUserLookupResponse } from "../mapper/user";

export class UserHandler {
  constructor(private readonly userApplication: UserApplication) {}

  static register(app: OpenAPIHono, userApplication: UserApplication) {
    const handler = new UserHandler(userApplication);
    app.openapi(bulkUserLookupRoute, handler.bulkLookup);
  }

  bulkLookup: RouteHandler<typeof bulkUserLookupRoute> = async (c) => {
    const { ids } = c.req.valid("json");

    try {
      const users = await this.userApplication.getUsersByIds(ids);
      return c.json(toBulkUserLookupResponse(users), 200);
    } catch (error) {
      if (error instanceof UserLookupError && error.code === "too_many_ids") {
        return c.json({ error: "too many ids" }, 400);
      }

      throw error;
    }
  };
}
