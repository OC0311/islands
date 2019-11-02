package main

import (
	"github.com/jiangjincc/islands/block"
)

func main() {
	b := block.CreateBlockchainWithGenesisBlock()
	_ = b.AddBlockToBlockChain([]byte("trade 100 RMB"))
	_ = b.AddBlockToBlockChain([]byte("trade 200 RMB"))
	_ = b.AddBlockToBlockChain([]byte("trade 500 RMB"))
	b.PrintBlocks()
	//fmt.Println(b)

}
