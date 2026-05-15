import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import path from "path";
import { defineConfig, loadEnv } from "vite";

const manualChunks = {
  vendor: ["react", "react-dom"],
  router: ["@tanstack/react-router"],
  query: ["@tanstack/react-query"],
  ui: ["radix-ui"],
  forms: ["react-hook-form", "@hookform/resolvers", "zod"],
  icons: ["lucide-react"],
};

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, "../");

  return {
    plugins: [react(), tailwindcss()],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },
    server: {
      proxy: {
        "/api": {
          target: env.VITE_BASE_URL || "http://localhost:3001",
          changeOrigin: true,
        },
      },
    },
    build: {
      rollupOptions: {
        output: {
          manualChunks(id) {
            const normalizedId = id.replaceAll("\\", "/");

            for (const [chunkName, packages] of Object.entries(manualChunks)) {
              if (
                packages.some((packageName) =>
                  normalizedId.includes(`/node_modules/${packageName}/`),
                )
              ) {
                return chunkName;
              }
            }
          },
        },
      },
    },
  };
});
