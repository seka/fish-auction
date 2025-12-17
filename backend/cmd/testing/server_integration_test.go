package testing

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server"
	"github.com/seka/fish-auction/backend/internal/server/handler"
	"golang.org/x/crypto/bcrypt"
)

func TestServerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// 1. テスト設定を読み込む
	cfg := config.LoadTest()

	// 2. テスト用 DB 名を生成
	testDBName := fmt.Sprintf("test_fish_auction_%d", time.Now().Unix())

	// 3. 管理用 DB に接続
	adminDB, err := sql.Open("postgres", cfg.AdminConnStr())
	if err != nil {
		t.Fatalf("Failed to connect to admin database: %v", err)
	}
	defer adminDB.Close()

	// 4. テスト用 DB を作成
	if err := createTestDatabase(adminDB, testDBName); err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer func() {
		if err := dropTestDatabase(adminDB, testDBName); err != nil {
			t.Errorf("Failed to drop test database: %v", err)
		}
	}()

	// 5. Registry を初期化（DB 接続、Redis 接続、マイグレーション）
	// テストでは Redis をローカルホストに接続
	repoReg, db, err := registry.NewRepositoryRegistry(
		cfg.TestDBConnStr(testDBName),
		"localhost:6379",
		5*time.Minute,
	)
	if err != nil {
		t.Fatalf("Failed to initialize registry: %v", err)
	}
	defer db.Close()

	appCfg := &config.Config{
		DBHost:     cfg.DBHost,
		DBPort:     cfg.DBPort,
		DBUser:     cfg.DBUser,
		DBPassword: cfg.DBPassword,
		DBName:     testDBName,
		RedisAddr:  "localhost:6379",
		CacheTTL:   5 * time.Minute,
		AppEnv:     "test",
		SMTPHost:   "localhost",
		SMTPPort:   "1025",
		SMTPFrom:   "test@example.com",
	}

	useCaseReg := registry.NewUseCaseRegistry(repoReg, appCfg)

	// 6. Handlers を初期化
	healthHandler := handler.NewHealthHandler()
	fishermanHandler := handler.NewFishermanHandler(useCaseReg)
	buyerHandler := handler.NewBuyerHandler(useCaseReg)
	itemHandler := handler.NewItemHandler(useCaseReg)
	bidHandler := handler.NewBidHandler(useCaseReg)
	invoiceHandler := handler.NewInvoiceHandler(useCaseReg)
	authHandler := handler.NewAuthHandler(useCaseReg)
	venueHandler := handler.NewVenueHandler(useCaseReg)
	auctionHandler := handler.NewAuctionHandler(useCaseReg)
	adminHandler := handler.NewAdminHandler(useCaseReg)
	authResetHandler := handler.NewAuthResetHandler(useCaseReg)
	adminAuthResetHandler := handler.NewAdminAuthResetHandler(useCaseReg)

	// 7. Server を起動
	srv := server.NewServer(
		healthHandler,
		fishermanHandler,
		buyerHandler,
		itemHandler,
		bidHandler,
		invoiceHandler,
		authHandler,
		venueHandler,
		auctionHandler,
		adminHandler,
		authResetHandler,
		adminAuthResetHandler,
	)

	// 8. サーバーを goroutine で起動
	serverAddr := ":18080" // テスト用ポート
	errChan := make(chan error, 1)
	go func() {
		if err := srv.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// 9. サーバーの準備完了を待機
	serverURL := "http://localhost:18080"
	if err := waitForServer(serverURL + "/api/health"); err != nil {
		t.Fatalf("Server failed to start: %v", err)
	}

	// 10. Health エンドポイントをテスト
	t.Run("Health", func(t *testing.T) {
		resp, err := http.Get(serverURL + "/api/health")
		if err != nil {
			t.Fatalf("Failed to call health endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	// 11. Full Auction Flow Test
	t.Run("FullAuctionFlow", func(t *testing.T) {
		jar, _ := cookiejar.New(nil)
		client := &http.Client{
			Jar: jar,
		}

		// 1. Register Fisherman
		fishermanID := registerUser(t, client, serverURL+"/api/fishermen", `{"name": "Captain Ahab"}`)

		// 2. Register Buyer
		// We use buyerID for nothing currently, but we use buyer account for login
		_ = registerUser(t, client, serverURL+"/api/buyers", `{"name": "Ishmael", "email": "ishmael@example.com", "password": "password123", "organization": "Pequod", "contact_info": "sea"}`)

		// 3. Seed Admin (Direct DB)
		// Assuming "admin" table/model exists and we can insert.
		seedAdmin(t, db, "admin@example.com", "admin123")

		// 4. Login Admin
		adminCookies := login(t, client, serverURL+"/api/login", `{"email": "admin@example.com", "password": "admin123"}`)

		// 5. Login Buyer
		buyerCookies := login(t, client, serverURL+"/api/buyers/login", `{"email": "ishmael@example.com", "password": "password123"}`)

		// 6. Create Venue (as Admin)
		// POST /api/admin/venues
		venueID := createResource(t, client, serverURL+"/api/admin/venues", `{"name": "Nantucket Harbor", "location": "Nantucket", "description": "Main harbor"}`, adminCookies)

		// 7. Create Auction (as Admin)
		// POST /api/admin/auctions
		// Links to Venue.
		// Make auction active now so we can bid.
		auctionDate := time.Now().Format("2006-01-02")
		// StartTime 00:00, EndTime 23:59
		auctionID := createResource(t, client, serverURL+"/api/admin/auctions", fmt.Sprintf(`{"venue_id": %d, "auction_date": "%s", "start_time": "00:00:00", "end_time": "23:59:59", "status": "in_progress"}`, venueID, auctionDate), adminCookies)

		// 8. Create Item (as Admin)
		// POST /api/admin/items
		// Note: Item creation needs FishermanID and AuctionID
		itemID := createResource(t, client, serverURL+"/api/admin/items", fmt.Sprintf(`{"auction_id": %d, "fisherman_id": %d, "fish_type": "Whale", "quantity": 1, "unit": "whole"}`, auctionID, fishermanID), adminCookies)

		// Create another item for bidding test not needed if we use the first one, but let's keep logic simple
		// reusing itemID for bidding if possible, or create duplicate?
		// The previous code created itemID2. Let's just use itemID for bidding.
		// Wait, previous code used itemID2.
		// Let's just use the one item we created.

		// 9. List Auctions (Public GET /api/auctions)
		// Verify valid JSON response
		listResources(t, client, serverURL+"/api/auctions")

		// 10. Place Bid (as Buyer)
		// POST /api/buyer/bids
		bidBody := fmt.Sprintf(`{"item_id": %d, "price": 5000}`, itemID)
		postResource(t, client, serverURL+"/api/buyer/bids", bidBody, buyerCookies)

		// 11. Verify Bid via Auction Details or Item Details
		verifyBid(t, client, serverURL+fmt.Sprintf("/api/auctions/%d/items", auctionID), itemID, 5000)
	})

	// 11. サーバーをシャットダウン
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		t.Errorf("Failed to shutdown server: %v", err)
	}
}

// createTestDatabase はテスト用 DB を作成
func createTestDatabase(db *sql.DB, dbName string) error {
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	return err
}

// dropTestDatabase はテスト用 DB を削除
func dropTestDatabase(db *sql.DB, dbName string) error {
	// アクティブな接続を切断
	_, _ = db.Exec(fmt.Sprintf(`
		SELECT pg_terminate_backend(pg_stat_activity.pid)
		FROM pg_stat_activity
		WHERE pg_stat_activity.datname = '%s'
		AND pid <> pg_backend_pid()
	`, dbName))

	_, err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	return err
}

// waitForServer はサーバーが起動するまで待機
func waitForServer(url string) error {
	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for server to start")
		case <-ticker.C:
			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					return nil
				}
			}
		}
	}
}

// Helpers

// Helper functions (registerUser, login, etc) follow...

func registerUser(t *testing.T, client *http.Client, urlStr, jsonBody string) int {
	resp, err := client.Post(urlStr, "application/json", strings.NewReader(jsonBody))
	if err != nil {
		t.Fatalf("Failed to register: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200/201, got %d: %s", resp.StatusCode, string(body))
	}
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
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
	resp, err := client.Post(urlStr, "application/json", strings.NewReader(jsonBody))
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, string(body))
	}
	return resp.Cookies()
}

func createResource(t *testing.T, client *http.Client, urlStr, jsonBody string, cookies []*http.Cookie) int {
	req, _ := http.NewRequest("POST", urlStr, strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to create resource at %s: %v", urlStr, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200/201 at %s, got %d: %s", urlStr, resp.StatusCode, string(body))
	}

	var res map[string]interface{}
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
	req, _ := http.NewRequest("POST", urlStr, strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to POST to %s: %v", urlStr, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200/201 at %s, got %d: %s", urlStr, resp.StatusCode, string(body))
	}
}

func listResources(t *testing.T, client *http.Client, urlStr string) {
	resp, err := client.Get(urlStr)
	if err != nil {
		t.Fatalf("Failed to list resources at %s: %v", urlStr, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}
}

func verifyBid(t *testing.T, client *http.Client, urlStr string, itemID, expectedPrice int) {
	resp, err := client.Get(urlStr)
	if err != nil {
		t.Fatalf("Failed to get details at %s: %v", urlStr, err)
	}
	defer resp.Body.Close()

	var items []map[string]interface{}
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

func seedAdmin(t *testing.T, db *sql.DB, email, password string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	// Correct column is password_hash
	_, err = db.Exec("INSERT INTO admins (email, password_hash, created_at) VALUES ($1, $2, NOW())", email, string(hash))
	if err != nil {
		t.Fatalf("Failed to seed admin: %v", err)
	}
}

// HttpCookieClient wrapper for standard client to simplify logic if needed
// Actually, standard http.Client with cookiejar is enough.
type CookieClient = http.Client // Alias for simplicity in signatures
