package main

import (
	"context"
	"fmt"

	"github.com/InjectiveLabs/injective-guilds-service/internal/db/mongoimpl"
	"github.com/InjectiveLabs/injective-guilds-service/internal/exchange"
	log "github.com/xlab/suplog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// this will support the guilds process to get correct denom price in USD from asset-price service
var denomToCoinID = map[string]string{
	"peggy0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2": "weth",
	"peggy0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48": "usd-coin",
	"inj": "injective-protocol",
	"peggy0xdAC17F958D2ee523a2206206994597C13D831ec7": "tether",
	"peggy0x514910771AF9Ca656af840dff83E8264EcF986CA": "chainlink",
	"peggy0x7Fc66500c84A76Ad7e9c93437bFc5Ac33E2DDaE9": "aave",
	"peggy0x7D1AfA7B718fb893dB30A3aBc0Cfc608AaCfeBB0": "matic-network",
	"peggy0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984": "uniswap",
	"peggy0x6B3595068778DD592e39A122f4f5a5cF09C90fE2": "sushi",
	"peggy0xc944E90C64B2c07662A292be6244BDf05Cda44a7": "the-graph",
	"peggy0xC011a73ee8576Fb46F5E1c5751cA3B9Fe0af2a6F": "havven",
	// TODO: Add quant-network to support list
	// "peggy0x4a220E6096B25EADb88358cb44068A3248254675":                      "quant-network",
	"peggy0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599":                      "wrapped-bitcoin",
	"peggy0xBB0E17EF65F82Ab018d8EDd776e8DD940327B28b":                      "axie-infinity",
	"ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9": "cosmos",
	"peggy0xAaEf88cEa01475125522e117BFe45cF32044E238":                      "guildfi",
	"ibc/B448C0CA358B958301D328CCDC5D5AD642FC30A6D3AE106FF721DB315F3DDE5C": "terrausd",
	"ibc/B8AF5D92165F35AB31F3FC7C7B444B9D240760FA5D406C49D24862BD0284E395": "terra-luna",
	"ibc/E7807A46C0B7B44B350DA58F51F278881B863EC4DCA94635DAB39E52C30766CB": "chihuahua-token",
}

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

func actionUpdateDenom() {
	ctx := context.Background()
	dbSvc, err := mongoimpl.NewService(ctx, *dbURL, "guilds")
	panicIf(err)

	defer dbSvc.Disconnect(ctx)

	// don't want to mess up internal db code -> implement write steps here
	mgo := dbSvc.(*mongoimpl.MongoImpl)
	denomsColl := mgo.GetClient().Database("guilds").Collection("denoms")

	for denom, coinId := range denomToCoinID {
		filter := bson.M{"denom": denom}
		upd := bson.M{"$set": bson.M{"coin_id": coinId}}

		opt := &options.UpdateOptions{}
		opt.SetUpsert(true)
		_, err := denomsColl.UpdateOne(ctx, filter, upd, opt)
		panicIf(err)
	}

	allDenoms, err := dbSvc.ListDenomCoinID(ctx)
	panicIf(err)

	for _, d := range allDenoms {
		log.Info(fmt.Sprintf("updated: %+v\n", d))
	}

	log.Info("Double checking if asset-price supports all added coins...")
	coinIDs := make([]string, 0)
	for _, v := range denomToCoinID {
		coinIDs = append(coinIDs, v)
	}

	provider, err := exchange.NewExchangeProvider("", "", *assetPriceURL)
	panicIf(err)

	prices := getAllCoinPrices(ctx, provider, coinIDs)
	notFoundCoins := make([]string, 0)

	existMap := make(map[string]bool)
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
	fmt.Println("!!! Bravo, all coins price are supported")
}
