import type { Context } from "hono";
import {
  UserApplication,
  UserLookupError,
} from "~/application/user";
import { parseBulkUserLookupRequest } from "~/presentation/http/dto/request/user";
import { newBulkUserLookupResponse } from "~/presentation/http/dto/response/user";

export class UserHandler {
  constructor(private readonly userApplication: UserApplication) {}

  bulkLookup = async (c: Context) => {
    let body: unknown;
    try {
      body = await c.req.json();
    } catch {
      return c.json({ error: "invalid json" }, 400);
    }

    const request = parseBulkUserLookupRequest(body);
    if (!request) {
      return c.json({ error: "invalid json" }, 400);
    }

    try {
      const users = await this.userApplication.getUsersByIds(request.ids);
      return c.json(newBulkUserLookupResponse(users));
    } catch (error) {
      if (error instanceof UserLookupError && error.code === "too_many_ids") {
        return c.json({ error: "too many ids" }, 400);
      }
      throw error;
    }
  };
}
