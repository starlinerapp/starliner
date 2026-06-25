export function appendSseDataLines(
  buffer: string,
  onLine: (payload: string) => void,
): string {
  const lines = buffer.split("\n");
  const remaining = lines.pop() ?? "";

  for (const line of lines) {
    if (!line.startsWith("data: ")) {
      continue;
    }

    const payload = line.slice(6).replace(/\r$/, "");
    if (!payload) {
      continue;
    }

    onLine(payload);
  }

  return remaining;
}
