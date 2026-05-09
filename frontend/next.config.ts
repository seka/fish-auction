import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  // ECS Fargate などコンテナ運用時に最小ランタイム成果物を生成する。
  output: 'standalone',
  rewrites: async () => {
    const backendUrl = process.env.API_BASE_URL || 'http://127.0.0.1:8080';
    return [
      {
        source: '/api/:path*',
        destination: `${backendUrl}/api/:path*`,
      },
    ];
  },
};

import createNextIntlPlugin from 'next-intl/plugin';

const withNextIntl = createNextIntlPlugin('./src/core/i18n/request.ts');

export default withNextIntl(nextConfig);
