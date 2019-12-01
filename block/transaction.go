package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"math/big"
	"time"

	"github.com/jiangjincc/islands/utils"
	"github.com/jiangjincc/islands/wallet"
)

type Transaction struct {
	// 交易hash
	TxHash []byte

	// in 消费
	In []*TXInput
	// out 未花费
	Out []*TXOutput
}

// 1、创世区块产生的交易
func NewCoinBaseTransaction(address string) *Transaction {
	// 代表消费
	txInput := &TXInput{
		TxHash:    []byte{},
		Vout:      -1,
		PublicKey: nil,
		Signature: nil,
	}

	txOutPut := NewTxOutput(10, address)

	coinBase := &Transaction{
		TxHash: []byte{},
		In:     []*TXInput{txInput},
		Out:    []*TXOutput{txOutPut},
	}
	coinBase.SetHash()
	return coinBase
}

// 2、转账产生的交易
func NewSimpleTransaction(from, to string, amount int64, bc *UTXOSet, tx []*Transaction) *Transaction {
	var (
		txInputs  []*TXInput
		txOutputs []*TXOutput
	)
	ws, _ := wallet.NewWallets()
	w := ws.Wallets[from]
	//// 找到from 所有未花费的tx

	// 找到需要花费的utxo
	//money, spendableUTXODic := bc.Blockchain.FindSpendableUTXOS(from, amount, tx)
	money, spendableUTXODic := bc.Blockchain.FindSpendableUTXOS(from, amount, tx)
	// 代表消费

	for hash, indexArray := range spendableUTXODic {
		hashByte, _ := hex.DecodeString(hash)
		for _, index := range indexArray {
			txInput := &TXInput{
				TxHash:    hashByte,
				Vout:      index,
				PublicKey: w.PublicKey,
				Signature: nil,
			}
			txInputs = append(txInputs, txInput)
		}
	}

	txOutPut := NewTxOutput(amount, to)

	// 找零 - 转回给发送方
	txOutPut2 := NewTxOutput(int64(money)-amount, from)

	txOutputs = append(txOutputs, txOutPut, txOutPut2)
	transaction := &Transaction{
		TxHash: []byte{},
		In:     txInputs,
		Out:    txOutputs,
	}

	transaction.SetHash()
	// 对交易进行签名
	bc.Blockchain.SignTransaction(transaction, w.PrivateKey, tx)
	return transaction
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

	hash := sha256.Sum256(bytes.Join([][]byte{utils.IntToHex(time.Now().Unix()), result.Bytes()}, []byte{}))

	t.TxHash = hash[:]
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {

	if tx.IsCoinbaseTransaction() {
		return
	}

	for _, in := range tx.In {
		if prevTXs[hex.EncodeToString(in.TxHash)].TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.Copy()

	for inID, vin := range txCopy.In {
		prevTx := prevTXs[hex.EncodeToString(vin.TxHash)]
		txCopy.In[inID].Signature = nil
		txCopy.In[inID].PublicKey = prevTx.Out[vin.Vout].PublicKey
		txCopy.TxHash = txCopy.Hash()
		txCopy.In[inID].PublicKey = nil

		// 签名代码
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.TxHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.In[inID].Signature = signature
	}
}

func (t *Transaction) Sign2(priKey ecdsa.PrivateKey, txs map[string]Transaction) {
	// 判断是coinbase 交易
	if t.IsCoinbaseTransaction() {
		return
	}

	// 判断是否查询到了input
	for _, in := range t.In {
		if txs[hex.EncodeToString(in.TxHash)].TxHash == nil {
			log.Panic("根据当前区块的交易输入未找到交易")
		}
	}

	// 将当前交易进行拷贝进行签名
	txCopy := t.Copy()

	// 签名
	for inHash, in := range txCopy.In {
		prevTx := txs[hex.EncodeToString(in.TxHash)]
		txCopy.In[inHash].Signature = nil
		txCopy.In[inHash].PublicKey = prevTx.Out[in.Vout].PublicKey
		txCopy.TxHash = txCopy.Hash()
		txCopy.In[inHash].PublicKey = nil

		// 使用私钥对hash 进行签名
		r, s, err := ecdsa.Sign(rand.Reader, &priKey, txCopy.TxHash)
		if err != nil {
			log.Panic(err)
		}

		signature := append(r.Bytes(), s.Bytes()...)
		// 设置每笔花费的签名
		t.In[inHash].Signature = signature
	}
}

func (t *Transaction) Verify(txMap map[string]Transaction) bool {
	if t.IsCoinbaseTransaction() {
		return true
	}

	// 判断是否查询到了input
	for _, in := range t.In {
		if txMap[hex.EncodeToString(in.TxHash)].TxHash == nil {
			log.Panic("根据当前区块的交易输入未找到交易")
		}
	}

	// 将当前交易进行拷贝进行签名
	txCopy := t.Copy()
	curve := elliptic.P256()
	// 签名
	for inHash, in := range t.In {
		prevTx := txMap[hex.EncodeToString(in.TxHash)]
		txCopy.In[inHash].Signature = nil
		txCopy.In[inHash].PublicKey = prevTx.Out[in.Vout].PublicKey
		txCopy.TxHash = txCopy.Hash()
		txCopy.In[inHash].PublicKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(in.Signature)
		r.SetBytes(in.Signature[:(sigLen / 2)])
		s.SetBytes(in.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		publickLen := len(in.PublicKey)
		x.SetBytes(in.PublicKey[:(publickLen / 2)])
		y.SetBytes(in.PublicKey[(publickLen / 2):])

		rawPubkey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubkey, txCopy.TxHash, &r, &s) == false {
			return false
		}
	}

	return true
}

// 深拷贝当前区块
func (t *Transaction) Copy() Transaction {
	var (
		outs   []*TXOutput
		inputs []*TXInput
	)

	for _, in := range t.In {
		inputs = append(inputs, &TXInput{TxHash: in.TxHash, Vout: in.Vout, PublicKey: nil, Signature: nil})
	}

	for _, out := range t.Out {
		outs = append(outs, &TXOutput{Value: out.Value, PublicKey: out.PublicKey})
	}

	copyTransaction := Transaction{TxHash: t.TxHash, In: inputs, Out: outs}
	return copyTransaction
}

func (t *Transaction) Hash() []byte {
	txCopy := t
	txCopy.TxHash = []byte{}
	hash := sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

func (t *Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(t)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}
