package utils

import (
	"fmt"
	"os"
	"encoding/hex"
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

func Embed(inputFile string, outputFile string, container string, key []byte, encryption bool, verbose bool) error {

	var ciphertext, nonce []byte
			
	// Read Input File (Data to be embedded)
	if verbose{
		fmt.Println("[DEBUG] Reading input file for hiding process.")
	}
	file, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		os.Exit(1)
	}

	if encryption {
		// Encrypt and compress the input file to GDP
		ciphertext, nonce, err = CompressAndEncrypt(file, key)
		if err != nil {
			fmt.Println("Encryption failed:", err)
			os.Exit(1)
		}
	} else {
		// Compress the input file to GDP without encryption
		ciphertext, err = CompressXZ(file)
		if err != nil {
			fmt.Println("Compression failed:", err)
			os.Exit(1)
		}
		nonce = nil
	}

	// Verbose output for cipher and nonce
	if verbose {
		// fmt.Println("[DEBUG] Ciphertext (hex):", hex.EncodeToString(ciphertext))
		fmt.Printf("[DEBUG] Ciphertext length: %d bytes\n", len(ciphertext))
		if nonce != nil {
			fmt.Println("[DEBUG] Nonce (hex):", hex.EncodeToString(nonce))
		}
	}

	// Create GDP File Structure
	gdpFile, err := MakeGDPFile(encryption, nonce, ciphertext)
	if err != nil {
		fmt.Println("Error creating GDP file:", err)
		os.Exit(1)
	}

	// Verbose output for GDP file size
	if verbose {
		fmt.Printf("[DEBUG] GDP file size: %d bytes\n", len(gdpFile))
	}

	// Read container WAV file
	containerData, metadata, err := WAVToPCM(container)
	if err != nil {
		fmt.Println("Error reading container file:", err)
		os.Exit(1)
	}

	// Verbose output for container WAV size
	if verbose {
		fmt.Printf("[DEBUG] Container WAV file size: %d bytes\n", len(containerData))
	}

	// Embed GDP file into container WAV file using LSB encoding
	embeddedWAV, err := EmbedToLSB(containerData, gdpFile)
	if err != nil {
		fmt.Println("Error embedding GDP into WAV:", err)
		os.Exit(1)
	}

	// Write the embedded WAV data to output file
	err = PCMToWAV(outputFile, embeddedWAV, *metadata)
	if err != nil {
		fmt.Println("Error writing to output file:", err)
		os.Exit(1)
	}

	return nil
}


// Extract retrieves hidden data from a WAV file
func Extract(container, outputFile string, key []byte, encryption bool, verbose bool) error {

	// Read container WAV to PCM file
	if verbose {
		fmt.Println("[DEBUG] Reading container WAV to PCM file")
	}
	containerData, _, err := WAVToPCM(container)
	if err != nil {
		fmt.Println("Error reading container file:", err)
		os.Exit(1)
	}

	// Verbose output for container WAV size
	if verbose {
		fmt.Printf("[DEBUG] Container WAV file size: %d bytes\n", len(containerData))
	}

	// Extract GDP file from LSB of container WAV file
	gdpFile, err := ExtractGDPFromLSB(containerData)
	if err != nil {
		fmt.Println("Error extracting GDP file from container:", err)
		os.Exit(1)
	}

	// Verbose output for GDP file size
	if verbose {
		fmt.Printf("[DEBUG] Extracted GDP file size: %d bytes\n", len(gdpFile))
	}

	// Parse the GDP file to get encryption flag, nonce, and ciphertext
	encryption, _, nonce, _, ciphertext, err := ParseGDPFile(gdpFile, false)
	if err != nil {
		fmt.Println("Error parsing GDP file:", err)
		os.Exit(1)
	}

	var plaintext []byte
	if encryption {
		plaintext, err = DecryptAndDecompress(ciphertext, key, nonce)
	} else {
		plaintext, err = DecompressXZ(ciphertext)
	}

	if err != nil {
		fmt.Println("Error decrypting/decompressing:", err)
		os.Exit(1)
	}

	// Verbose output for decrypted plaintext size
	if verbose {
		fmt.Printf("[DEBUG] Decrypted plaintext size: %d bytes\n", len(plaintext))
	}

	// Print the extracted plaintext
	if !verbose {
		fmt.Println("Success: Plaintext extracted and written to output file")
	}

	// Write to Output File
	err = os.WriteFile(outputFile, plaintext, 0644)
	if err != nil {
		fmt.Println("Error writing to output file:", err)
		os.Exit(1)
	}
	return nil
}

func DeriveKey(password string) []byte {
	return pbkdf2.Key([]byte(password), []byte("GoDeepSalt"), 100000, 32, sha256.New)
}