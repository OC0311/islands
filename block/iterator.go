package block

import (
	"log"
	"math/big"

	"github.com/boltdb/bolt"
)

type Iterator struct {
	DB          *bolt.DB
	CurrentHash []byte
}

func NewBlockIterator(db *bolt.DB, current []byte) *Iterator {
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_blockBucketName))
		current = b.Get([]byte(_topHash))
		return nil
	})

	return &Iterator{
		DB:          db,
		CurrentHash: current,
	}
}

func (i *Iterator) Next() (*Block, bool) {
	var (
		block  *Block
		isNext = true
	)

	err := i.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(_blockBucketName))

		blockBytes := bucket.Get(i.CurrentHash)
		block = UnSerialize(blockBytes)

		i.CurrentHash = block.PrevBlockHash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	// 判断是否是创世区块
	var hashInt big.Int
	if big.NewInt(0).Cmp(hashInt.SetBytes(i.CurrentHash)) == 0 {
		isNext = false
	}

	return block, isNext
}
