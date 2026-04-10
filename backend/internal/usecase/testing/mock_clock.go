package testing

import (
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockClock provides a mock implementation of the Clock interface for testing.
type MockClock struct {
	NowTime time.Time
}

// NewMockClock creates a new MockClock with a fixed time.
func NewMockClock(t time.Time) *MockClock {
	return &MockClock{NowTime: t}
}

func (m *MockClock) Now() time.Time {
	return m.NowTime
}

func (m *MockClock) NowIn(name model.LocationName) time.Time {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	return m.NowTime.In(jst)
}

func (m *MockClock) SetNow(t time.Time) {
	m.NowTime = t
}
