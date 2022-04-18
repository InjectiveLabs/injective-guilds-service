package main

import (
	"context"
	"fmt"

	"github.com/InjectiveLabs/injective-guilds-service/internal/db/mongoimpl"
	cli "github.com/jawher/mow.cli"
	log "github.com/xlab/suplog"
)

func parseSetCapacityArg(c *cli.Cmd) {
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

	capacity = c.Int(cli.IntOpt{
		Name:  "capacity",
		Desc:  "capacity to set",
		Value: 150,
	})
}

func setCapacityAction() {
	log.Info("connecting database")
	ctx := context.Background()
	dbSvc, err := mongoimpl.NewService(ctx, *dbURL, "guilds")
	panicIf(err)

	guild, err := dbSvc.GetSingleGuild(ctx, *guildID)
	panicIf(err)

	if guild.MemberCount > *capacity {
		err = fmt.Errorf("cannot set guild capacity lower than current member count")
		panic(err)
	}

	err = dbSvc.SetGuildCap(ctx, guild.ID.Hex(), *capacity)
	panicIf(err)

	log.Infof("🍺 updated guild %s (%s) member capacity to %d", guild.Name, guild.ID, *capacity)
}

func cmdSetCapacity(c *cli.Cmd) {
	// inputs:
	// guild id: --guild-id
	// db url: --db-url
	// capacity: --capacity
	parseSetCapacityArg(c)
	c.Action = setCapacityAction
}
