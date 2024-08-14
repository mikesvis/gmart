package handler

import (
	_context "context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mikesvis/gmart/internal/context"
	"github.com/mikesvis/gmart/internal/service/accural"
)

func (h *Handler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := _context.WithCancel(r.Context())
	defer cancel()

	userID := ctx.Value(context.UserIDContextKey).(uint64)

	balance, err := h.accural.GetUserBalance(ctx, userID)

	if err != nil && errors.Is(err, accural.ErrBadRequest) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := struct {
		Current   uint64 `json:"current"`
		Withdrawn uint64 `json:"withdrawn"`
	}{
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(response)
}
