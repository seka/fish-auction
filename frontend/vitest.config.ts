import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';
import path from 'path';

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './'),
      '@libs': path.resolve(__dirname, './src/libs'),
      '@atoms': path.resolve(__dirname, './src/components/atoms'),
      '@molecules': path.resolve(__dirname, './src/components/molecules'),
      '@functionals': path.resolve(__dirname, './src/components/functionals'),
      '@organisms': path.resolve(__dirname, './src/components/organisms'),
      '@templates': path.resolve(__dirname, './src/components/templates'),
      'styled-system': path.resolve(__dirname, './src/libs/styled-system'),
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
        'vitest.config.ts',
      ],
    },
  },
});
