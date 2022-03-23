package guildsprocess

import (
	"context"
	"errors"
	"fmt"

	"github.com/InjectiveLabs/injective-guilds-service/internal/db"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db/model"
	"github.com/InjectiveLabs/injective-guilds-service/internal/exchange"
	"github.com/shopspring/decimal"
	log "github.com/xlab/suplog"
)

// PortfolioHelper supports capture account portfolio (default subaccount)
type PortfolioHelper struct {
	denomToCoinID    map[string]string
	dbSvc            db.DBService
	exchangeProvider exchange.DataProvider
	logger           log.Logger
}

func NewPortfolioHelper(
	ctx context.Context,
	dbSvc db.DBService,
	provider exchange.DataProvider,
) (*PortfolioHelper, error) {
	helper := &PortfolioHelper{
		dbSvc:            dbSvc,
		exchangeProvider: provider,
	}

	if err := helper.UpdateDenomToCoinIDMap(ctx); err != nil {
		return nil, err
	}

	return helper, nil
}

func (p *PortfolioHelper) CaptureSingleMemberPortfolio(
	ctx context.Context,
	guild *model.Guild,
	member *model.GuildMember,
	addDenomPrices bool,
) (*model.AccountPortfolio, error) {
	denoms := model.GetGuildDenoms(guild)
	defaultSubaccountID := defaultSubaccountIDFromInjAddress(member.InjectiveAddress)

	// get balances
	balances, err := p.getSubaccountBalances(ctx, denoms, defaultSubaccountID)
	if err != nil {
		return nil, fmt.Errorf("get balance error: %w", err)
	}

	// get all positions
	positions, err := p.exchangeProvider.GetPositions(ctx, defaultSubaccountID)
	if err != nil {
		// TODO: Put intro retry queue
		return nil, fmt.Errorf("get position err: %w", err)
	}

	// compute pnl
	pnl := p.getUnrealizedPNL(guild, positions)

	// compute margin hold
	marginHold, err := p.getMarginHold(ctx, guild, positions, member)
	if err != nil {
		return nil, fmt.Errorf("get margin hold err: %w", err)
	}

	var prices map[string]float64
	if addDenomPrices {
		prices, err = p.GetDenomPrices(ctx, model.GetGuildDenoms(guild))
		if err != nil {
			return nil, fmt.Errorf("get margin hold err: %w", err)
		}
	}

	return buildPortfolio(member, balances, pnl, marginHold, prices), nil
}

func (p *PortfolioHelper) getSubaccountBalances(ctx context.Context, denoms []string, defaultSubaccountID string) (result []*exchange.Balance, err error) {
	balances, err := p.exchangeProvider.GetSubaccountBalances(ctx, defaultSubaccountID)
	if err != nil {
		return nil, err
	}

	denomToBalance := make(map[string]*exchange.Balance)
	for _, b := range balances {
		denomToBalance[b.Denom] = b
	}

	// filter denoms, add 0 if such denom not exists
	for _, denom := range denoms {
		b, exist := denomToBalance[denom]
		if !exist {
			result = append(result, &exchange.Balance{
				Denom:            denom,
				TotalBalance:     decimal.Zero,
				AvailableBalance: decimal.Zero,
			})
			continue
		}
		result = append(result, b)
	}

	return result, nil
}

// getMemberUnrealizedPNL returns pnl[denom]decimal.Decimal
// pnl in quoteDenom for now (but assume that usdc, usdt are around $1)
func (p *PortfolioHelper) getUnrealizedPNL(
	guild *model.Guild,
	positions []*exchange.DerivativePosition,
) map[string]decimal.Decimal {
	idToMarket := getIDToMarket(guild)
	pnl := make(map[string]decimal.Decimal)
	for _, pos := range positions {
		market, exist := idToMarket[pos.MarketID]
		if !exist {
			continue
		}
		quoteDenom := market.QuoteDenom
		// pnl[quoteDenom] += (markPrice - entryPrice) * quantity * signOf(direction)
		pnl[quoteDenom] = pnl[quoteDenom].Add(pos.MarkPrice.Sub(pos.EntryPrice).Mul(pos.Quantity).Mul(signOf(pos.Direction)))
	}

	return pnl
}

func (p *PortfolioHelper) getMarginHold(
	ctx context.Context,
	guild *model.Guild,
	positions []*exchange.DerivativePosition,
	member *model.GuildMember,
) (marginHolds map[string]decimal.Decimal, err error) {
	defaultSubaccountID := defaultSubaccountIDFromInjAddress(member.InjectiveAddress)

	idToMarket := getIDToMarket(guild)
	// Bojan: marginHold = sumOf(positions.margin) + sumOf(orders.margin) where orders.margin = notionalValue + fees
	marginHolds = make(map[string]decimal.Decimal)
	for _, p := range positions {
		market, exist := idToMarket[p.MarketID]
		if !exist {
			continue
		}

		quoteDenom := market.QuoteDenom
		marginHolds[quoteDenom] = marginHolds[quoteDenom].Add(p.Margin)
	}

	// margins from derivative orders
	derivOrders, err := p.exchangeProvider.GetDerivativeOrders(ctx, marketsFromGuild(guild, true), defaultSubaccountID)
	if err != nil {
		p.logger.WithError(err).Errorln("cannot get derivaitve orders")
		return nil, err
	}

	for _, o := range derivOrders {
		market, exist := idToMarket[o.MarketID]
		if !exist {
			// TODO: Optimization: Put into queue to remove this person
			// Reason: we don't support market which is not in guild
			continue
		}

		// we only have marginHold of quoteDenom in perp markets
		quoteDenom := market.QuoteDenom
		marginHolds[quoteDenom] = marginHolds[quoteDenom].Add(o.Margin)
	}

	// margins from spot orders
	spotOrders, err := p.exchangeProvider.GetSpotOrders(ctx, marketsFromGuild(guild, false), defaultSubaccountID)
	if err != nil {
		p.logger.WithError(err).Errorln("cannot get spot orders")
		return nil, err
	}

	// Albert:
	// for a limit buy in the ETH/USDT market, denom is USDT and balanceHold is (1 + takerFee)*(price * quantity)
	// for a limit sell in the ETH/USDT market, denom is ETH and balanceHold is just quantity
	for _, o := range spotOrders {
		market, exist := idToMarket[o.MarketID]
		if !exist {
			// TODO: Optimization: Put into queue to remove this person
			// Reason: we don't support market which is not in guild
			continue
		}

		if o.OrderSide == OrderSideBuy {
			// expected to parse successfully -> skip error
			takerFee, _ := decimal.NewFromString(market.TakerFeeRate.String())
			fee := takerFee.Mul(o.Price).Mul(o.UnfilledQuantity)
			margin := o.Price.Mul(o.UnfilledQuantity)

			// price * unfilled * (1 + takerFee)
			marginHolds[market.QuoteDenom] = marginHolds[market.QuoteDenom].Add(margin.Add(fee))
		} else {
			marginHolds[market.BaseDenom] = marginHolds[market.BaseDenom].Add(o.UnfilledQuantity)
		}
	}

	return marginHolds, nil
}

// getDenomPrices returns map[denom]priceInUSD
func (p *PortfolioHelper) GetDenomPrices(ctx context.Context, denoms []string) (map[string]float64, error) {
	result := make(map[string]float64)
	coinIDs := make([]string, 0)

	for _, d := range denoms {
		id, exist := p.denomToCoinID[d]
		if !exist {
			return nil, errors.New("not all denoms have coinIDs")
		}

		coinIDs = append(coinIDs, id)
	}

	prices, err := p.exchangeProvider.GetPriceUSD(ctx, coinIDs)
	if err != nil {
		return nil, err
	}

	for _, d := range denoms {
		id := p.denomToCoinID[d]
		found := false
		for _, price := range prices {
			if price.ID == id {
				result[d] = price.CurrentPrice
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("coin id have no price: %s", id)
		}
	}

	return result, nil
}

// basically, we need Denom -> CoinID to get coin price
// decide to store usd price in db at the moment, will use historical price later
func (p *PortfolioHelper) UpdateDenomToCoinIDMap(ctx context.Context) error {
	denomCoinID, err := p.dbSvc.ListDenomCoinID(ctx)
	if err != nil {
		return err
	}

	// add a check to make sure denoms are added
	if len(denomCoinID) == 0 {
		return errors.New("cannot find denom to coin-id, should run update-denom first")
	}

	p.denomToCoinID = make(map[string]string)
	for _, d := range denomCoinID {
		p.denomToCoinID[d.Denom] = d.CoinID
	}
	return nil
}
