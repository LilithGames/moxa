package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/samber/lo"

	"github.com/LilithGames/moxa/service"
)

var cmdShard = &cli.Command{
	Name: "shard",
	Aliases: []string{"s"},
	Usage: "manage shards",
	Action: cli.ShowSubcommandHelp,
	Subcommands: []*cli.Command{
		{
			Name: "list",
			Aliases: []string{"ls"},
			Action: actionListShard,
			Flags: []cli.Flag{},
		},
		{
			Name: "check",
			Aliases: []string{},
			Action: actionCheckShard,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "shard_name", Aliases: []string{"name", "n"}, Value: "", Usage: "shard name"},
				&cli.BoolFlag{Name: "all", Aliases: []string{"a"}, Value: false, Usage: "all shard"},
			},
		},
		{
			Name: "get",
			Aliases: []string{"g"},
			Action: actionGetShard,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "shard_name", Aliases: []string{"name", "n"}, Value: "", Usage: "shard name", Required: true},
			},
		},
		{
			Name: "add",
			Aliases: []string{},
			Action: actionAddShard,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "shard_name", Aliases: []string{"name", "n"}, Value: "", Usage: "shard name", Required: true},
				&cli.StringFlag{Name: "profile_name", Aliases: []string{"profile", "p"}, Value: "", Usage: "profile name", Required: true},
				&cli.IntFlag{Name: "replica", Aliases: []string{}, Value: 3, Usage: "shard replica size"},
			},
		},
		{
			Name: "rm",
			Aliases: []string{},
			Action: actionRemoveShard,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "shard_name", Aliases: []string{"name", "n"}, Value: "", Usage: "shard name", Required: true},
			},
		},
		{
			Name: "rebalance",
			Aliases: []string{},
			Action: actionRebalanceShard,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "shard_name", Aliases: []string{"name", "n"}, Value: "", Usage: "shard name", Required: true},
				&cli.IntFlag{Name: "replica", Aliases: []string{}, Value: 3, Usage: "shard replica size"},
			},
		},
	},
}

func actionListShard(cCtx *cli.Context) error {
	helper := NewHelper(cCtx)
	client := helper.MustClient()
	defer client.Close()

	resp, err := client.Spec().ListShardSpec(helper.Ctx(), &service.ListShardSpecRequest{Status: true})
	if err != nil {
		return fmt.Errorf("client.Spec().ListShardSpec err: %w", err)
	}
	tw := table.NewWriter()
	ShardViewTableSiri(tw, resp.Shards...)
	DefaultTableStyle(tw)
	fmt.Printf("%s\n", tw.Render())
	return nil
}

func actionCheckShard(cCtx *cli.Context) error {
	helper := NewHelper(cCtx)
	client := helper.MustClient()
	defer client.Close()

	var shards []*service.ShardView

	if cCtx.IsSet("all") && cCtx.Bool("all") {
		resp, err := client.Spec().ListShardSpec(helper.Ctx(), &service.ListShardSpecRequest{Status: true})
		if err != nil {
			return fmt.Errorf("client.Spec().ListShardSpec err: %w", err)
		}
		shards = resp.Shards
	} else if cCtx.IsSet("shard_name") {
		name := cCtx.String("shard_name")
		resp, err := client.Spec().GetShardSpec(helper.Ctx(), &service.GetShardSpecRequest{ShardName: name})
		if err != nil {
			return fmt.Errorf("client.Spec().GetShardSpec err: %w", err)
		}
		shards = append(shards, resp.Shard)
	} else {
		return fmt.Errorf("--shard_name or --all required")
	}

	for _, shard := range shards {
		if shard.Type != service.ShardStatusType_Ready {
			detail := ""
			if shard.TypeDetail != nil {
				detail = *shard.TypeDetail 
			}
			return fmt.Errorf("shard %s not ready: %s %s", shard.Spec.ShardName, shard.Type.String(), detail)
		}
		if !shard.Healthz.Healthz {
			detail := ""
			if shard.Healthz.HealthzDetail != nil {
				detail = *shard.Healthz.HealthzDetail
			}
			return fmt.Errorf("shard %s not healthz: %s", shard.Spec.ShardName, detail)
		}
	}
	fmt.Printf("all shard ok\n")
	return nil
}

func actionGetShard(cCtx *cli.Context) error {
	helper := NewHelper(cCtx)
	client := helper.MustClient()
	defer client.Close()

	name := cCtx.String("shard_name")
	resp, err := client.Spec().GetShardSpec(helper.Ctx(), &service.GetShardSpecRequest{ShardName: name})
	if err != nil {
		return fmt.Errorf("client.Spec().GetShardSpec err: %w", err)
	}
	fmt.Printf("%s\n", MustJsonSiri(resp.Shard))
	return nil
}

func actionAddShard(cCtx *cli.Context) error {
	helper := NewHelper(cCtx)
	client := helper.MustClient()
	defer client.Close()

	name := cCtx.String("shard_name")
	profile := cCtx.String("profile_name")
	spec := &service.AddShardSpecRequest{ShardName: name, ProfileName: profile}
	if cCtx.IsSet("replica") {
		spec.Replica = lo.ToPtr(int32(cCtx.Int("replica")))
	}
	resp, err := client.Spec().AddShardSpec(helper.Ctx(), spec)
	if err != nil {
		return fmt.Errorf("client.Spec().AddShardSpec err: %w", err)
	}
	fmt.Printf("%s\n", MustJsonSiri(resp.Shard))
	return nil
}

func actionRemoveShard(cCtx *cli.Context) error {
	helper := NewHelper(cCtx)
	client := helper.MustClient()
	defer client.Close()

	name := cCtx.String("shard_name")
	resp, err := client.Spec().RemoveShardSpec(helper.Ctx(), &service.RemoveShardSpecRequest{ShardName: name})
	if err != nil {
		return fmt.Errorf("client.Spec().RemoveShardSpec: %w", err)
	}
	if resp.Updated > 0 {
		fmt.Printf("shard spec %s removed\n", name)
	} else {
		fmt.Printf("no change\n")
	}
	return nil
}

func actionRebalanceShard(cCtx *cli.Context) error {
	helper := NewHelper(cCtx)
	client := helper.MustClient()
	defer client.Close()

	name := cCtx.String("shard_name")
	spec := &service.RebalanceShardSpecRequest{ShardName: name}
	if cCtx.IsSet("replica") {
		spec.Replica = lo.ToPtr(int32(cCtx.Int("replica")))
	}
	resp, err := client.Spec().RebalanceShardSpec(helper.Ctx(), spec)
	if err != nil {
		return fmt.Errorf("client.Spec().RebalanceShardSpec err: %w", err)
	}
	if resp.Updated > 0 {
		fmt.Printf("shard spec %s updated\n", name)
	} else {
		fmt.Printf("no change\n")
	}
	return nil
}
