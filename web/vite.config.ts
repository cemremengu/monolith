import path from "path";

import tailwindcss from "@tailwindcss/vite";
import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react";

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
          manualChunks: {
            // Vendor chunks
            vendor: ["react", "react-dom"],
            router: ["@tanstack/react-router"],
            query: ["@tanstack/react-query"],
            ui: ["@base-ui/react"],
            forms: ["react-hook-form", "@hookform/resolvers", "zod"],
            icons: ["lucide-react"],
          },
        },
      },
    },
  };
});
