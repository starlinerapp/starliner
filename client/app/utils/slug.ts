export function formatSlugInput(input: string): string {
  return input
    .toLowerCase()
    .replace(/[^a-z0-9-]+/g, "-")
    .replace(/-+/g, "-")
    .slice(0, 50);
}

export function sanitizeSlug(input: string): string {
  return formatSlugInput(input).replace(/^-+|-+$/g, "");
}
