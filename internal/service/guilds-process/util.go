package guildsprocess

import (
	"time"

	"github.com/InjectiveLabs/injective-guilds-service/internal/config"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db/model"
	"github.com/InjectiveLabs/injective-guilds-service/internal/exchange"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getIDToMarket(g *model.Guild) map[string]*model.GuildMarket {
	result := make(map[string]*model.GuildMarket)
	for _, m := range g.Markets {
		result[m.MarketID.String()] = m
	}
	return result
}

func signOf(direction string) decimal.Decimal {
	if direction == DirectionShort {
		return decimal.NewFromInt(-1)
	}
	return decimal.NewFromInt(1)
}

func buildPortfolio(
	member *model.GuildMember,
	balances []*exchange.Balance,
	injBalance []*model.BankBalance,
	pnl map[string]decimal.Decimal,
	marginHolds map[string]decimal.Decimal,
	usdPrices map[string]float64,
) *model.AccountPortfolio {
	portfolio := &model.AccountPortfolio{
		GuildID:          member.GuildID,
		InjectiveAddress: member.InjectiveAddress,
		UpdatedAt:        time.Now(),
	}

	for _, b := range balances {
		pnlValue, ok := pnl[b.Denom]
		if !ok {
			pnlValue = decimal.Zero
		}

		marginHoldValue, ok := marginHolds[b.Denom]
		if !ok {
			marginHoldValue = decimal.Zero
		}

		aBalance := &model.Balance{
			Denom: b.Denom,
		}

		aBalance.TotalBalance, _ = primitive.ParseDecimal128(b.TotalBalance.String())
		aBalance.AvailableBalance, _ = primitive.ParseDecimal128(b.AvailableBalance.String())
		aBalance.UnrealizedPNL, _ = primitive.ParseDecimal128(pnlValue.StringFixed(18))
		aBalance.MarginHold, _ = primitive.ParseDecimal128(marginHoldValue.StringFixed(18))
		if usdPrices != nil {
			// TODO: Impl historical price on asset-price service to recompute it on UI
			aBalance.PriceUSD = usdPrices[b.Denom]
		}

		portfolio.Balances = append(portfolio.Balances, aBalance)
	}

	// currently only track INJ balance
	// use bank balance model for future, if we want to track balance in other denoms
	// TODO: Impl historical price on asset-price service to recompute it on UI
	if usdPrices != nil {
		for _, b := range injBalance {
			if b.Denom == config.DEMOM_INJ {
				b.PriceUSD = usdPrices[config.DEMOM_INJ]
			}
		}
	}
	portfolio.BankBalances = injBalance

	return portfolio
}

func marketsFromGuild(guild *model.Guild, isPerp bool) []string {
	var result []string
	for _, m := range guild.Markets {
		if m.IsPerpetual == isPerp {
			result = append(result, m.MarketID.Hex())
		}
	}
	return result
}

func defaultSubaccountIDFromInjAddress(injAddress model.Address) string {
	ethAddr := common.BytesToAddress(injAddress.Bytes())
	return ethAddr.Hex() + "000000000000000000000000"
}

func sum(a primitive.Decimal128, b primitive.Decimal128) primitive.Decimal128 {
	parsedA, _ := decimal.NewFromString(a.String())
	parsedB, _ := decimal.NewFromString(b.String())
	result, _ := primitive.ParseDecimal128(parsedA.Add(parsedB).String())
	return result
}
