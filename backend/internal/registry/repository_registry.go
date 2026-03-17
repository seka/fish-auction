package registry

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres"
	cacheStore "github.com/seka/fish-auction/backend/internal/infrastructure/datastore/redis"
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
	NewAdminRepository() repository.AdminRepository // ... other repositories
	NewPushRepository() repository.PushRepository
	PasswordReset() repository.PasswordResetRepository
	NewItemCacheInvalidator() repository.CacheInvalidator
}

// repositoryRegistry implements the Repository interface
type repositoryRegistry struct {
	db       datastore.Database
	cache    datastore.Cache
	cacheTTL time.Duration
}

// NewRepositoryRegistry creates a new Repository registry
// It handles DB connection, Redis connection, and migration initialization
func NewRepositoryRegistry(cfg *config.Config) (Repository, *sql.DB, error) {
	db, err := connectDB(cfg.DBConnectionURL())
	if err != nil {
		return nil, nil, err
	}

	if err := runMigrations(db); err != nil {
		_ = db.Close()
		return nil, nil, err
	}

	redisClient, err := connectRedis(cfg.RedisAddr)
	if err != nil {
		_ = db.Close()
		return nil, nil, err
	}

	return &repositoryRegistry{
		db:       postgres.NewClient(db),
		cache:    cacheStore.NewClient(redisClient),
		cacheTTL: cfg.CacheTTL,
	}, db, nil
}

func connectDB(postgresAddr string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for range 10 {
		db, err = sql.Open("postgres", postgresAddr)
		if err == nil {
			err = db.PingContext(context.Background())
		}
		if err == nil {
			return db, nil
		}
		log.Printf("Failed to connect to DB: %v. Retrying in 2s...", err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to database after retries: %w", err)
}

func runMigrations(db *sql.DB) error {
	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			migrationFiles = append(migrationFiles, entry.Name())
		}
	}
	sort.Strings(migrationFiles)

	for _, file := range migrationFiles {
		log.Printf("Applying migration: %s", file)
		migrationSQL, err := migrations.FS.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		_, err = db.ExecContext(context.Background(), string(migrationSQL))
		if err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", file, err)
		}
	}
	return nil
}

func connectRedis(redisAddr string) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	var redisErr error
	for range 10 {
		redisErr = redisClient.Ping(context.Background()).Err()
		if redisErr == nil {
			return redisClient, nil
		}
		log.Printf("Failed to connect to Redis: %v. Retrying in 2s...", redisErr)
		time.Sleep(2 * time.Second)
	}

	_ = redisClient.Close()
	return nil, fmt.Errorf("could not connect to redis after retries: %w", redisErr)
}

func (r *repositoryRegistry) newItemCompositeStore() *datastore.ItemCompositeStore {
	repo := postgres.NewItemStore(r.db)
	cache := cacheStore.NewItemStore(r.cache, r.cacheTTL)
	return datastore.NewItemCompositeStore(repo, cache)
}

func (r *repositoryRegistry) NewItemRepository() repository.ItemRepository {
	return r.newItemCompositeStore()
}

func (r *repositoryRegistry) NewItemCacheInvalidator() repository.CacheInvalidator {
	return r.newItemCompositeStore()
}

func (r *repositoryRegistry) NewBidRepository() repository.BidRepository {
	return postgres.NewBidStore(r.db)
}

func (r *repositoryRegistry) NewBuyerRepository() repository.BuyerRepository {
	repo := postgres.NewBuyerStore(r.db)
	cache := cacheStore.NewBuyerStore(r.cache, r.cacheTTL)
	return datastore.NewBuyerCompositeStore(repo, cache)
}

func (r *repositoryRegistry) NewAuthenticationRepository() repository.AuthenticationRepository {
	return postgres.NewAuthenticationStore(r.db)
}

func (r *repositoryRegistry) NewFishermanRepository() repository.FishermanRepository {
	repo := postgres.NewFishermanStore(r.db)
	cache := cacheStore.NewFishermanStore(r.cache, r.cacheTTL)
	return datastore.NewFishermanCompositeStore(repo, cache)
}

func (r *repositoryRegistry) NewTransactionManager() repository.TransactionManager {
	return r.db.TransactionManager()
}

func (r *repositoryRegistry) NewVenueRepository() repository.VenueRepository {
	return postgres.NewVenueStore(r.db)
}

func (r *repositoryRegistry) NewAuctionRepository() repository.AuctionRepository {
	return postgres.NewAuctionStore(r.db)
}

func (r *repositoryRegistry) NewAdminRepository() repository.AdminRepository {
	return postgres.NewAdminStore(r.db)
}

func (r *repositoryRegistry) NewPushRepository() repository.PushRepository {
	return postgres.NewPushStore(r.db)
}

func (r *repositoryRegistry) PasswordReset() repository.PasswordResetRepository {
	return postgres.NewPasswordResetStore(r.db)
}
