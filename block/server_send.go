package block

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func sendMessage(to string, from []byte) {
	fmt.Println("向服务器发送请求...")

	c, err := net.Dial("tcp", to)
	if err != nil {
		log.Panicf("connect to server [%s] failed %v\n", to, err)
	}

	_, err = io.Copy(c, bytes.NewReader([]byte(fmt.Sprintf("%s\n", from))))
	if err != nil {
		log.Panicf("copy bytes error:%v", err)
	}
}

func sendVersion(toAddress string, bc *Blockchain) {
	heigth := bc.GetHeight()
	versionData := Version{
		Heigth:  int(heigth),
		AddFrom: nodeAddress,
	}
	request := []byte{}
	data := Encode(versionData)
	request = append(CommandToBytes(_version), data...)
	sendMessage(toAddress, request)
}

// 同步数据
func sendGetBlocks(toAddress string) {
	request := []byte{}
	data := Encode(GetBlocks{AddFrom: nodeAddress})
	request = append(CommandToBytes(_getBlock), data...)
	sendMessage(toAddress, request)
}

// 展示数据
func sendInv(toAddress string, hash [][]byte) {
	request := []byte{}
	data := Encode(Inv{AddFrom: nodeAddress, Hashs: hash})
	request = append(CommandToBytes(_inv), data...)
	sendMessage(toAddress, request)
}

// 获取数据
func sendGetData(toAddress string, hash []byte) {
	request := []byte{}
	data := Encode(GetData{Addfrom: nodeAddress, Tx: hash})
	request = append(CommandToBytes(_getData), data...)
	sendMessage(toAddress, request)
}

func sendBlock(toAddress string, block []byte) {
	request := []byte{}
	data := Encode(BlockData{AddFrom: nodeAddress, Block: block})
	request = append(CommandToBytes(_block), data...)
	sendMessage(toAddress, request)
}
