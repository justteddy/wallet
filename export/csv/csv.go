package csv

import (
	"bytes"
	"encoding/csv"

	"github.com/justteddy/wallet/types"
	"github.com/pkg/errors"
)

var headers = []string{"wallet_id", "operation_id", "amount", "date"}

func Format(ops []types.ExportOperation) ([]byte, error) {
	if len(ops) == 0 {
		return []byte{}, nil
	}

	buffer := bytes.NewBuffer(nil)
	w := csv.NewWriter(buffer)
	if err := w.Write(headers); err != nil {
		return nil, errors.Wrap(err, "write csv headers")
	}

	if err := w.WriteAll(transformToStringSlice(ops)); err != nil {
		return nil, errors.Wrap(err, "write csv data")
	}

	return buffer.Bytes(), nil
}

func transformToStringSlice(ops []types.ExportOperation) [][]string {
	result := make([][]string, 0, len(ops))
	for _, op := range ops {
		result = append(result, []string{
			op.WalletID,
			op.OperationType,
			op.Amount,
			op.Date,
		})
	}
	return result
}
