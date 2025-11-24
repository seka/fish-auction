# 漁港のせり (Fish Auction)

漁港のせりシステムの実装プロジェクトです。

## 技術構成 (Tech Stack)

- **Frontend**: Next.js (React)
- **Backend**: Go
- **Database**: PostgreSQL
- **Web Server**: Nginx (Reverse Proxy)
- **Infrastructure**: Docker Compose

## 開発環境のセットアップ (Setup)

### Requirements

- Docker
- Docker Compose

### Setup

#### 1. HTTPS Certificates (mkcert)

開発環境ではHTTPSを使用します。`mkcert` を使用して証明書を作成してください。

```bash
# Install local CA (First time only)
mkcert -install

# Generate certificates
mkdir -p nginx/certs
mkcert -key-file nginx/certs/key.pem -cert-file nginx/certs/cert.pem localhost 127.0.0.1
```

#### 2. Build and Run

以下のコマンドを実行して、アプリケーションをビルド・起動します。

```bash
docker-compose up -d --build
```

### Access

起動後、以下のURLで各サービスにアクセスできます。

- **Frontend**: [http://localhost](http://localhost)
- **Backend API**: [http://localhost/api/](http://localhost/api/)
  - Health Check: [http://localhost/api/health](http://localhost/api/health)

## Directory Structure

- `frontend/`: Next.js アプリケーション
- `backend/`: Go API サーバー
- `nginx/`: Nginx 設定
- `docker-compose.yml`: コンテナ構成定義
- `AGENTS.md`: AI エージェント用ガイドライン

## Development Rules

開発の進め方やAIエージェントの利用については [AGENTS.md](./AGENTS.md) を参照してください。
