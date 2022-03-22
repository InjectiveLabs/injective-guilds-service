package guildsprocess

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/InjectiveLabs/injective-guilds-service/internal/config"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db/model"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db/mongoimpl"
	"github.com/InjectiveLabs/injective-guilds-service/internal/exchange"
	cosmtypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	log "github.com/xlab/suplog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	DirectionLong  = "long"
	DirectionShort = "short"

	OrderSideBuy  = "buy"
	OrderSideSell = "sell"
)

type GuildsProcess struct {
	dbSvc    db.DBService
	exchange exchange.DataProvider
	logger   log.Logger

	portfolioUpdateInterval time.Duration
	disqualifyInterval      time.Duration
}

func NewProcess(cfg config.GuildProcessConfig) (*GuildsProcess, error) {
	ctx := context.Background()
	dbService, err := mongoimpl.NewService(ctx, cfg.DBConnectionURL, cfg.DBName)
	if err != nil {
		return nil, err
	}

	// make index
	if err := dbService.(*mongoimpl.MongoImpl).EnsureIndex(ctx); err != nil {
		log.WithError(err).Warningln("cannot ensure index")
	}

	// won't use lcd endpoint here
	exchangeProvider, err := exchange.NewExchangeProvider(cfg.ExchangeGRPCURL, "", cfg.AssetPriceURL)
	if err != nil {
		return nil, err
	}

	cosmtypes.GetConfig().SetBech32PrefixForAccount("inj", "injpub")

	return &GuildsProcess{
		dbSvc:                   dbService,
		exchange:                exchangeProvider,
		logger:                  log.WithField("svc", "guilds_process"),
		portfolioUpdateInterval: cfg.PortfolioUpdateInterval,
		disqualifyInterval:      cfg.DisqualifyInterval,
	}, nil
}

func (p *GuildsProcess) Run(ctx context.Context) {
	// run 2 cron jobs
	go p.runWithInterval(ctx, p.portfolioUpdateInterval, func(ctx context.Context) error {
		return p.capturePorfolioSnapshot(ctx)
	})

	go p.runWithInterval(ctx, p.disqualifyInterval, func(ctx context.Context) error {
		return p.processDisqualification(ctx)
	})
}

func (p *GuildsProcess) runWithInterval(ctx context.Context, interval time.Duration, fn func(ctx context.Context) error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			timeMarker := time.Now()
			err := fn(ctx)
			if err != nil {
				p.logger.WithError(err).Errorln("error while doing cronjob")
			}
			timeElasped := time.Since(timeMarker)

			if timeElasped < interval {
				time.Sleep(interval - timeElasped)
			}
		}
	}
}

// TODO: Improvement: Implement retryable mechanism
func (p *GuildsProcess) capturePorfolioSnapshot(ctx context.Context) error {
	guilds, err := p.dbSvc.ListAllGuilds(ctx)
	if err != nil {
		return fmt.Errorf("list guild err: %w", err)
	}

	now := time.Now()
	for _, guild := range guilds {
		guildID := guild.ID.Hex()

		members, err := p.dbSvc.ListGuildMembers(ctx, model.MemberFilter{
			GuildID: &guildID,
		})
		if err != nil {
			err = fmt.Errorf("list non-default member err: %w", err)
			p.logger.
				WithField("guildID", guildID).
				WithError(err).Warningln("skip this guild")
			continue
		}

		// for each guild, we get all denom prices once
		// eliminate failure + save time
		priceMap, err := p.getDenomPrices(ctx, model.GetGuildDenoms(guild))
		if err != nil {
			// should return an error
			return err
		}

		portfolios := make([]*model.AccountPortfolio, 0)

		// TODO: Create queue to re-try when failure happens
		// TODO: Create bulk accounts balances query on injective-exchange
		for _, member := range members {
			portfolioSnapshot, err := p.captureSingleMemberPortfolio(ctx, guild, member)
			if err != nil {
				p.logger.
					WithField("guildID", guildID).
					WithField("memberAddr", member.InjectiveAddress.String()).
					WithError(err).Warningln("capture snapshot error")
				continue
			}

			// fill denom price
			// priceMap has all denom prices
			for _, b := range portfolioSnapshot.Balances {
				b.PriceUSD = priceMap[b.Denom]
			}

			portfolioSnapshot.UpdatedAt = now
			portfolios = append(portfolios, portfolioSnapshot)
		}

		err = p.dbSvc.AddAccountPortfolios(ctx, guildID, portfolios)
		if err != nil {
			p.logger.
				WithField("guildID", guildID).
				WithError(err).Warningln("skip this guild")
		}
	}
	return nil
}

// getDenomPrices returns map[denom]priceInUSD
func (p *GuildsProcess) getDenomPrices(ctx context.Context, denoms []string) (map[string]float64, error) {
	result := make(map[string]float64)
	coinIDs := make([]string, 0)

	for _, d := range denoms {
		id, exist := exchange.DenomToCoinID[d]
		if !exist {
			p.logger.WithField("denom", d).Warning("coinID not found")
			return nil, errors.New("not all denoms have coinIDs")
		}

		coinIDs = append(coinIDs, id)
	}

	prices, err := p.exchange.GetPriceUSD(ctx, coinIDs)
	if err != nil {
		return nil, err
	}

	for _, d := range denoms {
		id := exchange.DenomToCoinID[d]
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

func defaultSubaccountIDFromInjAddress(injAddress model.Address) string {
	ethAddr := common.BytesToAddress(injAddress.Bytes())
	return ethAddr.Hex() + "000000000000000000000000"
}

func (p *GuildsProcess) captureSingleMemberPortfolio(
	ctx context.Context,
	guild *model.Guild,
	member *model.GuildMember,
) (*model.AccountPortfolio, error) {
	denoms := model.GetGuildDenoms(guild)
	defaultSubaccountID := defaultSubaccountIDFromInjAddress(member.InjectiveAddress)

	// get balances
	balances, err := p.getSubaccountBalances(ctx, denoms, defaultSubaccountID)
	if err != nil {
		return nil, err
	}

	// get all positions
	positions, err := p.exchange.GetPositions(ctx, defaultSubaccountID)
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

	return buildPortfolio(member, balances, pnl, marginHold), nil
}

func buildPortfolio(
	member *model.GuildMember,
	balances []*exchange.Balance,
	pnl map[string]decimal.Decimal,
	marginHolds map[string]decimal.Decimal,
) *model.AccountPortfolio {
	portfolio := &model.AccountPortfolio{
		GuildID: member.GuildID,
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
		aBalance.UnrealizedPNL, _ = primitive.ParseDecimal128(pnlValue.String())
		aBalance.MarginHold, _ = primitive.ParseDecimal128(marginHoldValue.String())
		portfolio.Balances = append(portfolio.Balances, aBalance)
	}

	return portfolio
}

func (p *GuildsProcess) getSubaccountBalances(ctx context.Context, denoms []string, defaultSubaccountID string) (result []*exchange.Balance, err error) {
	balances, err := p.exchange.GetSubaccountBalances(ctx, defaultSubaccountID)
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

// getMemberUnrealizedPNL returns pnl[denom]decimal.Decimal
// pnl in quoteDenom for now (but assume that usdc, usdt are around $1)
func (p *GuildsProcess) getUnrealizedPNL(
	guild *model.Guild,
	positions []*exchange.DerivativePosition,
) map[string]decimal.Decimal {
	idToMarket := getIDToMarket(guild)
	pnl := make(map[string]decimal.Decimal)
	for _, p := range positions {
		market, exist := idToMarket[p.MarketID]
		if !exist {
			continue
		}
		quoteDenom := market.QuoteDenom
		// pnl[quoteDenom] += (markPrice - entryPrice) * quantity * signOf(direction)
		pnl[quoteDenom] = pnl[quoteDenom].Add(p.MarkPrice.Sub(p.EntryPrice).Mul(p.Quantity).Mul(signOf(p.Direction)))
	}

	return pnl
}

func (p *GuildsProcess) getMarginHold(
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
	derivOrders, err := p.exchange.GetDerivativeOrders(ctx, marketsFromGuild(guild, true), defaultSubaccountID)
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
	spotOrders, err := p.exchange.GetSpotOrders(ctx, marketsFromGuild(guild, false), defaultSubaccountID)
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

// processDisqualify get orders from guild's markets
// then remove members whose orders' fee_recipient is not masterAddress
func (p *GuildsProcess) processDisqualification(ctx context.Context) error {
	guilds, err := p.dbSvc.ListAllGuilds(ctx)
	if err != nil {
		return fmt.Errorf("list guild err: %w", err)
	}

	// TODO: clarify disqualify reason
	for _, g := range guilds {
		guildID := g.ID.Hex()

		members, err := p.dbSvc.ListGuildMembers(ctx, model.MemberFilter{
			GuildID: &guildID,
		})
		if err != nil {
			err = fmt.Errorf("list non-default member err: %w", err)
			p.logger.
				WithField("guildID", guildID).
				WithError(err).Warningln("skip this guild")
			continue
		}

		markets := make([]string, 0)
		for _, market := range g.Markets {
			markets = append(markets, market.MarketID.Hex())
		}

		for _, member := range members {
			disqualify, err := p.shouldDisqualify(ctx, g, member.InjectiveAddress)
			if err != nil {
				continue
			}
			// we don't expect this to regularly happen,
			// so decided to delete each document this way
			if disqualify {
				err = p.dbSvc.RemoveMember(ctx, g.ID.Hex(), member.InjectiveAddress)
				if err != nil {
					log.WithField("memberAddress", member.InjectiveAddress.String()).
						WithError(err).Errorln("cannot delete member")
				}
			}
		}
	}
	return nil
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

// shouldDisqualify tries to disqualify a person if deriv/spot orders has fee recipient != master address
func (p *GuildsProcess) shouldDisqualify(
	ctx context.Context,
	guild *model.Guild,
	address model.Address,
) (bool, error) {
	defaultSubaccountID := defaultSubaccountIDFromInjAddress(address)
	spotOrders, err := p.exchange.GetSpotOrders(
		ctx, marketsFromGuild(guild, false),
		defaultSubaccountID,
	)
	if err != nil {
		p.logger.WithField("subaccountID", defaultSubaccountID).
			WithError(err).Warningln("get spot orders error")
	}

	for _, o := range spotOrders {
		if o.FeeRecipient != guild.MasterAddress.String() {
			return true, nil
		}
	}

	derivativeOrders, err := p.exchange.GetDerivativeOrders(
		ctx, marketsFromGuild(guild, true),
		defaultSubaccountID,
	)
	if err != nil {
		p.logger.WithField("subaccountID", defaultSubaccountID).
			WithError(err).Warningln("get derivative orders error")
		return false, err
	}

	for _, o := range derivativeOrders {
		if o.FeeRecipient != guild.MasterAddress.String() {
			return true, nil
		}
	}

	return false, nil
}

func (p *GuildsProcess) GracefullyShutdown(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// close db
	log.Info("closing db connection")
	if err := p.dbSvc.Disconnect(shutdownCtx); err != nil {
		log.WithError(err).Error("cannot close db connection")
	}

	// close exchange grpc
	log.Info("closing exchange grpc connection")
	if err := p.exchange.Close(); err != nil {
		log.WithError(err).Error("cannot close exchange grpc connection")
	}
}
