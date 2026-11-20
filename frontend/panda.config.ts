import { defineConfig } from '@pandacss/dev';
import { colors, spacing, typography, radii, shadows, animations } from './src/core/styles/tokens';

export default defineConfig({
  // Panda v2 requires explicit presets for utilities, conditions, and default tokens.
  presets: ['@pandacss/preset-base', '@pandacss/preset-panda'],

  // Enable CSS reset
  preflight: true,

  // Include all TypeScript/TSX files
  include: ['./app/**/*.{ts,tsx}', './src/**/*.{ts,tsx}'],

  // Exclude node_modules and .next
  exclude: ['src/libs/styled-system'],

  // Theme configuration
  theme: {
    extend: {
      tokens: {
        colors,
        spacing,
        radii,
        shadows,
        fonts: typography.fonts,
        fontWeights: typography.fontWeights,
        fontSizes: typography.fontSizes,
        lineHeights: typography.lineHeights,
      },
      keyframes: animations.keyframes,
    },
  },

  // Output directory for generated files
  outdir: 'src/libs/styled-system',

  // Configure JSX framework
  jsxFramework: 'react',
});
