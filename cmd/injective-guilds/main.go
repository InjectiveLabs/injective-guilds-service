package main

import (
	"os"

	cli "github.com/jawher/mow.cli"
)

var app = cli.App("injective-guilds", "A microserivce for trading guilds queries")

func main() {
	_ = app.Run(os.Args)
}
