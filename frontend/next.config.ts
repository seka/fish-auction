import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  rewrites: async () => {
    // Default to localhost for local development (yarn dev)
    const apiBaseUrl = process.env.API_BASE_URL || 'http://localhost:8080';
    return [
      {
        source: '/api/:path*',
        destination: `${apiBaseUrl}/api/:path*`,
      },
    ];
  },
};

import createNextIntlPlugin from 'next-intl/plugin';

const withNextIntl = createNextIntlPlugin('./src/core/i18n/request.ts');

export default withNextIntl(nextConfig);
