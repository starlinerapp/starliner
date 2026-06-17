import Redis from "ioredis";
import { serverEnv } from "~/env.server";

const client = new Redis({
  host: serverEnv.REDIS_ADDR.split(":")[0],
  port: parseInt(serverEnv.REDIS_ADDR.split(":")[1] ?? "6379", 10),
  password: serverEnv.REDIS_PASSWORD,
  db: 0,
});

client.on("error", (err) => {
  console.error("[cache] Redis connection error:", err);
});

export const cache = {
  async set(key: string, value: string, ttlSeconds: number): Promise<void> {
    await client.set(key, value, "EX", ttlSeconds);
    return undefined;
  },

  get(key: string): Promise<string | null> {
    return client.get(key);
  },

  async delete(key: string): Promise<void> {
    await client.del(key);
    return undefined;
  },
};
