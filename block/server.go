package block

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

var nodeAddress string
var knowAddress = []string{
	"localhost:3000",
}

func StartServer(nodeID string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	l, err := net.Listen("tcp", nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer l.Close()
	// 两个节点，主节点负责保存数据，钱包节点负责发送请求同步数据
	if nodeAddress != knowAddress[0] {
		sendMessage(knowAddress[0], []byte(nodeAddress))
	}
	bc := GetBlockchain()
	for {
		c, err := l.Accept()
		if err != nil {
			log.Printf("accept error:%v", err)
		}

		go handleConn(c, bc)
	}
}

// 处理请求
func handleConn(conn net.Conn, bc *Blockchain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}

	cmd := BytesToCommand(request[:_commandLen])
	switch cmd {
	case _version:
		handleSendVersion(request, bc)
	case _block:
		handleSendBlock(request, bc)
	case _getData:
		handleSendGetData(request, bc)
	case _getBlock:
		handleSendGetBlocks(request, bc)
	case _inv:
		handleSendInv(request, bc)
	}
}

func Encode(data interface{}) []byte {
	b, err := json.Marshal(data)
	if err != nil {
		log.Panic(err)
	}
	return b
}

func CommandToBytes(command string) []byte {
	var bytes [_commandLen]byte
	for i, e := range command {
		bytes[i] = byte(e)
	}
	return bytes[:]
}

func BytesToCommand(bytes []byte) string {
	var command []byte
	for _, b := range bytes {
		if b != 0x00 {
			command = append(command, b)
		}
	}

	return string(command)
}
