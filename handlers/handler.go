package handlers

//go:generate mockgen -source=handler.go -destination=mocks/handler.go -package=mocks

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
	Operations(ctx context.Context, wallet types.WalletID, opType types.OperationType, from, to time.Time) ([]types.Operation, error)
}

type exporter interface {
	Export(format types.ExportFormat, ops []types.Operation) ([]byte, error)
}

type Handler struct {
	wg walletGenerator
	s  storage
	e  exporter
}

func New(wg walletGenerator, s storage, e exporter) *Handler {
	return &Handler{
		wg: wg,
		s:  s,
		e:  e,
	}
}

func writeErrorResponse(w http.ResponseWriter, code int, err error) {
	log.WithError(err).Error("handler error")

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err = w.Write(errorResponse(err)); err != nil {
		log.WithError(err).Error("failed to write erroneous response")
	}
}

func errorResponse(err error) []byte {
	return []byte(`{"error":"` + strings.TrimSpace(err.Error()) + `"}`)
}
