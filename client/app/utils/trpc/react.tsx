import React from "react";
import superjson from "superjson";
import {
  defaultShouldDehydrateQuery,
  QueryClient,
  QueryClientProvider,
} from "@tanstack/react-query";
import {
  createTRPCClient,
  httpBatchStreamLink,
  httpSubscriptionLink,
  loggerLink,
  splitLink,
} from "@trpc/client";
import { type PropsWithChildren, useState } from "react";
import type { inferRouterInputs } from "@trpc/server";
import { createTRPCContext } from "@trpc/tanstack-react-query";
import type { AppRouter } from "~/server/main";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

function createQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 60 * 1000, // 1 minute
      },
      dehydrate: {
        serializeData: superjson.serialize,
        shouldDehydrateQuery: (query) =>
          defaultShouldDehydrateQuery(query) ||
          query.state.status === "pending",
      },
      hydrate: {
        deserializeData: superjson.deserialize,
      },
    },
  });
}

let browserQueryClient: QueryClient | undefined = undefined;

export const getQueryClient = () => {
  if (typeof window === "undefined") {
    return createQueryClient();
  }
  browserQueryClient ??= createQueryClient();
  return browserQueryClient;
};

const links = [
  loggerLink({
    enabled: (op) =>
      process.env.NODE_ENV === "development" ||
      (op.direction === "down" && op.result instanceof Error),
  }),
  splitLink({
    condition: (op) => op.type === "subscription",
    true: httpSubscriptionLink({
      transformer: superjson,
      url: "/api/trpc",
    }),
    false: httpBatchStreamLink({
      transformer: superjson,
      url: "/api/trpc",
      maxURLLength: 2083,
    }),
  }),
];

export const { TRPCProvider, useTRPC } = createTRPCContext<AppRouter>();

interface TRPCReactProviderProps {
  queryClient?: QueryClient;
}

export function TRPCReactProvider({
  children,
}: TRPCReactProviderProps & PropsWithChildren) {
  const queryClient = getQueryClient();
  const [trpcClient] = useState(() =>
    createTRPCClient<AppRouter>({
      links: links,
    }),
  );

  return (
    <QueryClientProvider client={queryClient}>
      <TRPCProvider trpcClient={trpcClient} queryClient={queryClient}>
        {children}
      </TRPCProvider>
      <ReactQueryDevtools />
    </QueryClientProvider>
  );
}

export type RouterInputs = inferRouterInputs<AppRouter>;
export type RouterOutputs = inferRouterInputs<AppRouter>;
