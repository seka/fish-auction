# syntax=docker/dockerfile:1.7

ARG NODE_IMAGE=node:25-alpine

# ============================================
# deps: 依存解決のキャッシュレイヤー
# ============================================
FROM ${NODE_IMAGE} AS deps
WORKDIR /app
COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile --ignore-scripts && yarn cache clean

# ============================================
# dev: ローカル開発用（next dev）
# ============================================
FROM ${NODE_IMAGE} AS dev
WORKDIR /app
COPY package.json yarn.lock panda.config.ts ./
RUN yarn install --frozen-lockfile --ignore-scripts
COPY . .
RUN yarn prepare
EXPOSE 3000
CMD ["yarn", "dev"]

# ============================================
# builder: 本番ビルド（next build → standalone 出力）
# ============================================
FROM ${NODE_IMAGE} AS builder
WORKDIR /app

# NEXT_PUBLIC_* はビルド時にバンドルへ埋め込まれるため ARG で受ける。
ARG NEXT_PUBLIC_API_URL=""
ARG NEXT_PUBLIC_VAPID_PUBLIC_KEY=""
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL} \
    NEXT_PUBLIC_VAPID_PUBLIC_KEY=${NEXT_PUBLIC_VAPID_PUBLIC_KEY} \
    NEXT_TELEMETRY_DISABLED=1

COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN yarn build

# ============================================
# prod: 本番ランタイム（standalone server.js を起動）
# ============================================
FROM ${NODE_IMAGE} AS prod
WORKDIR /app

ENV NODE_ENV=production \
    NEXT_TELEMETRY_DISABLED=1 \
    PORT=3000 \
    HOSTNAME=0.0.0.0

# 非 root 実行
RUN addgroup -g 10001 -S nodejs && \
    adduser -u 10001 -S nextjs -G nodejs

COPY --from=builder --chown=nextjs:nodejs /app/public        ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static  ./.next/static

USER nextjs:nodejs
EXPOSE 3000

CMD ["node", "server.js"]
