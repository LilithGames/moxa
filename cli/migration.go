package cli

import (
	"fmt"
	"context"

	"github.com/urfave/cli/v2"

	"github.com/LilithGames/moxa/service"
	"github.com/jedib0t/go-pretty/v6/table"
)

var cmdMigration = &cli.Command{
	Name: "migration",
	Aliases: []string{"m"},
	Usage: "manage shard migration",
	Action: cli.ShowSubcommandHelp,
	Subcommands: []*cli.Command{
		{
			Name: "list",
			Aliases: []string{"ls"},
			Action: actionListMigration,
			Flags: []cli.Flag{},
		},
		{
			Name: "get",
			Aliases: []string{"g"},
			Action: actionGetMigration,
			Flags: []cli.Flag{
				&cli.Uint64Flag{Name: "shard_id", Aliases: []string{"sid"}, Value: 0, Usage: "shard id", Required: true},
			},
		},
		{
			Name: "check",
			Aliases: []string{"c"},
			Action: actionCheckMigration,
			Flags: []cli.Flag{
				&cli.Uint64Flag{Name: "shard_id", Aliases: []string{"sid"}, Value: 0, Usage: "shard id"},
				&cli.BoolFlag{Name: "all", Aliases: []string{"a"}, Value: false, Usage: "all shard"},
			},
		},
		{
			Name: "prepare",
			Aliases: []string{"p"},
			Action: actionPrepareMigration,
			Flags: []cli.Flag{
				&cli.Uint64Flag{Name: "shard_id", Aliases: []string{"sid"}, Value: 0, Usage: "shard id"},
				&cli.BoolFlag{Name: "all", Aliases: []string{"a"}, Value: false, Usage: "all shard"},
			},
		},
	},
}

func actionListMigration(cCtx *cli.Context) error {
	helper := NewHelper(cCtx)
	client := helper.MustClient()
	defer client.Close()
	resp, err := client.Spec().ListMigration(context.TODO(), &service.ListMigrationRequest{})
	if err != nil {
		return fmt.Errorf("client.Spec().ListMigration err: %w", err)
	}
	tw := table.NewWriter()
	MigrationTableSiri(tw, resp.Migrations...)
	DefaultTableStyle(tw)
	fmt.Printf("%s\n", tw.Render())
	return nil
}

func actionGetMigration(cCtx *cli.Context) error {
	helper := NewHelper(cCtx)
	client := helper.MustClient()
	defer client.Close()
	shardID := cCtx.Uint64("shard_id")
	resp, err := client.Spec().GetMigration(context.TODO(), &service.QueryMigrationRequest{ShardId: shardID})
	if err != nil {
		return fmt.Errorf("client.Spec().GetMigration err: %w", err)
	}
	fmt.Printf("%s\n", MustJsonSiri(resp.State))
	return nil
}

func actionCheckMigration(cCtx *cli.Context) error {
	helper := NewHelper(cCtx)
	client := helper.MustClient()
	defer client.Close()

	if cCtx.IsSet("all") && cCtx.Bool("all") {
		resp, err := client.Spec().ListMigration(context.TODO(), &service.ListMigrationRequest{})
		if err != nil {
			return fmt.Errorf("client.Spec().ListMigration err: %w", err)
		}
		for _, mi := range resp.Migrations {
			if err := mi.GetErr(); err != nil {
				return cli.Exit(fmt.Sprintf("shard %d: %s", err.ShardId, err.Error), 1)
			} else if data := mi.GetData(); data != nil {
				if data.Type != "Empty" {
					return cli.Exit(fmt.Sprintf("shard %d migration type != Empty: %s", data.ShardId, data.Type), 1)
				}
				continue
			} else {
				panic(fmt.Errorf("unknown MigrationStateType: %T", mi.Item))
			}
		}
		fmt.Printf("all migration status ok\n", )
		return nil
	} else if cCtx.IsSet("shard_id") {
		shardID := cCtx.Uint64("shard_id")
		resp, err := client.Spec().GetMigration(context.TODO(), &service.QueryMigrationRequest{ShardId: shardID})
		if err != nil {
			return fmt.Errorf("client.Spec().GetMigration err: %w", err)
		}
		if resp.State.Type != "Empty" {
			return cli.Exit(fmt.Sprintf("shard %d migration type != Empty: %s", shardID, resp.State.Type), 1)
		}
		fmt.Printf("migration %d status ok\n", shardID)
		return nil
	}
	return cli.Exit("shard_id or all required", 1)
}

func actionPrepareMigration(cCtx *cli.Context) error {
	helper := NewHelper(cCtx)
	client := helper.MustClient()
	defer client.Close()

	if cCtx.IsSet("all") && cCtx.Bool("all") {
		resp, err := client.Spec().ListMigration(context.TODO(), &service.ListMigrationRequest{})
		if err != nil {
			return fmt.Errorf("client.Spec().ListMigration err: %w", err)
		}
		migrations := make([]*service.MigrationState, 0, len(resp.Migrations))
		for _, mi := range resp.Migrations {
			if err := mi.GetErr(); err != nil {
				return fmt.Errorf("get migration %d err: %s", err.ShardId, err.Error)
			}
			data := mi.GetData()
			migrations = append(migrations, data)
		}
		for _, mi := range migrations {
			if _, err := client.Spec().CreateMigration(context.TODO(), &service.CreateMigrationRequest{ShardId: mi.ShardId}); err != nil {
				return fmt.Errorf("create migration %d err: %w", mi.ShardId, err)
			}
		}
		fmt.Printf("all migration prepared\n")
		return nil

	} else if cCtx.IsSet("shard_id") {
		shardID := cCtx.Uint64("shard_id")
		if _, err := client.Spec().CreateMigration(context.TODO(), &service.CreateMigrationRequest{ShardId: shardID}); err != nil {
			return fmt.Errorf("client.Spec().CreateMigration err: %w", err)
		}
		fmt.Printf("migration %d prepared\n")
		return nil
	}
	return cli.Exit("shard_id or all required", 1)
}
