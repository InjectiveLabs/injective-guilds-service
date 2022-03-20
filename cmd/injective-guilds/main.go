package main

import (
	"os"

	cli "github.com/jawher/mow.cli"
)

var app = cli.App("injective-guilds", "A microserivce for trading guilds queries")

func main() {
	app.Command("api", "start Guilds service HTTP API server", cmdApi)
	app.Command("process", "start Guilds process, which takes portfolios snapshots and handle disqualification", cmdProcess)
	_ = app.Run(os.Args)
}
