package block

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"

	"github.com/jiangjincc/islands/encryption"

	"github.com/jiangjincc/islands/wallet"

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
func CreateBlockchainWithGenesisBlock(address string) {
	// 判断数据库文件是否存在
	if dbIsExist(_dbName) {
		fmt.Println("区块已经存在")
		return
	}

	if !wallet.IsValidForAddress([]byte(address)) {
		log.Println("无效地址")
		os.Exit(1)
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

		genesisBlock := CreateGenesisBlock([]*Transaction{NewCoinBaseTransaction(address)})

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
		os.Exit(0)
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
func (bc *Blockchain) AddBlockToBlockChain(data []*Transaction) error {

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

func (bc *Blockchain) MineNewBlock(from, to, amount []string) {
	// 处理交易逻辑
	var (
		block *Block
		txs   []*Transaction
	)

	for index, address := range from {
		if !wallet.IsValidForAddress([]byte(address)) || !wallet.IsValidForAddress([]byte(to[index])) {
			log.Println("无效地址")
			os.Exit(1)
		}
	}

	for i, address := range from {
		// 构建多个交易
		a, _ := strconv.Atoi(amount[i])
		txs = append(txs, NewSimpleTransaction(address, to[i], int64(a), bc, txs))
	}

	// 添加矿工奖励
	coinBase := NewCoinBaseTransaction(from[0])
	txs = append(txs, coinBase)

	// 获取最新区块
	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_blockBucketName))
		if b != nil {
			hash := b.Get([]byte(_topHash))
			blockBytes := b.Get([]byte(hash))
			block = UnSerialize(blockBytes)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	// 对交易签名进行校验
	for _, tx := range txs {
		if !bc.VerifyTransaction(tx) {
			fmt.Println(hex.EncodeToString(tx.TxHash))
			log.Panic("无效的交易")
		}
	}

	newBlock := NewBlock(txs, block.Height+1, block.Hash)
	err = bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_blockBucketName))
		if b != nil {
			err := b.Put([]byte(newBlock.Hash), newBlock.Serialize())
			if err != nil {
				return err
			}

			err = b.Put([]byte(_topHash), []byte(newBlock.Hash))
			if err != nil {
				return err
			}

			bc.Tip = newBlock.Hash
		}
		// 更新最新区块的信息
		return nil
	})

}

// 找到需要花费的utxo
// TODO 需要找到最合适的utxo
func (bc *Blockchain) FindSpendableUTXOS(from string, amount int64, txs []*Transaction) (int64, map[string][]int) {

	// 找到合适的utxo 拿出来花费
	var (
		value         int64
		allUTXO       = bc.UTXOs(from, txs)
		spendableUTXO = make(map[string][]int)
	)

	for _, out := range allUTXO {
		value += out.OutPut.Value
		hash := hex.EncodeToString(out.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], out.Index)
		if value >= amount {
			break
		}
	}

	if value < amount {
		fmt.Printf("账户【%s】余额不足：%d \n", from, value)
		os.Exit(0)
	}

	return value, spendableUTXO
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

func (bc *Blockchain) GetBalance(address string) {
	var (
		amount int64 = 0
	)
	if !wallet.IsValidForAddress([]byte(address)) {
		log.Println("无效地址")
		os.Exit(1)
	}
	txs := bc.UTXOs(address, []*Transaction{})
	for _, v := range txs {
		amount += v.OutPut.Value
	}
	fmt.Printf("账户[ %s ]余额为: %d \n", address, amount)
}

func (bc *Blockchain) UTXOs(address string, txs []*Transaction) []*UTXO {
	var (
		currentHash []byte = bc.Tip
		spentTxs           = make(map[string][]int)
		unUTXOs     []*UTXO
	)

	// 处理未打包区块
	for _, tx := range txs {
		if !tx.IsCoinbaseTransaction() {
			// 是否是address 的花费
			for _, in := range tx.In {

				publickHash := encryption.Base58Decode([]byte(address))
				// 获取中间的public,标示是自己的花费
				ripemd160Hash := publickHash[1 : len(publickHash)-4]
				if in.UnLockWithAddress(ripemd160Hash) {
					key := hex.EncodeToString(in.TxHash)
					spentTxs[key] = append(spentTxs[key], in.Vout)
				}
			}
		}
	}

	for _, tx := range txs {
	work1:
		for index, out := range tx.Out {
			if out.UnLockWithAddress(address) {
				if len(spentTxs) == 0 {
					utxo := &UTXO{
						TxHash: tx.TxHash,
						Index:  index,
						OutPut: out,
					}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash, indexArray := range spentTxs {
						// 说明已经被花费
						if hash == hex.EncodeToString(tx.TxHash) {
							var isUnSpentUTXO = false
							for _, outIndex := range indexArray {
								if index == outIndex {
									isUnSpentUTXO = true
									continue work1
								}
								if !isUnSpentUTXO {
									utxo := &UTXO{
										TxHash: tx.TxHash,
										Index:  index,
										OutPut: out,
									}
									unUTXOs = append(unUTXOs, utxo)
								}

							}
						} else {
							utxo := &UTXO{
								TxHash: tx.TxHash,
								Index:  index,
								OutPut: out,
							}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}

			}

		}
	}

	iterator := NewBlockIterator(bc.DB, currentHash)
	for {
		block, isNext := iterator.Next()

		for i := len(block.Txs) - 1; i >= 0; i-- {
			tx := block.Txs[i]

			if !tx.IsCoinbaseTransaction() {
				// 是否是address 的花费
				for _, in := range tx.In {
					publickHash := encryption.Base58Decode([]byte(address))
					// 获取中间的public,标示是自己的花费
					ripemd160Hash := publickHash[1 : len(publickHash)-4]
					if in.UnLockWithAddress(ripemd160Hash) {
						key := hex.EncodeToString(in.TxHash)
						spentTxs[key] = append(spentTxs[key], in.Vout)
					}
				}
			}

			// 是否为自己的未花费
		work:
			for index, out := range tx.Out {
				if out.UnLockWithAddress(address) {
					// 判断是否被花费
					if len(spentTxs) != 0 {
						isSpend := false
						for hash, indexArray := range spentTxs {
							for _, i := range indexArray {
								// 说明已经被花费
								if i == index && hash == hex.EncodeToString(tx.TxHash) {
									isSpend = true
									continue work
								}
							}

						}
						if !isSpend {
							utxo := &UTXO{
								TxHash: tx.TxHash,
								Index:  index,
								OutPut: out,
							}
							unUTXOs = append(unUTXOs, utxo)
						}
					} else {
						utxo := &UTXO{
							TxHash: tx.TxHash,
							Index:  index,
							OutPut: out,
						}
						unUTXOs = append(unUTXOs, utxo)
					}

				}

			}
		}

		if !isNext {
			break
		}
	}

	return unUTXOs
}

func (bc *Blockchain) SignTransaction(tx *Transaction, priKey ecdsa.PrivateKey) {
	if tx.IsCoinbaseTransaction() {
		return
	}

	txMap := make(map[string]Transaction)
	for _, in := range tx.In {
		prevTx, err := bc.FindTransaction(in.TxHash)
		if err != nil {
			log.Panic(err)
		}
		txMap[hex.EncodeToString(prevTx.TxHash)] = prevTx
	}

	//  对每一笔交易签名
	tx.Sign(priKey, txMap)
}

// 根据input ID查找交易
func (bc *Blockchain) FindTransaction(id []byte) (Transaction, error) {
	var (
		currentHash []byte = bc.Tip
	)

	iterator := NewBlockIterator(bc.DB, currentHash)
	for {
		block, _ := iterator.Next()
		for _, tx := range block.Txs {
			if bytes.Compare(tx.TxHash, id) == 0 {
				return *tx, nil
			}
		}

		// 判断是否为创世区块
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}

	}

	return Transaction{}, nil
}

// 验证交易
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	txMap := make(map[string]Transaction)
	for _, in := range tx.In {
		prevTx, err := bc.FindTransaction(in.TxHash)
		if err != nil {
			log.Panic(err)
		}
		txMap[hex.EncodeToString(prevTx.TxHash)] = prevTx
	}

	return tx.Verify(txMap)
}
