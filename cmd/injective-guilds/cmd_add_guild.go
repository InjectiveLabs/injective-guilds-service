package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/InjectiveLabs/injective-guilds-service/internal/config"
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

const defaultMinRequirement = float64(250.0)

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

	lcdURL = c.String(cli.StringOpt{
		Name:  "lcd-url",
		Desc:  "lcd url to get bank balance",
		Value: "https://lcd.injective.network",
	})

	spotRequirements = c.Strings(cli.StringsOpt{
		Name:  "spot-require",
		Desc:  "minimum requirements to join this guild. BaseRequirement/QuoteRequirement",
		Value: []string{},
	})

	derivativeRequirements = c.Strings(cli.StringsOpt{
		Name:  "derivative-require",
		Desc:  "minimum requirements to join this guild. QuoteRequirement for derivative market",
		Value: []string{},
	})

	minStaking = c.Int(cli.IntOpt{
		Name:  "min-staking",
		Desc:  "min staking amount", // <- we can't get this atm
		Value: 250,
	})

	memberParams = c.String(cli.StringOpt{
		Name:  "params",
		Desc:  "default member's optional params",
		Value: "",
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

func toMinAmounts(s string) (res []float64, err error) {
	parts := strings.Split(s, "/")
	for _, p := range parts {
		f, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return nil, err
		}

		res = append(res, f)
	}
	return res, nil
}

func checkGrant(ctx context.Context, exchangeProvider exchange.DataProvider, address string, masterAddress string) error {
	grants, err := exchangeProvider.GetGrants(ctx, address, masterAddress)
	if err != nil {
		return fmt.Errorf("get grants err: %w", err)
	}

	grantMap := make(map[string]bool)
	for _, g := range grants.Grants {
		grantMap[g.Authorization.Msg] = true
	}

	for _, msg := range config.GrantRequirements {
		if _, exist := grantMap[msg]; !exist {
			return fmt.Errorf("missing grant msg: %s", msg)
		}
	}
	return nil
}

func addGuildAction() {
	validateAddGuildArgs()

	ctx := context.Background()
	log.Info("connect db service at ", *dbURL)
	dbSvc, err := mongoimpl.NewService(ctx, *dbURL, "guilds")
	panicIf(err)

	log.Info("connecting exchange api at ", *exchangeURL)
	exchangeProvider, err := exchange.NewExchangeProvider(*exchangeURL, *lcdURL, *assetPriceURL)
	panicIf(err)

	log.Info("initializing portfolio helper")
	helper, err := guildsprocess.NewPortfolioHelper(ctx, exchangeProvider, log.WithField("svc", "add_guild"))
	panicIf(err)

	spotClient := spotExchangePB.NewInjectiveSpotExchangeRPCClient(exchangeProvider.GetExchangeConn())
	derivativeClient := derivativeExchangePB.NewInjectiveDerivativeExchangeRPCClient(exchangeProvider.GetExchangeConn())

	denomIDToMinAmount := make(map[string]float64)
	srIndex := 0

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
				Address:   market.GetBaseTokenMeta().GetAddress(),
				Symbol:    market.GetBaseTokenMeta().GetSymbol(),
				Decimals:  int(market.GetBaseTokenMeta().GetDecimals()),
				UpdatedAt: market.GetBaseTokenMeta().GetUpdatedAt(),
			},
			QuoteDenom: market.GetQuoteDenom(),
			QuoteTokenMeta: &model.TokenMeta{
				Name:      market.GetQuoteTokenMeta().GetName(),
				Address:   market.GetQuoteTokenMeta().GetAddress(),
				Symbol:    market.GetQuoteTokenMeta().GetSymbol(),
				Decimals:  int(market.GetQuoteTokenMeta().GetDecimals()),
				UpdatedAt: market.GetQuoteTokenMeta().GetUpdatedAt(),
			},
			TakerFeeRate: takerFeeRate,
			MakerFeeRate: marketFeeRate,
		})

		floats, err := toMinAmounts((*spotRequirements)[srIndex])
		denomIDToMinAmount[market.GetBaseDenom()] += floats[0]
		denomIDToMinAmount[market.GetQuoteDenom()] += floats[1]

		srIndex++
	}

	drIndex := 0
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

		floats, err := toMinAmounts((*derivativeRequirements)[drIndex])
		denomIDToMinAmount[market.GetQuoteDenom()] += floats[0]
		drIndex++
	}

	denomRequirements := make([]*model.DenomRequirement, 0)
	for k, v := range denomIDToMinAmount {
		denomRequirements = append(denomRequirements, &model.DenomRequirement{
			Denom:        k,
			MinAmountUSD: v,
		})
	}

	master, _ := cosmtypes.AccAddressFromBech32(*masterAddr)
	defaultMember, _ := cosmtypes.AccAddressFromBech32(*defaultMemberAddr)
	guild := &model.Guild{
		Name:          *name,
		Description:   *description,
		MasterAddress: model.Address{AccAddress: master},

		Requirements: denomRequirements,
		Markets:      markets,
		Capacity:     *capacity,
	}

	log.Info("double check grants of default member")
	err = checkGrant(ctx, exchangeProvider, defaultMember.String(), master.String())
	panicIf(err)

	// everything goes right, add to db
	log.Info("adding guild: ", *name)
	id, err := dbSvc.AddGuild(ctx, guild)
	panicIf(err)

	log.Info("capturing default member portfolio")
	portfolio, err := helper.CaptureSingleMemberPortfolio(ctx, guild, &model.GuildMember{
		GuildID:          *id,
		InjectiveAddress: model.Address{AccAddress: defaultMember},
	}, true)
	if err != nil {
		// TODO: Use transaction
		log.Error(fmt.Sprintf("capture portfolio failed: %s. Going to revert ...", err.Error()))
		err = dbSvc.DeleteGuild(ctx, id.Hex())
		panicIf(err)

		log.Fatal("revert done. no guild added")
	}

	log.Info("adding default member")
	err = dbSvc.AddMember(ctx, id.Hex(), model.Address{AccAddress: defaultMember}, portfolio, true, *memberParams)
	if err != nil {
		// TODO: Use transaction
		log.Error(fmt.Sprintf("adding default member failed: %s. Going to revert ...", err.Error()))
		err = dbSvc.DeleteGuild(ctx, id.Hex())
		panicIf(err)

		log.Fatal("revert done. no guild added")
	}

	log.Info("update guild portfolio snapshot into db")
	// guild portfolio balance now is first member's balance
	guildPortfolio := &model.GuildPortfolio{
		GuildID:      *id,
		Balances:     portfolio.Balances,
		BankBalances: portfolio.BankBalances,
		UpdatedAt:    portfolio.UpdatedAt,
	}
	err = dbSvc.AddGuildPortfolios(ctx, []*model.GuildPortfolio{guildPortfolio})
	panicIf(err)

	log.Info("🍺 all done, guild created = ", *id)
}

func cmdAddGuild(c *cli.Cmd) {
	// inputs:
	// guild name: --name
	// guild description: --description
	// database url: --db-url
	// exchange url: --exchange-url
	// spot market (can supply many --spot-id): --spot-id
	// spot requirements: --requirement baseTokenAmountInUSD/quoteTokenAmountInUSD
	// derivative market (can supply many --derivative-id): --derivative-id
	// derivative requirements: --requirement quoteTokenAmountInUSD
	// min staking requirements: --min-staking
	// capacity: --capactity
	// master address: --master
	// default member address: --default-member
	// returns guild id
	parseAddGuildArgs(c)
	c.Action = addGuildAction
}
