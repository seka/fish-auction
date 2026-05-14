# docker buildx bake 用の build orchestration 定義。
#
# 使い方:
#   cd backend
#   docker buildx bake                       # 全 prod-* ターゲットを並列ビルド（VERSION=0.0.0 で local 用）
#   docker buildx bake prod-server           # 個別ターゲット
#   VERSION=1.2.3 COMMIT=abc1234 REGISTRY=ghcr.io/seka/fish-auction \
#       docker buildx bake --push            # リリース時（VERSION がそのまま image tag になる）
#
# Dockerfile 直叩き (`docker build ./backend`) も従来通り動作する
# （最終ステージが prod-server）。bake は複数イメージの一括ビルド／push 用。

# バイナリの ldflags 注入値かつ image tag。同一値を使い回し、不整合を避ける。
variable "VERSION" {
  default = "0.0.0"
}

variable "COMMIT" {
  default = "unknown"
}

variable "REGISTRY" {
  default = "fish-auction"
}

group "default" {
  targets = [
    "prod-server",
    "prod-worker",
    "prod-relay",
    "prod-migration",
    "prod-ops",
  ]
}

target "_common" {
  context    = "."
  dockerfile = "Dockerfile"
  args = {
    VERSION = VERSION
    COMMIT  = COMMIT
  }
}

target "prod-server" {
  inherits = ["_common"]
  target   = "prod-server"
  tags     = ["${REGISTRY}/server:${VERSION}"]
}

target "prod-worker" {
  inherits = ["_common"]
  target   = "prod-worker"
  tags     = ["${REGISTRY}/worker:${VERSION}"]
}

target "prod-relay" {
  inherits = ["_common"]
  target   = "prod-relay"
  tags     = ["${REGISTRY}/relay:${VERSION}"]
}

target "prod-migration" {
  inherits = ["_common"]
  target   = "prod-migration"
  tags     = ["${REGISTRY}/migration:${VERSION}"]
}

target "prod-ops" {
  inherits = ["_common"]
  target   = "prod-ops"
  tags     = ["${REGISTRY}/ops:${VERSION}"]
}
