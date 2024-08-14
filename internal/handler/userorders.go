package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	_context "context"

	"github.com/mikesvis/gmart/internal/context"
	"github.com/mikesvis/gmart/internal/service/order"
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

	w.WriteHeader(http.StatusAccepted)

	// push to process accural channel
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

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(orders)
}
