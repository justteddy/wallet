package types

import (
	"errors"

	"github.com/justteddy/wallet/currency"
)

const DateLayout = "2006-01-02"

var ErrUnavailableBalance = errors.New("insufficient funds in the account")

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

type DBOperation struct {
	WalletID      WalletID      `db:"wallet_id"`
	OperationType OperationType `db:"operation_type"`
	Amount        int           `db:"amount"`
	CreatedAt     string        `db:"created_at"`
}

type ExportOperation struct {
	WalletID      string `json:"wallet_id"`
	OperationType string `json:"operation_type"`
	Amount        string `json:"amount"`
	Date          string `json:"date"`
}

// TransformDBToExportOperation transforms DBOperation to ExportOperation
func TransformDBToExportOperation(ops []DBOperation) []ExportOperation {
	expOps := make([]ExportOperation, 0, len(ops))
	for _, op := range ops {
		expOps = append(expOps, ExportOperation{
			WalletID:      string(op.WalletID),
			OperationType: string(op.OperationType),
			Amount:        currency.Format(op.Amount),
			Date:          op.CreatedAt,
		})
	}

	return expOps
}
