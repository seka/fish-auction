---
name: fish-auction-dev-workflow
description: Use when working through a change in fish-auction and you need the expected development workflow for planning, branching, implementation, verification, git checks, pull requests, and CI confirmation.
---

# Fish Auction Dev Workflow

## When to use

- 作業の進め方そのものを揃えたいとき
- 実装前後の確認漏れを減らしたいとき
- ブランチ、PR、CI 確認まで含めた流れを整理したいとき

## Standard flow

1. 要件と既存コードを調査して、変更対象と影響範囲を決める
2. `master` から作業ブランチを切る
3. 既存パターンに寄せて最小差分で実装する
4. 必要な確認を行い、`git status` と `git diff` を見直す
5. レビューしやすい粒度でコミットを作成する
6. 作業ブランチを master にマージする

## Before editing

- 先に `rg` や `fd` で影響範囲を調べる
- どの層まで触る変更かを先に決める
- ドメインルールが絡むなら、対応する usecase と test を読む

## During editing

- 変更は最小限に保つ
- 既存ファイルの構成、命名、責務分離に合わせる
- unrelated な差分は巻き込まない

## Before finishing

- `git status` で意図しない変更がないか見る
- `git diff` で不要な差分、デバッグコード、説明不足を確認する
- 追加した変更に対して、必要なテストや CI 観点を整理する

## PR and CI

- Direct Push はしない
- PR は `gh pr create` を使う
- PR 作成後は GitHub Actions の結果を確認する
- CI が未確認なら、その旨を明示する

## Working rule

このワークフローは「順番の固定」が目的。調査を省略してすぐ編集しない。
