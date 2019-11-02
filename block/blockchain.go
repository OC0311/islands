package block

const(
	_genesisBlockHeight = 1
)

// 存储有序的区块
type Blockchain struct{
	Blocks []*Block
}

// 生成创世区块函数的blockchain
func CreateBlockchainWithGenesisBlock() *Blockchain {
	data := "Genesis Block"
	genesisBlock := CreateGenesisBlock([]byte(data))
	blockChain := new(Blockchain)
	blockChain.Blocks = append(blockChain.Blocks,genesisBlock)
	return blockChain
}


// 添加新区块到链中
func (bc *Blockchain)Add(data []byte){
	blockHeight := len(bc.Blocks)
	prevBlock := bc.Blocks[blockHeight-1]
	block := NewBlock(data, prevBlock.Height+1, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, block)
}