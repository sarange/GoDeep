package main

import (
	"os"
	"os/exec"
	"testing"
	"crypto/sha256"
	"golang.org/x/crypto/pbkdf2"

	"github.com/sarange/godeep/utils"
)

// Test files
const (
	testContainerWAV       = "samples/container.wav"
	testSecretFile         = "samples/secret.txt"
	testOutputWAV          = "samples/output.wav"
	testExtractedFile      = "samples/extracted_secret.txt"
	testOutputWAVNoPass    = "samples/output_no_pass.wav"
	testExtractedFileNoPass = "samples/extracted_no_pass.txt"
	testPassword           = "testpassword"
)

// Generate encryption key
func generateKey() []byte {
	return pbkdf2.Key([]byte(testPassword), []byte("GoDeepSalt"), 100000, 32, sha256.New)
}

// ✅ **Test 1: Embed using `utils/` (With Password)**
func TestEmbedWithPassword(t *testing.T) {
	key := generateKey()
	err := utils.Embed(testSecretFile, testOutputWAV, testContainerWAV, key, true, true)
	if err != nil {
		t.Fatalf("Embedding failed: %v", err)
	}

	if _, err := os.Stat(testOutputWAV); os.IsNotExist(err) {
		t.Fatalf("Output WAV file was not created")
	}
}

// ✅ **Test 2: Extract using `utils/` (With Password)**
func TestExtractWithPassword(t *testing.T) {
	key := generateKey()
	err := utils.Extract(testOutputWAV, testExtractedFile, key, true, true)
	if err != nil {
		t.Fatalf("Extraction failed: %v", err)
	}

	if _, err := os.Stat(testExtractedFile); os.IsNotExist(err) {
		t.Fatalf("Extracted file was not created")
	}

	originalData, err := os.ReadFile(testSecretFile)
	if err != nil {
		t.Fatalf("Failed to read original secret file: %v", err)
	}

	extractedData, err := os.ReadFile(testExtractedFile)
	if err != nil {
		t.Fatalf("Failed to read extracted file: %v", err)
	}

	if string(originalData) != string(extractedData) {
		t.Fatalf("Extracted data does not match original secret file")
	}
}

// ✅ **Test 3: Embed using `utils/` (No Password)**
func TestEmbedNoPassword(t *testing.T) {
	err := utils.Embed(testSecretFile, testOutputWAVNoPass, testContainerWAV, nil, false, true)
	if err != nil {
		t.Fatalf("Embedding (no password) failed: %v", err)
	}

	if _, err := os.Stat(testOutputWAVNoPass); os.IsNotExist(err) {
		t.Fatalf("Output WAV file (no password) was not created")
	}
}

// ✅ **Test 4: Extract using `utils/` (No Password)**
func TestExtractNoPassword(t *testing.T) {
	err := utils.Extract(testOutputWAVNoPass, testExtractedFileNoPass, nil, false, true)
	if err != nil {
		t.Fatalf("Extraction (no password) failed: %v", err)
	}

	if _, err := os.Stat(testExtractedFileNoPass); os.IsNotExist(err) {
		t.Fatalf("Extracted file (no password) was not created")
	}

	originalData, err := os.ReadFile(testSecretFile)
	if err != nil {
		t.Fatalf("Failed to read original secret file: %v", err)
	}

	extractedData, err := os.ReadFile(testExtractedFileNoPass)
	if err != nil {
		t.Fatalf("Failed to read extracted file: %v", err)
	}

	if string(originalData) != string(extractedData) {
		t.Fatalf("Extracted data (no password) does not match original secret file")
	}
}

// ✅ **Test 5: CLI - Embed (With Password)**
func TestCLI_EmbedWithPassword(t *testing.T) {
	cmd := exec.Command("./godeep", "embed", "-i", testSecretFile, "-o", testOutputWAV, "-c", testContainerWAV, "-p", testPassword)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("CLI Embed failed: %v", err)
	}

	if _, err := os.Stat(testOutputWAV); os.IsNotExist(err) {
		t.Fatalf("Output WAV file was not created by CLI")
	}
}

// ✅ **Test 6: CLI - Extract (With Password)**
func TestCLI_ExtractWithPassword(t *testing.T) {
	cmd := exec.Command("./godeep", "extract", "-c", testOutputWAV, "-o", testExtractedFile, "-p", testPassword)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("CLI Extract failed: %v", err)
	}

	if _, err := os.Stat(testExtractedFile); os.IsNotExist(err) {
		t.Fatalf("Extracted file was not created by CLI")
	}

	originalData, err := os.ReadFile(testSecretFile)
	if err != nil {
		t.Fatalf("Failed to read original secret file: %v", err)
	}

	extractedData, err := os.ReadFile(testExtractedFile)
	if err != nil {
		t.Fatalf("Failed to read extracted file: %v", err)
	}

	if string(originalData) != string(extractedData) {
		t.Fatalf("Extracted data from CLI does not match original secret file")
	}
}

// ✅ **Test 7: CLI - Embed (No Password)**
func TestCLI_EmbedNoPassword(t *testing.T) {
	cmd := exec.Command("./godeep", "embed", "-i", testSecretFile, "-o", testOutputWAVNoPass, "-c", testContainerWAV, "--noencryption")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("CLI Embed (no password) failed: %v", err)
	}

	if _, err := os.Stat(testOutputWAVNoPass); os.IsNotExist(err) {
		t.Fatalf("Output WAV file (no password) was not created by CLI")
	}
}

// ✅ **Test 8: CLI - Extract (No Password)**
func TestCLI_ExtractNoPassword(t *testing.T) {
	cmd := exec.Command("./godeep", "extract", "-c", testOutputWAVNoPass, "-o", testExtractedFileNoPass, "--noencryption")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("CLI Extract (no password) failed: %v", err)
	}

	if _, err := os.Stat(testExtractedFileNoPass); os.IsNotExist(err) {
		t.Fatalf("Extracted file (no password) was not created by CLI")
	}

	originalData, err := os.ReadFile(testSecretFile)
	if err != nil {
		t.Fatalf("Failed to read original secret file: %v", err)
	}

	extractedData, err := os.ReadFile(testExtractedFileNoPass)
	if err != nil {
		t.Fatalf("Failed to read extracted file: %v", err)
	}

	if string(originalData) != string(extractedData) {
		t.Fatalf("Extracted data (no password) from CLI does not match original secret file")
	}
}
