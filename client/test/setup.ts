import { setupServer } from "msw/node";
import { afterAll, afterEach, beforeAll } from "vitest";

export const restHandlers = [];

export const server = setupServer(...restHandlers);

beforeAll(() => server.listen());
afterAll(() => server.close());
afterEach(() => server.resetHandlers());
