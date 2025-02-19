package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"fmt"
)

// GDP File Structure
// ---------------------------------------------------------
// Magick bytes       : 3 (byte[3]) ("GDP" -> "\x47\x44\x50")
// Encryption         : 1 byte (bool)
// Size of Nonce      : 1 byte (uint8)
// Nonce              : 0-255 bytes (based on Size of Nonce)
// Size of Ciphertext : 8 bytes (uint64)
// Ciphertext         : 0-18,446,744,073,709,551,615 bytes
// ---------------------------------------------------------

// ParseGDPFile parses a GDP file format and extracts encryption flag, nonce, and ciphertext.
func ParseGDPFile(input []byte, dummy bool) (bool, int, []byte, uint64, []byte, error) {
	if len(input) < 13 {
		return false, 0, nil, 0, nil, errors.New("invalid GDP file: too short")
	}

	// Check magic bytes
	if !bytes.Equal(input[:3], []byte("GDP")) {
		return false, 0, nil, 0, nil, errors.New(fmt.Sprintf("invalid GDP file: incorrect magic bytes %x", input[:3]))
	}

	// Read encryption flag (1 byte)
	encryption := input[3] != 0

	// Read nonce size (1 byte)
	nonceSize := int(input[4])
	if len(input) < 13+nonceSize && !dummy {
		return false, 0, nil, 0, nil, errors.New("invalid GDP file: data too short for nonce and size field")
	}

	// Read nonce (variable length)
	nonce := input[5 : 5+nonceSize]

	// Read ciphertext size (8 bytes, uint64)
	ciphertextSize := binary.LittleEndian.Uint64(input[5+nonceSize : 5+nonceSize+8])
	if len(input) < 5+nonceSize+8+int(ciphertextSize) && !dummy {
		return false, 0, nil, 0, nil, errors.New("invalid GDP file: incomplete ciphertext")
	}

	var ciphertext []byte
	if !dummy {
		// Read ciphertext
		ciphertext = input[5+nonceSize+8 : 5+nonceSize+8+int(ciphertextSize)]
	}

	return encryption, nonceSize, nonce, ciphertextSize, ciphertext, nil
}

// MakeGDPFile constructs a GDP file with the given encryption flag, nonce, and ciphertext.
func MakeGDPFile(encryption bool, nonce []byte, ciphertext []byte) ([]byte, error) {
	if len(nonce) > 255 {
		return nil, errors.New("nonce size exceeds 255 bytes")
	}

	var buf bytes.Buffer

	// Write magic bytes "GDP"
	buf.Write([]byte("GDP"))

	// Write encryption flag (1 byte)
	if encryption {
		buf.WriteByte(1)
	} else {
		buf.WriteByte(0)
	}

	// Write nonce size (1 byte)
	buf.WriteByte(byte(len(nonce)))

	// Write nonce data
	buf.Write(nonce)

	// Write ciphertext size (8 bytes, uint64 in little-endian)
	ciphertextSize := uint64(len(ciphertext))
	binary.Write(&buf, binary.LittleEndian, ciphertextSize)

	// Write ciphertext
	buf.Write(ciphertext)

	return buf.Bytes(), nil
}

// ReadGDPFile reads a GDP file from disk and parses its contents.
func ReadGDPFile(input string) (bool, int, []byte, uint64, []byte, error) {
	file, err := os.ReadFile(input)
	if err != nil {
		return false, 0, nil, 0, nil, err
	}
	return ParseGDPFile(file, false)
}

// WriteGDPFile writes a GDP file to disk.
func WriteGDPFile(filename string, encryption bool, nonce []byte, ciphertext []byte) error {
	file, err := MakeGDPFile(encryption, nonce, ciphertext)
	if err != nil {
		return err
	}

	// Write to file
	err = os.WriteFile(filename, file, 0644)
	if err != nil {
		return err
	}

	return nil
}
