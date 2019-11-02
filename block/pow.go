package block

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/jiangjincc/islands/utils"
	"math/big"
)

const(
	_targetBit = 8
)

type ProofOfWork struct{
	Block *Block
	// 难度
	target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	// 定义难度
	target := big.NewInt(1)
	target = target.Lsh(target, 256-_targetBit)

	return &ProofOfWork{
		Block:block,
		target:target,
	}
}

func (pow *ProofOfWork)Run()([]byte, int64){
	// 生成hash 判断有效性
	nonce := 0
	var hashInt big.Int
	var hash [32]byte
	for{
		dataBytes := pow.prepareData(nonce)
		hash = sha256.Sum256(dataBytes)
		fmt.Printf("\r%x\n", hash )

		hashInt.SetBytes(hash[:])
		if pow.target.Cmp(&hashInt) == 1{
			break
		}
		nonce ++
	}

	return hash[:], int64(nonce)
}

func (pow *ProofOfWork)prepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.Block.PrevBlockHash,
		pow.Block.Data,
		utils.IntToHex(pow.Block.Timestamp),
		utils.IntToHex(int64(_targetBit)),
		utils.IntToHex(int64(nonce)),
		utils.IntToHex(int64(pow.Block.Height)),

	},[]byte{})

	return data
}

func (pow *ProofOfWork) IsValid() bool {
	var hashInt big.Int
	hashInt.SetBytes(pow.Block.Hash)
	if pow.target.Cmp(&hashInt) == 1{
		return true
	}
	return false
}
