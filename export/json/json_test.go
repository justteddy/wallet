package json_test

import (
	"testing"

	"github.com/justteddy/wallet/export/json"
	"github.com/justteddy/wallet/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	data, err := json.Format([]types.ExportOperation{
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

	expected := []byte(`[{"wallet_id":"wallet1","operation_type":"operation","amount":"100.00$","date":"2030-01-01"},{"wallet_id":"wallet2","operation_type":"operation2","amount":"200.00$","date":"2030-01-02"}]`)

	require.NoError(t, err)
	assert.Equal(t, expected, data)
}
