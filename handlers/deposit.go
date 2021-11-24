package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justteddy/wallet/types"
	"github.com/pkg/errors"
)

type depositRequest struct {
	Amount int `json:"amount"`
}

func (h *Handler) HandleDeposit(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	walletID := params.ByName("wallet")
	if walletID == "" {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("empty wallet id"))
		return
	}

	var depositReq depositRequest
	if err := json.NewDecoder(r.Body).Decode(&depositReq); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, errors.Wrap(err, "decode request"))
		return
	}

	if err := h.s.Deposit(r.Context(), types.WalletID(walletID), depositReq.Amount); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "save to storage"))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
