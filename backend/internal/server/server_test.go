package server

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/server/handler"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

func TestServer_SecurityRoutes(t *testing.T) {
	// Setup Mocks (Minimal needed to avoid nil pointers, actual logic not executed for 401 checks)
	mockReg := &mock.MockRegistry{
		CreateFishermanUC:   &mock.MockCreateFishermanUseCase{ExecuteFunc: func(ctx context.Context, name string) (*model.Fisherman, error) { return &model.Fisherman{ID: 1}, nil }},
		CreateBuyerUC:       &mock.MockCreateBuyerUseCase{ExecuteFunc: func(ctx context.Context, n, e, p, o, c string) (*model.Buyer, error) { return &model.Buyer{ID: 1}, nil }},
		ListFishermenUC:     &mock.MockListFishermenUseCase{ExecuteFunc: func(ctx context.Context) ([]model.Fisherman, error) { return []model.Fisherman{}, nil }},
		ListBuyersUC:        &mock.MockListBuyersUseCase{ExecuteFunc: func(ctx context.Context) ([]model.Buyer, error) { return []model.Buyer{}, nil }},
		GetBuyerPurchasesUC: &mock.MockGetBuyerPurchasesUseCase{ExecuteFunc: func(ctx context.Context, id int) ([]model.Purchase, error) { return []model.Purchase{}, nil }},
		GetBuyerAuctionsUC:  &mock.MockGetBuyerAuctionsUseCase{ExecuteFunc: func(ctx context.Context, id int) ([]model.Auction, error) { return []model.Auction{}, nil }},
		// Add mocks for other handlers if needed for the 200 OK checks, but for 401 checks they won't be called.
	}

	// Initialize Handlers
	hHealth := handler.NewHealthHandler()
	hFisherman := handler.NewFishermanHandler(mockReg)
	hBuyer := handler.NewBuyerHandler(mockReg)
	hItem := handler.NewItemHandler(mockReg)
	hBid := handler.NewBidHandler(mockReg)
	hInvoice := handler.NewInvoiceHandler(mockReg)
	hAuth := handler.NewAuthHandler(mockReg)
	hVenue := handler.NewVenueHandler(mockReg)
	hAuction := handler.NewAuctionHandler(mockReg)
	hAdmin := handler.NewAdminHandler(mockReg)
	hAuthReset := handler.NewAuthResetHandler(mockReg)
	hAdminAuthReset := handler.NewAdminAuthResetHandler(mockReg)

	// Initialize Server
	s := NewServer(
		hHealth,
		hFisherman,
		hBuyer,
		hItem,
		hBid,
		hInvoice,
		hAuth,
		hVenue,
		hAuction,
		hAdmin,
		hAuthReset,
		hAdminAuthReset,
	)

	tests := []struct {
		name           string
		method         string
		path           string
		cookieName     string
		cookieValue    string
		expectedStatus int
	}{
		// --------------------------------------------------------------------
		// 1. Public Routes Verification (Legacy/Public Insecure POST blocked)
		// --------------------------------------------------------------------
		{name: "Public_CreateFisherman_Blocked", method: http.MethodPost, path: "/api/fishermen", expectedStatus: http.StatusMethodNotAllowed},
		{name: "Public_CreateBuyer_Blocked", method: http.MethodPost, path: "/api/buyers", expectedStatus: http.StatusMethodNotAllowed},
		{name: "Public_ListFishermen_Allowed", method: http.MethodGet, path: "/api/fishermen", expectedStatus: http.StatusOK},
		{name: "Public_ListBuyers_Allowed", method: http.MethodGet, path: "/api/buyers", expectedStatus: http.StatusOK},

		// --------------------------------------------------------------------
		// 2. Admin Routes Security Verification (Must be 401 without cookie)
		// --------------------------------------------------------------------
		// Fishermen
		{name: "Admin_CreateFisherman_NoAuth", method: http.MethodPost, path: "/api/admin/fishermen", expectedStatus: http.StatusUnauthorized},
		// Buyers
		{name: "Admin_CreateBuyer_NoAuth", method: http.MethodPost, path: "/api/admin/buyers", expectedStatus: http.StatusUnauthorized},
		// Items
		{name: "Admin_CreateItem_NoAuth", method: http.MethodPost, path: "/api/admin/items", expectedStatus: http.StatusUnauthorized},
		// Auctions
		{name: "Admin_CreateAuction_NoAuth", method: http.MethodPost, path: "/api/admin/auctions", expectedStatus: http.StatusUnauthorized},
		{name: "Admin_UpdateAuction_NoAuth", method: http.MethodPut, path: "/api/admin/auctions/1", expectedStatus: http.StatusUnauthorized},
		{name: "Admin_DeleteAuction_NoAuth", method: http.MethodDelete, path: "/api/admin/auctions/1", expectedStatus: http.StatusUnauthorized},
		// Venues
		{name: "Admin_CreateVenue_NoAuth", method: http.MethodPost, path: "/api/admin/venues", expectedStatus: http.StatusUnauthorized},
		{name: "Admin_UpdateVenue_NoAuth", method: http.MethodPut, path: "/api/admin/venues/1", expectedStatus: http.StatusUnauthorized},
		{name: "Admin_DeleteVenue_NoAuth", method: http.MethodDelete, path: "/api/admin/venues/1", expectedStatus: http.StatusUnauthorized},
		// Password
		{name: "Admin_UpdatePassword_NoAuth", method: http.MethodPut, path: "/api/admin/password", expectedStatus: http.StatusUnauthorized},

		// --------------------------------------------------------------------
		// 3. Buyer Routes Security Verification (Must be 401 without cookie)
		// --------------------------------------------------------------------
		// My Page
		{name: "Buyer_GetMe_NoAuth", method: http.MethodGet, path: "/api/buyer/me", expectedStatus: http.StatusUnauthorized},
		{name: "Buyer_GetPurchases_NoAuth", method: http.MethodGet, path: "/api/buyer/me/purchases", expectedStatus: http.StatusUnauthorized},
		{name: "Buyer_GetAuctions_NoAuth", method: http.MethodGet, path: "/api/buyer/me/auctions", expectedStatus: http.StatusUnauthorized},
		// Bids
		{name: "Buyer_CreateBid_NoAuth", method: http.MethodPost, path: "/api/buyer/bids", expectedStatus: http.StatusUnauthorized},
		// Password
		{name: "Buyer_UpdatePassword_NoAuth", method: http.MethodPut, path: "/api/buyer/password", expectedStatus: http.StatusUnauthorized},

		// --------------------------------------------------------------------
		// 4. Authorized Access Verification (Sample check with cookie)
		// --------------------------------------------------------------------
		// Fishermen (Admin)
		{
			name:           "Admin_CreateFisherman_Authorized",
			method:         http.MethodPost,
			path:           "/api/admin/fishermen",
			cookieName:     "admin_session",
			cookieValue:    "authenticated",
			expectedStatus: http.StatusOK,
		},
		// Buyer Me (Buyer)
		{
			name:           "Buyer_GetMe_Authorized",
			method:         http.MethodGet,
			path:           "/api/buyer/me",
			cookieName:     "buyer_session",
			cookieValue:    "authenticated",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var body io.Reader
			// Provide dummy JSON body for POST/PUT to ensure handlers don't decode EOF
			if tc.method == http.MethodPost || tc.method == http.MethodPut {
				body = strings.NewReader(`{"name": "test", "email": "test@example.com", "password": "pass", "organization": "org", "contact_info": "info"}`)
			}
			req := httptest.NewRequest(tc.method, tc.path, body)

			if tc.cookieName != "" {
				req.AddCookie(&http.Cookie{Name: tc.cookieName, Value: tc.cookieValue})
				// Add ID cookies required by middleware/handlers
				req.AddCookie(&http.Cookie{Name: "admin_id", Value: "1"})
				req.AddCookie(&http.Cookie{Name: "buyer_id", Value: "1"})
			}
			w := httptest.NewRecorder()

			// Accessing unexported router internal to package server
			s.router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}
}
