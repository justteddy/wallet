package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/justteddy/wallet/handlers"
	"github.com/justteddy/wallet/handlers/mocks"
	"github.com/justteddy/wallet/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleTransfer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("storage error", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"from_wallet": "wallet1", "to_wallet": "wallet2","amount": 100}`))
		req, err := http.NewRequest(http.MethodPost, "/transfer", body)
		require.NoError(t, err)

		storageMock := mocks.NewMockstorage(ctrl)
		storageMock.EXPECT().
			Transfer(gomock.Any(), types.WalletID("wallet1"), types.WalletID("wallet2"), 100).
			Times(1).
			Return(errors.New("storage error"))

		rr := httptest.NewRecorder()
		handlers.New(nil, storageMock, nil).HandleTransfer(rr, req, nil)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, `{"error":"save to storage: storage error"}`, rr.Body.String())
	})

	t.Run("happy path", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"from_wallet": "wallet1", "to_wallet": "wallet2","amount": 100}`))
		req, err := http.NewRequest(http.MethodPost, "/transfer", body)
		require.NoError(t, err)

		storageMock := mocks.NewMockstorage(ctrl)
		storageMock.EXPECT().
			Transfer(gomock.Any(), types.WalletID("wallet1"), types.WalletID("wallet2"), 100).
			Times(1).
			Return(nil)

		rr := httptest.NewRecorder()
		handlers.New(nil, storageMock, nil).HandleTransfer(rr, req, nil)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
