package guildsprocess

import (
	"context"
	"time"

	"github.com/InjectiveLabs/injective-guilds-service/internal/config"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db/mongoimpl"
	"github.com/InjectiveLabs/injective-guilds-service/internal/exchange"
	log "github.com/xlab/suplog"
)

type GuildsProcess struct {
	dbSvc    db.DBService
	exchange exchange.DataProvider
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
	exchangeProvider, err := exchange.NewExchangeProvider(cfg.ExchangeGRPCURL, "")
	if err != nil {
		return nil, err
	}

	return &GuildsProcess{
		dbSvc:    dbService,
		exchange: exchangeProvider,
	}, nil
}

func (p *GuildsProcess) Run(ctx context.Context) {
	// do something
}

func (p *GuildsProcess) processPorfolio(ctx context.Context) {

}

func (p *GuildsProcess) processDisqualify(ctx context.Context) {

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
