package block

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
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
		Vout:      -1,
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

	coinBase.SetHash()
	return coinBase
}

// 2、转账产生的交易
func NewSimpleTransaction(from, to string, amount int64, bc *Blockchain, tx []*Transaction) *Transaction {
	var (
		txInputs  []*TXInput
		txOutputs []*TXOutput
	)
	//// 找到from 所有未花费的tx
	//utxo := bc.UTXOs(from)
	// 找到需要花费的utxo
	money, spendableUTXODic := bc.FindSpendableUTXOS(from, amount, tx)
	// 代表消费

	for hash, indexArray := range spendableUTXODic {
		hashByte, _ := hex.DecodeString(hash)
		for _, index := range indexArray {
			txInput := &TXInput{
				TxHash:    hashByte,
				Vout:      index,
				ScriptSig: from,
			}
			txInputs = append(txInputs, txInput)
		}
	}

	txOutPut := &TXOutput{
		Value:        amount,
		ScriptPubKey: to,
	}

	// 找零 - 转回给发送方
	txOutPut2 := &TXOutput{
		Value:        int64(money) - amount,
		ScriptPubKey: from,
	}

	txOutputs = append(txOutputs, txOutPut, txOutPut2)
	coinBase := &Transaction{
		TxHash: []byte{},
		In:     txInputs,
		Out:    txOutputs,
	}

	coinBase.SetHash()
	return coinBase
}

// 是否是创世区块的内容
func (t *Transaction) IsCoinbaseTransaction() bool {
	return len(t.In[0].TxHash) == 0 && t.In[0].Vout == -1
}

// TODO 使用protobuff
func (t *Transaction) SetHash() {

	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(t)

	if err != nil {
		log.Panic(err)
	}

	t.TxHash = result.Bytes()[:]
}
