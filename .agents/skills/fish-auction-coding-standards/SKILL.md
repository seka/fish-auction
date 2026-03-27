---
name: fish-auction-coding-standards
description: Use when implementing or reviewing code in fish-auction and you need the project coding standards for TypeScript, Go, naming, typing, formatting, and basic function design.
---

# Fish Auction Coding Standards

## When to use

- 実装前にこのリポジトリの基本的な書き方を揃えたいとき
- レビュー時に最低限の実装品質を確認したいとき
- TypeScript と Go の両方で、既存コードに寄せた判断をしたいとき

## General

- TypeScript は `yarn format`、Go は `gofmt` に従う
- 変数名、関数名、型名は役割が分かる名前を優先する
- 関数は小さく保ち、責務を増やしすぎない (30行を目安とする)
- 既存の隣接ファイルの命名や分割方針を優先する

## TypeScript

- strict な型前提で扱う
- `any` は避け、型定義か Zod schema を使う
- API の入出力には型を付ける
- UI、query、API 通信の責務を混ぜない
- 整形や lint の必要がある場合は既存スクリプトに従う
- `yarn lint:ci` のエラー件数は 0 件にする

## Go

- エラー処理は明示的に行う
- `context.Context` を受け渡す
- HTTP 層、usecase、repository の責務を跨いでロジックを置かない
- 既存の usecase / handler の書式と命名に合わせる
- `golangci-lint` のエラー件数は 0 件にする

## Review checklist

- 不要な抽象化を増やしていないか
- 型や DTO が欠けていないか
- エラー経路が抜けていないか
- 変更対象のレイヤーに対して実装場所が適切か
- 近い既存実装と比べて不自然な命名になっていないか
- テストは問題ないか

## Read next only when needed

- フロントの詳細: `frontend/README.md`, `frontend/docs/ARCHITECTURE.md`
- バックの詳細: `backend/README.md`
- 実装配置の判断: `fish-auction-frontend-patterns`, `fish-auction-backend-patterns`

## Working rule

ルールを一般論で適用しすぎない。迷ったら、このコードベースの近傍実装を優先して合わせる。
