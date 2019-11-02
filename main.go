package main

import (
	"fmt"
	"github.com/jiangjincc/islands/block"
	"github.com/boltdb/bolt"
	"log"
)

func init (){
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}


func main(){
	//b := block.CreateBlockchainWithGenesisBlock()
	//b.Add([]byte("trade 100"))
	//b.Add([]byte("trade 200"))
	//b.Add([]byte("trade 500"))
	//fmt.Println(b)


	// 验证区块的有效性
	bc := block.NewBlock([]byte("test"), 1, []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	pow := block.NewProofOfWork(bc)
	fmt.Println(bc.Serialize())
	fmt.Println(pow.IsValid())
}


