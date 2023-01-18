package cli

import (
	"testing"

	"github.com/jedib0t/go-pretty/v6/table"
)

func TestRenderTable(t *testing.T) {
	header := table.Row{"#", "Col1", "Col2", "Col3"}
	row1 := table.Row{"1", "2", "3", "4"}
	row2 := table.Row{"1", "2", "3", "4"}
	row3 := table.Row{"1", "2", "3", "4"}

	tw := table.NewWriter()

	tw.AppendHeader(header)
	tw.AppendRow(row1)
	tw.AppendRow(row2)
	tw.AppendRow(row3)
	tw.SetStyle(table.StyleColoredBright)

	println(tw.Render())
}
