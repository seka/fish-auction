import { defineConfig } from '@pandacss/dev'

export default defineConfig({
  // Enable CSS reset
  preflight: true,

  // Include all TypeScript/TSX files
  include: ['./app/**/*.{ts,tsx}', './src/**/*.{ts,tsx}'],
  
  // Exclude node_modules and .next
  exclude: [],

  // Output directory for generated files
  outdir: 'styled-system',

  // Configure JSX framework
  jsxFramework: 'react',
})
