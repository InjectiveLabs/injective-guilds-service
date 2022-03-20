package exchange

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubaccountBalance(t *testing.T) {
	// integration test
	provider, err := NewExchangeProvider("sentry2.injective.network:9910", "https://lcd.injective.network")
	if err != nil {
		panic(err)
	}
	defer provider.Close()

	ctx := context.Background()
	balances, err := provider.GetDefaultSubaccountBalances(
		ctx,
		"0xaf79152ac5df276d9a8e1e2e22822f9713474902000000000000000000000000",
	)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, len(balances) > 0, true)
}

func TestGrants(t *testing.T) {
	// integration test
	provider, err := NewExchangeProvider("sentry2.injective.network:9910", "https://testnet.lcd.injective.dev")
	if err != nil {
		panic(err)
	}
	defer provider.Close()

	ctx := context.Background()
	grants, err := provider.GetGrants(
		ctx,
		"inj14au322k9munkmx5wrchz9q30juf5wjgz2cfqku",
		"inj1hkhdaj2a2clmq5jq6mspsggqs32vynpk228q3r",
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(grants)
}
