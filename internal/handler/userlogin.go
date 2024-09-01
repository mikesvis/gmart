package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mikesvis/gmart/internal/service/user"
)

type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var req UserLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	userID, err := h.user.GetUserID(ctx, req.Login, req.Password)

	if err != nil && errors.Is(err, user.ErrWrongArgument) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil && errors.Is(err, user.ErrUnauthorizedUser) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.user.Login(ctx, w, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
