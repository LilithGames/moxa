package cli

import (
	"github.com/urfave/cli/v2"
)

var cmdNodehost = &cli.Command{
	Name:    "nodehost",
	Aliases: []string{"nh"},
	Usage:   "manage nodehosts",
	Action:  cli.ShowSubcommandHelp,
	Subcommands: []*cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Action:  actionListNodehosts,
			Flags:   []cli.Flag{},
		},
	},
}

func actionListNodehosts(cCtx *cli.Context) error {
	return nil
}

