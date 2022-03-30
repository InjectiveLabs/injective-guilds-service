package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	guildsapisvc "github.com/InjectiveLabs/injective-guilds-service/api/gen/guilds_service"
	guildsapisvr "github.com/InjectiveLabs/injective-guilds-service/api/gen/http/guilds_service/server"
	"github.com/InjectiveLabs/injective-guilds-service/internal/config"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db/mongoimpl"
	"github.com/InjectiveLabs/injective-guilds-service/internal/exchange"
	guildsapi "github.com/InjectiveLabs/injective-guilds-service/internal/service/guilds-api"
	cli "github.com/jawher/mow.cli"
	"github.com/xlab/closer"
	log "github.com/xlab/suplog"
	goahttp "goa.design/goa/v3/http"
)

type APIServer struct {
	cfg      config.GuildsAPIServerConfig
	dbSvc    db.DBService
	exchange exchange.DataProvider
	handlers http.Handler
	server   *http.Server
}

func NewServer(cfg config.GuildsAPIServerConfig) (*APIServer, error) {
	var err error
	s := &APIServer{cfg: cfg}

	ctx := context.Background()
	s.dbSvc, err = mongoimpl.NewService(ctx, cfg.DBConnectionURL, cfg.DBName)
	if err != nil {
		return nil, err
	}

	s.exchange, err = exchange.NewExchangeProvider(cfg.ExchangeGRPCURL, cfg.LcdURL, cfg.AssetPriceURL)
	if err != nil {
		return nil, err
	}

	// prepare service implementations
	guildsApi, err := guildsapi.NewService(ctx, s.dbSvc, s.exchange)
	if err != nil {
		return nil, err
	}

	// prepare endpoints
	guildsServiceEndpoints := guildsapisvc.NewEndpoints(guildsApi)

	var (
		dec                 = goahttp.RequestDecoder
		enc                 = goahttp.ResponseEncoder
		logger              = log.WithField("svc", "guilds_service")
		eh                  = errorHandler(logger)
		mux                 = goahttp.NewMuxer()
		guildsServiceServer = guildsapisvr.New(guildsServiceEndpoints, mux, dec, enc, eh, nil)
	)

	// mounts
	guildsapisvr.Mount(mux, guildsServiceServer)
	s.handlers = mux

	return s, nil
}

func (s *APIServer) ListenAndServe(ctx context.Context) error {
	var (
		address string
		tls     bool
	)

	switch {
	case strings.HasPrefix(s.cfg.ListenAddress, "http://"):
		address = strings.TrimPrefix(s.cfg.ListenAddress, "http://")
		tls = false
	case strings.HasPrefix(s.cfg.ListenAddress, "https://"):
		address = strings.TrimPrefix(s.cfg.ListenAddress, "https://")
		tls = true
	default:
		return fmt.Errorf("unsupported protocol with address: %s, need http or https", s.cfg.ListenAddress)
	}

	// new server + listenFn
	s.server = &http.Server{Addr: address, Handler: s.handlers}
	go func() {
		var err error
		if tls {
			log.Infoln("listening with tls:", s.server.Addr)
			err = s.server.ListenAndServeTLS(
				s.cfg.TLSCertFilePath, s.cfg.TLSCertFilePath,
			)
		} else {
			log.Infoln("listening no tls:", s.server.Addr)
			err = s.server.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			log.WithError(err).Errorln("listen and serve error")

			// call to gracefully close everything after an error occurs
			closer.Close()
		}
	}()

	return nil
}

func (s *APIServer) GracefullyShutdown() {
	log.Info("service is going to stop")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Info("shutting down api server")
	if err := s.server.Shutdown(shutdownCtx); err != nil && err != http.ErrServerClosed {
		log.WithError(err).Error("cannot shutdown server")
	}

	// close db
	log.Info("closing db connection")
	if err := s.dbSvc.Disconnect(shutdownCtx); err != nil {
		log.WithError(err).Error("cannot close db connection")
	}

	// close exchange grpc
	log.Info("closing exchange grpc connection")
	if err := s.exchange.Close(); err != nil {
		log.WithError(err).Error("cannot close exchange grpc connection")
	}

	log.Info("server stopped")
}

func cmdApi(c *cli.Cmd) {
	c.Action = func() {
		cfg := config.LoadGuildsAPIServerConfig()
		err := cfg.Validate()
		panicIf(err)

		// check prices
		doubleCheckDenomConfig(cfg.AssetPriceURL)
		if !cfg.StatsdConfig.Disabled {
			// set global stat and log
			err = connectStatServerWithRetry(cfg.EnvName, cfg.StatsdConfig, retryCount)
			panicIf(err)
		}

		// setup logger
		log.DefaultLogger.SetLevel(getLogLevel(cfg.LogLevel))

		apiServer, err := NewServer(cfg)
		panicIf(err)

		err = apiServer.ListenAndServe(context.Background())
		panicIf(err)

		closer.Bind(func() {
			apiServer.GracefullyShutdown()
		})
		// wait until os signal then gracefully shutdown
		closer.Hold()
	}
}
