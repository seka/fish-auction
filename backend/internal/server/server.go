package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/server/handler/admin"
	"github.com/seka/fish-auction/backend/internal/server/handler/buyer"
	"github.com/seka/fish-auction/backend/internal/server/handler/public"
	"github.com/seka/fish-auction/backend/internal/server/middleware"
)

const (
	defaultShutdownTimeout = 30 * time.Second
)

// Server serves the request.
type Server struct {
	router                *http.ServeMux
	httpServer            *http.Server
	healthHandler         *public.HealthHandler
	fishermanHandler      *admin.FishermanHandler
	buyerAuthHandler      *public.BuyerAuthHandler
	buyerHandler          *buyer.BuyerHandler
	adminBuyerHandler     *admin.BuyerHandler
	publicItemHandler     *public.ItemHandler
	adminItemHandler      *admin.ItemHandler
	bidHandler            *buyer.BidHandler
	invoiceHandler        *admin.InvoiceHandler
	adminAuthHandler      *public.AdminAuthHandler
	publicVenueHandler    *public.VenueHandler
	adminVenueHandler     *admin.VenueHandler
	publicAuctionHandler  *public.AuctionHandler
	adminAuctionHandler   *admin.AuctionHandler
	adminHandler          *admin.AdminHandler
	adminAuthResetHandler *admin.AuthResetHandler
	authResetHandler      *public.AuthResetHandler
	pushHandler           *buyer.PushHandler
	adminLoginRL          *middleware.RateLimiterMiddleware
	buyerLoginRL          *middleware.RateLimiterMiddleware
	adminResetRL          *middleware.RateLimiterMiddleware
	buyerResetRL          *middleware.RateLimiterMiddleware
	adminMe               *admin.MeHandler
	adminAuth             *middleware.AdminAuthMiddleware
	buyerAuth             *middleware.BuyerAuthMiddleware
	cors                  *middleware.CORSMiddleware
	securityHeaders       *middleware.SecurityHeadersMiddleware
	csrf                  *middleware.CSRFMiddleware
	cacheControl          *middleware.CacheControlMiddleware
	gzip                  *middleware.GzipMiddleware
	maxBody               *middleware.MaxBodyMiddleware
	recovery              *middleware.RecoveryMiddleware
	trustedProxy          *middleware.TrustedProxyMiddleware
	requestID             *middleware.RequestIDMiddleware
	readTimeout           time.Duration
	writeTimeout          time.Duration
	idleTimeout           time.Duration
}

// NewServer creates a new Server instance.
func NewServer(
	healthHandler *public.HealthHandler,
	fishermanHandler *admin.FishermanHandler,
	buyerAuthHandler *public.BuyerAuthHandler,
	buyerHandler *buyer.BuyerHandler,
	adminBuyerHandler *admin.BuyerHandler,
	publicItemHandler *public.ItemHandler,
	adminItemHandler *admin.ItemHandler,
	bidHandler *buyer.BidHandler,
	invoiceHandler *admin.InvoiceHandler,
	adminAuthHandler *public.AdminAuthHandler,
	publicVenueHandler *public.VenueHandler,
	adminVenueHandler *admin.VenueHandler,
	publicAuctionHandler *public.AuctionHandler,
	adminAuctionHandler *admin.AuctionHandler,
	adminHandler *admin.AdminHandler,
	authResetHandler *public.AuthResetHandler,
	adminAuthResetHandler *admin.AuthResetHandler,
	pushHandler *buyer.PushHandler,
	adminMeHandler *admin.MeHandler,
	sessionRepo repository.SessionRepository,
	rateLimitRepo repository.RateLimitRepository,
	allowedOrigins []string,
	trustedProxies []string,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	idleTimeout time.Duration,
) *Server {
	s := &Server{
		router:                http.NewServeMux(),
		healthHandler:         healthHandler,
		fishermanHandler:      fishermanHandler,
		buyerAuthHandler:      buyerAuthHandler,
		buyerHandler:          buyerHandler,
		adminBuyerHandler:     adminBuyerHandler,
		publicItemHandler:     publicItemHandler,
		adminItemHandler:      adminItemHandler,
		bidHandler:            bidHandler,
		invoiceHandler:        invoiceHandler,
		adminAuthHandler:      adminAuthHandler,
		publicVenueHandler:    publicVenueHandler,
		adminVenueHandler:     adminVenueHandler,
		publicAuctionHandler:  publicAuctionHandler,
		adminAuctionHandler:   adminAuctionHandler,
		adminHandler:          adminHandler,
		authResetHandler:      authResetHandler,
		adminAuthResetHandler: adminAuthResetHandler,
		pushHandler:           pushHandler,
		adminMe:               adminMeHandler,
		adminLoginRL:          middleware.NewRateLimiterMiddleware(rateLimitRepo.IncrementAdminLogin, middleware.LoginRateLimit, middleware.LoginRateWindow),
		buyerLoginRL:          middleware.NewRateLimiterMiddleware(rateLimitRepo.IncrementBuyerLogin, middleware.LoginRateLimit, middleware.LoginRateWindow),
		adminResetRL:          middleware.NewRateLimiterMiddleware(rateLimitRepo.IncrementAdminReset, middleware.ResetRateLimit, middleware.ResetRateWindow),
		buyerResetRL:          middleware.NewRateLimiterMiddleware(rateLimitRepo.IncrementBuyerReset, middleware.ResetRateLimit, middleware.ResetRateWindow),
		adminAuth:             middleware.NewAdminAuthMiddleware(sessionRepo),
		buyerAuth:             middleware.NewBuyerAuthMiddleware(sessionRepo),
		cors:                  middleware.NewCORSMiddleware(allowedOrigins),
		securityHeaders:       middleware.NewSecurityHeadersMiddleware(),
		csrf:                  middleware.NewCSRFMiddleware(allowedOrigins),
		cacheControl:          middleware.NewCacheControlMiddleware(),
		gzip:                  middleware.NewGzipMiddleware(),
		maxBody:               middleware.NewMaxBodyMiddleware(1024 * 1024), // 1MB
		recovery:              middleware.NewRecoveryMiddleware(),
		trustedProxy:          middleware.NewTrustedProxyMiddleware(trustedProxies),
		requestID:             middleware.NewRequestIDMiddleware(),
		readTimeout:           readTimeout,
		writeTimeout:          writeTimeout,
		idleTimeout:           idleTimeout,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.registerPublicRoutes()
	s.registerAdminRoutes()
	s.registerBuyerRoutes()
}

func (s *Server) registerPublicRoutes() {
	s.healthHandler.RegisterRoutes(s.router)

	// admin login は RL でラップし、logout はそのまま登録する。
	s.router.Handle("POST /api/login", s.adminLoginRL.Handle(http.HandlerFunc(s.adminAuthHandler.Login)))
	s.router.HandleFunc("POST /api/admin/logout", s.adminAuthHandler.Logout)

	// buyer login は RL でラップし、logout はそのまま登録する。
	s.router.Handle("POST /api/buyer/login", s.buyerLoginRL.Handle(http.HandlerFunc(s.buyerAuthHandler.Login)))
	s.router.HandleFunc("POST /api/buyer/logout", s.buyerAuthHandler.Logout)

	// パスワードリセット request エンドポイントのみ RL 対象。verify / confirm は RL なし。
	s.router.Handle("POST /api/auth/password-reset/request", s.buyerResetRL.Handle(http.HandlerFunc(s.authResetHandler.RequestReset)))
	s.router.HandleFunc("POST /api/auth/password-reset/verify", s.authResetHandler.VerifyToken)
	s.router.HandleFunc("POST /api/auth/password-reset/confirm", s.authResetHandler.ConfirmReset)

	s.router.Handle("POST /api/admin/password-reset/request", s.adminResetRL.Handle(http.HandlerFunc(s.adminAuthResetHandler.RequestReset)))
	s.router.HandleFunc("POST /api/admin/password-reset/verify", s.adminAuthResetHandler.VerifyToken)
	s.router.HandleFunc("POST /api/admin/password-reset/confirm", s.adminAuthResetHandler.ConfirmReset)

	s.publicItemHandler.RegisterRoutes(s.router)
	s.publicAuctionHandler.RegisterRoutes(s.router)
	s.publicVenueHandler.RegisterRoutes(s.router)
}

func (s *Server) registerAdminRoutes() {
	adminMux := http.NewServeMux()

	s.adminItemHandler.RegisterRoutes(adminMux)
	s.fishermanHandler.RegisterRoutes(adminMux)
	s.adminBuyerHandler.RegisterRoutes(adminMux)
	s.adminAuctionHandler.RegisterRoutes(adminMux)
	s.adminVenueHandler.RegisterRoutes(adminMux)
	s.adminHandler.RegisterRoutes(adminMux)
	s.invoiceHandler.RegisterRoutes(adminMux)
	s.adminMe.RegisterRoutes(adminMux)

	s.router.Handle("/api/admin/", s.adminAuth.Handle(http.StripPrefix("/api/admin", adminMux)))
}

func (s *Server) registerBuyerRoutes() {
	buyerMux := http.NewServeMux()

	s.bidHandler.RegisterRoutes(buyerMux)
	s.buyerHandler.RegisterRoutes(buyerMux)
	s.pushHandler.RegisterRoutes(buyerMux)

	s.router.Handle("/api/buyer/", s.buyerAuth.Handle(http.StripPrefix("/api/buyer", buyerMux)))
}

// Start starts the HTTP server and blocks until the context is canceled.
func (s *Server) Start(ctx context.Context, addr string) error {
	if addr == "" {
		addr = ":8080"
	}

	handlerWithCORS := s.cors.Handle(s.router)
	handlerWithHeaders := s.securityHeaders.Handle(handlerWithCORS)
	handlerWithCSRF := s.csrf.Handle(handlerWithHeaders)
	handlerWithCache := s.cacheControl.Handle(handlerWithCSRF)
	handlerWithGzip := s.gzip.Handle(handlerWithCache)
	handlerWithLimit := s.maxBody.Handle(handlerWithGzip)
	// requestID は trustedProxy 直後に配置し、ALB が付与する X-Request-ID を尊重しつつ
	// 後段ハンドラ／ミドルウェアが context から request_id を参照できるようにする。
	handlerWithRequestID := s.requestID.Handle(handlerWithLimit)
	// trustedProxy はリカバリ直下に配置し、後続のハンドラ全てが
	// 書き換え後の RemoteAddr / Scheme を見られるようにする。
	handlerWithProxy := s.trustedProxy.Handle(handlerWithRequestID)
	handlerWithRecovery := s.recovery.Handle(handlerWithProxy)

	s.httpServer = &http.Server{
		Addr:              addr,
		Handler:           handlerWithRecovery,
		ReadTimeout:       s.readTimeout,
		ReadHeaderTimeout: 20 * time.Second,
		WriteTimeout:      s.writeTimeout,
		IdleTimeout:       s.idleTimeout,
	}

	// listen failure と signal の両方を main goroutine で受け取り、
	// どちらの経路でも defer 経由のクリーンアップが走るようにする。
	listenErrCh := make(chan error, 1)
	go func() {
		slog.Info("server starting", "addr", addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			listenErrCh <- err
		}
		close(listenErrCh)
	}()

	select {
	case err, ok := <-listenErrCh:
		if ok && err != nil {
			return fmt.Errorf("server listen failed: %w", err)
		}
	case <-ctx.Done():
	}

	slog.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	slog.Info("server exiting")
	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}
