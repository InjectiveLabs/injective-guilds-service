package main

import (
	"context"
	"net/http"
	"strings"

	log "github.com/xlab/suplog"
	"goa.design/goa/middleware"
)

func panicIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func errorHandler(logger log.Logger) func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		id := ctx.Value(middleware.RequestIDKey).(string)
		_, _ = w.Write([]byte("[" + id + "] encoding: " + err.Error()))
		logger.Errorf("[%s] ERROR: %s", id, err.Error())
	}
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
