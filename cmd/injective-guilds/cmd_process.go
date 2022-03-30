package main

import (
	"context"

	"github.com/InjectiveLabs/injective-guilds-service/internal/config"
	guildsprocess "github.com/InjectiveLabs/injective-guilds-service/internal/service/guilds-process"
	cli "github.com/jawher/mow.cli"
	"github.com/xlab/closer"
	log "github.com/xlab/suplog"
)

func cmdProcess(c *cli.Cmd) {
	c.Action = func() {
		cfg := config.LoadGuildsProcessConfig()
		err := cfg.Validate()
		panicIf(err)

		// check denom prices
		doubleCheckDenomConfig(cfg.AssetPriceURL)

		if !cfg.StatsdConfig.Disabled {
			// set global stat and log
			err = connectStatServerWithRetry(cfg.EnvName, cfg.StatsdConfig, retryCount)
			panicIf(err)
		}

		// setup logger
		log.DefaultLogger.SetLevel(getLogLevel(cfg.LogLevel))

		guildsProcess, err := guildsprocess.NewProcess(cfg)
		panicIf(err)
		// run process(es) and hold until interrupt
		ctx := context.Background()
		cancelCtx, cancelFn := context.WithCancel(ctx)
		guildsProcess.Run(cancelCtx)

		closer.Bind(func() {
			cancelFn()
			guildsProcess.GracefullyShutdown(ctx)
		})

		closer.Hold()
	}
}
