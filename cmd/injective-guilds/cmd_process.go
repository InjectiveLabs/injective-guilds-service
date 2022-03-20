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
	cfg := config.LoadGuildsProcessConfig()
	err := cfg.Validate()
	panicIf(err)

	// setup logger
	log.DefaultLogger.SetLevel(getLogLevel(cfg.LogLevel))
	guildsProcess, err := guildsprocess.NewProcess(cfg)
	panicIf(err)

	// run process(es) and hold until interrupt
	ctx := context.Background()
	cancelCtx, cancelFn := context.WithCancel(ctx)
	go func(cancelableContext context.Context) {
		guildsProcess.Run(cancelableContext)
	}(cancelCtx)

	closer.Bind(func() {
		cancelFn()
		guildsProcess.GracefullyShutdown(ctx)
	})

	closer.Hold()
}
