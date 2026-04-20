#!/bin/bash
awslocal sqs create-queue \
  --queue-name notification-queue \
  --attributes '{"VisibilityTimeout":"30","MessageRetentionPeriod":"86400"}'

awslocal sqs create-queue \
  --queue-name notification-queue-dlq \
  --attributes '{"MessageRetentionPeriod":"604800"}'

# DLQ 設定
awslocal sqs set-queue-attributes \
  --queue-url http://localhost:4566/000000000000/notification-queue \
  --attributes '{"RedrivePolicy":"{\"deadLetterTargetArn\":\"arn:aws:sqs:ap-northeast-1:000000000000:notification-queue-dlq\",\"maxReceiveCount\":\"3\"}"}'
