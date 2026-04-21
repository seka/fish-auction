#!/bin/bash
set -x

# メインキューの作成
awslocal sqs create-queue \
  --region ap-northeast-1 \
  --queue-name notification-queue \
  --attributes '{"VisibilityTimeout":"30","MessageRetentionPeriod":"86400"}'

# DLQ の作成
awslocal sqs create-queue \
  --region ap-northeast-1 \
  --queue-name notification-queue-dlq \
  --attributes '{"MessageRetentionPeriod":"604800"}'

# メインキューに DLQ を紐付け
awslocal sqs set-queue-attributes \
  --region ap-northeast-1 \
  --queue-url http://sqs.ap-northeast-1.localhost.localstack.cloud:4566/000000000000/notification-queue \
  --attributes '{"RedrivePolicy":"{\"deadLetterTargetArn\":\"arn:aws:sqs:ap-northeast-1:000000000000:notification-queue-dlq\",\"maxReceiveCount\":\"3\"}"}'
