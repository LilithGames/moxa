package cli

import (
	"google.golang.org/protobuf/proto"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/protobuf/encoding/protojson"
	"github.com/samber/lo"

	"github.com/LilithGames/moxa/service"
)

func DefaultTableStyle(tw table.Writer) {
	// tw.SetStyle(table.StyleColoredBright)
	tw.Style().Options.DrawBorder = false
	tw.Style().Options.SeparateColumns = false
	tw.Style().Options.SeparateFooter = false
	tw.Style().Options.SeparateHeader = false
	tw.Style().Options.SeparateRows = false
}

func MigrationRowSiri(item *service.MigrationStateListItem) table.Row {
	if data := item.GetData(); data != nil {
		return table.Row{
			data.ShardId,
			data.Type,
			data.Version,
			data.StateIndex,
		}
	}
	if err := item.GetErr(); err != nil {
		return table.Row{
			err.ShardId,
			err.Error,
		}
	}
	panic("impossible")
}

func MustJsonSiri(item proto.Message) string {
	opts := protojson.MarshalOptions{Multiline: true, Indent: "  ", EmitUnpopulated: true}
	bs, err := opts.Marshal(item)
	if err != nil {
		panic("impossible")
	}
	return string(bs)
}

func MigrationTableSiri(tw table.Writer, items ...*service.MigrationStateListItem)  {
	header := table.Row{"ShardID", "Type", "Version", "StateIndex"}
	tw.AppendHeader(header)
	for _, item := range items {
		row := MigrationRowSiri(item)
		tw.AppendRow(row)
	}
}

func ShardViewRowSiri(view *service.ShardView) table.Row {
	nodeIDs := lo.Map(view.Status.Nodes, func(node *service.NodeStatusView, _ int) uint64 {
		return node.NodeId
	})
	return table.Row{
		view.Spec.ShardName,
		view.Spec.ShardId,
		view.Spec.Replica,
		nodeIDs,
		view.Type,
		view.Healthz.Healthz,
	}
}

func ShardViewTableSiri(tw table.Writer, views ...*service.ShardView) {
	header := table.Row{"Name", "ShardID", "Replica", "Nodes", "Type", "Healthz"}
	tw.AppendHeader(header)
	for _, view := range views {
		row := ShardViewRowSiri(view)
		tw.AppendRow(row)
	}
}
