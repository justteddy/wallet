//go:build integration

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/justteddy/wallet/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	host            = "http://localhost:8080/"
	createWalletURL = "/wallet"
	depositURL      = "/deposit/%s"
	transferURL     = "/transfer"
	reportURL       = "/report/%s/%s"
)

func TestIntegration(t *testing.T) {
	waitHostIsReady(t)

	httpClient := &http.Client{Timeout: time.Second * 3}

	// create wallet
	wallet1 := createWallet(t, httpClient)
	assert.NotEmpty(t, wallet1)

	// deposit concurrently 100 times by 1$
	payload := depositPayload(100)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			deposit(t, httpClient, wallet1, payload)
		}()
	}

	wg.Wait()

	// ensure in exactly 100 deposit operations by 1$ on wallet1
	ops := reportJSON(t, httpClient, wallet1, reportPayload("", "", ""))
	assert.Len(t, ops, 100)
	for _, op := range ops {
		assert.Equal(t, wallet1, op.WalletID)
		assert.Equal(t, "1.00$", op.Amount)
		assert.Equal(t, "deposit", op.OperationType)
		assert.Equal(t, time.Now().Format(types.DateLayout), op.Date)
	}

	// create another wallet
	wallet2 := createWallet(t, httpClient)
	assert.NotEmpty(t, wallet2)

	// concurrently transfer half amount from wallet1 to wallet2 - 1$ per transaction
	payload = transferPayload(wallet1, wallet2, 100)
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			transfer(t, httpClient, payload, http.StatusOK)
		}()
	}

	wg.Wait()

	// ensure in exactly 50 withdraw operations by 1$ on wallet1
	ops = reportJSON(t, httpClient, wallet1, reportPayload("", "", "withdraw"))
	assert.Len(t, ops, 50)
	for _, op := range ops {
		assert.Equal(t, wallet1, op.WalletID)
		assert.Equal(t, "1.00$", op.Amount)
		assert.Equal(t, "withdraw", op.OperationType)
		assert.Equal(t, time.Now().Format(types.DateLayout), op.Date)
	}

	// ensure in exactly 50 deposit operations by 1$ on wallet2
	ops = reportJSON(t, httpClient, wallet2, reportPayload("", "", "deposit"))
	assert.Len(t, ops, 50)
	for _, op := range ops {
		assert.Equal(t, wallet2, op.WalletID)
		assert.Equal(t, "1.00$", op.Amount)
		assert.Equal(t, "deposit", op.OperationType)
		assert.Equal(t, time.Now().Format(types.DateLayout), op.Date)
	}

	// ensure that wallet1 doesn't have enough amount to transfer 50.01$ (there are only 50.00$)
	payload = transferPayload(wallet1, wallet2, 5001)
	transfer(t, httpClient, payload, http.StatusBadRequest)

	// ensure that there is no operations for tomorrow
	tomorrow := time.Now().AddDate(0, 0, 1).Format(types.DateLayout)
	ops = reportJSON(t, httpClient, wallet2, reportPayload(tomorrow, tomorrow, ""))
	assert.Len(t, ops, 0)
}

func transfer(t *testing.T, httpClient *http.Client, payload []byte, expectedStatusCode int) {
	reader := bytes.NewReader(payload)
	req, err := http.NewRequest("POST", host+transferURL, reader)
	require.NoError(t, err)

	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, expectedStatusCode, resp.StatusCode)
}

func transferPayload(fromWallet, toWallet string, amount int) []byte {
	return []byte(
		fmt.Sprintf(`{"from_wallet": "%s", "to_wallet": "%s", "amount": %d}`,
			fromWallet, toWallet, amount,
		),
	)
}

func reportJSON(t *testing.T, httpClient *http.Client, wallet string, payload []byte) []types.ExportOperation {
	reader := bytes.NewReader(payload)
	req, err := http.NewRequest("POST", host+fmt.Sprintf(reportURL, "json", wallet), reader)
	require.NoError(t, err)

	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var ops []types.ExportOperation
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&ops))

	return ops
}

func reportPayload(from, to, opType string) []byte {
	return []byte(
		fmt.Sprintf(`{"from_date": "%s", "to_date": "%s", "operation_type": "%s"}`,
			from, to, opType,
		),
	)
}

func createWallet(t *testing.T, httpClient *http.Client) string {
	req, err := http.NewRequest("POST", host+createWalletURL, nil)
	require.NoError(t, err)

	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response struct {
		WalletID string `json:"wallet_id"`
	}

	require.NoError(t, json.Unmarshal(body, &response))

	return response.WalletID
}

func deposit(t *testing.T, httpClient *http.Client, wallet string, payload []byte) {
	reader := bytes.NewReader(payload)
	req, err := http.NewRequest("POST", host+fmt.Sprintf(depositURL, wallet), reader)
	require.NoError(t, err)

	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func depositPayload(amount int) []byte {
	return []byte(fmt.Sprintf(`{"amount": %d}`, amount))
}

func waitHostIsReady(t *testing.T) {
	var ready bool
	for i := 0; i < 5; i++ {
		<-time.After(time.Second * 2)
		_, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			fmt.Println("host is unreachable, trying again")
			continue
		}
		fmt.Println("host is ready")
		ready = true
		break
	}

	require.True(t, ready)
}
