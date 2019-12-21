package cmd

import (
	"github.com/jiangjincc/islands/block"
	"github.com/spf13/cobra"
)

var (
	node string

	nodeCmd = &cobra.Command{
		Use:   "node",
		Short: "节点相关命令",
		Long:  ``,
	}

	nodeStartCmd = &cobra.Command{
		Use:   "start",
		Short: "节点启动",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			block.StartServer(node)
		},
	}
)

func nodeCmdExecute(rootCmd *cobra.Command) {
	nodeStartCmd.Flags().StringVarP(&node, "node", "n", "3000", "节点端口号")
	_ = txSendCmd.MarkFlagRequired("from")

	rootCmd.AddCommand(nodeCmd)
	nodeCmd.AddCommand(nodeStartCmd)
}
