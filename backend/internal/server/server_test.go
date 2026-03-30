package server

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	adminHandler "github.com/seka/fish-auction/backend/internal/server/handler/admin"
	buyerHandler "github.com/seka/fish-auction/backend/internal/server/handler/buyer"
	publicHandler "github.com/seka/fish-auction/backend/internal/server/handler/public"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

func TestServer_SecurityRoutes(t *testing.T) {
	// Setup Mocks (Minimal needed to avoid nil pointers, actual logic not executed for 401 checks)
	mockReg := &mock.MockRegistry{
		CreateFishermanUC:   &mock.MockCreateFishermanUseCase{ExecuteFunc: func(_ context.Context, _ string) (*model.Fisherman, error) { return &model.Fisherman{ID: 1}, nil }},
		CreateBuyerUC:       &mock.MockCreateBuyerUseCase{ExecuteFunc: func(_ context.Context, _, _, _, _, _ string) (*model.Buyer, error) { return &model.Buyer{ID: 1}, nil }},
		ListFishermenUC:     &mock.MockListFishermenUseCase{ExecuteFunc: func(_ context.Context) ([]model.Fisherman, error) { return []model.Fisherman{}, nil }},
		ListBuyersUC:        &mock.MockListBuyersUseCase{ExecuteFunc: func(_ context.Context) ([]model.Buyer, error) { return []model.Buyer{}, nil }},
		GetBuyerPurchasesUC: &mock.MockGetBuyerPurchasesUseCase{ExecuteFunc: func(_ context.Context, _ int) ([]model.Purchase, error) { return []model.Purchase{}, nil }},
		GetBuyerAuctionsUC:  &mock.MockGetBuyerAuctionsUseCase{ExecuteFunc: func(_ context.Context, _ int) ([]model.Auction, error) { return []model.Auction{}, nil }},
		CreateAuctionUC: &mock.MockCreateAuctionUseCase{ExecuteFunc: func(_ context.Context, auction *model.Auction) (*model.Auction, error) {
			return &model.Auction{ID: 1, VenueID: auction.VenueID}, nil
		}},
		UpdateAuctionStatusUC: &mock.MockUpdateAuctionStatusUseCase{ExecuteFunc: func(_ context.Context, _ int, _ model.AuctionStatus) error {
			return nil
		}},
		GetBuyerUC: &mock.MockGetBuyerUseCase{ExecuteFunc: func(_ context.Context, _ int) (*model.Buyer, error) {
			return &model.Buyer{ID: 1, Name: "Test Buyer"}, nil
		}},
	}

	// Initialize Handlers
	sessionRepo := &mock.MockSessionRepository{
		Sessions: map[string]*model.Session{
			"admin-session-1": {ID: "admin-session-1", UserID: 1, Role: model.SessionRoleAdmin},
			"buyer-session-1": {ID: "buyer-session-1", UserID: 1, Role: model.SessionRoleBuyer},
		},
	}
	hHealth := publicHandler.NewHealthHandler()
	hFisherman := adminHandler.NewFishermanHandler(mockReg)
	hBuyerAuth := publicHandler.NewBuyerAuthHandler(mockReg, sessionRepo)
	hBuyer := buyerHandler.NewBuyerHandler(mockReg)
	hAdminBuyer := adminHandler.NewBuyerHandler(mockReg)
	hPublicItem := publicHandler.NewItemHandler(mockReg)
	hAdminItem := adminHandler.NewItemHandler(mockReg)
	hBid := buyerHandler.NewBidHandler(mockReg)
	hInvoice := adminHandler.NewInvoiceHandler(mockReg)
	hAdminAuth := publicHandler.NewAdminAuthHandler(mockReg, sessionRepo)
	hPublicVenue := publicHandler.NewVenueHandler(mockReg)
	hAdminVenue := adminHandler.NewVenueHandler(mockReg)
	hPublicAuction := publicHandler.NewAuctionHandler(mockReg)
	hAdminAuction := adminHandler.NewAuctionHandler(mockReg)
	hAdmin := adminHandler.NewAdminHandler(mockReg)
	hAuthReset := publicHandler.NewAuthResetHandler(mockReg)
	hAdminAuthReset := adminHandler.NewAuthResetHandler(mockReg)
	hPush := buyerHandler.NewPushHandler(mockReg)

	// Initialize Server
	s := NewServer(
		hHealth,
		hFisherman,
		hBuyerAuth,
		hBuyer,
		hAdminBuyer,
		hPublicItem,
		hAdminItem,
		hBid,
		hInvoice,
		hAdminAuth,
		hPublicVenue,
		hAdminVenue,
		hPublicAuction,
		hAdminAuction,
		hAdmin,
		hAuthReset,
		hAdminAuthReset,
		hPush,
		sessionRepo,
		[]string{"https://localhost", "http://localhost:3000"},
		10*time.Second,
		10*time.Second,
		10*time.Second,
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
		// 1. Admin Routes Security Verification (Must be 401 without cookie)
		// --------------------------------------------------------------------
		// Fishermen
		{name: "Admin_ListFishermen_NoAuth", method: http.MethodGet, path: "/api/admin/fishermen", expectedStatus: http.StatusUnauthorized},
		{name: "Admin_CreateFisherman_NoAuth", method: http.MethodPost, path: "/api/admin/fishermen", expectedStatus: http.StatusUnauthorized},
		// Buyers
		{name: "Admin_ListBuyers_NoAuth", method: http.MethodGet, path: "/api/admin/buyers", expectedStatus: http.StatusUnauthorized},
		{name: "Admin_CreateBuyer_NoAuth", method: http.MethodPost, path: "/api/admin/buyers", expectedStatus: http.StatusUnauthorized},
		// Items
		{name: "Admin_CreateItem_NoAuth", method: http.MethodPost, path: "/api/admin/items", expectedStatus: http.StatusUnauthorized},
		// Auctions
		{name: "Admin_CreateAuction_NoAuth", method: http.MethodPost, path: "/api/admin/auctions", expectedStatus: http.StatusUnauthorized},
		{name: "Admin_UpdateAuction_NoAuth", method: http.MethodPut, path: "/api/admin/auctions/1", expectedStatus: http.StatusUnauthorized},
		{name: "Admin_UpdateAuctionStatus_NoAuth", method: http.MethodPatch, path: "/api/admin/auctions/1/status", expectedStatus: http.StatusUnauthorized},
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
			cookieValue:    "admin-session-1",
			expectedStatus: http.StatusCreated,
		},
		// Auction Status (Admin)
		{
			name:           "Admin_UpdateAuctionStatus_Authorized",
			method:         http.MethodPatch,
			path:           "/api/admin/auctions/1/status",
			cookieName:     "admin_session",
			cookieValue:    "admin-session-1",
			expectedStatus: http.StatusOK,
		},
		// List Fishermen (Admin)
		{
			name:           "Admin_ListFishermen_Authorized",
			method:         http.MethodGet,
			path:           "/api/admin/fishermen",
			cookieName:     "admin_session",
			cookieValue:    "admin-session-1",
			expectedStatus: http.StatusOK,
		},
		// List Buyers (Admin)
		{
			name:           "Admin_ListBuyers_Authorized",
			method:         http.MethodGet,
			path:           "/api/admin/buyers",
			cookieName:     "admin_session",
			cookieValue:    "admin-session-1",
			expectedStatus: http.StatusOK,
		},
		// Buyer Me (Buyer)
		{
			name:           "Buyer_GetMe_Authorized",
			method:         http.MethodGet,
			path:           "/api/buyer/me",
			cookieName:     "buyer_session",
			cookieValue:    "buyer-session-1",
			expectedStatus: http.StatusOK,
		},
		// Auction Create (Admin)
		{
			name:           "Admin_CreateAuction_Authorized",
			method:         http.MethodPost,
			path:           "/api/admin/auctions",
			cookieName:     "admin_session",
			cookieValue:    "admin-session-1",
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var body io.Reader
			// Provide dummy JSON body for POST/PUT to ensure handlers don't decode EOF
			if tc.method == http.MethodPost || tc.method == http.MethodPut || tc.method == http.MethodPatch {
				payload := `{"name": "test", "email": "test@example.com", "password": "pass", "organization": "org", "contact_info": "info", "status": "in_progress"}`
				if strings.Contains(tc.path, "auctions") {
					payload = `{"venue_id": 1, "auction_date": "2025-01-01", "start_time": "10:00:00", "end_time": "11:00:00", "status": "scheduled"}`
				}
				body = strings.NewReader(payload)
			}
			req := httptest.NewRequestWithContext(context.Background(), tc.method, tc.path, body)

			if tc.cookieName != "" {
				req.AddCookie(&http.Cookie{Name: tc.cookieName, Value: tc.cookieValue})
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
