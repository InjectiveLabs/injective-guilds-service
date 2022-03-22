package guildsprocess

import (
	"context"
	"fmt"
	"time"

	"github.com/InjectiveLabs/injective-guilds-service/internal/config"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db/model"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db/mongoimpl"
	"github.com/InjectiveLabs/injective-guilds-service/internal/exchange"
	cosmtypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/xlab/suplog"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// TODO: Implement retryable mechanism
func (p *GuildsProcess) capturePorfolioSnapshot(ctx context.Context) error {
	guilds, err := p.dbSvc.ListAllGuilds(ctx)
	if err != nil {
		return fmt.Errorf("list guild err: %w", err)
	}

	now := time.Now()
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

		portfolios := make([]*model.AccountPortfolio, 0)
		denoms := model.GetGuildDenoms(g)
		// TODO: Create queue to re-try when failure happens
		// TODO: Create bulk accounts balances query on injective-exchange
		for _, member := range members {
			balances, err := p.getMemberBalances(ctx, denoms, member)
			if err != nil {
				// try updating another guilds
				err = fmt.Errorf("get member balance err: %w", err)
				p.logger.
					WithField("guildID", g.ID.Hex()).
					WithError(err).Warningln("skip this guild")
				continue
			}

			portfolioSnapshot := &model.AccountPortfolio{
				GuildID:          g.ID,
				InjectiveAddress: member.InjectiveAddress,
				Balances:         balances,
				UpdatedAt:        now,
			}
			portfolios = append(portfolios, portfolioSnapshot)
		}

		err = p.dbSvc.AddAccountPortfolios(ctx, g.ID.Hex(), portfolios)
		if err != nil {
			p.logger.
				WithField("guildID", g.ID.Hex()).
				WithError(err).Warningln("skip this guild")
		}
	}
	return nil
}

func defaultSubaccountIDFromInjAddress(injAddress model.Address) string {
	ethAddr := common.BytesToAddress(injAddress.Bytes())
	return ethAddr.Hex() + "000000000000000000000000"
}

func (p *GuildsProcess) getMemberBalances(ctx context.Context, denoms []string, member *model.GuildMember) (result []*model.Balance, err error) {
	defaultSubaccountID := defaultSubaccountIDFromInjAddress(member.InjectiveAddress)

	balances, err := p.exchange.GetSubaccountBalances(ctx, defaultSubaccountID)
	if err != nil {
		return nil, err
	}

	denomToBalance := make(map[string]*exchange.Balance)
	for _, b := range balances {
		denomToBalance[b.Denom] = b
	}
	// filter denoms and add to result
	for _, denom := range denoms {
		b, exist := denomToBalance[denom]
		if !exist {
			result = append(result, &model.Balance{
				Denom:            denom,
				TotalBalance:     primitive.NewDecimal128(0, 0),
				AvailableBalance: primitive.NewDecimal128(0, 0),
			})
			continue
		}

		// we parsed this total balance successfully -> no error expected
		totalBalance, _ := primitive.ParseDecimal128(b.TotalBalance.String())
		availBalance, _ := primitive.ParseDecimal128(b.AvailableBalance.String())
		result = append(result, &model.Balance{
			Denom:            denom,
			TotalBalance:     totalBalance,
			AvailableBalance: availBalance,
		})
	}

	return result, nil
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
