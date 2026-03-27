---
name: fish-auction-backend-patterns
description: Use when changing the fish-auction Go backend, especially handlers, DTOs, usecases, repository boundaries, datastore implementations, migrations, or tests.
---

# Fish Auction Backend Patterns

## When to use

- Go バックエンドを修正するとき
- handler / usecase / repository の責務を崩さず実装したいとき
- DB 変更や新規 API をどこまで追加すべきか判断したいとき

## Layer rules

- `backend/internal/server/handler`: HTTP 入出力、DTO 変換、ルーティング
- `backend/internal/server/dto`: request / response DTO
- `backend/internal/usecase/<domain>`: 業務ロジック
- `backend/internal/domain/model`: ドメインモデル
- `backend/internal/domain/repository`: repository interface
- `backend/internal/infrastructure/datastore/*`: 永続化実装
- `backend/internal/registry`: handler や usecase の組み立て
- `backend/migrations`: スキーマ変更

## Handler pattern to preserve

近い handler に合わせる。典型的には次の流れ。

1. DTO へ decode
2. domain model へ詰め替え
3. usecase を呼ぶ
4. `util.HandleError` でエラー処理
5. DTO か JSON を返す

## Usecase pattern to preserve

- package はドメイン単位
- 1 ファイル 1 役割のユースケースとする
- interface と実装を同ファイルに置く構成が多い
- repository 経由で外部依存を隠す
- 近い既存ユースケースの命名を流用する

## Change checklist

- API 入出力が変わる: handler + dto + frontend の API 呼び出しを確認
- 業務ルールが変わる: usecase と対応テストを確認
- 永続化が変わる: domain repository + infrastructure + migration を確認
- 新規エンドポイント追加: handler の route 登録と registry 配線を確認
- 認証や role に触れる: session / middleware / cookie 周辺を確認

## Tests

- usecase テストは同じ package 配下の `*_test.go`
- handler テストは `backend/internal/server/handler/*_test.go`
- mock は `backend/internal/usecase/testing` と `backend/internal/server/testing`

## Important nuance

ルーティング定義は handler ごとに少し差がある。`mux.HandleFunc("/api/venues", ...)` 形式と、`mux.HandleFunc("POST /api/admin/password-reset/request", ...)` 形式が混在しているので、触るファイルの近傍スタイルに合わせる。

## Search patterns

- `rg "HandleFunc|RegisterRoutes" backend/internal/server/handler`
- `rg "type .*UseCase interface|New.*UseCase|Execute\\(" backend/internal/usecase`
- `rg "repository\\." backend/internal/usecase backend/internal/domain`

## Working rule

バックエンド変更では、対象 domain の usecase と handler の両方を先に読む。新規実装はそこに合わせ、repository 境界をまたぐロジックを handler に置かない。
