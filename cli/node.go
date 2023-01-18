package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/LilithGames/moxa/service"
	"github.com/samber/lo"
)

var cmdNode = &cli.Command{
	Name:    "node",
	Aliases: []string{"n"},
	Usage:   "manage nodes",
	Action:  cli.ShowSubcommandHelp,
	Subcommands: []*cli.Command{
		{
			Name:    "cordon",
			Aliases: []string{},
			Action:  actionCordonNode,
			Flags:   []cli.Flag{
				&cli.StringFlag{Name: "node_host_id", Aliases: []string{"nhid"}, Value: "", Usage: "node host id"},
				&cli.IntFlag{Name: "node_host_index", Aliases: []string{"index"}, Value: 0, Usage: "node host index"},
			},
		},
		{
			Name:    "uncordon",
			Aliases: []string{},
			Action:  actionUncordonNode,
			Flags:   []cli.Flag{
				&cli.StringFlag{Name: "node_host_id", Aliases: []string{"nhid"}, Value: "", Usage: "node host id"},
				&cli.IntFlag{Name: "node_host_index", Aliases: []string{"index"}, Value: 0, Usage: "node host index"},
			},
		},
		{
			Name:    "drain",
			Aliases: []string{},
			Action:  actionDrainNode,
			Flags:   []cli.Flag{
				&cli.StringFlag{Name: "node_host_id", Aliases: []string{"nhid"}, Value: "", Usage: "node host id"},
				&cli.IntFlag{Name: "node_host_index", Aliases: []string{"index"}, Value: 0, Usage: "node host index"},
			},
		},
	},
}

func actionCordonNode(cCtx *cli.Context) error {
	helper := NewHelper(cCtx)
	client := helper.MustClient()
	defer client.Close()
	spec := &service.CordonNodeSpecRequest{}
	if cCtx.IsSet("node_host_id") {
		spec.NodeHostId = lo.ToPtr(cCtx.String("node_host_id"))
	} else if cCtx.IsSet("node_host_index") {
		spec.NodeIndex = lo.ToPtr(int32(cCtx.Int("node_host_index")))
	}
	resp, err := client.Spec().CordonNodeSpec(helper.Ctx(), spec)
	if err != nil {
		return fmt.Errorf("client.Spec().CordonNodeSpec: %w", err)
	}
	if resp.Updated > 0 {
		fmt.Printf("node spec %v updated\n", spec)
	} else {
		fmt.Printf("no change\n")
	}
	return nil
}

func actionUncordonNode(cCtx *cli.Context) error {
	helper := NewHelper(cCtx)
	client := helper.MustClient()
	defer client.Close()
	spec := &service.UncordonNodeSpecRequest{}
	if cCtx.IsSet("node_host_id") {
		spec.NodeHostId = lo.ToPtr(cCtx.String("node_host_id"))
	} else if cCtx.IsSet("node_host_index") {
		spec.NodeIndex = lo.ToPtr(int32(cCtx.Int("node_host_index")))
	}
	resp, err := client.Spec().UncordonNodeSpec(helper.Ctx(), spec)
	if err != nil {
		return fmt.Errorf("client.Spec().UncordonNodeSpec: %w", err)
	}
	if resp.Updated > 0 {
		fmt.Printf("node spec %v updated\n", spec)
	} else {
		fmt.Printf("no change\n")
	}
	return nil
}

func actionDrainNode(cCtx *cli.Context) error {
	helper := NewHelper(cCtx)
	client := helper.MustClient()
	defer client.Close()
	spec := &service.DrainNodeSpecRequest{}
	if cCtx.IsSet("node_host_id") {
		spec.NodeHostId = lo.ToPtr(cCtx.String("node_host_id"))
	} else if cCtx.IsSet("node_host_index") {
		spec.NodeIndex = lo.ToPtr(int32(cCtx.Int("node_host_index")))
	}
	resp, err := client.Spec().DrainNodeSpec(helper.Ctx(), spec)
	if err != nil {
		return fmt.Errorf("client.Spec().DrainNodeSpec: %w", err)
	}
	if resp.Updated > 0 {
		fmt.Printf("node spec %v updated\n", spec)
	} else {
		fmt.Printf("no change\n")
	}
	return nil
}
