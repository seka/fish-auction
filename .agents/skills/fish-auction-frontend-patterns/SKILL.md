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

- `frontend/app/*/components`: Container (Page orchestration) components
- `frontend/src/features/*/components`: Domain-specific UI widgets (List, Form, etc.)
- `frontend/src/components`: UI components shared across multiple features
- `frontend/src/features/*/states`: UI logic or domain-specific hooks
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
- `rg "Schema|FormData" frontend/src/schemas frontend/src/features`

## Working rule

フロント変更では、まず同じ feature の隣接実装を 2〜3 ファイル読む。新しいページ要素を追加する場合は `app/**/components/` に Container を作成し、`src/features/**/components` から提供される再利用可能なコンポーネントや `queries`, `states` と組み合わせてオーケストレーションを行う。
