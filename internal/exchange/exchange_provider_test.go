package exchange

import (
	"context"
	"fmt"
	"testing"
)

func TestSubaccountBalance(t *testing.T) {
	// integration test
	provider, err := NewExchangeProvider("sentry2.injective.network:9910")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	balances, err := provider.GetDefaultSubaccountBalances(ctx, "0xaf79152ac5df276d9a8e1e2e22822f9713474902000000000000000000000000")
	if err != nil {
		panic(err)
	}

	for _, b := range balances {
		fmt.Println("denom:", b.Denom)
		fmt.Println("total:", b.TotalBalance.String())
		fmt.Println("avail:", b.AvailableBalance.String())
	}
}
