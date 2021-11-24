package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/justteddy/wallet/types"
	log "github.com/sirupsen/logrus"
)

type walletGenerator interface {
	Generate() (types.WalletID, error)
}

type storage interface {
	CreateWallet(ctx context.Context, wallet types.WalletID) error
	Deposit(ctx context.Context, wallet types.WalletID, amount int) error
	Transfer(ctx context.Context, fromWallet, toWallet types.WalletID, amount int) error
}

type reporter interface {
	Report(
		ctx context.Context,
		format types.ReportFormat,
		opType types.OperationType,
		wallet types.WalletID,
		dateFrom, dateTo time.Time,
	)
}

type handler struct {
	wg walletGenerator
	s  storage
	r  reporter
}

func New(wg walletGenerator, s storage, r reporter) *handler {
	return &handler{
		wg: wg,
		s:  s,
		r:  r,
	}
}

func writeErrorResponse(w http.ResponseWriter, code int, err error) {
	log.WithError(err).Error("handler error")

	w.WriteHeader(code)
	if _, err = w.Write(errorResponse(err)); err != nil {
		log.WithError(err).Error("failed to write erroneous response")
	}
}

func errorResponse(err error) []byte {
	return []byte(`{"error":"` + strings.TrimSpace(err.Error()) + `"}`)
}
