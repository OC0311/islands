package block

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Transaction struct {
	// 交易hash
	TxHash []byte

	// in 消费
	In []*TXInput
	//out 未花费
	Out []*TXOutput
}

// 1、创世区块产生的交易
func NewCoinBaseTransaction(address string) *Transaction {
	// 代表消费
	txInput := &TXInput{
		TxHash:    []byte{},
		Vout:      0,
		ScriptSig: "",
	}

	txOutPut := &TXOutput{
		Value:        10,
		ScriptPubKey: address,
	}

	coinBase := &Transaction{
		TxHash: []byte{},
		In:     []*TXInput{txInput},
		Out:    []*TXOutput{txOutPut},
	}

	coinBase.Serialize()
	return coinBase
}

// 2、转账产生的交易

// TODO 使用protobuff
func (t *Transaction) Serialize() {

	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(t)

	if err != nil {
		log.Panic(err)
	}

	t.TxHash = result.Bytes()[:]
}
