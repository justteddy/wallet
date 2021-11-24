package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestName(t *testing.T) {
	t.Skip()
	data := []byte(`
		{
			"amount": 1
		}
	`)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reader := bytes.NewReader(data)
			req, err := http.NewRequest("POST", "http://localhost:8080/deposit/34f3e0346472e0f2ac7ccb5a930fcdc20474a227136c516f96104b6427eb0419", reader)
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}()
	}

	wg.Wait()
}

func TestName2(t *testing.T) {
	t.Skip()
	data := []byte(`
		{
			"from_wallet": "34f3e0346472e0f2ac7ccb5a930fcdc20474a227136c516f96104b6427eb0419",
			"to_wallet": "26bb890c2e4ec72bafefbf95814e16e77767a5eb1ff1c784fe52edc9ae9b463c",
			"amount": 1
		}
	`)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reader := bytes.NewReader(data)
			req, err := http.NewRequest("POST", "http://localhost:8080/transfer", reader)
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}()
	}

	wg.Wait()
}

func TestX(t *testing.T) {
	val := 11252

	dollars := val / 100
	cents := val % 100

	fmt.Printf("%d.%d$\n", dollars, cents)

	fmt.Println(time.Time{}.IsZero())
}
