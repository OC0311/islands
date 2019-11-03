package block

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const (
	_genesisBlockHeight = 1
	_dbName             = "/Users/admin/workspace/project/github.com/jiangjincc/islands/blockchain.db"
	_blockBucketName    = "blocks"
	_topHash            = "top_hash"
)

// 存储有序的区块
type Blockchain struct {
	Tip []byte // 最新区块的hash
	DB  *bolt.DB
}

// 生成创世区块函数的blockchain
func CreateBlockchainWithGenesisBlock() {
	// 判断数据库文件是否存在
	if dbIsExist(_dbName) {
		fmt.Println("区块已经存在")
		return
	}

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
		// save last hash
		err = bucket.Put([]byte(_topHash), genesisBlock.Hash)

		return err
	})

	if err != nil {
		log.Panic(err)
	}

}

func GetBlockchain() *Blockchain {
	var (
		blockchain *Blockchain
	)

	if !dbIsExist(_dbName) {
		fmt.Println("请初始化区块链")
		return nil
	}

	db, err := bolt.Open(_dbName, 0600, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	_ = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(_blockBucketName))
		topHash := bucket.Get([]byte(_topHash))

		blockchain = &Blockchain{
			Tip: topHash,
			DB:  db,
		}

		return nil
	})

	return blockchain
}

func dbIsExist(dbName string) bool {
	_, err := os.Open(dbName)

	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
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

func (bc *Blockchain) PrintBlocks() {
	var (
		currentHash []byte = bc.Tip
	)

	iterator := NewBlockIterator(bc.DB, currentHash)
	for {
		block, isNext := iterator.Next()
		block.PrintBlock()
		if !isNext {
			break
		}
	}
}
