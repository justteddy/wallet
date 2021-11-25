package json

import (
	"encoding/json"

	"github.com/justteddy/wallet/types"
)

func Format(ops []types.ExportOperation) ([]byte, error) {
	return json.Marshal(ops)
}
