package cli

import (
	"log"
	"fmt"
	"syscall"

	"github.com/urfave/cli/v2"
	"github.com/lni/goutils/syncutil"
	"github.com/samber/lo"

	"github.com/LilithGames/moxa/service"
	"github.com/LilithGames/moxa/utils"
)

var cmdDebug = &cli.Command{
	Name: "debug",
	Aliases: []string{"d"},
	Usage: "debug",
	Action: cli.ShowSubcommandHelp,
	Subcommands: []*cli.Command{
		{
			Name: "subscribe",
			Aliases: []string{"sub"},
			Action: cli.ShowSubcommandHelp,
			Flags: []cli.Flag{},
			Subcommands: []*cli.Command{
				{
					Name: "members",
					Aliases: []string{"m"},
					Action: actionSubscribeMembers,
					Flags: []cli.Flag{},
				},
			},
		},
	},
}

func actionSubscribeMembers(cCtx *cli.Context) error {
	helper := NewHelper(cCtx)
	client := helper.MustClient()
	defer client.Close()

	resp, err := client.NodeHost().ListMemberState(helper.Ctx(), &service.ListMemberStateRequest{})
	if err != nil {
		return fmt.Errorf("client.NodeHost().ListMember err: %w", err)
	}
	for _, member := range resp.Members {
		fmt.Printf("member: %s %v\n", member.Meta.NodeHostId, lo.Keys(member.Shards))
	}

	session, err := client.NodeHost().SubscribeMemberState(helper.Ctx(), &service.SubscribeMemberStateRequest{Version: resp.Version})
	defer session.CloseSend()
	if err != nil {
		return fmt.Errorf("client.NodeHost().SubscribeMember err: %w", err)
	}
	stopper := syncutil.NewStopper()
	utils.SignalHandler(stopper, syscall.SIGINT, syscall.SIGTERM)
	stopper.RunWorker(func() {
		for {
			resp2, err := session.Recv()
			if err != nil {
				log.Println("[ERROR]", fmt.Errorf("session.Recv err: %w", err))
				return 
			}
			state := resp2.State
			fmt.Printf("member: %d %s %s %v\n", resp2.Version, resp2.Type.String(), state.Meta.NodeHostId, lo.Keys(state.Shards))
		}
	})
	stopper.RunWorker(func() {
		<-stopper.ShouldStop()
		helper.Cancel()
	})
	stopper.Wait()
	return nil
}
