# 漁港のせり (Fish Auction)

漁港における「せり（オークション）」をデジタル化し、効率的な取引を支援するシステムです。
荷主（漁協・漁師）、仲買人（バイヤー）、管理者それぞれの業務を統合的に管理します。

## 主要機能 (Features)

- **オークション管理**: せり情報の登録、リアルタイムな入札、落札結果の管理。
- **複数ロール対応**: 管理者、漁師、仲買人それぞれの専用機能を提供。
- **請求・精算**: 落札結果に基づいた請求書（Invoice）の自動生成。
- **リアルタイム通知**: 重要イベントのプッシュ通知機能。

## 技術構成 (Tech Stack)

このプロジェクトは以下の技術スタックで構成されています。詳細は各ディレクトリの README を参照してください。

- **[Frontend (Next.js)](./frontend/README.md)**: Next.js 16, TypeScript, Panda CSS, TanStack Query, Yarn
- **[Backend (Go)](./backend/README.md)**: Go, gorilla/mux, database/sql, PostgreSQL, Redis
- **Infrastructure**: Docker / Docker Compose, Nginx (Reverse Proxy), Mailhog

---

## セットアップ (Setup)

このプロジェクトは Docker Compose を使用して素早く開発環境を構築できます。各サービスを個別に実行したい場合は、それぞれのディレクトリの README を参照してください。

- **[Frontend 開発ガイド](./frontend/README.md)**
- **[Backend 開発ガイド](./backend/README.md)**

### クイックスタート (Docker Compose)

#### 1. 前提条件の確認

以下のツールがインストールされていることを確認してください。

- **Docker Desktop** / **Docker Compose**
- **mkcert**: ローカルで HTTPS を使用するために必要です。

#### 2. SSL 証明書の生成

Proxy (Nginx) で使用する証明書を生成します。

```bash
mkdir -p nginx/certs
mkcert -key-file nginx/certs/key.pem -cert-file nginx/certs/cert.pem localhost 127.0.0.1
```

#### 3. アプリケーションの起動

```bash
docker-compose up -d --build
```

起動後、以下のURLでアクセス可能です。

- **Frontend**: [https://localhost](https://localhost)
- **Backend API**: [https://localhost/api/](https://localhost/api/)
- **Mailhog (メール確認)**: [http://localhost:8025](http://localhost:8025)

## ディレクトリ構造 (Directory Structure)

```text
.
├── frontend/           # Next.js アプリケーション ([README](./frontend/README.md))
├── backend/            # Go API サーバー ([README](./backend/README.md))
├── nginx/              # Nginx 設定 (Reverse Proxy / SSL)
└── docker-compose.yml
```

## 開発ルール (Development Rules)

開発の進め方やAIエージェントの利用については [AGENTS.md](./AGENTS.md) を参照してください。
