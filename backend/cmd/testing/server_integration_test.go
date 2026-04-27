package testing

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/service"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server"
	adminHandler "github.com/seka/fish-auction/backend/internal/server/handler/admin"
	buyerHandler "github.com/seka/fish-auction/backend/internal/server/handler/buyer"
	publicHandler "github.com/seka/fish-auction/backend/internal/server/handler/public"
)

func TestServerIntegration(t *testing.T) {
	requireIntegrationTests(t)

	// Define a root context for the entire test
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// 1. 設定を読み込む
	cfg, err := config.LoadAppServerConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 2. Registry を初期化（DB 接続、Redis 接続、マイグレーション）
	repoReg, err := registry.NewRepositoryRegistry(cfg, cfg, cfg, cfg, true)
	if err != nil {
		t.Fatalf("Failed to initialize registry: %v", err)
	}
	defer func() { _ = repoReg.Cleanup() }()

	// 3. Service を初期化
	realServiceReg, err := registry.NewServiceRegistry(cfg, config.NoWebpushConfig, cfg)
	if err != nil {
		t.Fatalf("Failed to initialize service registry: %v", err)
	}

	// 統合テスト用にプッシュ通知サービスをモックに差し替え
	mockPush := &mockPushService{}
	serviceReg := &wrappedServiceRegistry{
		Service:  realServiceReg,
		mockPush: mockPush,
	}

	useCaseReg := registry.NewUseCaseRegistry(repoReg, serviceReg, cfg)

	// 4. Worker を初期化して起動
	workerReg, err := registry.NewWorkerRegistry(cfg, repoReg, serviceReg)
	if err != nil {
		t.Fatalf("Failed to initialize worker registry: %v", err)
	}
	w, err := workerReg.NewWorker()
	if err != nil {
		t.Fatalf("Failed to create worker: %v", err)
	}

	go func() {
		if err := w.Start(ctx); err != nil && ctx.Err() == nil {
			t.Errorf("Worker failed: %v", err)
		}
	}()

	// 3. Handlers を初期化
	healthHandler := publicHandler.NewHealthHandler()
	fishermanHandler := adminHandler.NewFishermanHandler(useCaseReg)
	sessionRepo := repoReg.NewSessionRepository()
	buyerAuthHandler := publicHandler.NewBuyerAuthHandler(useCaseReg, sessionRepo)
	buyerAccountHandler := buyerHandler.NewBuyerHandler(useCaseReg)
	adminBuyerHandler := adminHandler.NewBuyerHandler(useCaseReg)
	publicItemHandler := publicHandler.NewItemHandler(useCaseReg)
	adminItemHandler := adminHandler.NewItemHandler(useCaseReg)
	bidHandler := buyerHandler.NewBidHandler(useCaseReg)
	invoiceHandler := adminHandler.NewInvoiceHandler(useCaseReg)
	adminAuthHandler := publicHandler.NewAdminAuthHandler(useCaseReg, sessionRepo)
	publicVenueHandler := publicHandler.NewVenueHandler(useCaseReg)
	adminVenueHandler := adminHandler.NewVenueHandler(useCaseReg)
	publicAuctionHandler := publicHandler.NewAuctionHandler(useCaseReg)
	adminAuctionHandler := adminHandler.NewAuctionHandler(useCaseReg)
	adminAccountHandler := adminHandler.NewAdminHandler(useCaseReg)
	authResetHandler := publicHandler.NewAuthResetHandler(useCaseReg)
	adminAuthResetHandler := adminHandler.NewAuthResetHandler(useCaseReg)
	pushHandler := buyerHandler.NewPushHandler(useCaseReg)

	// 4. Server を起動
	srv := server.NewServer(
		healthHandler,
		fishermanHandler,
		buyerAuthHandler,
		buyerAccountHandler,
		adminBuyerHandler,
		publicItemHandler,
		adminItemHandler,
		bidHandler,
		invoiceHandler,
		adminAuthHandler,
		publicVenueHandler,
		adminVenueHandler,
		publicAuctionHandler,
		adminAuctionHandler,
		adminAccountHandler,
		authResetHandler,
		adminAuthResetHandler,
		pushHandler,
		sessionRepo,
		[]string{"https://localhost", "http://localhost:3000"},
		time.Minute,
		time.Minute,
		time.Minute,
	)

	// 5. サーバーを goroutine で起動
	errChan := make(chan error, 1)
	go func() {
		if err := srv.Start(ctx, cfg.ServerAddr()); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// 6. サーバーの準備完了を待機
	serverURL := "http://" + cfg.ServerAddr()
	if err := waitForServer(serverURL + "/api/health"); err != nil {
		t.Fatalf("Server failed to start: %v", err)
	}

	// 7. Health エンドポイントをテスト
	t.Run("Health", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+"/api/health", http.NoBody)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to call health endpoint: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	// 8. Full Auction Flow Test
	t.Run("FullAuctionFlow", func(t *testing.T) {
		jar, _ := cookiejar.New(nil)
		client := &http.Client{
			Jar: jar,
		}

		// 1. Seed Admin (using UseCase)
		seedAdmin(t, useCaseReg, "admin@example.com", "Admin-Password123")

		// 2. Login Admin
		adminCookies := login(t, client, serverURL+"/api/login", `{"email": "admin@example.com", "password": "Admin-Password123"}`)

		// 3. Register Fisherman (using Admin URL)
		fishermanID := registerUser(t, client, serverURL+"/api/admin/fishermen", `{"name": "Captain Ahab"}`, adminCookies)

		// 4. Register Buyer (using Admin URL)
		_ = registerUser(t, client, serverURL+"/api/admin/buyers", `{"name": "Ishmael", "email": "ishmael@example.com", "password": "Password123", "organization": "Pequod", "contact_info": "sea"}`, adminCookies)

		// 5. Login Buyer
		_ = login(t, client, serverURL+"/api/buyer/login", `{"email": "ishmael@example.com", "password": "Password123"}`)

		// 6. Create Venue (as Admin)
		// POST /api/admin/venues
		venueID := createResource(t, client, serverURL+"/api/admin/venues", `{"name": "Nantucket Harbor", "location": "Nantucket", "description": "Main harbor"}`, adminCookies)

		// 7. Create Auction (as Admin)
		// POST /api/admin/auctions
		// Links to Venue.
		// StartAt 00:00, EndAt 23:59 (JST)
		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		now := time.Now().In(jst)
		startAt := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jst).Format(time.RFC3339)
		endAt := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, jst).Format(time.RFC3339)
		auctionID := createResource(t, client, serverURL+"/api/admin/auctions", fmt.Sprintf(`{"venue_id": %d, "start_at": %q, "end_at": %q, "status": "in_progress"}`, venueID, startAt, endAt), adminCookies)

		// 8. Create Item (as Admin)
		// POST /api/admin/items
		// Note: Item creation needs FishermanID and AuctionID
		itemID := createResource(t, client, serverURL+"/api/admin/items", fmt.Sprintf(`{"auction_id": %d, "fisherman_id": %d, "fish_type": "Whale", "quantity": 1, "unit": "whole"}`, auctionID, fishermanID), adminCookies)
		// Update Item status to Available (Approval step)
		putResource(t, client, serverURL+fmt.Sprintf("/api/admin/items/%d", itemID), fmt.Sprintf(`{"auction_id": %d, "fisherman_id": %d, "fish_type": "Whale", "quantity": 1, "unit": "whole", "status": "Available"}`, auctionID, fishermanID), adminCookies)

		// Create another item for bidding test not needed if we use the first one, but let's keep logic simple
		// reusing itemID for bidding if possible, or create duplicate?
		// The previous code created itemID2. Let's just use itemID for bidding.
		// Wait, previous code used itemID2.
		// Let's just use the one item we created.

		// 9. List Auctions (Public GET /api/auctions)
		// 9. Public Listing (Items)
		listResources(t, client, serverURL+"/api/auctions")

		// 10. Place Bid (as Buyer)
		// POST /api/buyer/bids
		buyerCookies := login(t, client, serverURL+"/api/buyer/login", `{"email": "ishmael@example.com", "password": "Password123"}`)
		bidBody := fmt.Sprintf(`{"item_id": %d, "price": 5000}`, itemID)
		postResource(t, client, serverURL+"/api/buyer/bids", bidBody, buyerCookies)

		// 11. Verify Bid via Auction Details or Item Details
		verifyBid(t, client, serverURL+fmt.Sprintf("/api/auctions/%d/items", auctionID), itemID, 5000)
	})

	// 9. Asynchronous Notification Flow Test
	t.Run("AsynchronousNotificationFlow", func(t *testing.T) {
		jarA, _ := cookiejar.New(nil)
		clientA := &http.Client{Jar: jarA}
		jarB, _ := cookiejar.New(nil)
		clientB := &http.Client{Jar: jarB}

		// 1. Login Admin first (required to register buyers)
		adminCookies := login(t, clientA, serverURL+"/api/login", `{"email": "admin@example.com", "password": "Admin-Password123"}`)

		// 2. Create and Login Buyer A (Previous Bidder)
		_ = registerUser(t, clientA, serverURL+"/api/admin/buyers", `{"name": "Buyer A", "email": "buyera@example.com", "password": "Password123", "organization": "Org A", "contact_info": "email"}`, adminCookies)
		cookiesA := login(t, clientA, serverURL+"/api/buyer/login", `{"email": "buyera@example.com", "password": "Password123"}`)

		// 3. Subscribe Buyer A to Push Notifications
		subscribePush(t, clientA, serverURL+"/api/buyer/push/subscribe", `{"endpoint": "https://fcm.googleapis.com/fcm/send/fake-token", "keys": {"p256dh": "fake-p256dh", "auth": "fake-auth"}}`, cookiesA)

		// 4. Create and Login Buyer B (New Bidder)
		_ = registerUser(t, clientB, serverURL+"/api/admin/buyers", `{"name": "Buyer B", "email": "buyerb@example.com", "password": "Password123", "organization": "Org B", "contact_info": "email"}`, adminCookies)
		cookiesB := login(t, clientB, serverURL+"/api/buyer/login", `{"email": "buyerb@example.com", "password": "Password123"}`)

		// 5. Setup Item for Bidding
		fishermanID := registerUser(t, clientA, serverURL+"/api/admin/fishermen", `{"name": "Notification Fisherman"}`, adminCookies)
		venueID := createResource(t, clientA, serverURL+"/api/admin/venues", `{"name": "Notification Venue"}`, adminCookies)
		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		now := time.Now().In(jst)
		startAt := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jst).Format(time.RFC3339)
		endAt := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, jst).Format(time.RFC3339)
		auctionID := createResource(t, clientA, serverURL+"/api/admin/auctions", fmt.Sprintf(`{"venue_id": %d, "start_at": %q, "end_at": %q, "status": "in_progress"}`, venueID, startAt, endAt), adminCookies)
		itemID := createResource(t, clientA, serverURL+"/api/admin/items", fmt.Sprintf(`{"auction_id": %d, "fisherman_id": %d, "fish_type": "Tuna", "quantity": 10, "unit": "kg"}`, auctionID, fishermanID), adminCookies)
		putResource(t, clientA, serverURL+fmt.Sprintf("/api/admin/items/%d", itemID), fmt.Sprintf(`{"auction_id": %d, "fisherman_id": %d, "fish_type": "Tuna", "quantity": 10, "unit": "kg", "status": "Available"}`, auctionID, fishermanID), adminCookies)

		// 5. Buyer A Places Initial Bid (10,000)
		postResource(t, clientA, serverURL+"/api/buyer/bids", fmt.Sprintf(`{"item_id": %d, "price": 10000}`, itemID), cookiesA)

		// 6. Buyer B Places Higher Bid (15,000)
		// This triggers an outbid notification for Buyer A.
		postResource(t, clientB, serverURL+"/api/buyer/bids", fmt.Sprintf(`{"item_id": %d, "price": 15000}`, itemID), cookiesB)

		// 7. Verify Notification via Worker (Async)
		t.Log("Waiting for worker to process outbid notification...")
		found := false
		for i := 0; i < 30; i++ { // Try for 15 seconds (30 * 500ms)
			calls := mockPush.getCalls()
			for _, call := range calls {
				// Buyer A (ID 2 or similar) should have received a notification
				// We checking by email or just knowing it's Buyer A.
				// In this test, Buyer A is the first buyer we created in this t.Run.
				// Their ID depends on DB state, but we know it's not Buyer B.
				// Let's just check if ANY notification reached Buyer A.
				if call.payload != nil {
					payloadMap, ok := call.payload.(map[string]any)
					if ok {
						title, ok := payloadMap["title"].(string)
						if ok && strings.Contains(title, "高値更新") {
							found = true
							break
						}
					}
				}
			}
			if found {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}

		if !found {
			t.Error("Timed out waiting for outbid notification to be sent by worker")
		} else {
			t.Log("Successfully verified asynchronous outbid notification!")
		}
	})

	// 11. サーバーをシャットダウン
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		t.Errorf("Failed to shutdown server: %v", err)
	}
}

// waitForServer はサーバーが起動するまで待機
func waitForServer(targetURL string) error {
	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Use background context for waiting
	ctx := context.Background()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for server to start")
		case <-ticker.C:
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, http.NoBody)
			if err != nil {
				return fmt.Errorf("failed to create request: %w", err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err == nil {
				_ = resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					return nil
				}
			}
		}
	}
}

// Helpers

// Helper functions (registerUser, login, etc) follow...

func registerUser(t *testing.T, client *http.Client, urlStr, jsonBody string, cookies []*http.Cookie) int {
	req, err := http.NewRequestWithContext(context.Background(), "POST", urlStr, strings.NewReader(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for _, c := range cookies {
		req.AddCookie(c)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to register: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200/201, got %d: %s", resp.StatusCode, string(body))
	}
	var res map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&res)
	id, ok := res["id"].(float64)
	if !ok {
		// Try to parse from response if ID is not top level or different format
		// Fisherman/Buyer response: {ID: int, Name: string}
		// JSON numbers are float64
		t.Fatal("Could not parse ID from response")
	}
	return int(id)
}

func login(t *testing.T, client *http.Client, urlStr, jsonBody string) []*http.Cookie {
	req, err := http.NewRequestWithContext(context.Background(), "POST", urlStr, strings.NewReader(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, string(body))
	}
	return resp.Cookies()
}

func createResource(t *testing.T, client *http.Client, urlStr, jsonBody string, cookies []*http.Cookie) int {
	req, _ := http.NewRequestWithContext(context.Background(), "POST", urlStr, strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to create resource at %s: %v", urlStr, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200/201 at %s, got %d: %s", urlStr, resp.StatusCode, string(body))
	}

	var res map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	id, ok := res["id"].(float64)
	if !ok {
		t.Fatal("Could not parse ID from response")
	}
	return int(id)
}

func postResource(t *testing.T, client *http.Client, urlStr, jsonBody string, cookies []*http.Cookie) {
	req, _ := http.NewRequestWithContext(context.Background(), "POST", urlStr, strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to POST to %s: %v", urlStr, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200/201 at %s, got %d: %s", urlStr, resp.StatusCode, string(body))
	}
}

func subscribePush(t *testing.T, client *http.Client, urlStr, jsonBody string, cookies []*http.Cookie) {
	req, _ := http.NewRequestWithContext(context.Background(), "POST", urlStr, strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to subscribe at %s: %v", urlStr, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200 at %s, got %d: %s", urlStr, resp.StatusCode, string(body))
	}
}

func putResource(t *testing.T, client *http.Client, urlStr, jsonBody string, cookies []*http.Cookie) {
	req, _ := http.NewRequestWithContext(context.Background(), "PUT", urlStr, strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to PUT to %s: %v", urlStr, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200/204 at %s, got %d: %s", urlStr, resp.StatusCode, string(body))
	}
}

func listResources(t *testing.T, client *http.Client, urlStr string) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", urlStr, http.NoBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to list resources at %s: %v", urlStr, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}
}

func verifyBid(t *testing.T, client *http.Client, urlStr string, itemID, expectedPrice int) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", urlStr, http.NoBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to get details at %s: %v", urlStr, err)
	}
	defer func() { _ = resp.Body.Close() }()

	var items []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		t.Fatalf("Failed to decode items: %v", err)
	}

	found := false
	for _, item := range items {
		id, _ := item["id"].(float64)
		if int(id) == itemID {
			found = true
			bid, ok := item["highest_bid"].(float64)
			if !ok {
				t.Errorf("Item %d has no highest_bid", itemID)
			} else if int(bid) != expectedPrice {
				t.Errorf("Expected highest bid %d, got %d", expectedPrice, int(bid))
			}
		}
	}
	if !found {
		t.Errorf("Item %d not found in auction items", itemID)
	}
}

func seedAdmin(t *testing.T, useCaseReg registry.UseCase, email, password string) {
	_, err := useCaseReg.NewCreateAdminUseCase().Execute(context.Background(), email, password)
	if err != nil {
		t.Fatalf("Failed to seed admin: %v", err)
	}
}

// HttpCookieClient wrapper for standard client to simplify logic if needed
// Actually, standard http.Client with cookiejar is enough.
type CookieClient = http.Client // Alias for simplicity in signatures

// Mocks for Worker Integration Testing

type mockPushService struct {
	sentCalls []pushCall
	mu        sync.Mutex
}

type pushCall struct {
	buyerID int
	payload any
}

func (m *mockPushService) Send(_ context.Context, sub *model.PushSubscription, payload any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sentCalls = append(m.sentCalls, pushCall{
		buyerID: sub.BuyerID,
		payload: payload,
	})
	return nil
}

func (m *mockPushService) getCalls() []pushCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	calls := make([]pushCall, len(m.sentCalls))
	copy(calls, m.sentCalls)
	return calls
}

type wrappedServiceRegistry struct {
	registry.Service
	mockPush service.PushNotificationService
}

func (w *wrappedServiceRegistry) NewPushNotificationService() service.PushNotificationService {
	return w.mockPush
}
