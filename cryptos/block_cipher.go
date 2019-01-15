package cryptos

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/sha3"
)

//https://golang.org/pkg/crypto/cipher/

func createHash32Byte(key string) []byte {
	hasher := sha3.New256()
	hasher.Write([]byte(key))
	return hasher.Sum(nil)
}

func GcmDecrypt(nonce, ciphertext, key []byte, nonceIncluded bool) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	var plaintext []byte
	if nonceIncluded {
		nonceSize := aesgcm.NonceSize()
		plaintext, err = aesgcm.Open(nil, ciphertext[:nonceSize], ciphertext[nonceSize:], nil)
	} else {
		plaintext, err = aesgcm.Open(nil, nonce, ciphertext, nil)

	}
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("plaintext: %s\n", plaintext)

}

func GcmEncrypt(plaintext string, key []byte, nonceIncluded bool) ([]byte, []byte) {
	// Load your secret key from a safe place and reuse it across multiple
	// Seal/Open calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256).

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	var ciphertext []byte
	if nonceIncluded {
		ciphertext = aesgcm.Seal(nonce, nonce, []byte(plaintext), nil)
	} else {
		ciphertext = aesgcm.Seal(nil, nonce, []byte(plaintext), nil)
	}
	fmt.Printf("nonce: %x\n", nonce)
	fmt.Printf("ciphertext: %x\n", ciphertext)
	return nonce, ciphertext

}

func Ctr(plaintext string, key []byte) {
	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	// CTR mode is the same for both encryption and decryption, so we can
	// also decrypt that ciphertext with NewCTR.

	plaintext2 := make([]byte, len(plaintext))
	stream = cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext2, ciphertext[aes.BlockSize:])

	fmt.Printf("%s\n", plaintext2)
}

func BlockCipherSample() {
	nonce, ciphertext := GcmEncrypt("This is plaintext", createHash32Byte("password"), false)
	GcmDecrypt(nonce, ciphertext, createHash32Byte("password"), false)
	nonce, ciphertext = GcmEncrypt("This is plaintext", createHash32Byte("password"), true)
	GcmDecrypt(nonce, ciphertext, createHash32Byte("password"), true)

	Ctr("This is plaintext", createHash32Byte("password"))
	// AesSample()
}

//https://golang.org/pkg/crypto/cipher/
