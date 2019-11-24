package block

import (
	"bytes"

	"github.com/jiangjincc/islands/encryption"
)

type TXOutput struct {
	Value     int64
	PublicKey []byte
}

func (out *TXOutput) Lock(address string) {
	publickHash := encryption.Base58Decode([]byte(address))
	// 获取中间的public,标示是自己的花费
	out.PublicKey = publickHash[1 : len(publickHash)-4]
}

func (out *TXOutput) UnLockWithAddress(address string) bool {
	publickHash := encryption.Base58Decode([]byte(address))
	// 获取中间的public
	hash160 := publickHash[1 : len(publickHash)-4]
	return bytes.Compare(out.PublicKey, hash160) == 0
}

func NewTxOutput(value int64, address string) *TXOutput {
	// 表示是谁的金额
	txOutput := &TXOutput{
		Value:     value,
		PublicKey: nil,
	}
	// 设置hash
	txOutput.Lock(address)
	return txOutput
}
