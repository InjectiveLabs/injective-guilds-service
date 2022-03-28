package main

import (
	"context"

	"github.com/InjectiveLabs/injective-guilds-service/internal/db/mongoimpl"
	cli "github.com/jawher/mow.cli"
	log "github.com/xlab/suplog"
)

func parseDeleteGuildArgs(c *cli.Cmd) {
	guildID = c.String(cli.StringOpt{
		Name:  "guild-id",
		Desc:  "guild ID to delete",
		Value: "",
	})

	dbURL = c.String(cli.StringOpt{
		Name:  "db-url",
		Desc:  "database url",
		Value: "mongodb://localhost:27017",
	})
}

func deleteGuildAction() {
	log.Info("connecting database")
	ctx := context.Background()
	dbSvc, err := mongoimpl.NewService(ctx, *dbURL, "guilds")
	panicIf(err)

	err = dbSvc.DeleteGuild(ctx, *guildID)
	panicIf(err)
	log.Info("🍺 delete done, please check db")
}

func cmdDeleteGuild(c *cli.Cmd) {
	// inputs:
	// guild id: --guild-id
	parseDeleteGuildArgs(c)
	c.Action = deleteGuildAction
}
