package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/InjectiveLabs/injective-guilds-service/internal/config"
	"github.com/InjectiveLabs/metrics"
	log "github.com/xlab/suplog"
	"goa.design/goa/middleware"
)

const retryCount = 5

func panicIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func errorHandler(logger log.Logger) func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		id, ok := ctx.Value(middleware.RequestIDKey).(string)
		if !ok {
			id = "nil"
		}
		_, _ = w.Write([]byte("[" + id + "] encoding: " + err.Error()))
		logger.Errorf("[%s] ERROR: %s", id, err.Error())
	}
}

func connectStatServerWithRetry(
	envName string,
	cfg config.StatsdConfig,
	retryCount int,
) error {
	var err error

	hostName, _ := os.Hostname()
	for i := 0; i < retryCount; i++ {
		if err = metrics.Init(cfg.Addr, cfg.Prefix, &metrics.StatterConfig{
			Agent:                cfg.Agent,
			HostName:             hostName,
			EnvName:              envName,
			StuckFunctionTimeout: cfg.StuckDur,
			MockingEnabled:       cfg.Mocking,
		}); err == nil {
			return nil
		}
		log.WithError(err).Warningln("stat server connect failed, retrying...")
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("failed to connect stat server, last error: %w", err)
}

func getLogLevel(s string) log.Level {
	switch strings.ToLower(s) {
	case "1", "error":
		return log.ErrorLevel
	case "2", "warn":
		return log.WarnLevel
	case "3", "info":
		return log.InfoLevel
	case "4", "debug":
		return log.DebugLevel
	default:
		return log.FatalLevel
	}
}
