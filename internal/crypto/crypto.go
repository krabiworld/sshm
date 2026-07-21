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
	saltLength = 16
	keyLength  = 32
	timeCost   = 1
	memoryCost = 64 * 1024
	threads    = 4
)

type Cipher struct {
	masterPassword string
}

func NewCipher(masterPassword string) Cipher {
	return Cipher{masterPassword}
}

func (s Cipher) Encrypt(password string) (string, error) {
	salt := make([]byte, saltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	key := argon2.IDKey([]byte(s.masterPassword), salt, timeCost, memoryCost, threads, keyLength)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
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

	minLen := saltLength + 12
	if len(encryptedBytes) < minLen {
		return "", errors.New("incorrect data")
	}

	salt := encryptedBytes[:saltLength]
	nonceSize := 12
	nonce := encryptedBytes[saltLength : saltLength+nonceSize]
	ciphertext := encryptedBytes[saltLength+nonceSize:]

	key := argon2.IDKey([]byte(s.masterPassword), salt, timeCost, memoryCost, threads, keyLength)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("bad master password or corrupted data")
	}

	return string(plaintext), nil
}
