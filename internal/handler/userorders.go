package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	_context "context"

	"github.com/mikesvis/gmart/internal/context"
	"github.com/mikesvis/gmart/internal/domain"
	"github.com/mikesvis/gmart/internal/service/order"
	jsonutils "github.com/mikesvis/gmart/pkg/json"
)

func (h *Handler) CreateUserOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := _context.WithCancel(r.Context())
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	orderID, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	userID := ctx.Value(context.UserIDContextKey).(uint64)

	err = h.order.CreateOrder(ctx, uint64(orderID), userID)

	if err != nil && errors.Is(err, order.ErrBadRequest) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil && errors.Is(err, order.ErrConflict) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	if err != nil && errors.Is(err, order.ErrUnprocessableEntity) {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err != nil && errors.Is(err, order.ErrExisted) {
		w.WriteHeader(http.StatusOK)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	h.accrual.PushToAccural(uint64(orderID))

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := _context.WithCancel(r.Context())
	defer cancel()

	userID := ctx.Value(context.UserIDContextKey).(uint64)

	orders, err := h.order.GetOrdersByUser(ctx, userID)

	if err != nil && errors.Is(err, order.ErrNoContent) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err != nil && errors.Is(err, order.ErrBadRequest) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	type responseItem struct {
		ID        uint64             `json:"number,string"`
		Status    domain.Status      `json:"status"`
		Accrual   jsonutils.Rubles   `json:"accrual,omitempty"`
		CreatedAt jsonutils.JSONTime `json:"uploaded_at"`
	}

	var result []responseItem

	for _, v := range orders {
		result = append(result, responseItem{
			ID:        v.ID,
			Status:    v.Status,
			Accrual:   jsonutils.Rubles(v.Accrual),
			CreatedAt: jsonutils.JSONTime(v.CreatedAt),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(result)
}
