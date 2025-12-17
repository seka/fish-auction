import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";
import path from "path";

export default defineConfig({
    plugins: [react()],
    test: {
        environment: "jsdom",
        globals: true,
        setupFiles: "./vitest.setup.ts",
        alias: {
            "@": path.resolve(__dirname, "./"),
            "styled-system": path.resolve(__dirname, "./styled-system"),
        },
        coverage: {
            provider: "v8",
            exclude: [
                "styled-system/**",
                "panda.config.ts",
                "postcss.config.mjs",
                "eslint.config.mjs",
                "next.config.ts", // Next.js config
                "**/*.d.ts",      // Types
                "**/*.test.tsx",  // Tests
                "**/*.setup.ts",  // Setup files
                "vitest.config.ts",
            ],
        },
    },
});
