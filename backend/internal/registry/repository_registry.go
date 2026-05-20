package registry

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres"
	cacheStore "github.com/seka/fish-auction/backend/internal/infrastructure/datastore/redis"
	"github.com/seka/fish-auction/backend/internal/migration"
)

// Repository defines the interface for creating repositories and managing lifecycle.
type Repository interface {
	NewItemRepository() repository.ItemRepository
	NewBidRepository() repository.BidRepository
	NewBuyerRepository() repository.BuyerRepository
	NewAuthenticationRepository() repository.AuthenticationRepository
	NewFishermanRepository() repository.FishermanRepository
	NewTransactionManager() repository.TransactionManager
	NewVenueRepository() repository.VenueRepository
	NewAuctionRepository() repository.AuctionRepository
	NewAdminRepository() repository.AdminRepository
	NewPushRepository() repository.PushRepository
	PasswordReset() repository.PasswordResetRepository
	NewItemCacheInvalidator() repository.CacheInvalidator
	NewSessionRepository() repository.SessionRepository
	NewOutboxRepository() repository.OutboxRepository
	NewRateLimitRepository() repository.RateLimitRepository
	// Cleanup closes underlying connections (DB, Redis, etc.) via their interfaces.
	Cleanup() error
}

// repositoryRegistry implements the Repository interface
type repositoryRegistry struct {
	db         datastore.Database
	cache      datastore.Cache
	cacheTTL   time.Duration
	sessionTTL time.Duration
}

// NewRepositoryRegistry creates a new Repository registry.
// Migrations are no longer run here; invoke cmd/migration before starting
// processes that depend on schema state.
func NewRepositoryRegistry(
	dbCfg config.DatabaseConfig,
	redisCfg config.RedisConfig,
	cacheCfg config.CacheConfig,
	sessionCfg config.SessionConfig,
) (Repository, error) {
	db, err := migration.Connect(context.Background(), dbCfg.DBConnectionURL())
	if err != nil {
		return nil, err
	}

	// RedisAddr が空のときは Redis 接続をスキップする（relay のような Redis 不要プロセス向け）。
	var cache datastore.Cache
	var redisClient *redis.Client
	if redisCfg.RedisAddr() != "" {
		var err error
		redisClient, err = connectRedis(redisCfg.RedisAddr(), redisCfg.GetRedisDB())
		if err != nil {
			_ = db.Close()
			return nil, err
		}
		cache = cacheStore.NewClient(redisClient)
	}

	return &repositoryRegistry{
		db:         postgres.NewClient(db),
		cache:      cache,
		cacheTTL:   cacheCfg.GetCacheTTL(),
		sessionTTL: sessionCfg.GetSessionTTL(),
	}, nil
}

func (r *repositoryRegistry) NewRateLimitRepository() repository.RateLimitRepository {
	return postgres.NewRateLimitStore(r.db)
}

func (r *repositoryRegistry) Cleanup() error {
	var errs []string
	if err := r.db.Close(); err != nil {
		errs = append(errs, fmt.Sprintf("database close error: %v", err))
	}
	if r.cache != nil {
		if err := r.cache.Close(); err != nil {
			errs = append(errs, fmt.Sprintf("cache close error: %v", err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("cleanup errors: %s", strings.Join(errs, "; "))
	}
	return nil
}

func connectRedis(redisAddr string, db int) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   db,
	})

	var redisErr error
	for range 10 {
		redisErr = redisClient.Ping(context.Background()).Err()
		if redisErr == nil {
			return redisClient, nil
		}
		slog.Warn("failed to connect to Redis; retrying", "err", redisErr, "interval", "2s")
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

func (r *repositoryRegistry) NewSessionRepository() repository.SessionRepository {
	return cacheStore.NewSessionStore(r.cache, r.sessionTTL)
}

func (r *repositoryRegistry) NewOutboxRepository() repository.OutboxRepository {
	return postgres.NewOutboxStore(r.db)
}
