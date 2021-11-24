package handlers_test

import (
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

func TestHandleCreateWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("wallet generator error", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/wallet", nil)
		require.NoError(t, err)

		generatorMock := mocks.NewMockwalletGenerator(ctrl)
		generatorMock.EXPECT().
			Generate().
			Times(1).
			Return(types.WalletID(""), errors.New("some error"))

		rr := httptest.NewRecorder()
		handlers.New(generatorMock, nil, nil).HandleCreateWallet(rr, req, nil)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, `{"error":"generate wallet id: some error"}`, rr.Body.String())
	})

	t.Run("storage error", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/wallet", nil)
		require.NoError(t, err)

		generatorMock := mocks.NewMockwalletGenerator(ctrl)
		generatorMock.EXPECT().Generate().Times(1).Return(types.WalletID("walletID"), nil)

		storageMock := mocks.NewMockstorage(ctrl)
		storageMock.EXPECT().
			CreateWallet(gomock.Any(), types.WalletID("walletID")).
			Times(1).
			Return(errors.New("storage error"))

		rr := httptest.NewRecorder()
		handlers.New(generatorMock, storageMock, nil).HandleCreateWallet(rr, req, nil)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, `{"error":"create wallet: storage error"}`, rr.Body.String())
	})

	t.Run("happy path", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/wallet", nil)
		require.NoError(t, err)

		generatorMock := mocks.NewMockwalletGenerator(ctrl)
		generatorMock.EXPECT().Generate().Times(1).Return(types.WalletID("walletID"), nil)

		storageMock := mocks.NewMockstorage(ctrl)
		storageMock.EXPECT().
			CreateWallet(gomock.Any(), types.WalletID("walletID")).
			Times(1).
			Return(nil)

		rr := httptest.NewRecorder()
		handlers.New(generatorMock, storageMock, nil).HandleCreateWallet(rr, req, nil)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, `{"wallet_id":"walletID"}`, rr.Body.String())
	})
}
