package wallet

import (
	"bytes"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ripemd160"
)

const (
	walletFile = "wallets.dat"
)

// 钱包
type Wallets struct {
	Wallets map[string]*Wallet
}

func NewWallets() (*Wallets, error) {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets := &Wallets{Wallets: make(map[string]*Wallet)}
		return wallets, nil
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallet Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallet)
	if err != nil {
		log.Panic(err)
	}
	return &wallet, nil
}

func (w *Wallets) CreateNewWallet() {
	wallet := NewWallet()
	fmt.Printf("Address: %s \n", string(wallet.GetAddress()))
	w.Wallets[string(wallet.GetAddress())] = wallet
	w.SaveToFile()
}

func (w *Wallets) Ripemd160Hash(publicKey []byte) []byte {
	hash256 := sha256.New()
	hash256.Write(publicKey)
	hash := hash256.Sum(nil)

	ripemd := ripemd160.New()
	ripemd.Write(hash)
	ripemdHash := ripemd.Sum(nil)
	return ripemdHash
}

func (w *Wallets) SaveToFile() {
	var content bytes.Buffer
	// 为了可以序列化任何类型
	// 因为有可能会存放接口类型
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(w)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

func (w *Wallets) WalletList() {
	for address, _ := range w.Wallets {
		fmt.Printf("Address:%s\n", address)
	}
}
