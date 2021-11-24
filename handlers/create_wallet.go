package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justteddy/wallet/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type createWalletResponse struct {
	WalletID types.WalletID `json:"wallet_id"`
}

func (h *Handler) HandleCreateWallet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	walletID, err := h.wg.Generate()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "generate wallet id"))
		return
	}

	if err := h.s.CreateWallet(r.Context(), walletID); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "create wallet"))
		return
	}

	resp, err := json.Marshal(&createWalletResponse{WalletID: walletID})
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "marshal response"))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		log.WithError(err).Error("failed to write successful response")
	}
}
