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
	_, err := s.conn.ExecContext(ctx, queryInsertWallet, wallet)

	return errors.Wrap(err, "create wallet query error")
}

func (s *storage) Deposit(ctx context.Context, wallet types.WalletID, amount int) error {
	tx, err := s.conn.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}

	var balance int
	if err := tx.QueryRowContext(ctx, queryLockWalletForUpdate, wallet).Scan(&balance); err != nil {
		return completeTx(tx, errors.Wrap(err, "lock wallet row"))
	}

	if _, err := tx.ExecContext(ctx, queryInsertOperation, wallet, types.OperationTypeDeposit, amount); err != nil {
		return completeTx(tx, errors.Wrap(err, "create operation"))
	}

	if _, err := tx.ExecContext(ctx, queryUpdateWallet, balance+amount, wallet); err != nil {
		return completeTx(tx, errors.Wrap(err, "update balance"))
	}

	return completeTx(tx, nil)
}
func (s *storage) Transfer(ctx context.Context, fromWallet, toWallet types.WalletID, amount int) error {
	return nil
}
