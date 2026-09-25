import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";
import path from "path";

export default defineConfig({
  plugins: [react()],
  test: {
    environment: "jsdom",
    globals: false,
    setupFiles: ["./vitest.setup.ts"],
    // Exclude Playwright e2e tests from the vitest run
    exclude: ["tests/e2e/**", "**/node_modules/**"],
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "."),
      "@sorolens/ui": path.resolve(__dirname, "../../packages/ui/src/index.ts"),
      "@sorolens/xdr": path.resolve(__dirname, "../../packages/xdr/src/index.ts"),
    },
  },
});

