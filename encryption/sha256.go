package encryption

import (
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/ripemd160"
)

func main() {
	ripe160()
}

// hash 不是加密加密，因为无法解密
// 256位
func Sha256() {
	hasher := sha256.New()
	hasher.Write([]byte("raojiangjin"))
	bytes := hasher.Sum(nil)
	fmt.Printf("%x", bytes)
}

// hash
func ripe160() {
	ripemd := ripemd160.New()
	ripemd.Write([]byte("raojiangjin"))
	bytes := ripemd.Sum(nil)
	fmt.Printf("%x", bytes)

}
