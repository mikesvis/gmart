package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mikesvis/gmart/internal/service/user"
)

type UserRegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *Handler) UserRegister(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var req UserRegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	userID, err := h.user.RegisterUser(ctx, req.Login, req.Password)

	if err != nil && errors.Is(err, user.ErrBadRequest) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil && errors.Is(err, user.ErrConflict) {
		http.Error(w, err.Error(), http.StatusConflict)
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
