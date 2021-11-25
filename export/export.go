package export

import (
	"github.com/justteddy/wallet/export/csv"
	"github.com/justteddy/wallet/export/json"
	"github.com/justteddy/wallet/types"
	"github.com/pkg/errors"
)

type exportFunc func(ops []types.ExportOperation) ([]byte, error)

type exporter struct {
	toJSON exportFunc
	toCSV  exportFunc
}

func New() *exporter {
	return &exporter{
		toJSON: json.Format,
		toCSV:  csv.Format,
	}
}

// Export marshals []types.ExportOperation to []byte using different formats
func (e *exporter) Export(format types.ExportFormat, ops []types.ExportOperation) ([]byte, error) {
	switch format {
	case types.ExportFormatJSON:
		return e.toJSON(ops)
	case types.ExportFormatCSV:
		return e.toCSV(ops)
	default:
		return nil, errors.New("unexpected export format")
	}
}
