import type { Context } from "hono";
import type { AuthService } from "~/domain/port/auth";

export class AuthHandler {
  constructor(private readonly auth: AuthService) {}

  handle = (c: Context) => {
    return this.auth.handler(c.req.raw);
  };
}
