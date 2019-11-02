package utils

import (
	"bytes"
	"encoding/binary"
	"log"
)

func IntToHex(i int64) []byte {
	buff := new(bytes.Buffer)
	// 大端序
	err := binary.Write(buff, binary.BigEndian, i)
	if err != nil{
		log.Panic(err)
	}
	return buff.Bytes()
}
