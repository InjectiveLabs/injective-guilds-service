package main

import (
	"context"

	"github.com/InjectiveLabs/injective-guilds-service/internal/config"
	"github.com/InjectiveLabs/injective-guilds-service/internal/exchange"
	log "github.com/xlab/suplog"
)

func getAllCoinPrices(ctx context.Context, provider exchange.DataProvider, coinIDs []string) []*exchange.CoinPrice {
	var (
		err    error
		prices = make([]*exchange.CoinPrice, 0)
		ids    = coinIDs[:]
	)

	for len(ids) > 0 {
		var tmpPrices []*exchange.CoinPrice
		if len(ids) > 10 {
			tmpPrices, err = provider.GetPriceUSD(ctx, ids[:10])
			panicIf(err)

			ids = ids[10:]
		} else {
			tmpPrices, err = provider.GetPriceUSD(ctx, ids[:])
			panicIf(err)

			ids = []string{}
		}
		prices = append(prices, tmpPrices...)
	}
	return prices
}

func doubleCheckDenomConfig(assetPriceURL string) {
	log.Info("double checking if asset-price supports all configured coins...")

	ctx := context.Background()
	coinIDs := make([]string, 0)
	for _, v := range config.DenomConfigs {
		coinIDs = append(coinIDs, v.CoinID)
	}

	provider, err := exchange.NewExchangeProvider("", "", assetPriceURL)
	panicIf(err)

	var (
		prices        = getAllCoinPrices(ctx, provider, coinIDs)
		notFoundCoins = make([]string, 0)
		existMap      = make(map[string]bool)
	)

	for _, price := range prices {
		existMap[price.ID] = true
	}

	for _, id := range coinIDs {
		if _, exist := existMap[id]; !exist {
			notFoundCoins = append(notFoundCoins, id)
		}
	}

	if len(notFoundCoins) > 0 {
		log.Fatal("asset-price desn't support: ", notFoundCoins)
	}
	log.Info("!!! Bravo, all coins price are supported")
}
