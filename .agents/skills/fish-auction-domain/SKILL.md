---
name: fish-auction-domain
description: Use when a task depends on fish-auction business language such as admin, buyer, fisherman, auction, bid, invoice, item, venue, notification, or role-specific behavior.
---

# Fish Auction Domain

## When to use

- 業務用語の意味を揃えたいとき
- ロールごとの差分を見落としたくないとき
- API や UI の変更がどの業務領域に属するか判断したいとき

## Actors seen in the codebase

- `admin`: 管理系操作、マスタ管理、管理者認証
- `buyer`: 仲買人向けの入札、購入、認証、通知
- `public`: 未認証向けの公開情報

## Main domain modules

- `auth`: 管理者認証とパスワードリセット
- `buyer`: 仲買人認証、一覧、購入・オークション参照
- `admin`: 管理者作成、パスワード更新・再設定
- `auction`: せり本体、状態更新、一覧・詳細
- `bid`: 入札
- `item`: 出品物の作成、更新、並び替え
- `venue`: 会場管理
- `invoice`: 請求一覧
- `notification`: Push 通知
- `fisherman`: 荷主管理

## Interpretation workflow

変更要求を読んだら、次の 5 点に分解する。

1. 誰の操作か: `admin` / `buyer` / `public`
2. 何の集約か: `auction` / `item` / `venue` / `invoice` など
3. 操作種別か参照系か: create / update / delete / list / get / login
4. 副作用はあるか: session、password reset、並び替え、通知、請求
5. どこに表示されるか: 管理画面だけか、公開画面にも波及するか

## First places to inspect

- 業務ロジック: `backend/internal/usecase/<domain>`
- ハンドラー: `backend/internal/server/handler`
- フロントの query: `frontend/src/data/queries`, `frontend/src/features/*/queries`
- 既存テスト: `*_test.go`, `*.test.ts`

## Guardrails

- ロール名が違う処理を混同しない
- `admin` と `public` の query key / invalidation を分けて考える
- ドメイン用語の意味は README だけで決めず、対応する usecase と test で確認する
- 新しい仕様を足すときは、近い既存ユースケース名に合わせて命名する

## Useful searches

- `rg "package admin|package buyer|package auction|package bid" backend/internal/usecase`
- `rg "admin|buyer|public" frontend/src/data/queries frontend/src/features`

## Working rule

このスキルは「業務語彙の入口」を与えるためのもの。仕様の断定が必要なら、同じドメインの usecase テストを優先して読む。
