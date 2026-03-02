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
          ...prefix(":environment", [
            layout("routes/dashboard/projects/[id]/[environment]/layout.tsx", [
              layout(
                "routes/dashboard/projects/[id]/[environment]/architecture/layout.tsx",
                [
                  ...prefix("architecture", [
                    index(
                      "routes/dashboard/projects/[id]/[environment]/architecture/index.tsx",
                    ),
                    route(
                      "image",
                      "routes/dashboard/projects/[id]/[environment]/architecture/image.tsx",
                    ),
                    route(
                      "ingress",
                      "routes/dashboard/projects/[id]/[environment]/architecture/ingress.tsx",
                    ),
                    route(
                      "database",
                      "routes/dashboard/projects/[id]/[environment]/architecture/database.tsx",
                    ),
                  ]),
                ],
              ),
              route(
                "observability",
                "routes/dashboard/projects/[id]/[environment]/observability.tsx",
              ),
              route(
                "logs",
                "routes/dashboard/projects/[id]/[environment]/logs.tsx",
              ),
              route(
                "settings",
                "routes/dashboard/projects/[id]/[environment]/settings.tsx",
              ),
            ]),
          ]),
        ]),
      ]),

      ...prefix("clusters", [
        index("routes/dashboard/clusters/index.tsx"),
        route("all", "routes/dashboard/clusters/all.tsx"),
        route("new", "routes/dashboard/clusters/new.tsx"),
        ...prefix(":id", [
          index("routes/dashboard/clusters/[id]/index.tsx"),
          layout("routes/dashboard/clusters/[id]/layout.tsx", [
            route("general", "routes/dashboard/clusters/[id]/general.tsx"),
            route("settings", "routes/dashboard/clusters/[id]/settings.tsx"),
          ]),
          ...prefix("resources", [
            route(
              "private-key",
              "routes/dashboard/clusters/[id]/resources/private-key.tsx",
            ),
          ]),
        ]),
      ]),

      ...prefix("settings", [
        index("routes/dashboard/settings/index.tsx"),
        route("organization", "routes/dashboard/settings/organization.tsx"),
      ]),
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
