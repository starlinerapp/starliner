export function isDockerfile(text: string): boolean {
  const lines = text
    .split("\n")
    .map((l) => l.trim())
    .filter((l) => l && !l.startsWith("#"));
  return lines.some((l) => /^FROM\s+\S+/i.test(l));
}

export function parseDockerfile(
  content: string,
): { name: string; value: string }[] {
  const lines = content.split("\n");
  const result: { name: string; value: string }[] = [];

  for (let i = 0; i < lines.length; i++) {
    let line = lines[i].trim();
    if (!line || line.startsWith("#")) continue;

    while (line.endsWith("\\") && i + 1 < lines.length) {
      line = `${line.slice(0, -1).trim()} ${lines[++i].trim()}`;
    }

    if (!/^ARG\s+/i.test(line)) continue;

    const argBody = line.replace(/^ARG\s+/i, "").trim();
    const eqIndex = argBody.indexOf("=");

    if (eqIndex === -1) {
      if (argBody) result.push({ name: argBody, value: "" });
    } else {
      const name = argBody.slice(0, eqIndex).trim();
      let value = argBody.slice(eqIndex + 1).trim();
      if (
        (value.startsWith('"') && value.endsWith('"')) ||
        (value.startsWith("'") && value.endsWith("'"))
      ) {
        value = value.slice(1, -1);
      }
      if (name) result.push({ name, value });
    }
  }

  return result;
}
