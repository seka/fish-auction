import { defineConfig } from '@pandacss/dev'

export default defineConfig({
  // Enable CSS reset
  preflight: true,

  // Include all TypeScript/TSX files
  include: ['./app/**/*.{ts,tsx}', './src/**/*.{ts,tsx}'],

  // Exclude node_modules and .next
  exclude: [],

  // Theme configuration
  theme: {
    extend: {
      tokens: {
        colors: {
          // Primary colors
          primary: {
            50: { value: '#eef2ff' },   // indigo-50
            100: { value: '#e0e7ff' },  // indigo-100
            600: { value: '#4f46e5' },  // indigo-600
            900: { value: '#312e81' },  // indigo-900
          },
          // Secondary colors
          secondary: {
            50: { value: '#fff7ed' },   // orange-50
            100: { value: '#ffedd5' },  // orange-100
            600: { value: '#ea580c' },  // orange-600
          },
          // UI colors
          blue: {
            50: { value: '#eff6ff' },
            100: { value: '#dbeafe' },
            200: { value: '#bfdbfe' },
            600: { value: '#2563eb' },
            800: { value: '#1e40af' },
            900: { value: '#1e3a8a' },
          },
          green: {
            50: { value: '#f0fdf4' },
            100: { value: '#dcfce7' },
            600: { value: '#16a34a' },
          },
          purple: {
            50: { value: '#faf5ff' },
            100: { value: '#f3e8ff' },
            600: { value: '#9333ea' },
          },
          yellow: {
            50: { value: '#fefce8' },
            100: { value: '#fef9c3' },
            600: { value: '#ca8a04' },
          },
          // Gray scale
          gray: {
            50: { value: '#f9fafb' },
            100: { value: '#f3f4f6' },
            500: { value: '#6b7280' },
            600: { value: '#4b5563' },
            800: { value: '#1f2937' },
            900: { value: '#111827' },
          },
          // Background and foreground
          background: { value: '#ffffff' },
          foreground: { value: '#171717' },
        },
        spacing: {
          // Extend default spacing
          '18': { value: '4.5rem' },
        },
        radii: {
          '2xl': { value: '1rem' },
          '3xl': { value: '1.5rem' },
        },
        shadows: {
          'xl': { value: '0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1)' },
          '2xl': { value: '0 25px 50px -12px rgb(0 0 0 / 0.25)' },
        },
      },
      keyframes: {
        slideIn: {
          '0%': { transform: 'translateX(100%)', opacity: '0' },
          '100%': { transform: 'translateX(0)', opacity: '1' },
        },
      },
    },
  },

  // Output directory for generated files
  outdir: 'styled-system',

  // Configure JSX framework
  jsxFramework: 'react',
})
