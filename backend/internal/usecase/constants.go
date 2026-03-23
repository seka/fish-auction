package usecase

import "time"

const (
	// AuctionExtensionThreshold is the time remaining before the auction ends
	// during which a new bid will trigger an extension.
	AuctionExtensionThreshold = 5 * time.Minute

	// AuctionExtensionDuration is the duration by which the auction will be
	// extended when ShouldExtend is true.
	AuctionExtensionDuration = 5 * time.Minute

	// MaxFailedLoginAttempts is the number of consecutive failed login
	// attempts before an account is locked.
	MaxFailedLoginAttempts = 5

	// AccountLockDuration is the duration for which an account is locked
	// after exceeding MaxFailedLoginAttempts.
	AccountLockDuration = 30 * time.Minute
)
