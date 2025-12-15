package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
)

type BuyerHandler struct {
	createUseCase         buyer.CreateBuyerUseCase
	listUseCase           buyer.ListBuyersUseCase
	loginUseCase          buyer.LoginBuyerUseCase
	getPurchasesUseCase   buyer.GetBuyerPurchasesUseCase
	getAuctionsUseCase    buyer.GetBuyerAuctionsUseCase
	updatePasswordUseCase buyer.UpdatePasswordUseCase
}

func NewBuyerHandler(r registry.UseCase) *BuyerHandler {
	return &BuyerHandler{
		createUseCase:         r.NewCreateBuyerUseCase(),
		listUseCase:           r.NewListBuyersUseCase(),
		loginUseCase:          r.NewLoginBuyerUseCase(),
		getPurchasesUseCase:   r.NewGetBuyerPurchasesUseCase(),
		getAuctionsUseCase:    r.NewGetBuyerAuctionsUseCase(),
		updatePasswordUseCase: r.NewBuyerUpdatePasswordUseCase(),
	}
}

func (h *BuyerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBuyerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}
	buyer, err := h.createUseCase.Execute(r.Context(), req.Name, req.Email, req.Password, req.Organization, req.ContactInfo)
	if err != nil {
		util.HandleError(w, err)
		return
	}
	resp := dto.BuyerResponse{ID: buyer.ID, Name: buyer.Name}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *BuyerHandler) List(w http.ResponseWriter, r *http.Request) {
	buyers, err := h.listUseCase.Execute(r.Context())
	if err != nil {
		util.HandleError(w, err)
		return
	}
	resp := make([]dto.BuyerResponse, len(buyers))
	for i, b := range buyers {
		resp[i] = dto.BuyerResponse{ID: b.ID, Name: b.Name}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *BuyerHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginBuyerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	buyer, err := h.loginUseCase.Execute(r.Context(), req.Email, req.Password)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	// Set session cookie
	// In a real app, generate a secure token and store it in Redis/DB
	// For now, we'll use a simple signed token or just the ID (insecure but simple for MVP)
	// Let's use "buyer_session" cookie with value "authenticated" and maybe another for ID?
	// Or better, use a JWT or signed cookie.
	// Given the previous "admin_session" implementation, let's stick to simple cookie for now.
	// But we need the BuyerID for bidding.
	// So let's set "buyer_id" cookie as well (HttpOnly).

	http.SetCookie(w, &http.Cookie{
		Name:     "buyer_session",
		Value:    "authenticated",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteLaxMode,
	})

	// Store Buyer ID in a separate cookie or encoded in the session
	// For simplicity, let's use a separate cookie "buyer_id"
	http.SetCookie(w, &http.Cookie{
		Name:     "buyer_id",
		Value:    fmt.Sprintf("%d", buyer.ID),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	resp := dto.BuyerResponse{ID: buyer.ID, Name: buyer.Name}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *BuyerHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "buyer_session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "buyer_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out"})
}

func (h *BuyerHandler) GetCurrentBuyer(w http.ResponseWriter, r *http.Request) {
	// Check for buyer_session cookie
	cookie, err := r.Cookie("buyer_session")
	if err != nil || cookie.Value != "authenticated" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get buyer ID from cookie
	idCookie, err := r.Cookie("buyer_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": true,
		"buyer_id":      idCookie.Value,
	})
}

func (h *BuyerHandler) GetMyPurchases(w http.ResponseWriter, r *http.Request) {
	// Check for buyer_session cookie
	cookie, err := r.Cookie("buyer_session")
	if err != nil || cookie.Value != "authenticated" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get buyer ID from cookie
	idCookie, err := r.Cookie("buyer_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Convert buyer_id to int
	var buyerID int
	if _, err := fmt.Sscanf(idCookie.Value, "%d", &buyerID); err != nil {
		http.Error(w, "Invalid buyer ID", http.StatusBadRequest)
		return
	}

	purchases, err := h.getPurchasesUseCase.Execute(r.Context(), buyerID)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	// Convert to DTOs
	resp := make([]dto.PurchaseResponse, len(purchases))
	for i, p := range purchases {
		resp[i] = dto.PurchaseResponse{
			ID:          p.ID,
			ItemID:      p.ItemID,
			FishType:    p.FishType,
			Quantity:    p.Quantity,
			Unit:        p.Unit,
			Price:       p.Price,
			BuyerID:     p.BuyerID,
			AuctionID:   p.AuctionID,
			AuctionDate: p.AuctionDate,
			CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *BuyerHandler) GetMyAuctions(w http.ResponseWriter, r *http.Request) {
	// Check for buyer_session cookie
	cookie, err := r.Cookie("buyer_session")
	if err != nil || cookie.Value != "authenticated" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get buyer ID from cookie
	idCookie, err := r.Cookie("buyer_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Convert buyer_id to int
	var buyerID int
	if _, err := fmt.Sscanf(idCookie.Value, "%d", &buyerID); err != nil {
		http.Error(w, "Invalid buyer ID", http.StatusBadRequest)
		return
	}

	auctions, err := h.getAuctionsUseCase.Execute(r.Context(), buyerID)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	// Convert to DTOs
	resp := make([]dto.AuctionResponse, len(auctions))
	for i, a := range auctions {
		var startTime, endTime *string
		if a.StartTime != nil {
			s := a.StartTime.Format("15:04:05")
			startTime = &s
		}
		if a.EndTime != nil {
			e := a.EndTime.Format("15:04:05")
			endTime = &e
		}

		resp[i] = dto.AuctionResponse{
			ID:          a.ID,
			VenueID:     a.VenueID,
			AuctionDate: a.AuctionDate.Format("2006-01-02"),
			StartTime:   startTime,
			EndTime:     endTime,
			Status:      string(a.Status),
			CreatedAt:   a.CreatedAt,
			UpdatedAt:   a.UpdatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *BuyerHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	cookie, err := r.Cookie("buyer_session")
	if err != nil || cookie.Value != "authenticated" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	idCookie, err := r.Cookie("buyer_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var buyerID int
	if _, err := fmt.Sscanf(idCookie.Value, "%d", &buyerID); err != nil {
		http.Error(w, "Invalid buyer ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	if err := h.updatePasswordUseCase.Execute(r.Context(), buyerID, req.CurrentPassword, req.NewPassword); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully"})
}

func (h *BuyerHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/buyers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.Create(w, r)
		case http.MethodGet:
			h.List(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/buyers/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.Login(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/buyers/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.Logout(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/buyers/me", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetCurrentBuyer(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/buyers/me/purchases", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetMyPurchases(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/buyers/me/auctions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetMyAuctions(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/buyers/password", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			h.UpdatePassword(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
