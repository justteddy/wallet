package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/justteddy/wallet/handlers"
	"github.com/justteddy/wallet/handlers/mocks"
	"github.com/justteddy/wallet/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleDeposit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("empty wallet", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/deposit", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handlers.New(nil, nil, nil).HandleDeposit(rr, req, []httprouter.Param{})

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, `{"error":"empty wallet id"}`, rr.Body.String())
	})

	t.Run("storage error", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"amount": 100}`))
		req, err := http.NewRequest(http.MethodPost, "/deposit/walletID", body)
		require.NoError(t, err)

		storageMock := mocks.NewMockstorage(ctrl)
		storageMock.EXPECT().
			Deposit(gomock.Any(), types.WalletID("walletID"), 100).
			Times(1).
			Return(errors.New("storage error"))

		params := []httprouter.Param{
			{
				Key:   "wallet",
				Value: "walletID",
			},
		}
		rr := httptest.NewRecorder()
		handlers.New(nil, storageMock, nil).HandleDeposit(rr, req, params)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, `{"error":"save to storage: storage error"}`, rr.Body.String())
	})

	t.Run("happy path", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"amount": 100}`))
		req, err := http.NewRequest(http.MethodPost, "/deposit/walletID", body)
		require.NoError(t, err)

		storageMock := mocks.NewMockstorage(ctrl)
		storageMock.EXPECT().
			Deposit(gomock.Any(), types.WalletID("walletID"), 100).
			Times(1).
			Return(nil)

		params := []httprouter.Param{
			{
				Key:   "wallet",
				Value: "walletID",
			},
		}
		rr := httptest.NewRecorder()
		handlers.New(nil, storageMock, nil).HandleDeposit(rr, req, params)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
