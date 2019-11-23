package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
)

func IntToHex(i int64) []byte {
	buff := new(bytes.Buffer)
	// 大端序
	err := binary.Write(buff, binary.BigEndian, i)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func JsonToArray(jsonStr string) []string {
	var (
		res []string
	)

	if err := json.Unmarshal([]byte(jsonStr), &res); err != nil {
		panic(err)
	}
	return res
}

// 字节数组反转
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
