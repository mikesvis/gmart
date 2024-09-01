package handler

import (
	_context "context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/mikesvis/gmart/internal/context"
	"github.com/mikesvis/gmart/internal/service/accrual"
	"github.com/mikesvis/gmart/internal/service/order"
	jsonutils "github.com/mikesvis/gmart/pkg/json"
	"github.com/mikesvis/gmart/pkg/luhn"
)

func (h *Handler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := _context.WithCancel(r.Context())
	defer cancel()

	userID := ctx.Value(context.UserIDContextKey).(uint64)

	balance, err := h.accrual.GetUserBalance(ctx, userID)

	if err != nil && errors.Is(err, accrual.ErrWrongArgument) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := struct {
		Current   jsonutils.Rubles `json:"current"`
		Withdrawn jsonutils.Rubles `json:"withdrawn"`
	}{
		Current:   jsonutils.Rubles(balance.Current),
		Withdrawn: jsonutils.Rubles(balance.Withdrawn),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(response)
}

func (h *Handler) WithdrawForOrder(w http.ResponseWriter, r *http.Request) {
	// У меня есть вопросы!!! Передайте пожалуста товарщиу кто составлял задание
	// что есть такое понятие как "нормализация таблиц"!!!!
	// Я потратил несколько дней чтобы разобраться: от чегой-та когда приходит списание, оно не фиксируется у пользователя?!
	// По заданию выходит так что нужно просто зарегиcтрировать списание !внимание! за заказ!
	// При этом такой заказ может быть не создан в этом проекте!
	// Я бы с удовольствием послушал как это сделать учитывая тесты?
	// Я тогда тоже буду выкручиваться: просто у себя создавать такой заказ (который указан в списании)
	// потому что начисление за него "по идее" тоже может быть (правда там будут сложности на какую сумму начислять)
	ctx, cancel := _context.WithCancel(r.Context())
	defer cancel()

	var request struct {
		OrderID string  `json:"order"`
		Sum     float64 `json:"sum"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderID, err := strconv.ParseInt(request.OrderID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !luhn.IsValid(uint64(orderID)) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	userID := ctx.Value(context.UserIDContextKey).(uint64)

	balance, err := h.accrual.GetUserBalance(ctx, userID)
	if err != nil && errors.Is(err, accrual.ErrWrongArgument) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sum := int64(request.Sum * 100)
	if balance.Current < uint64(sum) {
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}

	// собсно вот: создаем заказ при списании, либо он уже существует
	err = h.order.CreateOrder(ctx, uint64(orderID), userID)
	if err != nil && !errors.Is(err, order.ErrOwnedByAnother) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	err = h.accrual.WithdrawToOrderID(ctx, uint64(orderID), sum*-1)

	if err != nil && errors.Is(err, accrual.ErrWrongArgument) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetUserWithdrawals(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := _context.WithCancel(r.Context())
	defer cancel()

	userID := ctx.Value(context.UserIDContextKey).(uint64)

	withdrawals, err := h.accrual.GetUserWithdrawals(ctx, userID)

	if err != nil && errors.Is(err, accrual.ErrResultIsEmpty) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err != nil && errors.Is(err, accrual.ErrWrongArgument) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type responseItem struct {
		OrderID     uint64             `json:"order,string"`
		Sum         jsonutils.Rubles   `json:"sum"`
		ProcessedAt jsonutils.JSONTime `json:"processed_at"`
	}
	var response []responseItem
	for _, v := range withdrawals {
		response = append(response, responseItem{
			OrderID:     v.OrderID,
			Sum:         jsonutils.Rubles(v.Sum),
			ProcessedAt: jsonutils.JSONTime(v.ProcessedAt),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(response)
}
