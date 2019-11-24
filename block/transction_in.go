package block

import (
	"bytes"

	"github.com/jiangjincc/islands/utils"
)

type TXInput struct {
	// 1. 交易的Hash
	TxHash []byte
	// 2. 存储TXOutput在Vout里面的索引
	Vout int

	// 数字签名可以防止花费别人的钱
	Signature []byte
	PublicKey []byte
}

func (in *TXInput) UnLockWithAddress(ripemd160Hash []byte) bool {
	publicKey := utils.Ripemd160Hash(in.PublicKey)
	return bytes.Compare(publicKey, ripemd160Hash) == 0
}
