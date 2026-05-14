# syntax=docker/dockerfile:1.7

ARG NODE_IMAGE=node:25-alpine

# ============================================
# base: 依存解決のキャッシュレイヤー
# ============================================
FROM ${NODE_IMAGE} AS base
WORKDIR /app
COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile --ignore-scripts && yarn cache clean

# ============================================
# dev: ローカル開発用（next dev）
# ============================================
FROM base AS dev
COPY . .
RUN yarn prepare
EXPOSE 3000
CMD ["yarn", "dev"]

# ============================================
# builder: 本番ビルド（next build → standalone 出力）
# ============================================
FROM base AS builder

# NEXT_PUBLIC_* はビルド時にバンドルへ埋め込まれるため ARG で受ける。
ARG NEXT_PUBLIC_API_URL=""
ARG NEXT_PUBLIC_VAPID_PUBLIC_KEY=""
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL} \
    NEXT_PUBLIC_VAPID_PUBLIC_KEY=${NEXT_PUBLIC_VAPID_PUBLIC_KEY} \
    NEXT_TELEMETRY_DISABLED=1

COPY . .
RUN yarn build

# ============================================
# prod-base: 本番ランタイム共通（非 root、共通 ENV）
# ============================================
FROM ${NODE_IMAGE} AS prod-base
WORKDIR /app
ENV NODE_ENV=production \
    NEXT_TELEMETRY_DISABLED=1 \
    PORT=3000 \
    HOSTNAME=0.0.0.0
RUN addgroup -g 10001 -S nodejs && \
    adduser  -u 10001 -S nextjs -G nodejs
USER nextjs:nodejs

# ============================================
# prod: 本番ランタイム（standalone server.js を起動）
# `--target` 未指定の docker build がデフォルトでこのステージを生成するよう、
# 最終ステージに配置する。
# ============================================
FROM prod-base AS prod
COPY --chown=nextjs:nodejs --from=builder /app/public           ./public
COPY --chown=nextjs:nodejs --from=builder /app/.next/standalone ./
COPY --chown=nextjs:nodejs --from=builder /app/.next/static     ./.next/static
EXPOSE 3000
CMD ["node", "server.js"]
