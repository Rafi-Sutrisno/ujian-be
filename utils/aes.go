package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

func getAESKey() ([]byte, error) {
	keyHex := os.Getenv("AES_KEY")
	if keyHex == "" {
		return nil, errors.New("AES_KEY environment variable not set")
	}
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode AES_KEY: %w", err)
	}
	if len(key) != 32 {
		return nil, errors.New("AES key must be 32 bytes (64 hex characters)")
	}
	return key, nil
}

func AESEncrypt(stringToEncrypt string) (string, error) {
	key, err := getAESKey()
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(stringToEncrypt), nil)
	return hex.EncodeToString(ciphertext), nil
}

func AESDecrypt(encryptedString string) (string, error) {
	key, err := getAESKey()
	if err != nil {
		return "", err
	}
	enc, err := hex.DecodeString(encryptedString)
	if err != nil {
		return "", fmt.Errorf("invalid hex input: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(enc) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}
