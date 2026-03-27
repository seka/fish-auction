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
| /api/health | GET | 20 | ヘルスチェック |
| /api/items | GET | 25 | アイテム一覧 |
| /api/auctions | GET | 20 | 競り一覧 |
| /api/venues | GET | 15 | 会場一覧 |
| /api/invoices | GET | 10 | 請求一覧 |
| /api/items?status=Available | GET | 10 | 出品中アイテム一覧 |

`LOAD_TEST_AUCTION_IDS` を指定すると、次の detail 系 endpoint も自動で対象に含みます。

- `/api/auctions/{id}`
- `/api/auctions/{id}/items`

`LOAD_TEST_VENUE_IDS` を指定すると、次の endpoint も自動で対象に含みます。

- `/api/venues/{id}`

## 環境変数

| 変数名 | デフォルト値 | 説明 |
|---|---|---|
| LOAD_TEST_CONCURRENCY | 1000 | 同時接続数 |
| LOAD_TEST_REQUESTS | 10000 | リクエスト総数 |
| LOAD_TEST_DURATION | 0 | 実行時間（秒）、0=無制限 |
| LOAD_TEST_TARGET_URL | http://localhost:8080 | テスト対象URL |
| LOAD_TEST_AUCTION_IDS |  | カンマ区切りの auction ID。detail endpoint を対象に含める |
| LOAD_TEST_VENUE_IDS |  | カンマ区切りの venue ID。detail endpoint を対象に含める |

## 使用例

### ローカル環境でのテスト
```bash
# デフォルト設定で実行（1000並行、10000リクエスト）
go test ./cmd/testing/... -run TestLoadTest -v

# 軽量テスト（10並行、100リクエスト）
LOAD_TEST_CONCURRENCY=10 \
LOAD_TEST_REQUESTS=100 \
go test ./cmd/testing/... -run TestLoadTest -v

# detail endpoint も含めて実行
LOAD_TEST_AUCTION_IDS=1,2,3 \
LOAD_TEST_VENUE_IDS=1,2 \
LOAD_TEST_CONCURRENCY=50 \
LOAD_TEST_REQUESTS=1000 \
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

# 有効な detail ID を使って実運用に近い read-heavy テスト
LOAD_TEST_TARGET_URL=https://staging.example.com \
LOAD_TEST_AUCTION_IDS=101,102,103 \
LOAD_TEST_VENUE_IDS=10,11 \
LOAD_TEST_CONCURRENCY=500 \
LOAD_TEST_REQUESTS=20000 \
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

1. **対象API**: 現在の負荷試験は public の read-only endpoint を中心に実行します
2. **detail ID**: `LOAD_TEST_AUCTION_IDS`, `LOAD_TEST_VENUE_IDS` には実在する ID を指定してください
3. **リソース**: 高並行数（5000+）の場合、システムのファイルディスクリプタ制限に注意してください
4. **本番環境**: 本番環境でテストする場合は、事前に影響を確認してください
