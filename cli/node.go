package cli

import (
	"fmt"
	"strings"
	"strconv"

	"github.com/urfave/cli/v2"

	"github.com/LilithGames/moxa/service"
	"github.com/LilithGames/moxa/cluster"
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
				&cli.IntFlag{Name: "node_host_index", Aliases: []string{"index", "i"}, Value: 0, Usage: "node host index"},
			},
		},
		{
			Name:    "uncordon",
			Aliases: []string{},
			Action:  actionUncordonNode,
			Flags:   []cli.Flag{
				&cli.StringFlag{Name: "node_host_id", Aliases: []string{"nhid"}, Value: "", Usage: "node host id"},
				&cli.IntFlag{Name: "node_host_index", Aliases: []string{"index", "i"}, Value: 0, Usage: "node host index"},
			},
		},
		{
			Name:    "drain",
			Aliases: []string{},
			Action:  actionDrainNode,
			Flags:   []cli.Flag{
				&cli.StringFlag{Name: "node_host_id", Aliases: []string{"nhid"}, Value: "", Usage: "node host id"},
				&cli.IntFlag{Name: "node_host_index", Aliases: []string{"index", "i"}, Value: 0, Usage: "node host index"},
			},
		},
		{
			Name:    "delete",
			Aliases: []string{},
			Action:  actionDeleteNode,
			Flags:   []cli.Flag{
				&cli.StringFlag{Name: "node_host_id", Aliases: []string{"nhid"}, Value: "", Usage: "node host id"},
				&cli.IntFlag{Name: "node_host_index", Aliases: []string{"index", "i"}, Value: 0, Usage: "node host index"},
				&cli.BoolFlag{Name: "force", Aliases: []string{"f"}, Value: false, Usage: "force delete, may cause unhealth"},
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

func actionDeleteNode(cCtx *cli.Context) error {
	helper := NewHelper(cCtx)
	client := helper.MustClient()
	defer client.Close()
	req := &service.SyncRemoveNodeRequest{ShardId: cluster.MasterShardID}
	if cCtx.IsSet("node_host_id") {
		nodeID, err := parseNodeHostID(cCtx.String("node_host_id"))
		if err != nil {
			return fmt.Errorf("parseNodeHostID err: %w", err)
		}
		req.NodeId = lo.ToPtr(nodeID)
	} else if cCtx.IsSet("node_host_index") {
		req.NodeIndex = lo.ToPtr(int32(cCtx.Int("node_host_index")))
	}
	if cCtx.IsSet("force") {
		req.Force = lo.ToPtr(cCtx.Bool("force"))
	}
	resp, err := client.NodeHost().SyncRemoveNode(helper.Ctx(), req)
	if err != nil {
		return fmt.Errorf("client.NodeHost().SyncRemoveNode err: %w", err)
	}
	if resp.Updated > 0 {
		fmt.Printf("node deleted\n")
	} else {
		fmt.Printf("no change\n")
	}
	return nil
}

func parseNodeHostID(nhid string) (uint64, error) {
	items := strings.Split(nhid, "-")
	if len(items) != 2 || items[0] != "nhid" {
		return 0, fmt.Errorf("invalid nodeHostID %s", nhid)
	}
	id, err := strconv.ParseUint(items[len(items)-1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid nodeHostID %s, ParseUint err: %w", nhid, err)
	}
	return id, nil
}

