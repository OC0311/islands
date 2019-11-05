package block

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/jiangjincc/islands/utils"
)

const (
	// 挖矿难度
	_targetBit = 16
)

type ProofOfWork struct {
	Block *Block
	// 难度
	target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	// 定义难度
	target := big.NewInt(1)
	target = target.Lsh(target, 256-_targetBit)

	return &ProofOfWork{
		Block:  block,
		target: target,
	}
}

func (pow *ProofOfWork) Run() ([]byte, int64) {

	var (
		nonce   = 0
		hashInt big.Int
		hash    [32]byte
	)

	for {
		dataBytes := pow.prepareData(nonce)
		hash = sha256.Sum256(dataBytes)
		fmt.Printf("\r正在挖矿: %x", hash)

		hashInt.SetBytes(hash[:])
		if pow.target.Cmp(&hashInt) == 1 {
			fmt.Print("\n")
			break
		}

		nonce++
	}

	return hash[:], int64(nonce)
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.Block.PrevBlockHash,
		pow.Block.HashTransaction(),
		utils.IntToHex(pow.Block.Timestamp),
		utils.IntToHex(int64(_targetBit)),
		utils.IntToHex(int64(nonce)),
		utils.IntToHex(int64(pow.Block.Height)),
	}, []byte{})

	return data
}

func (pow *ProofOfWork) IsValid() bool {
	var hashInt big.Int
	hashInt.SetBytes(pow.Block.Hash)
	if pow.target.Cmp(&hashInt) == 1 {
		return true
	}
	return false
}
