# docker buildx bake 用の build orchestration 定義。
#
# 使い方:
#   cd frontend
#   docker buildx bake                  # prod ターゲットをビルド（VERSION=0.0.0 で local 用）
#   VERSION=1.2.3 REGISTRY=ghcr.io/seka/fish-auction \
#       docker buildx bake --push       # リリース時（VERSION がそのまま image tag になる）
#
# Dockerfile 直叩き (`docker build ./frontend`) も従来通り動作する
# （最終ステージが prod）。bake は本番イメージのビルド／push 用。

variable "VERSION" {
  default = "0.0.0"
}

variable "REGISTRY" {
  default = "fish-auction"
}

# next.config.ts の rewrites destination・クライアントバンドル埋め込み値は
# next build 時点で確定するため、bake 起動時に env で渡す。デフォルトは空文字。
variable "API_BASE_URL" {
  default = ""
}

variable "NEXT_PUBLIC_API_URL" {
  default = ""
}

variable "NEXT_PUBLIC_VAPID_PUBLIC_KEY" {
  default = ""
}

group "default" {
  targets = ["prod"]
}

target "_common" {
  context    = "."
  dockerfile = "Dockerfile"
  args = {
    API_BASE_URL                 = API_BASE_URL
    NEXT_PUBLIC_API_URL          = NEXT_PUBLIC_API_URL
    NEXT_PUBLIC_VAPID_PUBLIC_KEY = NEXT_PUBLIC_VAPID_PUBLIC_KEY
  }
}

target "prod" {
  inherits = ["_common"]
  target   = "prod"
  tags     = ["${REGISTRY}/frontend:${VERSION}"]
}
