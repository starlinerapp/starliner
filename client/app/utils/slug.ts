export function toSlug(input: string): string {
  return input
    .toLowerCase()
    .trim()
    .normalize("NFD") // decompose accented chars
    .replace(/[\u0300-\u036f]/g, "") // remove diacritic marks
    .replace(/[^a-z0-9\s-]/g, "") // remove non-alphanumeric chars
    .replace(/[\s_]+/g, "-") // spaces/underscores → hyphens
    .replace(/-+/g, "-") // collapse multiple hyphens
    .replace(/^-|-$/g, ""); // trim leading/trailing hyphens
}
