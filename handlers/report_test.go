package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/justteddy/wallet/handlers"
	"github.com/justteddy/wallet/handlers/mocks"
	"github.com/justteddy/wallet/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("validation error - empty wallet", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"from_date": "2030-01-01", "to_date": "2030-01-01","operation_type": "deposit"}`))
		req, err := http.NewRequest(http.MethodPost, "/report", body)
		require.NoError(t, err)

		params := []httprouter.Param{
			{
				Key:   "format",
				Value: "json",
			},
			{
				Key:   "wallet",
				Value: "",
			},
		}
		rr := httptest.NewRecorder()
		handlers.New(nil, nil, nil).HandleReport(rr, req, params)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, `{"error":"empty wallet id"}`, rr.Body.String())
	})

	t.Run("validation error - invalid export format", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"from_date": "2030-01-01", "to_date": "2030-01-01","operation_type": "deposit"}`))
		req, err := http.NewRequest(http.MethodPost, "/report", body)
		require.NoError(t, err)

		params := []httprouter.Param{
			{
				Key:   "format",
				Value: "UnexpectedFormat",
			},
			{
				Key:   "wallet",
				Value: "walletID",
			},
		}
		rr := httptest.NewRecorder()
		handlers.New(nil, nil, nil).HandleReport(rr, req, params)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, `{"error":"unexpected export format"}`, rr.Body.String())
	})

	t.Run("validation error - invalid operation type", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"from_date": "2030-01-01", "to_date": "2030-01-01", "operation_type": "UnexpectedOperation"}`))
		req, err := http.NewRequest(http.MethodPost, "/report", body)
		require.NoError(t, err)

		params := []httprouter.Param{
			{
				Key:   "format",
				Value: "json",
			},
			{
				Key:   "wallet",
				Value: "walletID",
			},
		}
		rr := httptest.NewRecorder()
		handlers.New(nil, nil, nil).HandleReport(rr, req, params)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, `{"error":"unexpected operation type"}`, rr.Body.String())
	})

	t.Run("validation error - invalid from date", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"from_date": "01-01-2030", "to_date": "2030-01-01","operation_type": "deposit"}`))
		req, err := http.NewRequest(http.MethodPost, "/report", body)
		require.NoError(t, err)

		params := []httprouter.Param{
			{
				Key:   "format",
				Value: "json",
			},
			{
				Key:   "wallet",
				Value: "walletID",
			},
		}
		rr := httptest.NewRecorder()
		handlers.New(nil, nil, nil).HandleReport(rr, req, params)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid date format in from_date")
	})

	t.Run("validation error - invalid to date", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"from_date": "2030-01-01", "to_date": "01-01-2030","operation_type": "deposit"}`))
		req, err := http.NewRequest(http.MethodPost, "/report", body)
		require.NoError(t, err)

		params := []httprouter.Param{
			{
				Key:   "format",
				Value: "json",
			},
			{
				Key:   "wallet",
				Value: "walletID",
			},
		}
		rr := httptest.NewRecorder()
		handlers.New(nil, nil, nil).HandleReport(rr, req, params)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid date format in to_date")
	})

	t.Run("validation error - from_date greater than to_date", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"from_date": "2030-01-02", "to_date": "2030-01-01","operation_type": "deposit"}`))
		req, err := http.NewRequest(http.MethodPost, "/report", body)
		require.NoError(t, err)

		params := []httprouter.Param{
			{
				Key:   "format",
				Value: "json",
			},
			{
				Key:   "wallet",
				Value: "walletID",
			},
		}
		rr := httptest.NewRecorder()
		handlers.New(nil, nil, nil).HandleReport(rr, req, params)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, `{"error":"from_date is greater than to_date"}`, rr.Body.String())
	})

	t.Run("storage error", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"from_date": "2030-01-01", "to_date": "2030-01-01","operation_type": "deposit"}`))
		req, err := http.NewRequest(http.MethodPost, "/report", body)
		require.NoError(t, err)

		fromDate, _ := time.Parse("2006-02-01", "2030-01-01")
		toDate := fromDate

		storageMock := mocks.NewMockstorage(ctrl)
		storageMock.EXPECT().
			Operations(gomock.Any(), types.WalletID("walletID"), types.OperationTypeDeposit, fromDate, toDate).
			Times(1).
			Return(nil, errors.New("storage error"))

		params := []httprouter.Param{
			{
				Key:   "format",
				Value: "json",
			},
			{
				Key:   "wallet",
				Value: "walletID",
			},
		}
		rr := httptest.NewRecorder()
		handlers.New(nil, storageMock, nil).HandleReport(rr, req, params)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, `{"error":"fetch operations: storage error"}`, rr.Body.String())
	})

	t.Run("exporter error", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"from_date": "2030-01-01", "to_date": "2030-01-01","operation_type": "deposit"}`))
		req, err := http.NewRequest(http.MethodPost, "/report", body)
		require.NoError(t, err)

		fromDate, _ := time.Parse("2006-02-01", "2030-01-01")
		toDate := fromDate

		storageMock := mocks.NewMockstorage(ctrl)
		storageMock.EXPECT().
			Operations(gomock.Any(), types.WalletID("walletID"), types.OperationTypeDeposit, fromDate, toDate).
			Times(1).
			Return([]types.Operation{}, nil)

		exporterMock := mocks.NewMockexporter(ctrl)
		exporterMock.EXPECT().
			Export(types.ExportFormatJSON, []types.Operation{}).
			Times(1).
			Return(nil, errors.New("exporter error"))

		params := []httprouter.Param{
			{
				Key:   "format",
				Value: "json",
			},
			{
				Key:   "wallet",
				Value: "walletID",
			},
		}
		rr := httptest.NewRecorder()
		handlers.New(nil, storageMock, exporterMock).HandleReport(rr, req, params)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, `{"error":"export operations: exporter error"}`, rr.Body.String())
	})

	t.Run("happy path", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"from_date": "2030-01-01", "to_date": "2030-01-01","operation_type": "deposit"}`))
		req, err := http.NewRequest(http.MethodPost, "/report", body)
		require.NoError(t, err)

		fromDate, _ := time.Parse("2006-02-01", "2030-01-01")
		toDate := fromDate

		storageMock := mocks.NewMockstorage(ctrl)
		storageMock.EXPECT().
			Operations(gomock.Any(), types.WalletID("walletID"), types.OperationTypeDeposit, fromDate, toDate).
			Times(1).
			Return([]types.Operation{}, nil)

		exporterMock := mocks.NewMockexporter(ctrl)
		exporterMock.EXPECT().
			Export(types.ExportFormatJSON, []types.Operation{}).
			Times(1).
			Return([]byte(`success`), nil)

		params := []httprouter.Param{
			{
				Key:   "format",
				Value: "json",
			},
			{
				Key:   "wallet",
				Value: "walletID",
			},
		}
		rr := httptest.NewRecorder()
		handlers.New(nil, storageMock, exporterMock).HandleReport(rr, req, params)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, `success`, rr.Body.String())
	})
}
