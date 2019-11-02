package main

import (
	"log"

	"github.com/boltdb/bolt"
	"github.com/jiangjincc/islands/block"
)

const (
	_blockBucketName = "blocks"
)

var (
	db *bolt.DB
)

func init() {
	var err error

	db, err = bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	defer db.Close()
	b := block.CreateBlockchainWithGenesisBlock()
	_ = b.AddBlockToBlockChain([]byte("trade 100 RMB"))
	_ = b.AddBlockToBlockChain([]byte("trade 200 RMB"))
	_ = b.AddBlockToBlockChain([]byte("trade 500 RMB"))
	b.PrintBlocks()
	//fmt.Println(b)

}
