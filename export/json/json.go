package json

import (
	"encoding/json"

	"github.com/justteddy/wallet/types"
)

func Format(ops []types.Operation) ([]byte, error) {
	return json.Marshal(ops)
}
