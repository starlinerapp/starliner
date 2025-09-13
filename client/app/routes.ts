import {
  type RouteConfig,
  index,
  route,
  prefix,
  layout,
} from "@react-router/dev/routes";

export default [
  index("routes/home.tsx"),
  route("settings", "routes/settings.tsx"),

  layout("routes/auth/layout.tsx", [
    route("signup", "routes/auth/signup.tsx"),
    route("login", "routes/auth/login.tsx"),
  ]),

  ...prefix("api", [route("auth/*", "routes/api/auth.ts")]),
] satisfies RouteConfig;
