import {
  type RouteConfig,
  route,
  prefix,
  layout,
  index,
} from "@react-router/dev/routes";

export default [
  layout("routes/dashboard/layout.tsx", [
    ...prefix(":slug?", [
      index("routes/dashboard/index.tsx"),

      ...prefix("projects", [
        index("routes/dashboard/projects/index.tsx"),
        route("all", "routes/dashboard/projects/all.tsx"),
        route("new", "routes/dashboard/projects/new.tsx"),
        ...prefix(":id", [
          index("routes/dashboard/projects/[id]/index.tsx"),
          layout("routes/dashboard/projects/[id]/layout.tsx", [
            route(
              "architecture",
              "routes/dashboard/projects/[id]/architecture.tsx",
            ),
            route(
              "observability",
              "routes/dashboard/projects/[id]/observability.tsx",
            ),
            route("logs", "routes/dashboard/projects/[id]/logs.tsx"),
            route("settings", "routes/dashboard/projects/[id]/settings.tsx"),
          ]),
        ]),
      ]),

      route("settings", "routes/dashboard/settings.tsx"),
    ]),
  ]),

  layout("routes/auth/layout.tsx", [
    route("signup", "routes/auth/signup.tsx"),
    route("login", "routes/auth/login.tsx"),
  ]),

  layout("routes/organizations/layout.tsx", [
    route("/organizations/new", "routes/organizations/new.tsx"),
  ]),

  ...prefix("api", [
    route("auth/*", "routes/api/auth.ts"),

    // tRPC routes
    route("trpc/*", "routes/api/trpc.ts"),
  ]),
] satisfies RouteConfig;
