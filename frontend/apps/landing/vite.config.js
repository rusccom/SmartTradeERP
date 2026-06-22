import { fileURLToPath } from "url";
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

const pkg = (name) => fileURLToPath(new URL(`../../packages/${name}/src`, import.meta.url));

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@smarterp/ui": pkg("ui"),
      "@smarterp/api": pkg("api"),
      "@smarterp/auth": pkg("auth"),
      "@smarterp/i18n": pkg("i18n"),
      "@smarterp/model": pkg("model"),
    },
  },
});
