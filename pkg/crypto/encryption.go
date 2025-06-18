package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

func Encrypt(plaintext, masterPassword string) (string, error) {
	// 1. Convert master password to 32-byte key
	key := sha256.Sum256([]byte(masterPassword))

	// 2. Create AES cipher
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	// 3. Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 4. Gen random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 5. Encrypt
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// 6. Encode to base64 for storage
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(ciphertext, masterPassword string) (string, error) {
	// 1. Decode from base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// 2. Same key derivation as encrypt
	key := sha256.Sum256([]byte(masterPassword))

	// 3. Create same AES cipher
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	// 4. Create same GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 5. Extract nonce (first 12 bytes)
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce := data[:nonceSize]
	ciphertext_bytes := data[nonceSize:]

	// 6. Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext_bytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
