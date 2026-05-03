package model

import (
	"fmt"
)

// JobType represents the type of an asynchronous job.
type JobType string

const (
	// JobTypePushOutbid is the job type for notifying a buyer that their bid was outbid.
	JobTypePushOutbid JobType = "push.outbid"
	// JobTypePushAuctionStatusChanged is the job type for notifying buyers that an auction status changed.
	JobTypePushAuctionStatusChanged JobType = "push.auction_status_changed"
	// JobTypeEmail is the job type for sending emails.
	JobTypeEmail JobType = "email"
)

// NewJobType creates a JobType from a string and validates it.
func NewJobType(s string) (JobType, error) {
	switch JobType(s) {
	case JobTypePushOutbid:
		return JobTypePushOutbid, nil
	case JobTypePushAuctionStatusChanged:
		return JobTypePushAuctionStatusChanged, nil
	case JobTypeEmail:
		return JobTypeEmail, nil
	default:
		return "", fmt.Errorf("unsupported job type: %s", s)
	}
}
