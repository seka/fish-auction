package admin

// SetRandRead replaces randRead for testing.
// Returns a function to restore the original value.
func SetRandRead(f func(b []byte) (n int, err error)) func() {
	original := randRead
	randRead = f
	return func() {
		randRead = original
	}
}

// MockRandRead creates a mock random reader that returns values or error
type MockRandRead struct {
	Err error
}

func (m *MockRandRead) Read(b []byte) (n int, err error) {
	if m.Err != nil {
		return 0, m.Err
	}
	// Fill with dummy data
	for i := range b {
		b[i] = 0x01
	}
	return len(b), nil
}

func GetRandReadFunc(err error) func([]byte) (int, error) {
	return func(b []byte) (int, error) {
		if err != nil {
			return 0, err
		}
		for i := range b {
			b[i] = 0x01
		}
		return len(b), nil
	}
}
