package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
)

var MySecret = []byte("WellItIsOkayIfYouSawThis")

func EncryptAES(plaintext string) (string, error) {
	block, err := aes.NewCipher(MySecret)
	if err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize)
	_, err = rand.Read(iv)
	if err != nil {
		return "", err
	}

	paddedText := padText(plaintext, aes.BlockSize)

	mode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(paddedText))
	mode.CryptBlocks(encrypted, []byte(paddedText))

	ivHex := hex.EncodeToString(iv)
	encStr := hex.EncodeToString(encrypted)

	encryptedHex := ivHex + encStr
	return encryptedHex, nil
}

func padText(text string, blockSize int) []byte {
	padding := blockSize - (len(text) % blockSize)
	paddingBytes := make([]byte, padding)
	for i := range paddingBytes {
		paddingBytes[i] = byte(padding)
	}
	return append([]byte(text), paddingBytes...)
}

func DecryptServer(encrypted []byte) (decrypted []byte) {
	key := []byte(MySecret)
	block, _ := aes.NewCipher(key)
	if len(encrypted) < aes.BlockSize {
		panic("ciphertext too short")
	}
	encrypted, err := hex.DecodeString(string(encrypted))
	if err != nil {
		return []byte("")
	}
	iv := encrypted[:aes.BlockSize]
	decrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decrypted, decrypted)
	return decrypted
}
