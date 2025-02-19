package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/ulikunitz/xz"
)

// CompressAndEncrypt compresses the input data and then encrypts it with AES-GCM
func CompressAndEncrypt(plaintext, key []byte) ([]byte, []byte, error) {
	// Compress the data using XZ
	compressed, err := CompressXZ(plaintext)
	if err != nil {
		return nil, nil, fmt.Errorf("compression failed: %w", err)
	}

	// Encrypt the compressed data using AES-GCM
	ciphertext, nonce, err := EncryptAESGCM(compressed, key)
	if err != nil {
		return nil, nil, fmt.Errorf("encryption failed: %w", err)
	}

	return ciphertext, nonce, nil
}

// DecryptAndDecompress decrypts the data with AES-GCM and decompresses it using XZ
func DecryptAndDecompress(ciphertext, key, nonce []byte) ([]byte, error) {
	// Decrypt the data using AES-GCM
	decrypted, err := DecryptAESGCM(ciphertext, nonce, key)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	// Decompress the decrypted data
	decompressed, err := DecompressXZ(decrypted)
	if err != nil {
		return nil, fmt.Errorf("decompression failed: %w", err)
	}

	return decompressed, nil
}

func CompressXZ(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := xz.NewWriter(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create XZ writer: %w", err)
	}

	_, err = writer.Write(data)
	if err != nil {
		writer.Close() // Ensure writer is closed before returning
		return nil, fmt.Errorf("failed to write to XZ writer: %w", err)
	}

	err = writer.Close() // Explicitly check Close() error
	if err != nil {
		return nil, fmt.Errorf("failed to close XZ writer: %w", err)
	}

	return buf.Bytes(), nil
}


// DecompressXZ decompresses XZ data
func DecompressXZ(compressedData []byte) ([]byte, error) {
	reader, err := xz.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, fmt.Errorf("failed to create XZ reader: %w", err)
	}

	return io.ReadAll(reader)
}

// EncryptAESGCM encrypts data using AES-GCM
func EncryptAESGCM(plaintext, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

// DecryptAESGCM decrypts data using AES-GCM
func DecryptAESGCM(ciphertext, nonce, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plaintext, nil
}
