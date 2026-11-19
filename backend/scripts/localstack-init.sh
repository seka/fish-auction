#!/bin/bash
# LocalStack 上に SQS キューを作成する初期化スクリプト。
# - docker-compose: LocalStack コンテナの /etc/localstack/init/ready.d から起動時に実行される
# - CI (backend-test.yml): ランナー上から実行される
set -euxo pipefail

# このスクリプトは LocalStack 専用。実行環境に本物の AWS 認証情報やエンドポイントが
# 設定されていても使わないよう、既定値ではなく常に上書きする。
# 接続先を変えたい場合は AWS 共通の AWS_ENDPOINT_URL ではなく LOCALSTACK_ENDPOINT を使う。
export AWS_ENDPOINT_URL="${LOCALSTACK_ENDPOINT:-http://localhost:4566}"
# LocalStack はアクセスキー ID が 12 桁の数字だとそれをアカウント ID として扱う。
# キュー URL / ARN は既定アカウント 000000000000 を前提にしているため、数字以外の値に固定する。
# 例: test -> 000000000000（OK） / 111111111111 -> 別アカウントに作成され URL が一致しない（NG）
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=ap-northeast-1
unset AWS_SESSION_TOKEN AWS_PROFILE

# 接続先が LocalStack でない場合は何も作成せずに中止する
if ! curl -sf "${AWS_ENDPOINT_URL}/_localstack/health" >/dev/null; then
  echo "LocalStack に接続できないため中止します: ${AWS_ENDPOINT_URL}" >&2
  exit 1
fi

# メインキューの作成
aws sqs create-queue \
  --queue-name notification-queue \
  --attributes '{"VisibilityTimeout":"30","MessageRetentionPeriod":"86400"}'

# DLQ の作成
aws sqs create-queue \
  --queue-name notification-queue-dlq \
  --attributes '{"MessageRetentionPeriod":"604800"}'

# メインキューに DLQ を紐付け
# キュー URL はエンドポイントの形式によって変わるため、固定値ではなく取得した値を使う
main_queue_url="$(aws sqs get-queue-url --queue-name notification-queue --query QueueUrl --output text)"
aws sqs set-queue-attributes \
  --queue-url "${main_queue_url}" \
  --attributes '{"RedrivePolicy":"{\"deadLetterTargetArn\":\"arn:aws:sqs:ap-northeast-1:000000000000:notification-queue-dlq\",\"maxReceiveCount\":\"3\"}"}'

# 統合テスト用キュー（コンテナ側 worker と分離するため、本番 notification-queue とは別）
aws sqs create-queue \
  --queue-name notification-queue-test \
  --attributes '{"VisibilityTimeout":"30","MessageRetentionPeriod":"86400"}'
