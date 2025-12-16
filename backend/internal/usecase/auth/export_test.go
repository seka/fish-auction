package auth

// SetRandRead sets the randRead function for testing
// Returns a function to restore the original value
func SetRandRead(f func(b []byte) (n int, err error)) func() {
	orig := randRead
	randRead = f
	return func() {
		randRead = orig
	}
}
