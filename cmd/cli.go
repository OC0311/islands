package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "islands",
		Short: "islands blockchain cli tool",
	}
)

func CmdExecute() {
	// 链相关命令
	chainCmdExecute(rootCmd)
	txCmdExecute(rootCmd)
	walletCmdExecute(rootCmd)
	_ = rootCmd.Execute()
}
