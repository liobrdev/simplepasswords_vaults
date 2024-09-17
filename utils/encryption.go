package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

func Encrypt(plaintext, hexEncodedKey string) (hexEncodedCiphertext string, err error) {
	var key []byte

	if key, err = hex.DecodeString(hexEncodedKey); err != nil {
		return "", err
	}

	var block cipher.Block

	if block, err = aes.NewCipher(key); err != nil {
		return "", err
	}

	var gcm cipher.AEAD

	if gcm, err = cipher.NewGCM(block); err != nil {
    return "", err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return hex.EncodeToString(ciphertext), nil 
}

func Decrypt(hexEncodedCiphertext, hexEncodedKey string) (plaintext string, err error) {
	var key []byte

	if key, err = hex.DecodeString(hexEncodedKey); err != nil {
		return "", err
	}

	var block cipher.Block

	if block, err = aes.NewCipher(key); err != nil {
		return "", err
	}

	var gcm cipher.AEAD

	if gcm, err = cipher.NewGCM(block); err != nil {
    return "", err
	}

	var ciphertext []byte

	if ciphertext, err = hex.DecodeString(hexEncodedCiphertext); err != nil {
		return "", err
	}

	nonce, textToDecrypt := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]

	var decryptedText []byte

	if decryptedText, err = gcm.Open(nil, nonce, textToDecrypt, nil); err != nil {
		return "", err
	}

	return string(decryptedText), nil
}
