package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"github.com/jiangjincc/islands/utils"

	"github.com/jiangjincc/islands/encryption"
)

var (
	version  = byte(0x00)
	checkLen = 4
)

type Wallet struct {
	PublicKey  []byte
	PrivateKey ecdsa.PrivateKey
}

func (w *Wallet) GetAddress() []byte {
	// 1、先将PuKey 256 -> 160 = 20byte
	ripemdHash := utils.Ripemd160Hash(w.PublicKey)
	versionRipemd160Hash := append([]byte{version}, ripemdHash...)
	checkSumBytes := CheckSum(versionRipemd160Hash)

	bytes := append(versionRipemd160Hash, checkSumBytes...)
	return encryption.Base58Encode(bytes)
}

func IsValidForAddress(address []byte) bool {
	versionPublicCheckSum := encryption.Base58Decode(address)
	checkSumBytes := versionPublicCheckSum[len(versionPublicCheckSum)-checkLen:]
	ripemd160CheckSum := versionPublicCheckSum[0 : len(versionPublicCheckSum)-checkLen]

	checkBytes := CheckSum(ripemd160CheckSum)
	if bytes.Compare(checkSumBytes, checkBytes) == 0 {
		return true
	}
	return false
}

func NewWallet() *Wallet {
	priKey, pubKey := newKeyPair()
	return &Wallet{PrivateKey: priKey, PublicKey: pubKey}
}

// 通过私钥产生公钥
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	publicKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, publicKey
}

func CheckSum(b []byte) []byte {
	hash1 := sha256.Sum256(b)
	hash2 := sha256.Sum256(hash1[:])

	return hash2[:checkLen]
}
