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
	name              *string
	description       *string
	capacity          *int
	masterAddr        *string
	defaultMemberAddr *string

	minSpotBase        *int
	minSpotQuote       *int
	minDerivativeQuote *int
	minStaking         *int

	dbURL         *string
	exchangeURL   *string
	assetPriceURL *string
)

func setConfig() {
	// config cosmos type address prefix
	cosmtypes.GetConfig().SetBech32PrefixForAccount("inj", "injpub")
}

func main() {
	setConfig()
	app.Command("api", "start Guilds service HTTP API server", cmdApi)
	app.Command("process", "start Guilds process, which takes portfolios snapshots and handle disqualification", cmdProcess)
	app.Command("update-denom", "update all denom-coinID map to database", cmdUpdateDenom)
	app.Command("add-guild", "add a guild", cmdAddGuild)

	_ = app.Run(os.Args)
}
