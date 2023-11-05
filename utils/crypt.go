package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
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
func unpadText(text []byte) ([]byte, error) {
	if len(text) == 0 {
		return nil, errors.New("Input text is empty")
	}

	padding := int(text[len(text)-1])

	if padding >= len(text) {
		return nil, errors.New("Invalid padding value")
	}

	if padding == 0 {
		return nil, errors.New("Invalid padding value: 0")
	}

	if padding > len(text) {
		return nil, errors.New("Padding value exceeds text length")
	}

	unpaddedText := text[:len(text)-padding]
	return unpaddedText, nil
}
func EncryptServer(origData []byte) (encrypted []byte) {
	key := []byte(MySecret)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	encrypted = make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
	encHex := hex.EncodeToString(encrypted)
	return []byte(encHex)
}

func DecryptAES(encrypted []byte) (decryptedRet []byte) {
	encryptedBytes, err := hex.DecodeString(string(encrypted))
	if err != nil {
		return nil
	}

	block, err := aes.NewCipher(MySecret)
	if err != nil {
		return nil
	}

	if len(encryptedBytes) < aes.BlockSize {
		return nil
	}
	iv := encryptedBytes[:aes.BlockSize]
	encStr := encryptedBytes[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encStr))
	mode.CryptBlocks(decrypted, encStr)

	decrypted, err = unpadText(decrypted)
	if err != nil {
		return nil
	}
	return decrypted
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
