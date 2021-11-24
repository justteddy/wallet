package handlers

import (
	"context"
	"time"

	"github.com/justteddy/wallet/types"
)

type storage interface {
	CreateWallet(ctx context.Context) (string, error)
	Deposit(ctx context.Context, wallet types.WalletID, amount int) error
	Transfer(ctx context.Context, fromWallet, toWallet types.WalletID, amount int) error
}

type reporter interface {
	Report(format types.ReportFormat, opType types.OperationType, wallet types.WalletID, dateFrom, dateTo time.Time)
}

type handler struct {
	s storage
	r reporter
}
