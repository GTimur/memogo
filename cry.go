// Package memogo - cry - string data encription/decription realization
package memogo

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"log"
)

/* Вектор инициализации (16 byte) */
var commonIV = []byte("TESTVERSION2018.")

/* Ключ шифрования для AES (32 byte) */
var word = "storeYourDataInSafePlace12312312"

// AesEncrypt - encript data string
func AesEncrypt(text string) ([]byte, error) {
	key := word
	IV := commonIV
	// Create the aes encryption algorithm
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Printf("Error: NewCipher(%d bytes) = %s\n", len(key), err)
		return nil, err
	}

	// Encrypted string
	cfb := cipher.NewCFBEncrypter(c, IV)
	ciphertext := make([]byte, len(text))
	cfb.XORKeyStream(ciphertext, []byte(text))
	return ciphertext, nil
}

// AesDecript - decript data string
func AesDecript(text []byte) (string, error) {
	key := word
	IV := commonIV
	// Create the aes encryption algorithm
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Printf("Error: NewCipher(%d bytes) = %s\n", len(key), err)
		return "", err
	}

	// Decrypt strings
	cfbdec := cipher.NewCFBDecrypter(c, IV)
	plaintextCopy := make([]byte, len(text))
	cfbdec.XORKeyStream(plaintextCopy, []byte(text))
	return fmt.Sprintf("%s", string(plaintextCopy)), nil
}
