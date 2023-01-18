package cli

import (
	"github.com/urfave/cli/v2"
)

// https://cli.urfave.org/v2/getting-started/
var App = &cli.App{
	Name: "moxactl",
	Usage: "manage moxa cluster",
	Action: cli.ShowAppHelp,
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "address", Aliases: []string{"addr"}, Value: "127.0.0.1:8001", Usage: "server address"},
	},
	Commands: []*cli.Command{cmdShard, cmdMigration, cmdNode},
}
