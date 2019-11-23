package wallet

import (
	"testing"
)

func TestNewWallet(t *testing.T) {
	w := NewWallet()

	w.IsValidForAddress(w.GetAddress())
}

func TestNewWallets(t *testing.T) {
	ws := NewWallets()
	ws.CreateNewWallet()
	ws.CreateNewWallet()
	ws.CreateNewWallet()
}
