package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	saltLen = 16
	time    = 2
	memory  = 19456
	threads = 1
	keyLen  = 32
)

type Cipher struct {
	masterPassword string
}

func NewCipher(masterPassword string) Cipher {
	return Cipher{masterPassword}
}

func (s Cipher) Encrypt(password string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	gcm, err := getGCM(s.masterPassword, salt)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(password), nil)

	result := make([]byte, 0, len(salt)+len(nonce)+len(ciphertext))
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return base64.StdEncoding.EncodeToString(result), nil
}

func (s Cipher) Decrypt(password string) (string, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		return "", err
	}

	minLen := saltLen + 12
	if len(encryptedBytes) < minLen {
		return "", errors.New("incorrect data")
	}

	salt := encryptedBytes[:saltLen]

	gcm, err := getGCM(s.masterPassword, salt)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce := encryptedBytes[saltLen : saltLen+nonceSize]
	ciphertext := encryptedBytes[saltLen+nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func getGCM(masterPassword string, salt []byte) (cipher.AEAD, error) {
	key := argon2.IDKey([]byte(masterPassword), salt, time, memory, threads, keyLen)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(block)
}
