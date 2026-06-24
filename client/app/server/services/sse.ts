import type { Readable } from "node:stream";
import type { AxiosResponse } from "axios";

type StreamRequest = () => Promise<AxiosResponse<Readable>>;

interface StreamSseOptions {
  trim?: boolean;
}

export async function* streamSse(
  requestStream: StreamRequest,
  signal: AbortSignal | undefined,
  options: StreamSseOptions = {},
): AsyncGenerator<string> {
  let response: AxiosResponse<Readable> | undefined;

  try {
    response = await requestStream();

    signal?.addEventListener("abort", () => {
      response?.data?.destroy();
    });

    const decoder = new TextDecoder();
    let buffer = "";

    for await (const chunk of response.data) {
      if (signal?.aborted) {
        break;
      }

      buffer += decoder.decode(chunk, { stream: true });
      const lines = buffer.split("\n");
      buffer = lines.pop() ?? "";

      for (const line of lines) {
        if (line.startsWith("data: ")) {
          const data = line.slice(6);
          yield options.trim ? data.trim() : data;
        }
      }
    }
  } catch (err) {
    if (signal?.aborted) {
      return;
    }
    throw err;
  } finally {
    response?.data?.destroy();
  }
}


export async function* streamSseJson<T = unknown>(
  requestStream: StreamRequest,
  signal: AbortSignal | undefined,
): AsyncGenerator<T | string> {
  for await (const data of streamSse(requestStream, signal, { trim: true })) {
    if (!data) {
      continue;
    }

    try {
      yield JSON.parse(data) as T;
    } catch {
      yield data;
    }
  }
}
