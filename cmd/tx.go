package cmd

import (
	"fmt"

	"github.com/jiangjincc/islands/utils"

	"github.com/jiangjincc/islands/block"

	"github.com/spf13/cobra"
)

var (
	// 参数
	from   string
	to     string
	amount string

	address string

	txCmd = &cobra.Command{
		Use:   "tx",
		Short: "交易相关命令",
		Long:  ``,
	}

	txSendCmd = &cobra.Command{
		Use:   "send",
		Short: "发送交易",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			b := block.GetBlockchain()
			b.MineNewBlock(utils.JsonToArray(from), utils.JsonToArray(to), utils.JsonToArray(amount))

			// 更新utxo表
			utxoSet := &block.UTXOSet{Blockchain: b}
			utxoSet.Update()
		},
	}

	getBalanceCmd = &cobra.Command{
		Use:   "balance",
		Short: "查询月余额",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			b := block.GetBlockchain()

			utxoSet := block.UTXOSet{Blockchain: b}
			amount := utxoSet.GetBalance(address)
			//b.GetBalance(address)
			fmt.Printf("账户[ %s ]余额为: %d \n", address, amount)
		},
	}
)

func txCmdExecute(rootCmd *cobra.Command) {
	// 解析参数
	txSendCmd.Flags().StringVarP(&from, "from", "f", "", "发送人")
	_ = txSendCmd.MarkFlagRequired("from")

	txSendCmd.Flags().StringVarP(&to, "to", "t", "", "接受人")
	_ = txSendCmd.MarkFlagRequired("to")

	txSendCmd.Flags().StringVarP(&amount, "amount", "a", "", "金额")
	_ = txSendCmd.MarkFlagRequired("amount")

	getBalanceCmd.Flags().StringVarP(&address, "address", "a", "", "地址")
	_ = getBalanceCmd.MarkFlagRequired("address")

	rootCmd.AddCommand(txCmd)
	txCmd.AddCommand(txSendCmd)
	txCmd.AddCommand(getBalanceCmd)
}
