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
	dbSvc           db.DBService
	exchange        exchange.DataProvider
	logger          log.Logger
	portfolioHelper *PortfolioHelper

	portfolioUpdateInterval time.Duration
	disqualifyInterval      time.Duration
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

	portfolioHelper, err := NewPortfolioHelper(ctx, dbService, exchangeProvider)
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

		// for each guild, don't need price be 100% accurate, we get all denom prices once
		// eliminate failure + save time
		denoms := model.GetGuildDenoms(guild)
		denoms = append(denoms, "inj")
		priceMap, err := p.portfolioHelper.GetDenomPrices(ctx, denoms)
		if err != nil {
			err = fmt.Errorf("get denom price err: %w", err)
			p.logger.
				WithField("guildID", guildID).
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
					WithField("guildID", guildID).
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

			if len(portfolioSnapshot.BankBalances) > 0 && (portfolioSnapshot.BankBalances[0].Denom == config.DEMOM_INJ) {
				sumInjBankBalance = sum(sumInjBankBalance, portfolioSnapshot.BankBalances[0].Balance)
			}

			portfolioSnapshot.UpdatedAt = now
			portfolios = append(portfolios, portfolioSnapshot)
		}

		if len(portfolios) > 0 {
			p.logger.
				WithField("count", len(portfolios)).
				WithField("guildID", guildID).Infoln("updated portfolios")
			err = p.dbSvc.AddAccountPortfolios(ctx, portfolios)
			if err != nil {
				p.logger.
					WithField("guildID", guildID).
					WithError(err).Warningln("skip this guild")
			}
		}

		// if no failed member then we are confident with writing this snapshot
		if len(portfolios) == len(members) {
			// update guild
			guildPortfolio := &model.GuildPortfolio{
				GuildID:   guild.ID,
				UpdatedAt: now,
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
					WithField("guildID", guildID).
					WithError(err).Warningln("cannot add guild portfolio")
			}
		}
	}
	return nil
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
		isDefaultMember := false

		members, err := p.dbSvc.ListGuildMembers(ctx, model.MemberFilter{
			GuildID:         &guildID,
			IsDefaultMember: &isDefaultMember,
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

		countDisqualifed := 0
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
					continue
				}

				countDisqualifed++
			}
		}

		p.logger.
			WithField("count", countDisqualifed).
			WithField("guildID", guildID).
			Info("disqualifed members")
	}
	return nil
}

// shouldDisqualify tries to disqualify a person if deriv/spot orders has fee recipient != master address
func (p *GuildsProcess) shouldDisqualify(
	ctx context.Context,
	guild *model.Guild,
	address model.Address,
) (bool, error) {
	defaultSubaccountID := defaultSubaccountIDFromInjAddress(address)
	masterAddress := strings.ToLower(guild.MasterAddress.String())

	spotOrders, err := p.exchange.GetSpotOrders(
		ctx, marketsFromGuild(guild, false),
		defaultSubaccountID,
	)
	if err != nil {
		p.logger.WithField("subaccountID", defaultSubaccountID).
			WithError(err).Warningln("get spot orders error")
	}

	for _, o := range spotOrders {
		if strings.ToLower(o.FeeRecipient) != masterAddress {
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
		if strings.ToLower(o.FeeRecipient) != masterAddress {
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
