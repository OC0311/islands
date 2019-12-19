package block

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
)

var nodeAddress string
var knowAddress = []string{
	"localhost:3000",
}

func startServer(nodeID string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	l, err := net.Listen("tcp", nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer l.Close()
	// 两个节点，主节点负责保存数据，钱包节点负责发送请求同步数据
	if nodeAddress != knowAddress[0] {
		sendMessage(knowAddress[0], nodeAddress)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			log.Printf("accept error:%v", err)
		}

		request, err := ioutil.ReadAll(c)
		if err != nil {
			log.Printf("revevc error:%s", err)
		}
		fmt.Printf("receive data:%v", request)
	}
}

func sendMessage(to, from string) {
	fmt.Println("向服务器发送请求...")

	c, err := net.Dial("tcp", to)
	if err != nil {
		log.Panicf("connect to server [%s] failed %v\n", to, err)
	}

	_, err = io.Copy(c, bytes.NewReader([]byte(from)))
	if err != nil {
		log.Panicf("copy bytes error:%v", err)
	}
}
