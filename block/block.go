package block

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"

	"github.com/jedib0t/go-pretty/table"
)

type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	Data          []byte
	Hash          []byte
	Height        int64
	Nonce         int64
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
	powIns := NewProofOfWork(block)
	hash, nonce := powIns.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func CreateGenesisBlock(data []byte) *Block {
	return NewBlock(data, 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

// TODO 使用protobuff
func (b *Block) Serialize() []byte {

	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)

	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func (b *Block) PrintBlock() {

	tb := table.NewWriter()
	tb.SetOutputMirror(os.Stdout)
	tb.AppendHeader(table.Row{"内容", "区块信息"})
	tb.AppendRows([]table.Row{
		{"Height", b.Height},
		{"Data", string(b.Data)},
		{"Timestamp", time.Unix(b.Timestamp, 0).Format("2006-01-02 15:04:05")},
		{"Nonce", b.Nonce},
		{"Hash", byteForString(b.Hash)},
		{"PrevHash", byteForString(b.PrevBlockHash)},
	})
	tb.SetStyle(table.StyleDefault)
	tb.Render()
}

func UnSerialize(bt []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(bt))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

func byteForString(b []byte) string {
	return fmt.Sprintf("%x", b)
}

type Iterator struct {
	DB          *bolt.DB
	CurrentHash []byte
}

func NewBlockIterator(db *bolt.DB, current []byte) *Iterator {
	return &Iterator{
		DB:          db,
		CurrentHash: current,
	}
}

func (i *Iterator) Next() *Block {
	var (
		block *Block
	)

	err := i.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(_blockBucketName))

		blockBytes := bucket.Get(i.CurrentHash)
		block = UnSerialize(blockBytes)

		// 打印区块信息
		// block.PrintBlock()

		i.CurrentHash = block.PrevBlockHash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return block
}