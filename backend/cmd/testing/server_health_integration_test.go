package testing

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server"
	"github.com/seka/fish-auction/backend/internal/server/handler"
)

func TestServerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// 1. テスト用 DB 名を生成
	testDBName := fmt.Sprintf("test_fish_auction_%d", time.Now().Unix())

	// 2. 管理用 DB に接続
	adminConnStr := getAdminConnStr()
	adminDB, err := sql.Open("postgres", adminConnStr)
	if err != nil {
		t.Fatalf("Failed to connect to admin database: %v", err)
	}
	defer adminDB.Close()

	// 3. テスト用 DB を作成
	if err := createTestDatabase(adminDB, testDBName); err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer func() {
		if err := dropTestDatabase(adminDB, testDBName); err != nil {
			t.Errorf("Failed to drop test database: %v", err)
		}
	}()

	// 4. テスト用 DB の接続文字列を構築
	testConnStr := buildTestConnStr(testDBName)

	// 5. Registry を初期化（DB 接続とマイグレーション）
	repoReg, db, err := registry.NewRepositoryRegistry(testConnStr)
	if err != nil {
		t.Fatalf("Failed to initialize registry: %v", err)
	}
	defer db.Close()

	useCaseReg := registry.NewUseCaseRegistry(repoReg)

	// 6. Handlers を初期化
	healthHandler := handler.NewHealthHandler()
	fishermanHandler := handler.NewFishermanHandler(useCaseReg)
	buyerHandler := handler.NewBuyerHandler(useCaseReg)
	itemHandler := handler.NewItemHandler(useCaseReg)
	bidHandler := handler.NewBidHandler(useCaseReg)
	invoiceHandler := handler.NewInvoiceHandler(useCaseReg)
	authHandler := handler.NewAuthHandler(useCaseReg)

	// 7. Server を起動
	srv := server.NewServer(
		healthHandler,
		fishermanHandler,
		buyerHandler,
		itemHandler,
		bidHandler,
		invoiceHandler,
		authHandler,
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
	t.Run("HealthEndpoint", func(t *testing.T) {
		resp, err := http.Get(serverURL + "/api/health")
		if err != nil {
			t.Fatalf("Failed to call health endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	// 11. サーバーをシャットダウン
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		t.Errorf("Failed to shutdown server: %v", err)
	}
}

// getAdminConnStr は管理用 DB の接続文字列を取得
func getAdminConnStr() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "postgres")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

// buildTestConnStr はテスト用 DB の接続文字列を構築
func buildTestConnStr(testDBName string) string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, testDBName)
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

// getEnv は環境変数を取得（デフォルト値あり）
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
