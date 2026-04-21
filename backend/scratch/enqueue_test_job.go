package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/queue/sqs"
	"github.com/seka/fish-auction/backend/internal/job/message"
)

func main() {
	// Configure SQS client for LocalStack
	region := os.Getenv("SQS_REGION")
	if region == "" {
		region = "ap-northeast-1"
	}
	queueURL := os.Getenv("SQS_QUEUE_URL")
	if queueURL == "" {
		queueURL = "http://localhost:4566/000000000000/notification-queue"
	}
	endpoint := os.Getenv("SQS_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:4566"
	}

	ctx := context.Background()
	client, err := sqs.NewClient(ctx, region, queueURL, endpoint)
	if err != nil {
		log.Fatalf("failed to create SQS client: %v", err)
	}

	// 1. Create a test payload
	// In the real app, this is what would be sent to the buyer's browser
	testPayload := map[string]string{
		"title": "テスト通知",
		"body":  "ジョブキューの正常動作を確認しました。",
		"url":   "/auctions",
	}

	// 2. Wrap it in a PushNotificationMessage
	jobMessage := message.PushNotificationMessage{
		BuyerID: 1, // '株式会社 魚河岸' in seed.sql
		Payload: testPayload,
	}

	fmt.Printf("Enqueuing test job to %s...\n", queueURL)
	fmt.Printf("Payload: %+v\n", jobMessage)

	// 3. Enqueue
	err = client.Enqueue(ctx, model.JobTypePushNotification, jobMessage)
	if err != nil {
		log.Fatalf("failed to enqueue job: %v", err)
	}

	fmt.Println("Successfully enqueued job!")
	fmt.Println("\nNext steps:")
	fmt.Println("1. Ensure 'docker-compose up localstack db worker' is running")
	fmt.Println("2. Check worker logs: 'docker-compose logs -f worker'")
	fmt.Println("3. Verify the worker processes the job and shows 'ReceiveCount = 1'")
}
