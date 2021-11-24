package types

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
