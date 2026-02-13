import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { cpSync, existsSync, mkdirSync, rmSync } from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const MONACO_SOURCE = path.resolve(__dirname, 'node_modules', 'monaco-editor', 'min', 'vs');
const MONACO_DESTINATION = path.resolve(__dirname, 'public', 'monaco', 'vs');

function copyMonacoAssets() {
  if (!existsSync(MONACO_SOURCE)) {
    console.warn('[vite] monaco-editor assets were not found, skipping copy step');
    return;
  }

  rmSync(MONACO_DESTINATION, { recursive: true, force: true });
  mkdirSync(MONACO_DESTINATION, { recursive: true });
  cpSync(MONACO_SOURCE, MONACO_DESTINATION, { recursive: true });
}

copyMonacoAssets();

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: 'dist',
    sourcemap: false,
  },
});
