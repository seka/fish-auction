---
name: fish-auction-test-ops
description: Use ONLY when explicitly instructed to run tests or heavyweight operations using Docker to bypass sandbox restrictions or verify complete system behavior.
---

# Fish Auction Test Operations

## When to use

- ユーザーから「テストを実行して」「テストを確認して」と明示的に指示されたとき
- 実装の最終確認として、クリーンな環境での整合性を検証したいとき
- サンドボックス環境の制限によりローカルコマンド（`make test`等）が失敗し、Docker での回避が必要なとき

## Docker-based Verification

ビルドキャッシュやディレクトリ権限の問題を回避するため、Docker コンテナ内でテストを実行します。これらのコマンドは実行に時間がかかる場合があるため、必要最小限の範囲で実行してください。

### Backend Tests
```bash
docker-compose run --rm backend make app.test
```
*統合テストも含める場合:*
```bash
docker-compose run --rm backend make app.integration-test
```

### Frontend Tests
```bash
docker-compose run --rm frontend yarn test
```

## Working rule

このスキルに含まれるコマンドはリソースを消費し、待ち時間が発生します。自動的に実行せず、必ず実行前にユーザーの指示があることを確認してください。
