package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	Data          []byte
	Hash          []byte
	Height        int64
	Nonce int64
}

// 创建新的区块
func NewBlock(data []byte, height int64, PrevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().UTC().Unix(),
		PrevBlockHash: PrevBlockHash,
		Data:          data,
		Height:        height,
	}

	// 工作量证明
	powIns :=  NewProofOfWork(block)
	hash, nonce := powIns.Run()

	block.Hash  = hash[:]
	block.Nonce = nonce
	//block.SetHash()
	return block
}


func CreateGenesisBlock(data []byte) *Block {
	return NewBlock(data, 1, []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}

func (b *Block)SetHash(){
	height := utils.IntToHex(b.Height)
	timeString := strconv.FormatInt(b.Timestamp,2)
	timeBytes := []byte(timeString)

	all :=bytes.Join([][]byte{
		height, timeBytes, b.Data, b.PrevBlockHash, b.Hash,
	}, []byte{})

	hash := sha256.Sum256(all)
	b.Hash = hash[:]
}


// TODO 使用protobuff
func (b *Block)Serialize() []byte {

	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)

	if err != nil{
		log.Panic(err)
	}

	return result.Bytes()
}

func UnSerialize(bt []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(bt))
	err := decoder.Decode(&block)
	if err != nil{
		log.Panic(err)
	}

	return &block
}



