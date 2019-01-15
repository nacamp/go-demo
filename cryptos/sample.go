package cryptos

import (
	"encoding/base64"
	"fmt"
	"log"

	"golang.org/x/crypto/scrypt"
)

//https://godoc.org/golang.org/x/crypto/scrypt
//adaptive key derivation function : scrypt
func ScryptSample() {
	salt := []byte{0xc8, 0x28, 0xf2, 0x58, 0xa7, 0x6a, 0xad, 0x7b}

	dk, err := scrypt.Key([]byte("some password"), salt, 1<<15, 8, 1, 32)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(dk))
}

func Sample() {
	ScryptSample()
}
