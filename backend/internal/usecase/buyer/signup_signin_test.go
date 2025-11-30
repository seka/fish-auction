package buyer_test

import (
	"context"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
	mock "github.com/seka/fish-auction/backend/internal/usecase/testing"
)

func TestSignupAndSigninFlow(t *testing.T) {
	// In-memory storage
	buyers := make(map[int]*model.Buyer)
	auths := make(map[string]*model.Authentication) // email -> auth
	nextID := 1

	// Setup Mocks
	buyerRepo := &mock.MockBuyerRepository{
		CreateFunc: func(ctx context.Context, b *model.Buyer) (*model.Buyer, error) {
			b.ID = nextID
			nextID++
			buyers[b.ID] = b
			return b, nil
		},
		FindByIDFunc: func(ctx context.Context, id int) (*model.Buyer, error) {
			if b, ok := buyers[id]; ok {
				return b, nil
			}
			return nil, nil // or error
		},
	}

	authRepo := &mock.MockAuthenticationRepository{
		CreateFunc: func(ctx context.Context, a *model.Authentication) (*model.Authentication, error) {
			a.ID = nextID // simple ID generation
			nextID++
			auths[a.Email] = a
			return a, nil
		},
		FindByEmailFunc: func(ctx context.Context, email string) (*model.Authentication, error) {
			if a, ok := auths[email]; ok {
				return a, nil
			}
			return nil, context.DeadlineExceeded // Using a dummy error, or create a custom one
		},
		UpdateLoginSuccessFunc: func(ctx context.Context, id int, loginAt time.Time) error {
			// Find auth by ID (inefficient but fine for test)
			for _, a := range auths {
				if a.ID == id {
					// In a real DB this would update the record
					return nil
				}
			}
			return nil
		},
		IncrementFailedAttemptsFunc: func(ctx context.Context, id int) error {
			for _, a := range auths {
				if a.ID == id {
					a.FailedAttempts++
					return nil
				}
			}
			return nil
		},
		LockAccountFunc: func(ctx context.Context, id int, until time.Time) error {
			for _, a := range auths {
				if a.ID == id {
					a.LockedUntil = &until
					return nil
				}
			}
			return nil
		},
	}

	createUC := buyer.NewCreateBuyerUseCase(buyerRepo, authRepo)
	loginUC := buyer.NewLoginBuyerUseCase(buyerRepo, authRepo)

	ctx := context.Background()
	email := "test@example.com"
	password := "securepassword123"

	// 1. Signup
	createdBuyer, err := createUC.Execute(ctx, "Test User", email, password, "Test Org", "Contact")
	if err != nil {
		t.Fatalf("Signup failed: %v", err)
	}
	if createdBuyer == nil {
		t.Fatal("Signup returned nil buyer")
	}

	// 2. Signin (Success)
	loggedInBuyer, err := loginUC.Execute(ctx, email, password)
	if err != nil {
		t.Fatalf("Signin failed with correct password: %v", err)
	}
	if loggedInBuyer.ID != createdBuyer.ID {
		t.Errorf("Signin returned wrong buyer ID: got %d, want %d", loggedInBuyer.ID, createdBuyer.ID)
	}

	// 3. Signin (Failure - Wrong Password)
	_, err = loginUC.Execute(ctx, email, "wrongpassword")
	if err == nil {
		t.Error("Signin should fail with wrong password, but it succeeded")
	}

	// 4. Signin (Failure - Wrong Email)
	_, err = loginUC.Execute(ctx, "wrong@example.com", password)
	if err == nil {
		t.Error("Signin should fail with wrong email, but it succeeded")
	}
}
