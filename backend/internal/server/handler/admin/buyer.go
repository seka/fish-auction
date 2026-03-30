package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/handler/admin/request"
	"github.com/seka/fish-auction/backend/internal/server/handler/admin/response"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
)

// BuyerHandler handles admin HTTP requests related to buyers.
type BuyerHandler struct {
	createUseCase buyer.CreateBuyerUseCase
	listUseCase   buyer.ListBuyersUseCase
	deleteUseCase buyer.DeleteBuyerUseCase
}

// NewBuyerHandler creates a new BuyerHandler instance.
func NewBuyerHandler(r registry.UseCase) *BuyerHandler {
	return &BuyerHandler{
		createUseCase: r.NewCreateBuyerUseCase(),
		listUseCase:   r.NewListBuyersUseCase(),
		deleteUseCase: r.NewDeleteBuyerUseCase(),
	}
}

// Create handles the buyer creation request.
func (h *BuyerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req request.CreateBuyer
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	buy, err := h.createUseCase.Execute(r.Context(), req.Name, req.Email, req.Password, req.Organization, req.ContactInfo)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusCreated, response.Buyer{ID: buy.ID, Name: buy.Name})
}

// List handles the request to list buyers.
func (h *BuyerHandler) List(w http.ResponseWriter, r *http.Request) {
	buyers, err := h.listUseCase.Execute(r.Context())
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := make([]response.Buyer, len(buyers))
	for i, v := range buyers {
		resp[i] = response.Buyer{
			ID:   v.ID,
			Name: v.Name,
		}
	}

	util.WriteJSON(w, http.StatusOK, resp)
}

// Delete handles the buyer deletion request.
func (h *BuyerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid buyer ID")
		return
	}

	if err := h.deleteUseCase.Execute(r.Context(), id); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes registers the admin buyer handler routes to the given mux.
func (h *BuyerHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /buyers", h.List)
	mux.HandleFunc("POST /buyers", h.Create)
	mux.HandleFunc("DELETE /buyers/{id}", h.Delete)
}
