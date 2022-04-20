package guildsapi

import (
	"math"
	"time"

	svc "github.com/InjectiveLabs/injective-guilds-service/api/gen/guilds_service"
	"github.com/InjectiveLabs/injective-guilds-service/internal/config"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db/model"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemberMessage struct {
	Action    string `json:"action"`
	ExpiredAt int64  `json:"expired_at"` // unix timestamp, second
}

type Period struct {
	StartTime time.Time
	EndTime   time.Time
}

func modelPortfolioToHTTP(p *model.AccountPortfolio) *svc.SingleAccountPortfolio {
	portfolio := p.Copy()

	if len(portfolio.BankBalances) > 0 && p.BankBalances[0].Denom == config.DEMOM_INJ {
		portfolio.Balances = addInjBankToBalance(p.Balances, p.BankBalances[0])
	}

	var balances []*svc.Balance
	for _, b := range portfolio.Balances {
		balances = append(balances, &svc.Balance{
			Denom:            b.Denom,
			PriceUsd:         b.PriceUSD,
			TotalBalance:     b.TotalBalance.String(),
			AvailableBalance: b.AvailableBalance.String(),
			UnrealizedPnl:    b.UnrealizedPNL.String(),
			MarginHold:       b.MarginHold.String(),
		})
	}
	return &svc.SingleAccountPortfolio{
		InjectiveAddress: p.InjectiveAddress.String(),
		Balances:         balances,
		UpdatedAt:        p.UpdatedAt.UnixMilli(),
	}
}

// list timestamp [startTime, ceilToMonth(endTime))
func monthlyPeriods(startTime, endTime time.Time) (result []*Period) {
	current := startTime
	endTime = endTime.AddDate(0, 1, 0)

	times := make([]time.Time, 0)
	for current.Before(endTime) {
		times = append(times, current)
		beginOfMonth := time.Date(current.Year(), current.Month(), 1, 0, 0, 0, 0, current.Location())
		current = beginOfMonth.AddDate(0, 1, 0)
	}

	for i := 0; i < len(times)-1; i++ {
		result = append(result, &Period{
			StartTime: times[i],
			EndTime:   times[i+1],
		})
	}

	return result
}

func addInjBankToBalance(balance []*model.Balance, inj *model.BankBalance) []*model.Balance {
	for _, b := range balance {
		if b.Denom == config.DEMOM_INJ {
			b.TotalBalance = sum(b.TotalBalance, inj.Balance)
			b.AvailableBalance = sum(b.AvailableBalance, inj.Balance)
			return balance
		}
	}

	// if not found then append inj denom
	balance = append(balance, &model.Balance{
		Denom:            config.DEMOM_INJ,
		PriceUSD:         inj.PriceUSD,
		TotalBalance:     inj.Balance,
		AvailableBalance: inj.Balance,
	})
	return balance
}

func modelGuildToResponse(m *model.Guild, portfolio *model.GuildPortfolio, defaultMember *model.GuildMember) *svc.Guild {
	var (
		requirements    []*svc.Requirement
		balances        []*svc.Balance
		denomToUsdPrice = make(map[string]float64)
	)

	if len(portfolio.BankBalances) > 0 && portfolio.BankBalances[0].Denom == config.DEMOM_INJ {
		portfolio.Balances = addInjBankToBalance(portfolio.Balances, portfolio.BankBalances[0])
	}

	for _, b := range portfolio.Balances {
		balances = append(balances, &svc.Balance{
			Denom:            b.Denom,
			TotalBalance:     b.TotalBalance.String(),
			AvailableBalance: b.AvailableBalance.String(),
			UnrealizedPnl:    b.UnrealizedPNL.String(),
			MarginHold:       b.MarginHold.String(),
			PriceUsd:         b.PriceUSD,
		})

		denomToUsdPrice[b.Denom] = b.PriceUSD
	}

	for _, req := range m.Requirements {
		displayDecimal := config.DenomConfigs[req.Denom].DisplayDecimal

		var priceUsd float64
		if _, isStableCoin := config.StableCoinDenoms[req.Denom]; isStableCoin {
			priceUsd = 1
		} else {
			priceUsd = denomToUsdPrice[req.Denom]
		}

		roundedFloat := math.Ceil(req.MinAmountUSD*math.Pow10(displayDecimal)/priceUsd) / math.Pow10(displayDecimal)
		// IMPORTANT: We want to return a min requirement, so that FE and BE will be sync-ed
		requirements = append(requirements, &svc.Requirement{
			Denom:        req.Denom,
			MinAmountUsd: req.MinAmountUSD,
			MinAmount:    roundedFloat,
		})
	}

	var currentPortfolio *svc.SingleGuildPortfolio
	if len(balances) != 0 {
		currentPortfolio = &svc.SingleGuildPortfolio{
			Balances:  balances,
			UpdatedAt: portfolio.UpdatedAt.UnixMilli(),
		}
	}

	return &svc.Guild{
		ID:                   m.ID.Hex(),
		Name:                 m.Name,
		Description:          m.Description,
		MasterAddress:        m.MasterAddress.String(),
		Requirements:         requirements,
		Capacity:             m.Capacity,
		MemberCount:          m.MemberCount,
		CurrentPortfolio:     currentPortfolio,
		DefaultMemberAddress: defaultMember.InjectiveAddress.String(),
	}
}

func sum(a primitive.Decimal128, b primitive.Decimal128) primitive.Decimal128 {
	parsedA, _ := decimal.NewFromString(a.String())
	parsedB, _ := decimal.NewFromString(b.String())
	result, _ := primitive.ParseDecimal128(parsedA.Add(parsedB).String())
	return result
}

func defaultSubaccountIDFromInjAddress(injAddress model.Address) string {
	ethAddr := common.BytesToAddress(injAddress.Bytes())
	return ethAddr.Hex() + "000000000000000000000000"
}
