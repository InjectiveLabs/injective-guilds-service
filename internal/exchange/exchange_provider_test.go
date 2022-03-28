package exchange

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubaccountBalance(t *testing.T) {
	// integration test
	provider, err := NewExchangeProvider("sentry2.injective.network:9910", "https://lcd.injective.network", "")
	if err != nil {
		panic(err)
	}
	defer provider.Close()

	ctx := context.Background()
	balances, err := provider.GetSubaccountBalances(
		ctx,
		"0xaf79152ac5df276d9a8e1e2e22822f9713474902000000000000000000000000",
	)
	assert.NoError(t, err)
	assert.Equal(t, len(balances) > 0, true)
}

func TestGrants(t *testing.T) {
	// integration test
	provider, err := NewExchangeProvider("sentry2.injective.network:9910", "https://testnet.lcd.injective.dev", "")
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
	assert.NoError(t, err)

	fmt.Println(grants)
}

func TestAssetPrice(t *testing.T) {
	// integration test
	provider, err := NewExchangeProvider("sentry2.injective.network:9910", "", "https://k8s.mainnet.asset.injective.network")
	if err != nil {
		panic(err)
	}
	defer provider.Close()

	ctx := context.Background()
	prices, err := provider.GetPriceUSD(
		ctx,
		[]string{"bitcoin", "ethereum"},
	)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, len(prices), 2)
	for _, p := range prices {
		fmt.Printf("Price of '%s': %.4f\n", p.ID, p.CurrentPrice)
	}
}

func TestBankAccount(t *testing.T) {
	// integration test
	provider, err := NewExchangeProvider("sentry2.injective.network:9910", "https://lcd.injective.network", "https://k8s.mainnet.asset.injective.network")
	if err != nil {
		panic(err)
	}
	defer provider.Close()

	ctx := context.Background()
	bank, err := provider.GetBankBalance(
		ctx,
		"inj14au322k9munkmx5wrchz9q30juf5wjgz2cfqku",
	)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, len(bank.Balances), 1)
	for _, b := range bank.Balances {
		assert.Equal(t, "inj", b.Denom)
	}

	fmt.Printf("%#v\n", bank)
}
