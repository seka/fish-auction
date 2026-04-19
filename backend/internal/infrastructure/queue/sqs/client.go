package sqs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/service"
	notificationMessage "github.com/seka/fish-auction/backend/internal/job/message"
)

// Client implements service.JobQueue using AWS SQS.
type Client struct {
	client   *sqs.Client
	queueURL string
}

var _ service.JobQueue = (*Client)(nil)

// NewClient creates a new SQS client.
func NewClient(ctx context.Context, region, queueURL, endpoint string) (*Client, error) {
	var opts []func(*config.LoadOptions) error
	if region != "" {
		opts = append(opts, config.WithRegion(region))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	sqsClient := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	})

	return &Client{
		client:   sqsClient,
		queueURL: queueURL,
	}, nil
}

// Enqueue sends a message to the SQS queue.
func (c *Client) Enqueue(ctx context.Context, jobType model.JobType, payload any) error {
	var body []byte
	var err error

	// Map domain payload to infrastructure DTO and marshal to JSON.
	switch jobType {
	case model.JobTypePushNotification:
		// Map the payload to the message. Use reflection or type assertion if more complex,
		// but for now we expect a compatible struct or map.
		msg := notificationMessage.PushNotificationMessage{
			// Since UseCase passes a tag-less struct or map, we can rely on json.Marshal/Unmarshal
			// or manual mapping here. For simplicity and to satisfy the user's "infra knowledge"
			// requirement, we can re-marshal/unmarshal if payload is already generic,
			// or directly populate if we know the domain structure.
			// Here we assume the UseCase passes a struct with matching field names.
			Payload: payload,
		}

		// Because the legacy implementation had BuyerID separately, we need to extract it if needed.
		// However, to keep UseCase clean, we can expect the payload passed to Enqueue to be the full data.
		// If payload is already the "parameters" from UseCase, we map it.
		if p, ok := payload.(map[string]any); ok {
			if bid, ok := p["BuyerID"].(int); ok {
				msg.BuyerID = bid
			}
			if bpayload, ok := p["Payload"]; ok {
				msg.Payload = bpayload
			}
		} else {
			// Fallback: Use json trick to fill DTO if payload is a tag-less struct
			tmp, _ := json.Marshal(payload)
			_ = json.Unmarshal(tmp, &msg)
		}

		body, err = json.Marshal(msg)
	default:
		body, err = json.Marshal(payload)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal job payload: %w", err)
	}

	_, err = c.client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(string(body)),
		QueueUrl:    aws.String(c.queueURL),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"JobType": {
				DataType:    aws.String("String"),
				StringValue: aws.String(string(jobType)),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send SQS message: %w", err)
	}
	return nil
}

// Dequeue polls for messages from the SQS queue.
func (c *Client) Dequeue(ctx context.Context, waitTimeSeconds int32) ([]*model.JobMessage, error) {
	output, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(c.queueURL),
		MaxNumberOfMessages:   10,
		WaitTimeSeconds:       waitTimeSeconds,
		MessageAttributeNames: []string{"All"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to receive SQS messages: %w", err)
	}

	messages := make([]*model.JobMessage, 0, len(output.Messages))
	for _, m := range output.Messages {
		jobType := model.JobType("")
		if attr, ok := m.MessageAttributes["JobType"]; ok && attr.StringValue != nil {
			jobType = model.JobType(*attr.StringValue)
		}

		messages = append(messages, &model.JobMessage{
			ID:            *m.MessageId,
			ReceiptHandle: *m.ReceiptHandle,
			JobType:       jobType,
			Payload:       []byte(*m.Body),
		})
	}

	return messages, nil
}

// DeleteMessage deletes a message from the SQS queue.
func (c *Client) DeleteMessage(ctx context.Context, message *model.JobMessage) error {
	_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: aws.String(message.ReceiptHandle),
	})
	if err != nil {
		return fmt.Errorf("failed to delete SQS message: %w", err)
	}
	return nil
}
