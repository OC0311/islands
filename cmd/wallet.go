package cmd

import (
	"github.com/jiangjincc/islands/wallet"
	"github.com/spf13/cobra"
)

var (
	walletCmd = &cobra.Command{
		Use:   "wallet",
		Short: "钱包相关命令",
		Long:  ``,
	}

	createWalletCmd = &cobra.Command{
		Use:   "create",
		Short: "创建钱包",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			wallet, _ := wallet.NewWallets()
			wallet.CreateNewWallet()
		},
	}

	listWalletCmd = &cobra.Command{
		Use:   "list",
		Short: "创建钱包",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			wallet, _ := wallet.NewWallets()
			wallet.WalletList()
		},
	}
)

func walletCmdExecute(rootCmd *cobra.Command) {
	rootCmd.AddCommand(walletCmd)
	walletCmd.AddCommand(createWalletCmd)
	walletCmd.AddCommand(listWalletCmd)
}
