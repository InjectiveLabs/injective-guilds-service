package main

import (
	"context"
	"fmt"
	"os"

	"github.com/InjectiveLabs/injective-guilds-service/internal/db/model"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db/mongoimpl"
	"github.com/InjectiveLabs/injective-guilds-service/internal/exchange"
	guildsprocess "github.com/InjectiveLabs/injective-guilds-service/internal/service/guilds-process"
	derivativeExchangePB "github.com/InjectiveLabs/sdk-go/exchange/derivative_exchange_rpc/pb"
	spotExchangePB "github.com/InjectiveLabs/sdk-go/exchange/spot_exchange_rpc/pb"
	cosmtypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	cli "github.com/jawher/mow.cli"
	log "github.com/xlab/suplog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func parseAddGuildArgs(c *cli.Cmd) {
	spotIDs = c.Strings(cli.StringsOpt{
		Name:  "spot-id",
		Desc:  "spot marketID",
		Value: []string{},
	})

	derivativeIDs = c.Strings(cli.StringsOpt{
		Name:  "derivative-id",
		Desc:  "spot marketID",
		Value: []string{},
	})

	name = c.String(cli.StringOpt{
		Name:  "name",
		Desc:  "guild name",
		Value: "Default guild name",
	})

	description = c.String(cli.StringOpt{
		Name:  "description",
		Desc:  "guild description",
		Value: "This is guild description",
	})

	capacity = c.Int(cli.IntOpt{
		Name:  "capacity",
		Desc:  "guild capacity",
		Value: 500,
	})

	masterAddr = c.String(cli.StringOpt{
		Name:  "master",
		Desc:  "guild's master address",
		Value: "",
	})

	defaultMemberAddr = c.String(cli.StringOpt{
		Name:  "default-member",
		Desc:  "guild's default member address",
		Value: "",
	})

	dbURL = c.String(cli.StringOpt{
		Name:  "db-url",
		Desc:  "database url",
		Value: "mongodb://localhost:27017",
	})

	exchangeURL = c.String(cli.StringOpt{
		Name:  "exchange-url",
		Desc:  "exchange grpc api url",
		Value: "localhost:9910",
	})

	assetPriceURL = c.String(cli.StringOpt{
		Name:  "asset-price-url",
		Desc:  "asset-price url",
		Value: "https://k8s.mainnet.asset.injective.network",
	})

	minSpotBase = c.Int(cli.IntOpt{
		Name:  "min-spot-base",
		Desc:  "min spot base usd amount",
		Value: 250,
	})

	minSpotQuote = c.Int(cli.IntOpt{
		Name:  "min-spot-quote",
		Desc:  "min spot quote usd amount",
		Value: 250,
	})

	minDerivativeQuote = c.Int(cli.IntOpt{
		Name:  "min-derivative-quote",
		Desc:  "min derivative base usd amount",
		Value: 250,
	})

	minStaking = c.Int(cli.IntOpt{
		Name:  "min-staking",
		Desc:  "min staking amount", // <- we can't get this atm
		Value: 250,
	})
}

func validateAddGuildArgs() {
	if len(*spotIDs) == 0 && len(*derivativeIDs) == 0 {
		log.Error("cannot create guild with nomarket")
		os.Exit(1)
	}

	_, err := cosmtypes.AccAddressFromBech32(*masterAddr)
	panicIf(err)

	_, err = cosmtypes.AccAddressFromBech32(*defaultMemberAddr)
	panicIf(err)
}

func addGuildAction() {
	validateAddGuildArgs()

	ctx := context.Background()
	log.Info("connect db service at", *dbURL)
	dbSvc, err := mongoimpl.NewService(ctx, *dbURL, "guilds")
	panicIf(err)

	log.Info("connecting exchange api at", *exchangeURL)
	exchangeProvider, err := exchange.NewExchangeProvider(*exchangeURL, "", *assetPriceURL)
	panicIf(err)

	spotClient := spotExchangePB.NewInjectiveSpotExchangeRPCClient(exchangeProvider.GetExchangeConn())
	derivativeClient := derivativeExchangePB.NewInjectiveDerivativeExchangeRPCClient(exchangeProvider.GetExchangeConn())

	log.Info("checking markets using exchange api service")

	markets := make([]*model.GuildMarket, 0)
	for _, m := range *spotIDs {
		log.Info(fmt.Sprintf(">>> checking spot market id: %s", m))

		req := &spotExchangePB.MarketRequest{
			MarketId: m,
		}
		marketResp, err := spotClient.Market(ctx, req)
		panicIf(err)

		market := marketResp.GetMarket()
		takerFeeRate, err := primitive.ParseDecimal128(market.GetTakerFeeRate())
		panicIf(err)
		marketFeeRate, err := primitive.ParseDecimal128(market.GetMakerFeeRate())
		panicIf(err)

		markets = append(markets, &model.GuildMarket{
			MarketID:    model.Hash{Hash: common.HexToHash(market.GetMarketId())},
			IsPerpetual: false, // to clarify this is spot market
			BaseDenom:   market.GetBaseDenom(),
			BaseTokenMeta: &model.TokenMeta{
				Name:      market.GetBaseTokenMeta().GetName(),
				Address:   market.GetBaseTokenMeta().Address,
				Symbol:    market.GetBaseTokenMeta().Symbol,
				Decimals:  int(market.GetBaseTokenMeta().GetDecimals()),
				UpdatedAt: market.GetBaseTokenMeta().GetUpdatedAt(),
			},
			QuoteDenom: market.GetQuoteDenom(),
			QuoteTokenMeta: &model.TokenMeta{
				Name:      market.GetQuoteTokenMeta().GetName(),
				Address:   market.GetQuoteTokenMeta().Address,
				Symbol:    market.GetQuoteTokenMeta().Symbol,
				Decimals:  int(market.GetQuoteTokenMeta().GetDecimals()),
				UpdatedAt: market.GetQuoteTokenMeta().GetUpdatedAt(),
			},
			TakerFeeRate: takerFeeRate,
			MakerFeeRate: marketFeeRate,
		})
	}

	for _, m := range *derivativeIDs {
		log.Info(fmt.Sprintf(">>> checking derivative market id: %s", m))
		req := &derivativeExchangePB.MarketRequest{
			MarketId: m,
		}

		marketResp, err := derivativeClient.Market(ctx, req)
		panicIf(err)

		market := marketResp.GetMarket()
		markets = append(markets, &model.GuildMarket{
			MarketID:      model.Hash{Hash: common.HexToHash(market.GetMarketId())},
			IsPerpetual:   true,
			BaseDenom:     "",
			BaseTokenMeta: nil,
			QuoteDenom:    market.GetQuoteDenom(),
			QuoteTokenMeta: &model.TokenMeta{
				Name:      market.GetQuoteTokenMeta().GetName(),
				Address:   market.GetQuoteTokenMeta().GetAddress(),
				Symbol:    market.GetQuoteTokenMeta().GetSymbol(),
				Decimals:  int(market.GetQuoteTokenMeta().GetDecimals()),
				UpdatedAt: market.GetQuoteTokenMeta().GetUpdatedAt(),
			},
		})
	}

	master, _ := cosmtypes.AccAddressFromBech32(*masterAddr)
	defaultMember, _ := cosmtypes.AccAddressFromBech32(*defaultMemberAddr)
	guild := &model.Guild{
		Name:                       *name,
		Description:                *description,
		MasterAddress:              model.Address{AccAddress: master},
		SpotBaseRequirement:        *minSpotBase,
		SpotQuoteRequirement:       *minSpotQuote,
		DerivativeQuoteRequirement: *minDerivativeQuote,
		StakingRequirement:         *minStaking,
		Markets:                    markets,
		Capacity:                   *capacity,
	}

	log.Info("adding guild: ", *name)
	id, err := dbSvc.AddGuild(ctx, guild)
	panicIf(err)

	log.Info("adding default member portfolio")
	err = dbSvc.AddMember(ctx, id.Hex(), model.Address{AccAddress: defaultMember})
	panicIf(err)

	log.Info("capturing default member portfolio")
	helper, err := guildsprocess.NewPortfolioHelper(ctx, dbSvc, exchangeProvider)
	panicIf(err)

	portfolio, err := helper.CaptureSingleMemberPortfolio(ctx, guild, &model.GuildMember{
		GuildID:          *id,
		InjectiveAddress: model.Address{AccAddress: defaultMember},
	}, true)
	panicIf(err)

	err = dbSvc.AddAccountPortfolios(ctx, id.Hex(), []*model.AccountPortfolio{portfolio})
	panicIf(err)

	log.Info("🍺 all done")
}

func cmdAddGuild(c *cli.Cmd) {
	// inputs:
	// guild name: --name
	// guild description: --description
	// database url: --db-url
	// exchange url: --exchange-url
	// spot market (can supply many --spot-id): --spot-id
	// derivative market (can supply many --derivative-id): --derivative-id
	// capacity: --capactity
	// master address: --master
	// default member address: --default-member
	// returns guild id
	// example: go run cmd/injective-guilds/*.go add-guild \
	// --derivative-id=0xc559df216747fc11540e638646c384ad977617d6d8f0ea5ffdfc18d52e58ab01 \
	// --spot-id=0xfbc729e93b05b4c48916c1433c9f9c2ddb24605a73483303ea0f87a8886b52af \
	// --name=testguild --description "a test guild" --master=inj1awx03zmnnlsjuvp7x8ac3lphw50p0nea6p2584 \
	// --default-member=inj1zggdm44ln2gu7c5d2ge4wyr4wfs0cfn5lyfw4k --exchange-url=sentry2.injective.network:9910
	parseAddGuildArgs(c)
	c.Action = addGuildAction
}
