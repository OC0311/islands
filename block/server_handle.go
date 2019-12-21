package block

import (
	"encoding/json"
	"fmt"
	"log"
)

// 同步数据
func handleSendGetBlocks(request []byte, bc *Blockchain) {
	fmt.Println("send block handle...")
	var data GetBlocks
	err := json.Unmarshal(request[13:], &data)
	if err != nil {
		log.Panic(err)
	}

	// 获取所有的区块hash
	hashs := bc.GetBlockHashs()
	sendInv(data.AddFrom, hashs)
}

// 展示数据
func handleSendInv(request []byte, bc *Blockchain) {
	fmt.Println("version handle...")
	var data Inv
	err := json.Unmarshal(request[13:], &data)
	if err != nil {
		log.Panic(err)
	}
	for _, hash := range data.Hashs {
		sendGetData(data.AddFrom, hash)
	}

}

// 获取数据
func handleSendGetData(request []byte, bc *Blockchain) {
	fmt.Println("version handle...")
	var data GetData
	err := json.Unmarshal(request[13:], &data)
	if err != nil {
		log.Panic(err)
	}
	// 通过传过来的hash 获取区块数据
	blockBytes := bc.GetBlock(data.Tx)
	sendBlock(data.Addfrom, blockBytes)
}

// 接受到新区块
func handleSendBlock(request []byte, bc *Blockchain) {
	fmt.Println("version handle...")
	var data BlockData
	err := json.Unmarshal(request[13:], &data)
	if err != nil {
		log.Panic(err)
	}
	// 添加到区块链中
	blockBytes := data.Block
	bc.AddBlock(UnSerialize(blockBytes))
}

func handleSendVersion(request []byte, bc *Blockchain) {
	fmt.Println("version handle...")
	var data Version
	err := json.Unmarshal(request[13:], &data)
	if err != nil {
		log.Panic(err)
	}

	versionHeigth := data.Heigth
	heigth := bc.GetHeight()
	if heigth > int64(versionHeigth) {
		sendVersion(data.AddFrom, bc)
	} else if heigth < int64((versionHeigth)) {
		sendGetBlocks(data.AddFrom)
	}
}
