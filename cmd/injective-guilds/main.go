package main

import (
	"os"

	cosmtypes "github.com/cosmos/cosmos-sdk/types"
	cli "github.com/jawher/mow.cli"
)

var (
	app               = cli.App("injective-guilds", "A microserivce for trading guilds queries")
	spotIDs           *[]string
	derivativeIDs     *[]string
	guildID           *string
	name              *string
	description       *string
	capacity          *int
	masterAddr        *string
	defaultMemberAddr *string
	memberParams      *string

	spotRequirements       *[]string
	derivativeRequirements *[]string
	minStaking             *int

	dbURL         *string
	exchangeURL   *string
	assetPriceURL *string
	lcdURL        *string
)

func setConfig() {
	// config cosmos type address prefix
	cosmtypes.GetConfig().SetBech32PrefixForAccount("inj", "injpub")
}

func main() {
	setConfig()
	app.Command("api", "start Guilds service HTTP API server", cmdApi)
	app.Command("process", "start Guilds process, which takes portfolios snapshots and handle disqualification", cmdProcess)
	app.Command("add-guild", "add a guild", cmdAddGuild)
	app.Command("delete-guild", "delete a guild", cmdDeleteGuild)
	app.Command("set-capacity", "set member capacity of a guild", cmdSetCapacity)

	_ = app.Run(os.Args)
}
