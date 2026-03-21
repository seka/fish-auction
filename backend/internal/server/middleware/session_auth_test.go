package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

func TestAdminAuthMiddleware_Success(t *testing.T) {
	sessionRepo := &mock.MockSessionRepository{
		Sessions: map[string]*model.Session{
			"admin-session-1": {ID: "admin-session-1", UserID: 1, Role: model.SessionRoleAdmin},
		},
	}
	mw := NewAdminAuthMiddleware(sessionRepo)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adminID, ok := AdminIDFromContext(r.Context())
		if !ok || adminID != 1 {
			t.Fatalf("expected admin id 1 in context, got %v %v", adminID, ok)
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/admin/fishermen", nil)
	req.AddCookie(&http.Cookie{Name: "admin_session", Value: "admin-session-1"})
	w := httptest.NewRecorder()

	mw.Handle(next).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestAdminAuthMiddleware_RoleMismatch(t *testing.T) {
	sessionRepo := &mock.MockSessionRepository{
		Sessions: map[string]*model.Session{
			"buyer-session-1": {ID: "buyer-session-1", UserID: 1, Role: model.SessionRoleBuyer},
		},
	}
	mw := NewAdminAuthMiddleware(sessionRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/fishermen", nil)
	req.AddCookie(&http.Cookie{Name: "admin_session", Value: "buyer-session-1"})
	w := httptest.NewRecorder()

	mw.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestBuyerAuthMiddleware_Success(t *testing.T) {
	sessionRepo := &mock.MockSessionRepository{
		Sessions: map[string]*model.Session{
			"buyer-session-1": {ID: "buyer-session-1", UserID: 7, Role: model.SessionRoleBuyer},
		},
	}
	mw := NewBuyerAuthMiddleware(sessionRepo)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buyerID, ok := BuyerIDFromContext(r.Context())
		if !ok || buyerID != 7 {
			t.Fatalf("expected buyer id 7 in context, got %v %v", buyerID, ok)
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/buyer/me", nil)
	req.AddCookie(&http.Cookie{Name: "buyer_session", Value: "buyer-session-1"})
	w := httptest.NewRecorder()

	mw.Handle(next).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}
