package block

import (
	"log"
	"math/big"

	"github.com/boltdb/bolt"
)

const (
	_genesisBlockHeight = 1
	_dbName             = "blockchain.db"
	_blockBucketName    = "blocks"
	_topHash            = "top_hash"
)

// 存储有序的区块
type Blockchain struct {
	Tip []byte // 最新区块的hash
	DB  *bolt.DB
}

// 生成创世区块函数的blockchain
func CreateBlockchainWithGenesisBlock() *Blockchain {
	var (
		blockHash []byte
	)
	db, err := bolt.Open(_dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(_blockBucketName))
		if err != nil {
			return err
		}

		data := "Genesis Block"
		genesisBlock := CreateGenesisBlock([]byte(data))

		err = bucket.Put([]byte(genesisBlock.Hash), genesisBlock.Serialize())
		if err != nil {
			return err
		}
		blockHash = genesisBlock.Hash
		// save last hash
		err = bucket.Put([]byte(_topHash), blockHash)

		return err
	})

	if err != nil {
		log.Panic(err)
	}

	blockChain := &Blockchain{
		Tip: blockHash,
		DB:  db,
	}
	return blockChain
}

// 添加新区块到链中
func (bc *Blockchain) AddBlockToBlockChain(data []byte) error {

	err := bc.DB.Update(func(tx *bolt.Tx) error {
		// 获取最新区块的信息
		bucket := tx.Bucket([]byte(_blockBucketName))

		topHash := bc.Tip
		if topHash == nil {
			topHash = bucket.Get([]byte(_topHash))
		}

		prevBlockBytes := bucket.Get(topHash)
		prevBlock := UnSerialize(prevBlockBytes)

		// 创建新的区块
		block := NewBlock(data, prevBlock.Height+1, prevBlock.Hash)

		// 存储新区块
		err := bucket.Put(block.Hash, block.Serialize())
		if err != nil {
			return err
		}

		bc.Tip = block.Hash
		err = bucket.Put([]byte(_topHash), bc.Tip)
		return err
	})

	return err
}

func (bc *Blockchain) PrintBlocks() error {
	var (
		currentHash []byte = bc.Tip
	)

	iterator := NewBlockIterator(bc.DB, currentHash)
	for {
		block := iterator.Next()
		block.PrintBlock()
		// 判断是否是创世去区块
		var hashInt big.Int
		if big.NewInt(0).Cmp(hashInt.SetBytes(block.PrevBlockHash)) == 0 {
			break
		}

		currentHash = block.PrevBlockHash
	}

	return nil
}
