# 負荷テスト使用方法

## 概要
`cmd/testing/integration_stress_test.go` は実環境での負荷テストや同時接続テストに使用できる統合テストです。

## 機能
- 複数のエンドポイントにランダムな重みでリクエスト送信
- 並行リクエスト（デフォルト: 1000 goroutine）
- メトリクス収集（成功/失敗率、レスポンスタイム、スループット）
- 環境変数で設定可能

## エンドポイントと重み

| エンドポイント | メソッド | 重み | 説明 |
|---|---|---|---|
| /api/health | GET | 50 | ヘルスチェック |
| /api/items | GET | 20 | アイテム一覧 |
| /api/items | POST | 5 | アイテム作成 |
| /api/fishermen | GET | 10 | 漁師一覧 |
| /api/fishermen | POST | 5 | 漁師作成 |
| /api/buyers | GET | 5 | 買い手一覧 |
| /api/buyers | POST | 5 | 買い手作成 |

## 環境変数

| 変数名 | デフォルト値 | 説明 |
|---|---|---|
| LOAD_TEST_CONCURRENCY | 1000 | 同時接続数 |
| LOAD_TEST_REQUESTS | 10000 | リクエスト総数 |
| LOAD_TEST_DURATION | 0 | 実行時間（秒）、0=無制限 |
| LOAD_TEST_TARGET_URL | http://localhost:8080 | テスト対象URL |

## 使用例

### ローカル環境でのテスト
```bash
# デフォルト設定で実行（1000並行、10000リクエスト）
go test ./cmd/testing/... -run TestLoadTest -v

# 軽量テスト（10並行、100リクエスト）
LOAD_TEST_CONCURRENCY=10 \
LOAD_TEST_REQUESTS=100 \
go test ./cmd/testing/... -run TestLoadTest -v

# 高負荷テスト（5000並行、100000リクエスト）
LOAD_TEST_CONCURRENCY=5000 \
LOAD_TEST_REQUESTS=100000 \
go test ./cmd/testing/... -run TestLoadTest -v
```

### 実環境でのテスト
```bash
# 本番環境に対して60秒間の負荷テスト
LOAD_TEST_TARGET_URL=https://production.example.com \
LOAD_TEST_CONCURRENCY=1000 \
LOAD_TEST_DURATION=60 \
go test ./cmd/testing/... -run TestLoadTest -v

# ステージング環境で同時接続テスト
LOAD_TEST_TARGET_URL=https://staging.example.com \
LOAD_TEST_CONCURRENCY=2000 \
LOAD_TEST_REQUESTS=50000 \
go test ./cmd/testing/... -run TestLoadTest -v
```

## メトリクス出力例

```
=== Load Test Results ===
Total Requests:     10000
Successful:         9950 (99.50%)
Failed:             50 (0.50%)
Duration:           7.3s
Throughput:         1369.86 req/s

Response Times:
  Min:              1ms
  Max:              250ms
  Mean:             30ms
  P50:              25ms
  P95:              80ms
  P99:              150ms

Errors:
  HTTP 500: 30
  Connection timeout: 20
```

## 注意事項

1. **テストデータ**: POST リクエストはランダムなテストデータを生成します
2. **エラー**: 存在しないIDを参照する場合があるため、一部エラーが発生する可能性があります
3. **リソース**: 高並行数（5000+）の場合、システムのファイルディスクリプタ制限に注意してください
4. **本番環境**: 本番環境でテストする場合は、事前に影響を確認してください
