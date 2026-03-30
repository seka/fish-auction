---
name: fish-auction-data-model
description: Use when changing fish-auction data structures, tracing relationships between domain models, database tables, frontend types, and validation schemas, or planning a schema/API change.
---

# Fish Auction Data Model

## When to use

- DB 変更、API 変更、型変更の影響範囲を見たいとき
- domain model / DB row / frontend model を混同したくないとき
- `auction`, `item`, `venue`, `buyer`, `invoice` まわりの関係を整理したいとき

## Source of truth

- DB の実体: `backend/migrations/001_init.sql`
- ドメイン表現: `backend/internal/domain/model`
- 永続化の row 表現: `backend/internal/infrastructure/entity`
- サーバー側データの型定義: `frontend/src/data/entities` (@entities)
- フロントの機能固有モデル: `frontend/src/features/*/types`
- フロントの入力検証: `frontend/src/schemas`

迷ったら、テーブル構造は migration、業務上の意味は domain model、画面都合の shape は frontend model を優先して読む。

## Main aggregates

- `Venue`
  会場。`auctions.venue_id` で参照される。
- `Auction`
  せり本体。domain では `AuctionPeriod` と `AuctionStatus` を持つ。
- `AuctionItem`
  出品物。`auction_id` と `fisherman_id` を持ち、入札状態と並び順を持つ。
- `Bid`
  入札。domain では `BidPrice` value object を使う。
- `Purchase` / `InvoiceItem`
  落札結果や請求の読み取りモデル。永続化テーブル名と 1 対 1 でない場合がある。
- `Buyer` / `Authentication`
  仲買人本体と認証情報は分かれている。
- `Admin` / `Session`
  管理者認証とセッション。
- `PushSubscription`
  buyer に紐づく通知購読。

## Relationship map

- `venues` 1 - n `auctions`
- `auctions` 1 - n `auction_items`
- `fishermen` 1 - n `auction_items`
- `buyers` 1 - 1..n `authentications`, `transactions`, `push_subscriptions`
- `auction_items` 1 - n `transactions`

`Purchase` と `InvoiceItem` は集計結果や参照用モデルとして扱われることがあるので、追加変更時は repository / usecase を先に読む。

## Important layer differences

- domain model は value object を持つ
  例: `Auction.Period`, `Bid.Price`
- infrastructure entity は DB 列に近い
  例: `entity.Bid.Price` は `int`, `entity.AuctionItem` は `db` tag を持つ
- frontend model は画面で使いやすい primitive に寄せている
  例: `Auction` は `auctionDate`, `startTime`, `endTime` を文字列で持つ
- frontend schema は入力制約を持つ
  例: `auctionSchema`, `venueSchema`, `bidSchema`

## Read in this order

1. `backend/migrations/001_init.sql`
2. `backend/internal/domain/model/<target>.go`
3. `backend/internal/infrastructure/entity/<target>.go`
4. `backend/internal/server/dto` と handler
5. `frontend/src/models` と `frontend/src/models/schemas`

## Guardrails

- domain model と DB row を同一視しない
- frontend の型を backend の真実の源泉にしない
- `Purchase` や `InvoiceItem` のような参照モデルは、元テーブルだけで判断しない
- ステータス値は domain model と schema の両方を確認する
- nullable な列は frontend 型と domain 型で表現差がないか確認する

## Search patterns

- `rg "type .* struct" backend/internal/domain/model backend/internal/infrastructure/entity`
- `rg "CREATE TABLE|REFERENCES|UNIQUE|CHECK" backend/migrations/001_init.sql`
- `rg "export interface|z\\.object|z\\.enum" frontend/src/models`

## Working rule

モデル変更では、migration、domain model、entity、frontend model、schema の 5 点セットで見る。どれか 1 つだけ直して終わりにしない。
