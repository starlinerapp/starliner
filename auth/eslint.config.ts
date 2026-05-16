import js from "@eslint/js";
import tseslint from "typescript-eslint";
import { defineConfig } from "eslint/config";

export default defineConfig([
  {
    ignores: ["src/infrastructure/api/client/generated/**", "dist/**"],
  },
  {
    files: ["**/*.{js,mjs,cjs,ts,mts,cts,jsx,tsx}"],
    plugins: { js },
    extends: ["js/recommended"],
    settings: {
      react: {
        version: "detect",
      },
    },
  },
  tseslint.configs.recommended,
]);
