package errors

import (
	"errors"

	domainErrors "github.com/seka/fish-auction/backend/internal/domain/errors"
)

// PushNotificationError represents an error that occurred while sending a push notification.
type PushNotificationError struct {
	StatusCode int
	Message    string
}

func (e *PushNotificationError) Error() string {
	return e.Message
}

// HandleError converts implementation-specific push errors into domain-specific errors.
func HandleError(err error, endpoint string) error {
	if err == nil {
		return nil
	}

	var pushErr *PushNotificationError
	// If it's a structural PushNotificationError
	if errors.As(err, &pushErr) {
		switch pushErr.StatusCode {
		case 410: // Gone
			return &domainErrors.GoneError{Resource: "Subscription", ID: endpoint, Message: pushErr.Message}
		case 404: // Not Found
			return &domainErrors.NotFoundError{Resource: "Subscription", ID: endpoint}
		}
	}

	// For other errors, return as is (could be network errors etc.)
	return err
}
