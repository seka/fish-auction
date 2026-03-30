---
name: fish-auction-frontend-patterns
description: Use when changing the fish-auction frontend built with Next.js, React, Panda CSS, TanStack Query, react-hook-form, and Zod, especially to place code in the correct layer and follow existing patterns.
---

# Fish Auction Frontend Patterns

## When to use

- フロントエンドの実装や修正をするとき
- 新しい画面要素、フォーム、query、API 呼び出しの置き場所を決めたいとき
- 既存パターンに合わせて最小差分で直したいとき

## Layer rules

- `frontend/app`: route と page の入口
- `frontend/src/features/*/components`: 機能固有の UI
- `frontend/src/components`: 複数機能で共有する UI
- `frontend/src/features/*/states`: UI ロジックや hook
- `frontend/src/features/*/queries`: 機能固有の query / mutation (データの型変換・キャストを行う境界)
- `frontend/src/features/*/types`: 機能固有の型定義（ドメインモデル + UI型）
- `frontend/src/data/api`: バックエンド API との 1 対 1 通信 (apiClient)
- `frontend/src/data/entities`: サーバー側データ構造の定義 (@entities)
- `frontend/src/schemas`: Zod スキーマ
- `frontend/src/core/api/client.ts`: 共通の fetch ラッパー

## Practical placement guide

- API エンドポイントを叩くだけの関数は `src/data/api`
- React Query の状態管理は `src/data/queries` か feature 配下の `queries`
- UI と query の組み合わせやイベント調停は `states`
- 共有不能な見た目は feature 配下、共有可能な見た目は `src/components`

## Existing conventions to preserve

- `apiClient` を経由して通信する
- Query key は domain ごとにまとめる
- mutation 後は影響する key を invalidate する
- 公開画面と管理画面の両方に効く更新では、両方の key を確認する
- Panda CSS と既存トークンを優先し、場当たりのインライン実装は増やさない

## Important nuance

hook 名は feature により異なる。`useXXXManagement` に固定せず、`useLogin` のような軽量な命名も含めて、同じ feature の隣接ファイルに合わせる。

## First files to inspect

- `frontend/docs/ARCHITECTURE.md`
- 変更対象 feature の `components`, `states`, `queries`
- `frontend/src/data/api/*`
- `frontend/src/core/api/client.ts`

## Search patterns

- `rg "useQuery|useMutation|invalidateQueries" frontend/src`
- `rg "apiClient\\.(get|post|put|delete)" frontend/src`
- `rg "Schema|FormData" frontend/src/models frontend/src/features`

## Working rule

フロント変更では、まず同じ feature の隣接実装を 2〜3 ファイル読む。新しい抽象化を作る前に、既存の query key、schema、container 構成に寄せる。
