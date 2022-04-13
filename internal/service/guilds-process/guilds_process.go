package guildsprocess

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/InjectiveLabs/injective-guilds-service/internal/config"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db/model"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db/mongoimpl"
	"github.com/InjectiveLabs/injective-guilds-service/internal/exchange"
	metrics "github.com/InjectiveLabs/metrics"
	log "github.com/xlab/suplog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	expirationTimeLayout = "2006-01-02T15:04:05Z"

	DirectionLong  = "long"
	DirectionShort = "short"

	OrderSideBuy  = "buy"
	OrderSideSell = "sell"
)

type GuildsProcess struct {
	dbSvc           db.DBService
	exchange        exchange.DataProvider
	logger          log.Logger
	portfolioHelper *PortfolioHelper

	portfolioUpdateInterval time.Duration
	disqualifyInterval      time.Duration

	grants  []string
	svcTags metrics.Tags
}

func NewProcess(cfg config.GuildProcessConfig) (*GuildsProcess, error) {
	ctx := context.Background()
	logger := log.WithField("svc", "guilds_process")

	logger.Infoln("initializing db service")
	dbService, err := mongoimpl.NewService(ctx, cfg.DBConnectionURL, cfg.DBName)
	if err != nil {
		return nil, err
	}

	logger.Infoln("ensuring db indexes...")
	// make index
	if err := dbService.(*mongoimpl.MongoImpl).EnsureIndex(ctx); err != nil {
		log.WithError(err).Warningln("cannot ensure index")
	}

	logger.Infoln("connecting exchange grpc api")
	// won't use lcd endpoint here
	exchangeProvider, err := exchange.NewExchangeProvider(cfg.ExchangeGRPCURL, cfg.LcdURL, cfg.AssetPriceURL)
	if err != nil {
		return nil, err
	}

	svcTags := metrics.Tags{
		"svc": "guilds_process",
	}
	portfolioHelper, err := NewPortfolioHelper(ctx, exchangeProvider, logger)
	if err != nil {
		return nil, err
	}

	return &GuildsProcess{
		dbSvc:                   dbService,
		exchange:                exchangeProvider,
		logger:                  logger,
		portfolioUpdateInterval: cfg.PortfolioUpdateInterval,
		disqualifyInterval:      cfg.DisqualifyInterval,
		portfolioHelper:         portfolioHelper,
		grants:                  config.GrantRequirements,
		svcTags:                 svcTags,
	}, nil
}

func (p *GuildsProcess) Run(ctx context.Context) {
	p.logger.Infoln("guilds process is running to update portfolio and check to disqualify members")
	// run 2 cron jobs
	go p.runWithInterval(ctx, p.portfolioUpdateInterval, func(ctx context.Context) error {
		return p.captureMemberPortfolios(ctx)
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
func (p *GuildsProcess) captureMemberPortfolios(ctx context.Context) error {
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(p.svcTags)

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
				WithField("guild_id", guildID).
				WithError(err).Warningln("skip this guild")
			continue
		}

		// for each guild, don't need price be 100% accurate, we get all denom prices once
		// eliminate failure + save time
		denoms := model.GetGuildDenoms(guild)
		denoms = append(denoms, "inj")
		priceMap, err := p.portfolioHelper.GetDenomPrices(ctx, denoms)
		if err != nil {
			err = fmt.Errorf("get denom price err: %w", err)
			p.logger.
				WithField("guild_id", guildID).
				WithError(err).Warningln("get price err, skip this guild")
			continue
		}

		portfolios := make([]*model.AccountPortfolio, 0)
		denomToBalance := make(map[string]*model.Balance)
		var sumInjBankBalance = primitive.NewDecimal128(0, 0)
		// TODO: Create queue to re-try when failure happens
		// TODO: Create bulk accounts balances query on injective-exchange
		for _, member := range members {
			portfolioSnapshot, err := p.portfolioHelper.CaptureSingleMemberPortfolio(ctx, guild, member, false)
			if err != nil {
				p.logger.
					WithField("guild_id", guildID).
					WithField("memberAddr", member.InjectiveAddress.String()).
					WithError(err).Warningln("capture snapshot error")
				continue
			}

			if portfolioSnapshot == nil {
				continue
			}

			// fill denom price, priceMap has all denom prices
			for _, b := range portfolioSnapshot.Balances {
				b.PriceUSD = priceMap[b.Denom]

				// add to denom to balances
				if _, exist := denomToBalance[b.Denom]; !exist {
					denomToBalance[b.Denom] = &model.Balance{
						Denom:            b.Denom,
						PriceUSD:         b.PriceUSD,
						TotalBalance:     b.TotalBalance,
						AvailableBalance: b.AvailableBalance,
						UnrealizedPNL:    b.UnrealizedPNL,
						MarginHold:       b.MarginHold,
					}
				} else {
					tmp := denomToBalance[b.Denom]
					tmp.TotalBalance = sum(tmp.TotalBalance, b.TotalBalance)
					tmp.AvailableBalance = sum(tmp.AvailableBalance, b.AvailableBalance)
					tmp.UnrealizedPNL = sum(tmp.UnrealizedPNL, b.UnrealizedPNL)
					tmp.MarginHold = sum(tmp.MarginHold, b.MarginHold)
				}
			}

			// also fill denom price in bank balance
			for _, b := range portfolioSnapshot.BankBalances {
				if b.Denom == config.DEMOM_INJ {
					sumInjBankBalance = sum(sumInjBankBalance, b.Balance)
				}

				priceUSD, exist := priceMap[b.Denom]
				if exist {
					b.PriceUSD = priceUSD
				}
			}

			portfolioSnapshot.UpdatedAt = now
			portfolios = append(portfolios, portfolioSnapshot)
		}

		if len(portfolios) > 0 {
			p.logger.
				WithField("count", len(portfolios)).
				WithField("guild_id", guildID).Infoln("updated portfolios")
			err = p.dbSvc.AddAccountPortfolios(ctx, portfolios)
			if err != nil {
				p.logger.
					WithField("guild_id", guildID).
					WithError(err).Warningln("skip this guild")
			}
		}

		// if no failed member then we are confident with writing this snapshot
		if len(portfolios) == len(members) {
			// update guild
			guildPortfolio := &model.GuildPortfolio{
				GuildID:     guild.ID,
				UpdatedAt:   now,
				MemberCount: len(members),
				BankBalances: []*model.BankBalance{
					{
						Denom:    config.DEMOM_INJ,
						PriceUSD: priceMap[config.DEMOM_INJ],
						Balance:  sumInjBankBalance,
					},
				},
			}
			for _, v := range denomToBalance {
				guildPortfolio.Balances = append(guildPortfolio.Balances, v)
			}

			err = p.dbSvc.AddGuildPortfolios(ctx, []*model.GuildPortfolio{guildPortfolio})
			if err != nil {
				p.logger.
					WithField("guild_id", guildID).
					WithError(err).Warningln("cannot add guild portfolio")
			}
		}
	}
	return nil
}

// processDisqualify get orders from guild's markets
// then remove members whose orders' fee_recipient is not masterAddress
func (p *GuildsProcess) processDisqualification(ctx context.Context) error {
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(p.svcTags)

	guilds, err := p.dbSvc.ListAllGuilds(ctx)
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		return fmt.Errorf("list guild err: %w", err)
	}

	// TODO: clarify disqualify reason
	for _, g := range guilds {
		guildID := g.ID.Hex()
		isDefaultMember := false

		members, err := p.dbSvc.ListGuildMembers(ctx, model.MemberFilter{
			GuildID:         &guildID,
			IsDefaultMember: &isDefaultMember,
		})
		if err != nil {
			err = fmt.Errorf("list non-default member err: %w", err)
			p.logger.
				WithField("guild_id", guildID).
				WithError(err).Warningln("skip this guild")
			continue
		}

		markets := make([]string, 0)
		for _, market := range g.Markets {
			markets = append(markets, market.MarketID.Hex())
		}

		countDisqualifed := 0
		for _, member := range members {
			needDisqualified, err := p.shouldDisqualify(ctx, g, member.InjectiveAddress)
			if err != nil {
				continue
			}

			// we don't expect this to regularly happen,
			// so decided to delete each document this way
			if needDisqualified {
				err = p.dbSvc.RemoveMember(ctx, g.ID.Hex(), member.InjectiveAddress)
				if err != nil {
					log.WithField("memberAddress", member.InjectiveAddress.String()).
						WithError(err).Errorln("cannot delete member")
					continue
				}

				countDisqualifed++
			}
		}

		p.logger.WithField("count", countDisqualifed).WithField("guild_id", guildID).Info("disqualifed members")
	}
	return nil
}

func (p *GuildsProcess) meetGrantRequirements(
	ctx context.Context,
	guild *model.Guild,
	address string,
) (bool, error) {
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(p.svcTags)

	grants, err := p.exchange.GetGrants(ctx, address, guild.MasterAddress.String())
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		return false, err
	}

	msgToExpiration := make(map[string]time.Time)

	for _, g := range grants.Grants {
		t, err := time.Parse(expirationTimeLayout, g.Expiration)
		if err != nil {
			return false, fmt.Errorf("time parse err: %w", err)
		}

		msgToExpiration[g.Authorization.Msg] = t
	}

	// all expected grants must be provided
	now := time.Now()
	for _, expectedMsg := range p.grants {
		expiration, ok := msgToExpiration[expectedMsg]
		if !ok {
			p.logger.WithFields(log.Fields{
				"address":        address,
				"guild_id":       guild.ID.Hex(),
				"missed_message": expectedMsg,
			}).Info("account missed a grant")

			return false, nil
		}

		if expiration.Before(now) {
			p.logger.WithFields(log.Fields{
				"address":         address,
				"guild_id":        guild.ID.Hex(),
				"expired_message": expectedMsg,
				"expired_at":      expiration.String(),
			}).Info("account missed a grant")
			return false, nil
		}
	}

	return true, nil
}

func (p *GuildsProcess) spotOrdersHaveInvalidFeeRecipient(
	ctx context.Context,
	guild *model.Guild,
	defaultSubaccountID string,
	masterAddress string,
) (bool, error) {
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(p.svcTags)

	// open orders
	spotOrders, err := p.exchange.GetSpotOrders(
		ctx, []string{},
		defaultSubaccountID,
	)
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		return false, fmt.Errorf("get spot orders err: %w", err)
	}

	for _, o := range spotOrders {
		if strings.ToLower(o.FeeRecipient) != masterAddress {
			// we can log here to trace
			p.logger.WithFields(log.Fields{
				"subaccount_id": defaultSubaccountID,
				"guild_id":      guild.ID.Hex(),
				"fee_recipient": o.FeeRecipient,
			}).Info("account has invalid spot order")
			return true, nil
		}
	}

	return false, nil
}

func (p *GuildsProcess) derivativeOrdersHaveInvalidFeeRecipient(
	ctx context.Context,
	guild *model.Guild,
	defaultSubaccountID string,
	masterAddress string,
) (bool, error) {
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(p.svcTags)

	derivativeOrders, err := p.exchange.GetDerivativeOrders(
		ctx, []string{},
		defaultSubaccountID,
	)
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		return false, fmt.Errorf("get derivative orders err: %w", err)
	}

	for _, o := range derivativeOrders {
		if strings.ToLower(o.FeeRecipient) != masterAddress {
			p.logger.WithFields(log.Fields{
				"subaccount_id": defaultSubaccountID,
				"guild_id":      guild.ID.Hex(),
				"fee_recipient": o.FeeRecipient,
			}).Info("account has invalid derivative order")
			return true, nil
		}
	}

	return false, nil
}

// shouldDisqualify disqualifies a person if:
// - not enough grant requirement (user revoked at least one of them)
// - deriv/spot orders has fee recipient != master address
func (p *GuildsProcess) shouldDisqualify(
	ctx context.Context,
	guild *model.Guild,
	address model.Address,
) (bool, error) {
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(p.svcTags)

	defaultSubaccountID := defaultSubaccountIDFromInjAddress(address)
	masterAddress := strings.ToLower(guild.MasterAddress.String())
	// check grants
	meetRequirement, err := p.meetGrantRequirements(ctx, guild, address.String())
	if err != nil {
		// even it's not a fatal error, we should have metrics to monitor them
		metrics.ReportFuncError(p.svcTags)
		p.logger.WithField("address", address.String()).
			WithError(err).Warningln("check grants error")
	}

	if err == nil && !meetRequirement {
		// TODO: Add disqualification reason
		return true, nil
	}

	isInvalid, err := p.spotOrdersHaveInvalidFeeRecipient(ctx, guild, defaultSubaccountID, masterAddress)
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		p.logger.WithField("address", address.String()).
			WithError(err).Warningln("check orders error")
	}

	if err == nil && isInvalid {
		return true, nil
	}

	isInvalid, err = p.derivativeOrdersHaveInvalidFeeRecipient(ctx, guild, defaultSubaccountID, masterAddress)
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		return false, err
	}

	if err == nil && isInvalid {
		return true, nil
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
