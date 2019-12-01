package block

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/jiangjincc/islands/utils"

	"github.com/jiangjincc/islands/encryption"

	"github.com/boltdb/bolt"
)

const (
	_utxoTableName = "utxoTableName.db"
)

type UTXOSet struct {
	Blockchain *Blockchain
}

func (u *UTXOSet) ResetUTXOSet() map[string]*TxOutOuts {

	err := u.Blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_utxoTableName))
		if b != nil {
			err := tx.DeleteBucket([]byte(_utxoTableName))
			if err != nil {
				log.Panic(err)
			}

		}

		b, _ = tx.CreateBucket([]byte(_utxoTableName))
		if b != nil {
			txOutputMap := u.Blockchain.FindUTXOMap()
			for keyHash, outs := range txOutputMap {
				txHash, _ := hex.DecodeString(keyHash)
				//for i := 0; i < len(outs.UTXOS); i++ {
				//	fmt.Println(outs.UTXOS[i].OutPut.PublicKey, outs.UTXOS[i].OutPut.Value)
				//	fmt.Println()
				//}

				b.Put(txHash, outs.Serialize())
			}
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	return nil
}

func (u *UTXOSet) findUTXOForAddress(address string) []*UTXO {
	var utxos []*UTXO

	u.Blockchain.DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(_utxoTableName))

		// 游标
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			txOutputs := UnSerializeTxOutOuts(v)

			for _, utxo := range txOutputs.UTXOS {

				if utxo.OutPut.UnLockWithAddress(address) {
					utxos = append(utxos, utxo)
				}
			}
		}

		return nil
	})

	return utxos
}

func (u *UTXOSet) GetBalance(address string) int64 {
	utxos := u.findUTXOForAddress(address)
	var amount int64
	for _, utxo := range utxos {
		amount += utxo.OutPut.Value
	}

	return amount
}

// 返回要凑多少钱
func (u *UTXOSet) FindUnPackageSpendableUTXOS(from string, txs []*Transaction) []*UTXO {
	var (
		spentTxs = make(map[string][]int)
		unUTXOs  []*UTXO
	)

	// 处理未打包区块
	for _, tx := range txs {
		if !tx.IsCoinbaseTransaction() {
			// 是否是address 的花费
			for _, in := range tx.In {

				publickHash := encryption.Base58Decode([]byte(from))
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
			if out.UnLockWithAddress(from) {
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

	return unUTXOs
}

func (u *UTXOSet) FindUnSpendUTXOS(from string, amount int64, txs []*Transaction) (int64, map[string][]int) {
	// 查找还没有打包的未花费
	unPackageSpendableUTXO := u.FindUnPackageSpendableUTXOS(from, txs)

	var money int64 = 0
	spendableUTXO := make(map[string][]int)
	for _, utxo := range unPackageSpendableUTXO {
		money += utxo.OutPut.Value
		txHash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[txHash] = append(spendableUTXO[txHash], utxo.Index)
		if money >= amount {
			return money, spendableUTXO
		}
	}

	u.Blockchain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_utxoTableName))
		if b == nil {
			log.Panic(b)
		}
		c := b.Cursor()
	out:
		for k, v := c.First(); k != nil; k, v = c.Next() {
			txOutputs := UnSerializeTxOutOuts(v)

			for _, out := range txOutputs.UTXOS {
				money += out.OutPut.Value
				txHash := hex.EncodeToString(out.TxHash)
				spendableUTXO[txHash] = append(spendableUTXO[txHash], out.Index)
				if money >= amount {
					break out

				}
			}

		}

		return nil
	})

	if money < amount {
		log.Panic("余额不足")
	}
	return money, spendableUTXO
}

// 更新
func (u *UTXOSet) Update() {
	var (
		currentHash []byte = u.Blockchain.Tip
	)

	iterator := NewBlockIterator(u.Blockchain.DB, currentHash)
	block, _ := iterator.Next()

	ins := []*TXInput{}
	outsMap := make(map[string]*TxOutOuts)
	// 找到所有我要删除的数据
	for _, tx := range block.Txs {
		for _, in := range tx.In {
			ins = append(ins, in)
		}
	}

	for _, tx := range block.Txs {
		utxos := []*UTXO{}
		for index, out := range tx.Out {
			isSpend := false
			for _, in := range ins {
				// 说明这个out 已经被消费
				if in.Vout == index && bytes.Compare(tx.TxHash, in.TxHash) == 0 && bytes.Compare(out.PublicKey, utils.Ripemd160Hash(in.PublicKey)) == 0 {
					// 说明out已经被花费
					isSpend = true
					continue
				}

			}
			if !isSpend {
				utxo := &UTXO{TxHash: tx.TxHash, Index: index, OutPut: out}
				utxos = append(utxos, utxo)
			}

		}

		if len(utxos) > 0 {

			txHash := hex.EncodeToString(tx.TxHash)
			outsMap[txHash] = &TxOutOuts{utxos}
		}
	}
	// 删除
	err := u.Blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_utxoTableName))
		if b != nil {
			for _, in := range ins {
				txOutputBytes := b.Get(in.TxHash)
				if len(txOutputBytes) == 0 {
					continue
				}
				txOutputs := UnSerializeTxOutOuts(txOutputBytes)
				utxos := []*UTXO{}
				isNeedDel := false
				for _, utxo := range txOutputs.UTXOS {

					if in.Vout == utxo.Index &&
						bytes.Compare(utxo.OutPut.PublicKey, utils.Ripemd160Hash(in.PublicKey)) == 0 {

						isNeedDel = true
					} else {

						utxos = append(utxos, utxo)
					}
				}

				if isNeedDel {
					err := b.Delete(in.TxHash)

					if err != nil {
						log.Panic(err)
					}
					if len(utxos) > 0 {

						preTxOutputs := outsMap[hex.EncodeToString(in.TxHash)]
						preTxOutputs.UTXOS = append(preTxOutputs.UTXOS, utxos...)

						outsMap[hex.EncodeToString(in.TxHash)] = preTxOutputs
						//b.Put(in.TxHash, txOutputs.Serialize())
					}
				}
			}

			// 新增
			for keyHash, outPuts := range outsMap {

				keyBytes, _ := hex.DecodeString(keyHash)
				fmt.Println(outPuts.UTXOS)
				err := b.Put(keyBytes, outPuts.Serialize())
				if err != nil {
					log.Panic(err)
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

// 更新
func (u *UTXOSet) UpdateV2() {
	var (
		currentHash []byte = u.Blockchain.Tip
	)

	iterator := NewBlockIterator(u.Blockchain.DB, currentHash)
	block, _ := iterator.Next()

	ins := []*TXInput{}
	outsMap := make(map[string]*TxOutOuts)
	// 找到所有我要删除的数据
	for _, tx := range block.Txs {
		for _, in := range tx.In {
			ins = append(ins, in)
		}
	}

	for _, tx := range block.Txs {
		utxos := []*UTXO{}
		for index, out := range tx.Out {
			isSpent := false
			for _, in := range ins {
				if in.Vout == index && bytes.Compare(tx.TxHash, in.TxHash) == 0 && bytes.Compare(out.PublicKey, utils.Ripemd160Hash(in.PublicKey)) == 0 {
					isSpent = true
					continue
				}
			}

			if isSpent == false {
				utxo := &UTXO{tx.TxHash, index, out}
				utxos = append(utxos, utxo)
			}

		}

		if len(utxos) > 0 {
			txHash := hex.EncodeToString(tx.TxHash)
			outsMap[txHash] = &TxOutOuts{utxos}
		}

	}

	err := u.Blockchain.DB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(_utxoTableName))

		if b != nil {

			// 删除
			for _, in := range ins {

				txOutputsBytes := b.Get(in.TxHash)

				if len(txOutputsBytes) == 0 {
					continue
				}

				txOutputs := UnSerializeTxOutOuts(txOutputsBytes)

				UTXOS := []*UTXO{}

				// 判断是否需要
				isNeedDelete := false

				for _, utxo := range txOutputs.UTXOS {

					if in.Vout == utxo.Index && bytes.Compare(utxo.OutPut.PublicKey, utils.Ripemd160Hash(in.PublicKey)) == 0 {

						isNeedDelete = true
					} else {
						UTXOS = append(UTXOS, utxo)
					}
				}

				if isNeedDelete {
					b.Delete(in.TxHash)
					if len(UTXOS) > 0 {

						preTXOutputs := outsMap[hex.EncodeToString(in.TxHash)]

						preTXOutputs.UTXOS = append(preTXOutputs.UTXOS, UTXOS...)

						outsMap[hex.EncodeToString(in.TxHash)] = preTXOutputs

					}
				}

			}

			// 新增

			for keyHash, outPuts := range outsMap {
				keyHashBytes, _ := hex.DecodeString(keyHash)
				b.Put(keyHashBytes, outPuts.Serialize())
			}

		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}
