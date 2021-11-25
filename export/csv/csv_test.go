package csv_test

import (
	"testing"

	"github.com/justteddy/wallet/export/csv"
	"github.com/justteddy/wallet/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	data, err := csv.Format([]types.ExportOperation{
		{
			WalletID:      "wallet1",
			OperationType: "operation",
			Amount:        "100.00$",
			Date:          "2030-01-01",
		},
		{
			WalletID:      "wallet2",
			OperationType: "operation2",
			Amount:        "200.00$",
			Date:          "2030-01-02",
		},
	})

	expected := []byte(`wallet_id,operation_id,amount,date
wallet1,operation,100.00$,2030-01-01
wallet2,operation2,200.00$,2030-01-02
`)

	require.NoError(t, err)
	assert.Equal(t, expected, data)
}
