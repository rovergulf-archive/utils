package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/rovergulf/utils/clog"
	"golang.org/x/crypto/sha3"
	"io"
)

const (
	key = "qhA2KF6q8twuPyXj"
)

func EncodeSHA3String(v string) string {
	hash := sha3.New512()
	hash.Write([]byte(v))
	return hex.EncodeToString(hash.Sum(nil))
}

func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		clog.Error(err)
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		clog.Error(err)
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		clog.Error(err)
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		clog.Error(err)
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		clog.Error(err)
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
