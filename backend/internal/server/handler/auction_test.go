package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/handler"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

func TestAuctionHandler_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCreateUC := &mock.MockCreateAuctionUseCase{
			ExecuteFunc: func(ctx context.Context, auction *model.Auction) (*model.Auction, error) {
				auction.ID = 1
				auction.CreatedAt = time.Now()
				auction.UpdatedAt = time.Now()
				return auction, nil
			},
		}
		mockReg := &mock.MockRegistry{CreateAuctionUC: mockCreateUC}
		h := handler.NewAuctionHandler(mockReg)

		reqBody := dto.CreateAuctionRequest{
			VenueID:     1,
			AuctionDate: "2023-01-01",
			Status:      "Scheduled",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/auctions", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.Create(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestAuctionHandler_List(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockListUC := &mock.MockListAuctionsUseCase{
			ExecuteFunc: func(ctx context.Context, filters *repository.AuctionFilters) ([]model.Auction, error) {
				now := time.Now()
				return []model.Auction{{ID: 1, AuctionDate: now, CreatedAt: now, UpdatedAt: now}}, nil
			},
		}
		mockReg := &mock.MockRegistry{ListAuctionsUC: mockListUC}
		h := handler.NewAuctionHandler(mockReg)

		req := httptest.NewRequest(http.MethodGet, "/api/auctions", nil)
		w := httptest.NewRecorder()

		h.List(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("UseCaseError", func(t *testing.T) {
		mockListUC := &mock.MockListAuctionsUseCase{
			ExecuteFunc: func(ctx context.Context, filters *repository.AuctionFilters) ([]model.Auction, error) {
				return nil, errors.New("db error")
			},
		}
		mockReg := &mock.MockRegistry{ListAuctionsUC: mockListUC}
		h := handler.NewAuctionHandler(mockReg)

		req := httptest.NewRequest(http.MethodGet, "/api/auctions", nil)
		w := httptest.NewRecorder()

		h.List(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}

func TestAuctionHandler_Get(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockGetUC := &mock.MockGetAuctionUseCase{
			ExecuteFunc: func(ctx context.Context, id int) (*model.Auction, error) {
				now := time.Now()
				return &model.Auction{ID: 1, AuctionDate: now, CreatedAt: now, UpdatedAt: now}, nil
			},
		}
		mockReg := &mock.MockRegistry{GetAuctionUC: mockGetUC}
		h := handler.NewAuctionHandler(mockReg)

		req := httptest.NewRequest(http.MethodGet, "/api/auctions/1", nil)
		w := httptest.NewRecorder()

		h.Get(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}
