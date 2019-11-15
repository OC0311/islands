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
