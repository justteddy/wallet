package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justteddy/wallet/types"
	"github.com/pkg/errors"
)

type transferRequest struct {
	FromWallet types.WalletID `json:"from_wallet"`
	ToWallet   types.WalletID `json:"to_wallet"`
	Amount     int            `json:"amount"`
}

func (h *Handler) HandleTransfer(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var transferReq transferRequest
	if err := json.NewDecoder(r.Body).Decode(&transferReq); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, errors.Wrap(err, "decode request"))
		return
	}

	if transferReq.FromWallet == "" {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("from_wallet is required"))
		return
	}

	if transferReq.ToWallet == "" {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("to_wallet is required"))
		return
	}

	if transferReq.Amount <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("invalid amount"))
		return
	}

	if err := h.s.Transfer(r.Context(), transferReq.FromWallet, transferReq.ToWallet, transferReq.Amount); err != nil {
		if errors.Cause(err) == types.ErrUnavailableBalance {
			writeErrorResponse(w, http.StatusBadRequest, types.ErrUnavailableBalance)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "save to storage"))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
