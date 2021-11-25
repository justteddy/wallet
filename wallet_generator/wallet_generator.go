package wallet_generator

import (
	"crypto/sha256"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/justteddy/wallet/types"
	"github.com/pkg/errors"
)

type generator struct {
	current int64
}

func New() *generator {
	return &generator{}
}

// Generate generates unique wallet identifier.
// It works properly in scope of single service instance.
// In case of multi instance service run should use some distributed mechanisms like - database sequence, etcd, etc.
func (g *generator) Generate() (types.WalletID, error) {
	current := atomic.AddInt64(&g.current, 1)
	key := fmt.Sprintf("%d", int64(time.Now().UTC().Nanosecond())+current)

	hasher := sha256.New()
	if _, err := hasher.Write([]byte(key)); err != nil {
		return "", errors.Wrap(err, "generate wallet id")
	}

	return types.WalletID(fmt.Sprintf("%x", hasher.Sum(nil))), nil
}
