package guildsprocess

import (
	"context"
	"sort"
	"testing"

	"github.com/InjectiveLabs/injective-guilds-service/internal/db/model"
	"github.com/InjectiveLabs/injective-guilds-service/internal/exchange"
	cosmtypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	log "github.com/xlab/suplog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	marketInjUstID   = "0x9b9980167ecc3645ff1a5517886652d94a0825e54a77d2057cbbe3ebee015963"
	marketWethUsdtID = "0x54d4505adef6a5cef26bc403a33d595620ded4e15b9e2bc3dd489b714813366a"

	denomInj  = "inj"
	denomUst  = "ibc/B448C0CA358B958301D328CCDC5D5AD642FC30A6D3AE106FF721DB315F3DDE5C"
	denomWeth = "peggy0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	denomUsdt = "peggy0xdAC17F958D2ee523a2206206994597C13D831ec7"
)

func parseDecimal128(s string) primitive.Decimal128 {
	p, _ := primitive.ParseDecimal128(s)
	return p
}

func TestCaptureSinglePortfolio(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	mockExchange := exchange.NewMockDataProvider(ctrl)
	cosmtypes.GetConfig().SetBech32PrefixForAccount("inj", "injpub")

	helper, err := NewPortfolioHelper(ctx, mockExchange, log.WithField("svc", "test"))
	assert.NoError(t, err)

	subaccountID := "0xEB8cf88b739fE12E303E31fb88fC37751E17cF3D000000000000000000000000"
	accAddress, _ := cosmtypes.AccAddressFromBech32("inj1awx03zmnnlsjuvp7x8ac3lphw50p0nea6p2584")
	masterAddress := "inj1wng2ucn0ak3aw5gq9j7m2z88m5aznwntqnekuv"

	subaccountBalances := []*exchange.Balance{
		{
			Denom:            denomInj,
			TotalBalance:     decimal.RequireFromString("2000000000000000"),
			AvailableBalance: decimal.RequireFromString("2000000000000000"),
		},
		{
			Denom:            denomUsdt,
			TotalBalance:     decimal.RequireFromString("50000000"),
			AvailableBalance: decimal.RequireFromString("50000000"),
		},
		{
			Denom:            denomWeth,
			TotalBalance:     decimal.RequireFromString("1000000000000000000"),
			AvailableBalance: decimal.RequireFromString("500000000000000000"),
		},
		{
			Denom:            denomUst,
			TotalBalance:     decimal.RequireFromString("30000000"),
			AvailableBalance: decimal.RequireFromString("30000000"),
		},
	}
	spotOrders := []*exchange.SpotOrder{
		{
			MarketID:         marketInjUstID,
			FeeRecipient:     masterAddress,
			OrderSide:        OrderSideSell,
			Price:            decimal.RequireFromString("0.000000056"),
			UnfilledQuantity: decimal.RequireFromString("120000000000000"),
		},
		{
			MarketID:         marketInjUstID,
			FeeRecipient:     masterAddress,
			OrderSide:        OrderSideBuy,
			Price:            decimal.RequireFromString("0.000000032"),
			UnfilledQuantity: decimal.RequireFromString("230000000000000"),
		},
		{
			MarketID:         marketInjUstID,
			FeeRecipient:     masterAddress,
			OrderSide:        OrderSideBuy,
			Price:            decimal.RequireFromString("0.000000016"),
			UnfilledQuantity: decimal.RequireFromString("240000000000000"),
		},
	}
	derivativeOrders := []*exchange.DerivativeOrder{
		{
			MarketID:     marketWethUsdtID,
			FeeRecipient: masterAddress,
			Margin:       decimal.RequireFromString("11200000000"),
		},
		{
			MarketID:     marketWethUsdtID,
			FeeRecipient: masterAddress,
			Margin:       decimal.RequireFromString("11400000000"),
		},
	}
	bankBalances := &exchange.BankAccountBalances{
		Balances: []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		}{
			{
				Denom:  denomInj,
				Amount: "500000000000000000",
			},
		},
	}
	positions := []*exchange.DerivativePosition{
		{
			MarketID:   marketWethUsdtID,
			Direction:  DirectionLong,
			Quantity:   decimal.RequireFromString("0.5"),
			Margin:     decimal.RequireFromString("15000000000"),
			EntryPrice: decimal.RequireFromString("3500000000"),
			MarkPrice:  decimal.RequireFromString("3200000000"),
		},
		{
			MarketID:   marketWethUsdtID,
			Direction:  DirectionLong,
			Quantity:   decimal.RequireFromString("0.8"),
			Margin:     decimal.RequireFromString("15000000000"),
			EntryPrice: decimal.RequireFromString("2800000000"),
			MarkPrice:  decimal.RequireFromString("3200000000"),
		},
	}
	priceUSD := []*exchange.CoinPrice{
		{
			ID:           "injective-protocol",
			CurrentPrice: 6.0,
		},
		{
			ID:           "tether",
			CurrentPrice: 1,
		},
		{
			ID:           "terrausd",
			CurrentPrice: 1,
		},
		{
			ID:           "weth",
			CurrentPrice: 3300,
		},
	}

	mockExchange.EXPECT().GetSubaccountBalances(gomock.Any(), subaccountID).Return(subaccountBalances, nil).Times(1)
	mockExchange.EXPECT().GetSpotOrders(gomock.Any(), []string{marketInjUstID}, subaccountID).
		Return(spotOrders, nil).Times(1)
	mockExchange.EXPECT().GetDerivativeOrders(gomock.Any(), []string{marketWethUsdtID}, subaccountID).
		Return(derivativeOrders, nil).Times(1)
	mockExchange.EXPECT().GetPositions(gomock.Any(), subaccountID).
		Return(positions, nil).Times(1)
	mockExchange.EXPECT().GetBankBalance(gomock.Any(), accAddress.String()).
		Return(bankBalances, nil).Times(1)
	// TODO: element order
	mockExchange.EXPECT().GetPriceUSD(gomock.Any(), gomock.Any()).Return(priceUSD, nil).Times(1)

	takerFee, _ := primitive.ParseDecimal128("0.001")
	portfolio, err := helper.CaptureSingleMemberPortfolio(
		ctx,
		&model.Guild{
			Markets: []*model.GuildMarket{
				{
					IsPerpetual: false,
					MarketID:    model.Hash{Hash: common.HexToHash(marketInjUstID)},
					BaseDenom:   denomInj,
					BaseTokenMeta: &model.TokenMeta{
						Decimals: 15,
					},
					QuoteDenom: denomUst,
					QuoteTokenMeta: &model.TokenMeta{
						Decimals: 6,
					},
					TakerFeeRate: takerFee,
				},
				{
					IsPerpetual: true,
					MarketID:    model.Hash{Hash: common.HexToHash(marketWethUsdtID)},
					QuoteDenom:  denomUsdt,
				},
			},
		},
		&model.GuildMember{
			InjectiveAddress: model.Address{AccAddress: accAddress},
		}, true,
	)
	assert.NoError(t, err)
	assert.NotNil(t, portfolio)
	sort.Slice(portfolio.Balances, func(i, j int) bool {
		return portfolio.Balances[i].Denom < portfolio.Balances[j].Denom
	})

	assert.Equal(t, model.Address{AccAddress: accAddress}, portfolio.InjectiveAddress)

	assert.Equal(t, []*model.Balance{
		{
			Denom:            "ibc/B448C0CA358B958301D328CCDC5D5AD642FC30A6D3AE106FF721DB315F3DDE5C",
			PriceUSD:         1,
			TotalBalance:     parseDecimal128("30000000"),
			AvailableBalance: parseDecimal128("30000000"),
			MarginHold:       parseDecimal128("11211200"),
			UnrealizedPNL:    parseDecimal128("0"),
		},
		{
			Denom:            "inj",
			PriceUSD:         6.0,
			TotalBalance:     parseDecimal128("2000000000000000"),
			AvailableBalance: parseDecimal128("2000000000000000"),
			MarginHold:       parseDecimal128("120000000000000"),
			UnrealizedPNL:    parseDecimal128("0"),
		},
		{
			Denom:            "peggy0xdAC17F958D2ee523a2206206994597C13D831ec7",
			PriceUSD:         1,
			TotalBalance:     parseDecimal128("50000000"),
			AvailableBalance: parseDecimal128("50000000"),
			MarginHold:       parseDecimal128("52600000000"),
			UnrealizedPNL:    parseDecimal128("170000000"),
		},
	}, portfolio.Balances)

	assert.Equal(t, []*model.BankBalance{
		{
			Denom:    "inj",
			PriceUSD: 6.0,
			Balance:  parseDecimal128("500000000000000000"),
		},
	}, portfolio.BankBalances)

	// calcuation steps

	// peggy0xdAC17F958D2ee523a2206206994597C13D831ec7 // usdt
	// Total balance: 50000000
	// Available balance: 50000000
	// Margin hold: 52600000000 = 11400000000 + 11200000000 + 15000000000 + 15000000000 (position margins + orders margin)
	// Unrealized PNL: 170000000 = (3200000000 - 3500000000) * 0.5 + (3200000000 - 2800000000) * 0.8

	// inj
	// Total balance: 2000000000000000
	// Available balance: 2000000000000000
	// Margin hold: 120000000000000 (a sell order)
	// UnrealizedPNL: 0 // no position

	// ibc/B448C0CA358B958301D328CCDC5D5AD642FC30A6D3AE106FF721DB315F3DDE5C // ust
	// Total balance: 30000000
	// Available balance: 30000000
	// Margin hold: 11211200 # 240000000000000 * 0,000000016 * (1 + 0,001) + 0,000000032 * 230000000000000 * (1 + 0,001)
	// UnrealizedPNL: 0
}
