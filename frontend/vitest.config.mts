import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';
import path from 'path';

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(import.meta.dirname, './'),
      '@atoms': path.resolve(import.meta.dirname, './src/core/components/atoms'),
      '@molecules': path.resolve(import.meta.dirname, './src/core/components/molecules'),
      '@organisms': path.resolve(import.meta.dirname, './src/core/components/organisms'),
      '@templates': path.resolve(import.meta.dirname, './src/core/components/templates'),
      '@bootstraps': path.resolve(import.meta.dirname, './src/bootstraps'),
      '@libs': path.resolve(import.meta.dirname, './src/libs'),
      '@entities': path.resolve(import.meta.dirname, './src/data/entities'),
      '@entities/*': path.resolve(import.meta.dirname, './src/data/entities/*'),
      '@schemas': path.resolve(import.meta.dirname, './src/schemas'),
      '@schemas/*': path.resolve(import.meta.dirname, './src/schemas/*'),
      '@testings': path.resolve(import.meta.dirname, './src/testings'),
      'styled-system': path.resolve(import.meta.dirname, './src/libs/styled-system'),
    },
  },
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: './vitest.setup.ts',
    coverage: {
      provider: 'v8',
      exclude: [
        'src/libs/styled-system/**',
        'panda.config.ts',
        'postcss.config.mjs',
        'eslint.config.mjs',
        'next.config.ts', // Next.js config
        '**/*.d.ts', // Types
        '**/*.test.tsx', // Tests
        '**/*.setup.ts', // Setup files
        'vitest.config.mts',
      ],
    },
  },
});
