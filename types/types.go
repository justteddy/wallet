package types

import "errors"

var (
	ErrUnavailableBalance = errors.New("insufficient funds in the account")
)

type WalletID string

type ReportFormat string

const (
	ReportFormatJSON = "json"
	ReportFormatCSV  = "csv"
)

type OperationType string

const (
	OperationTypeDeposit  = "deposit"
	OperationTypeWithdraw = "withdraw"
)
