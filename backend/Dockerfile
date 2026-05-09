# syntax=docker/dockerfile:1.7

# golang/alpine のメジャーマイナーを固定。本番運用時は CI で digest pin することを推奨。
ARG GO_IMAGE=golang:1.26-alpine3.21
ARG RUNTIME_IMAGE=alpine:3.21

# ============================================
# base: 依存解決のキャッシュレイヤー
# ============================================
FROM ${GO_IMAGE} AS base
RUN apk add --no-cache make
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# ============================================
# dev: ローカル開発用（air によるホットリロード）
# ============================================
FROM base AS dev
RUN apk add --no-cache postgresql-client
RUN go install github.com/cosmtrek/air@v1.49.0
COPY . .
CMD ["air", "-c", ".air.toml"]

# ============================================
# builder: 本番バイナリを全 cmd 分ビルド
# ============================================
FROM base AS builder
COPY . .
ARG VERSION=dev
ARG COMMIT=unknown
ENV CGO_ENABLED=0 GOOS=linux
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" \
        -o /out/server    ./cmd/server && \
    go build -trimpath -ldflags "-s -w" -o /out/worker    ./cmd/worker && \
    go build -trimpath -ldflags "-s -w" -o /out/relay     ./cmd/relay && \
    go build -trimpath -ldflags "-s -w" -o /out/migration ./cmd/migration

# ============================================
# prod-base: 本番ランタイム共通（非 root、CA 証明書、tzdata）
# ============================================
FROM ${RUNTIME_IMAGE} AS prod-base
RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -g 10001 -S app && \
    adduser  -u 10001 -S app -G app
USER app:app
WORKDIR /app

# ============================================
# prod-server: HTTP サーバー
# ============================================
FROM prod-base AS prod-server
COPY --from=builder /out/server /app/server
EXPOSE 8080
ENTRYPOINT ["/app/server"]

# ============================================
# prod-worker: SQS コンシューマー
# ============================================
FROM prod-base AS prod-worker
COPY --from=builder /out/worker /app/worker
ENTRYPOINT ["/app/worker"]

# ============================================
# prod-relay: outbox → SQS リレー
# ============================================
FROM prod-base AS prod-relay
COPY --from=builder /out/relay /app/relay
ENTRYPOINT ["/app/relay"]

# ============================================
# prod-migration: ECS RunTask 用ワンショット
# ============================================
FROM prod-base AS prod-migration
COPY --from=builder /out/migration /app/migration
ENTRYPOINT ["/app/migration"]
CMD ["up"]
