package wallet

import "fmt"

// 钱包
type Wallets struct {
	Wallets map[string]*Wallet
}

func NewWallets() *Wallets {
	return &Wallets{
		Wallets: make(map[string]*Wallet),
	}
}

func (w *Wallets) CreateNewWallet() {
	wallet := NewWallet()
	fmt.Printf("Address: %s \n", string(wallet.GetAddress()))
	w.Wallets[string(wallet.GetAddress())] = wallet
}
