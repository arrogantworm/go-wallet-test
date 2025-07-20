package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	custom_errors "wallet-app/pkg/errors"
	"wallet-app/pkg/models"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type (
	UpdateWalletJSON struct {
		WalletID  uuid.UUID `json:"walletId"`
		Operation string    `json:"operationType"`
		Amount    float64   `json:"amount"`
	}
)

func (h *Handler) getWalletInfo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	walletID, err := uuid.Parse(idStr)
	if err != nil {
		h.sendError(w, fmt.Sprintf("wrong uuid: %v", err), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	balance, err := h.service.GetBalance(ctx, walletID)
	if err != nil {
		if errors.Is(err, custom_errors.ErrWalletNotFound) {
			h.sendError(w, "wallet not found", http.StatusNotFound)
		} else {
			h.sendError(w, "could not get balance", http.StatusInternalServerError)
			return
		}
	}

	res := models.Wallet{
		ID:      walletID,
		Balance: balance,
	}

	h.sendJSON(w, res, http.StatusOK)
}

func (h *Handler) updateWalletBalance(w http.ResponseWriter, r *http.Request) {

	var req UpdateWalletJSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "bad request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	switch strings.ToLower(req.Operation) {
	case "deposit":
		if err := h.service.Deposit(ctx, req.WalletID, req.Amount); err != nil {
			if errors.Is(err, custom_errors.ErrWalletNotFound) {
				if err := h.service.NewWallet(ctx, req.WalletID, req.Amount); err != nil {
					h.sendError(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				h.sendError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		h.sendSuccess(w, "balance updated", http.StatusOK)
	case "withdraw":
		if err := h.service.Withdraw(ctx, req.WalletID, req.Amount); err != nil {
			if errors.Is(err, custom_errors.ErrNotEnoughFunds) {
				h.sendError(w, err.Error(), http.StatusBadRequest)
				return
			} else {
				h.sendError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		h.sendSuccess(w, "balance updated", http.StatusOK)
	default:
		h.sendError(w, "wrong operation type", http.StatusBadRequest)
	}
}
