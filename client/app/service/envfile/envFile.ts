export function isEnvFile(text: string): boolean {
  const lines = text
    .split("\n")
    .map((l) => l.trim())
    .filter((l) => l && !l.startsWith("#"));
  // All non-empty, non-comment lines must look like KEY=VALUE
  const kvLines = lines.filter((l) => /^[A-Z_][A-Z0-9_]*\s*=.*$/.test(l));
  return kvLines.length === lines.length;
}

export function parseEnvFile(
  content: string,
): { name: string; value: string }[] {
  return content
    .split("\n")
    .map((line) => line.trim())
    .filter((line) => line && !line.startsWith("#"))
    .flatMap((line) => {
      const eqIndex = line.indexOf("=");
      if (eqIndex === -1) return [];
      const name = line.slice(0, eqIndex).trim();
      let value = line.slice(eqIndex + 1).trim();
      if (
        (value.startsWith('"') && value.endsWith('"')) ||
        (value.startsWith("'") && value.endsWith("'"))
      ) {
        value = value.slice(1, -1);
      }
      if (!name) return [];
      return [{ name, value }];
    });
}
