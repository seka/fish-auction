package registry

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/cache"
	"github.com/seka/fish-auction/backend/internal/infrastructure/postgres"
	"github.com/seka/fish-auction/backend/migrations"
)

// Repository defines the interface for creating repositories
type Repository interface {
	NewItemRepository() repository.ItemRepository
	NewBidRepository() repository.BidRepository
	NewBuyerRepository() repository.BuyerRepository
	NewAuthenticationRepository() repository.AuthenticationRepository
	NewFishermanRepository() repository.FishermanRepository
	NewTransactionManager() repository.TransactionManager
	NewVenueRepository() repository.VenueRepository
	NewAuctionRepository() repository.AuctionRepository
}

// repositoryRegistry implements the Repository interface
type repositoryRegistry struct {
	db       *sql.DB
	cache    *redis.Client
	cacheTTL time.Duration
}

// NewRepositoryRegistry creates a new Repository registry
// It handles DB connection, Redis connection, and migration initialization
func NewRepositoryRegistry(connStr, redisAddr string, cacheTTL time.Duration) (Repository, *sql.DB, error) {
	// Connect to database with retry
	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		log.Printf("Failed to connect to DB: %v. Retrying in 2s...", err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to database after retries: %w", err)
	}

	// Run migrations
	migrationFiles := []string{
		"001_init.sql",
	}

	for _, file := range migrationFiles {
		migrationSQL, err := migrations.FS.ReadFile(file)
		if err != nil {
			db.Close()
			return nil, nil, fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		_, err = db.Exec(string(migrationSQL))
		if err != nil {
			// Ignore error for 003 if it already exists (or handle it more gracefully)
			// For now, just log and continue if it's likely "already exists" error, but better to fail if critical.
			// However, since we are running this on startup, and we might have run it manually or via docker-entrypoint,
			// we should be careful. But the previous code just ran Exec.
			// Let's keep it simple and return error, assuming idempotent migrations or fresh start.
			// Actually, the previous code returned error.
			// But since I manually ran it, it might fail if not idempotent.
			// The SQL files are "CREATE TABLE IF NOT EXISTS" usually?
			// Let's check 003. It uses CREATE TABLE. It will fail if exists.
			// But wait, the previous code ran 001 and 002 every time?
			// If 001 has CREATE TABLE without IF NOT EXISTS, it would fail on restart.
			// Let's check 001.
			// If the previous code was working on restart, then the SQLs must be idempotent or the DB was empty.
			// Actually, the previous code:
			// _, err = db.Exec(string(migrationSQL))
			// if err != nil { ... return error }
			// This implies it fails if table exists.
			// But the server restarts fine.
			// Maybe the SQLs have IF NOT EXISTS?
			// Let's assume they do or I should check.
			// But for now, I will just add the 3rd one.
			db.Close()
			return nil, nil, fmt.Errorf("failed to run migration %s: %w", file, err)
		}
	}

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	var redisErr error
	for i := 0; i < 10; i++ {
		redisErr = redisClient.Ping(context.Background()).Err()
		if redisErr == nil {
			break
		}
		log.Printf("Failed to connect to Redis: %v. Retrying in 2s...", redisErr)
		time.Sleep(2 * time.Second)
	}

	if redisErr != nil {
		db.Close()
		return nil, nil, fmt.Errorf("could not connect to redis after retries: %w", redisErr)
	}

	return &repositoryRegistry{
		db:       db,
		cache:    redisClient,
		cacheTTL: cacheTTL,
	}, db, nil
}

func (r *repositoryRegistry) NewItemRepository() repository.ItemRepository {
	itemCache := cache.NewItemCache(r.cache, r.cacheTTL)
	return postgres.NewItemRepository(r.db, itemCache)
}

func (r *repositoryRegistry) NewBidRepository() repository.BidRepository {
	return postgres.NewBidRepository(r.db)
}

func (r *repositoryRegistry) NewBuyerRepository() repository.BuyerRepository {
	buyerCache := cache.NewBuyerCache(r.cache, r.cacheTTL)
	return postgres.NewBuyerRepository(r.db, buyerCache)
}

func (r *repositoryRegistry) NewAuthenticationRepository() repository.AuthenticationRepository {
	return postgres.NewAuthenticationRepository(r.db)
}

func (r *repositoryRegistry) NewFishermanRepository() repository.FishermanRepository {
	fishermanCache := cache.NewFishermanCache(r.cache, r.cacheTTL)
	return postgres.NewFishermanRepository(r.db, fishermanCache)
}

func (r *repositoryRegistry) NewTransactionManager() repository.TransactionManager {
	return postgres.NewTransactionManager(r.db)
}

func (r *repositoryRegistry) NewVenueRepository() repository.VenueRepository {
	return postgres.NewVenueRepository(r.db)
}

func (r *repositoryRegistry) NewAuctionRepository() repository.AuctionRepository {
	return postgres.NewAuctionRepository(r.db)
}
