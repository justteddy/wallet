package types

import (
	"errors"
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
	WalletID      WalletID      `json:"wallet_id" csv:"wallet_id"`
	OperationType OperationType `json:"operation_type" csv:"operation_type"`
	Amount        int           `json:"amount" csv:"amount"`
	Date          string        `json:"date" csv:"date"`
}
