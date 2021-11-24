package storage

import (
	"context"
	"time"

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
	if err := tx.QueryRowContext(ctx, queryLockWalletForCreate, wallet).Scan(&balance); err != nil {
		return completeTx(tx, errors.Wrap(err, "lock wallet"))
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
	tx, err := s.conn.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}

	rows, err := tx.QueryContext(ctx, queryLockWalletsForTransfer, fromWallet, toWallet)
	if err != nil {
		return completeTx(tx, errors.Wrap(err, "lock wallets"))
	}

	var from, to struct {
		id      string
		balance int
	}

	for rows.Next() {
		var (
			id      string
			balance int
		)
		if err := rows.Scan(&id, &balance); err != nil {
			return completeTx(tx, errors.Wrap(err, "scan rows"))
		}

		if id == string(fromWallet) {
			from.id = id
			from.balance = balance
		} else {
			to.id = id
			to.balance = balance
		}
	}

	if from.balance-amount < 0 {
		return completeTx(tx, types.ErrUnavailableBalance)
	}

	// add withdraw operation on fromWallet
	if _, err := tx.ExecContext(ctx, queryInsertOperation, fromWallet, types.OperationTypeWithdraw, amount); err != nil {
		return completeTx(tx, errors.Wrap(err, "create operation withdraw"))
	}

	// add deposit operation on toWallet
	if _, err := tx.ExecContext(ctx, queryInsertOperation, toWallet, types.OperationTypeDeposit, amount); err != nil {
		return completeTx(tx, errors.Wrap(err, "create operation deposit"))
	}

	// change fromWallet balance
	if _, err := tx.ExecContext(ctx, queryUpdateWallet, from.balance-amount, fromWallet); err != nil {
		return completeTx(tx, errors.Wrap(err, "update fromWallet balance"))
	}

	// change toWallet balance
	if _, err := tx.ExecContext(ctx, queryUpdateWallet, to.balance+amount, toWallet); err != nil {
		return completeTx(tx, errors.Wrap(err, "update toWallet balance"))
	}

	return completeTx(tx, nil)
}

func (s *storage) Operations(ctx context.Context, wallet types.WalletID, opType types.OperationType, from, to time.Time) ([]types.Operation, error) {
	return nil, nil
}
