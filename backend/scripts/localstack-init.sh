#!/bin/bash
# LocalStack 上に SQS キューを作成する初期化スクリプト。
# - docker-compose: LocalStack コンテナの /etc/localstack/init/ready.d から起動時に実行される
# - CI (backend-test.yml): ランナー上から実行される
set -euxo pipefail

# このスクリプトは LocalStack 専用（本番のキューは対象外）。
# 実行環境に本物の AWS 認証情報やエンドポイントが設定されていても使わないよう、
# 環境変数からは受け取らず固定値で上書きする。aws コマンドに渡すため export する。
# 接続先は CI（ランナー）・docker-compose（LocalStack コンテナ内）のどちらでも localhost:4566。
export AWS_ENDPOINT_URL="http://localhost:4566"

# LocalStack はアクセスキー ID が 12 桁の数字だとそれをアカウント ID として扱う。
# キュー URL / ARN は既定アカウント 000000000000 を前提にしているため、数字以外の値に固定する。
# 例: test -> 000000000000（OK） / 111111111111 -> 別アカウントに作成され URL が一致しない（NG）
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=ap-northeast-1

# プロファイル側の endpoint_url / region や、本物の一時認証情報が混ざらないよう
# このスクリプト内部では PROFILE を見ないようにする
unset AWS_SESSION_TOKEN AWS_PROFILE

# LocalStack に接続できない場合は、中止する。
if ! curl -sf --max-time 5 "${AWS_ENDPOINT_URL}/_localstack/health" >/dev/null; then
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
