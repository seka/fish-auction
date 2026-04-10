package service

import (
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// Clock provides an interface for time-related operations.
type Clock interface {
	Now() time.Time
	NowIn(location model.LocationName) time.Time
}

type realClock struct{}

// NewRealClock creates a new Clock instance that uses the real system time.
func NewRealClock() Clock {
	return &realClock{}
}

func (c *realClock) Now() time.Time {
	return time.Now()
}

func (c *realClock) NowIn(name model.LocationName) time.Time {
	return model.NewTimeZone(name).Now()
}
