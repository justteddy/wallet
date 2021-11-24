package wallet_generator_test

import (
	"sync"
	"testing"

	"github.com/justteddy/wallet/types"
	"github.com/justteddy/wallet/wallet_generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateWalletID(t *testing.T) {
	g := wallet_generator.New()

	wallets := make(map[types.WalletID]struct{})
	walletsCh := make(chan types.WalletID, 1)

	done := make(chan struct{})
	go func() {
		defer close(done)
		for walletID := range walletsCh {
			wallets[walletID] = struct{}{}
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < 100_000; i++ {
		wg.Add(1)
		go func(counter int) {
			defer wg.Done()
			walletID, err := g.Generate()
			require.NoError(t, err)
			walletsCh <- walletID
		}(i)
	}

	wg.Wait()
	close(walletsCh)

	// wait until all wallets written
	<-done

	// check if all wallet identifiers are unique
	assert.Len(t, wallets, 100_000)
}
