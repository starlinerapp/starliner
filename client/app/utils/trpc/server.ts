// Defines bridge between tRPC server and client.

import { createTRPCOptionsProxy } from "@trpc/tanstack-react-query";
import type { LoaderFunctionArgs } from "react-router";
import { appRouter } from "~/server/main";
import { createCallerFactory, createTRPCContext } from "~/server/trpc";
import { getQueryClient } from "~/utils/trpc/react";

const createContext = (opts: { headers: Headers }) => {
  return createTRPCContext({
    headers: opts.headers,
  });
};

const createCaller = createCallerFactory(appRouter);

export const caller = async (loaderArgs: LoaderFunctionArgs) =>
  createCaller(await createContext({ headers: loaderArgs.request.headers }));

export async function createTRPC(loaderArgs: LoaderFunctionArgs) {
  return createTRPCOptionsProxy({
    ctx: () => createContext({ headers: loaderArgs.request.headers }),
    queryClient: getQueryClient,
    router: appRouter,
  });
}
