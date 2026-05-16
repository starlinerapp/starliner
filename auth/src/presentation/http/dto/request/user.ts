export type BulkUserLookupRequest = {
  ids: string[];
};

export function parseBulkUserLookupRequest(
  body: unknown,
): BulkUserLookupRequest | null {
  if (typeof body !== "object" || body === null) {
    return null;
  }

  const rawIds = (body as { ids?: unknown }).ids;
  const ids = Array.isArray(rawIds)
    ? rawIds.filter((id): id is string => typeof id === "string")
    : [];

  return { ids };
}
