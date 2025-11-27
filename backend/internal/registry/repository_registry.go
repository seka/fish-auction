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
	NewFishermanRepository() repository.FishermanRepository
	NewTransactionManager() repository.TransactionManager
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
	migrationSQL, err := migrations.FS.ReadFile("001_init.sql")
	if err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to read migration file: %w", err)
	}

	_, err = db.Exec(string(migrationSQL))
	if err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to connect to redis: %w", err)
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

func (r *repositoryRegistry) NewFishermanRepository() repository.FishermanRepository {
	fishermanCache := cache.NewFishermanCache(r.cache, r.cacheTTL)
	return postgres.NewFishermanRepository(r.db, fishermanCache)
}

func (r *repositoryRegistry) NewTransactionManager() repository.TransactionManager {
	return postgres.NewTransactionManager(r.db)
}
