package types

import (
	"errors"
)

const (
	DateLayout = "2006-01-02"
)

var (
	ErrUnavailableBalance = errors.New("insufficient funds in the account")
)

type WalletID string

type ExportFormat string

const (
	ExportFormatJSON ExportFormat = "json"
	ExportFormatCSV  ExportFormat = "csv"
)

type OperationType string

const (
	OperationTypeDeposit  OperationType = "deposit"
	OperationTypeWithdraw OperationType = "withdraw"
)

var (
	AllExportFormats = map[ExportFormat]struct{}{
		ExportFormatJSON: {},
		ExportFormatCSV:  {},
	}

	AllOperationTypes = map[OperationType]struct{}{
		OperationTypeDeposit:  {},
		OperationTypeWithdraw: {},
	}
)

type Operation struct {
	WalletID      WalletID      `db:"wallet_id"`
	OperationType OperationType `db:"operation_type"`
	Amount        int           `db:"amount"`
	Date          string        `db:"date"`
}
