package storage

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/justteddy/wallet/types"
	"github.com/pkg/errors"
)

type storage struct {
	conn *sqlx.DB
}

func New(conn *sqlx.DB) *storage {
	return &storage{
		conn: conn,
	}
}

func (s *storage) CreateWallet(ctx context.Context, wallet types.WalletID) error {
	_, err := s.conn.ExecContext(ctx, queryCreateWallet, wallet)

	return errors.Wrap(err, "create wallet query error")
}

func (s *storage) Deposit(ctx context.Context, wallet types.WalletID, amount int) error { return nil }
func (s *storage) Transfer(ctx context.Context, fromWallet, toWallet types.WalletID, amount int) error {
	return nil
}
