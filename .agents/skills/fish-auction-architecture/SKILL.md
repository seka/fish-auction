---
name: fish-auction-architecture
description: Use when working in fish-auction and you need to decide which files or layers to change, trace a request across frontend and backend, or estimate the impact of a feature or bugfix.
---

# Fish Auction Architecture

## When to use

- 変更箇所の当たりを最初に付けたいとき
- フロントからバックまでの処理経路を追いたいとき
- 影響範囲を洗い出してから実装したいとき

## Project map

- 構成は Next.js フロントエンドと Go バックエンドの 2 層
- `frontend/app`: App Router のルーティングと画面入口
- `frontend/src/features`: 機能単位の UI / state / query
- `frontend/src/data`: API 呼び出しと TanStack Query の共有層
- `frontend/src/core`: API クライアント、i18n、スタイルなどの他に依存関係のない開発の基盤
- `backend/internal/server`: handler / dto / middleware
- `backend/internal/usecase`: ユースケース単位の業務ロジック
- `backend/internal/domain`: モデル、repository interface、domain error
- `backend/internal/infrastructure`: PostgreSQL / Redis / Mailhog などの実装
- `backend/migrations`: DB 変更

## Communication

- クライアントとサーバーの通信は REST API を前提に見る
- 公開系は `http.ServeMux` ベース、管理系の一部は `gorilla/mux` のサブルーターを使う

## Request trace

通常は次の順で追う。

1. `frontend/app` か `frontend/src/features/*/components`
2. `frontend/src/features/*/states` または `queries`
3. `frontend/src/data/api`
4. `frontend/src/core/api/client.ts`
5. `backend/internal/server/handler`
6. `backend/internal/usecase`
7. `backend/internal/domain/repository`
8. `backend/internal/infrastructure/datastore`

## File selection rules

- 画面表示だけの変更: `frontend/src/features/*/components` か `frontend/src/components`
- フロントの送受信変更: `frontend/src/data/api`, `frontend/src/data/queries`, `frontend/src/data/entities`, `frontend/src/features/*/types`
- API の入出力変更: `backend/internal/server/dto`, `backend/internal/server/handler`
- 業務ルール変更: `backend/internal/usecase`
- 保存形式変更: `backend/internal/infrastructure/entity`, `datastore`, `migrations`
- ロールや認可に関わる変更: handler だけでなく session / middleware / usecase を確認

## Read next only when needed

- 全体像: `README.md`
- フロント詳細: `frontend/README.md`, `frontend/docs/ARCHITECTURE.md`
- バック詳細: `backend/README.md`

## Search patterns

- 画面起点で探す: `rg "router.push|useQuery|useMutation|apiClient" frontend`
- API 起点で探す: `rg "/api/" frontend backend`
- ドメイン起点で探す: `rg "auction|bid|buyer|fisherman|invoice|item|venue" frontend backend`

## Working rule

最初に「どの層まで触る変更か」を 1 行で決めてから読む。迷ったら handler と usecase を先に見て、UI 側は後追いする。
