package block

import (
	"os"
	"testing"

	"github.com/jedib0t/go-pretty/table"
)

func TestBlockchain_PrintBlocks(t *testing.T) {
	b := &Block{
		Height: 1,
		Txs:    []*Transaction{},
	}

	tb := table.NewWriter()
	tb.SetOutputMirror(os.Stdout)
	tb.AppendHeader(table.Row{"内容", "区块信息"})
	tb.AppendRows([]table.Row{
		{"区块高度", b.Height},
		{"区块数据", b.HashTransaction()},
	})
	tb.SetStyle(table.StyleDefault)
	tb.Render()
}
