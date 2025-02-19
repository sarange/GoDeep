package utils

import (
	"encoding/binary"
	"errors"
	"os"

	"github.com/go-audio/wav"
	"github.com/go-audio/audio"
)

type wavMetadata struct {
	SampleRate uint32
	BitDepth uint16
	NumChans uint16
	AudioFormat uint16
}

// WAVToPCM reads a WAV file and extracts PCM data and its header.
func WAVToPCM(inputFile string) ([]byte, *wavMetadata, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	decoder := wav.NewDecoder(file)
	if !decoder.IsValidFile() {
		return nil, nil, errors.New("invalid WAV file")
	}

	pcm, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, nil, err
	}

	// Convert []int to []byte (assuming 16-bit PCM)
	byteData := make([]byte, len(pcm.Data)*2) // 16-bit samples → 2 bytes each
	for i, sample := range pcm.Data {
		s := int16(sample)
		byteData[i*2] = byte(s & 0xFF)         // Lower byte
		byteData[i*2+1] = byte((s >> 8) & 0xFF) // Upper byte
	}

	metadata := &wavMetadata{
		SampleRate:  decoder.SampleRate,
		BitDepth:    decoder.BitDepth,
		NumChans:    decoder.NumChans,
		AudioFormat: decoder.WavAudioFormat,
	}

	return byteData, metadata, nil
}



// PCMToWAV converts a PCM byte array to a WAV byte array using a header.
func PCMToWAV(outputFile string, pcmData []byte, metadata wavMetadata) error {
	// Create or overwrite the file
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Convert []byte PCM data into []int PCM samples
	intData := make([]int, len(pcmData)/2)
	for i := 0; i < len(pcmData); i += 2 {
		intData[i/2] = int(binary.LittleEndian.Uint16(pcmData[i : i+2]))
	}

	// Create an *audio.IntBuffer to store PCM data
	pcmBuffer := &audio.IntBuffer{
		Data:   intData,
		Format: &audio.Format{SampleRate: int(metadata.SampleRate), NumChannels: int(metadata.NumChans)},
	}

	// Create WAV encoder
	encoder := wav.NewEncoder(file, int(metadata.SampleRate), int(metadata.BitDepth), int(metadata.NumChans), int(metadata.AudioFormat))

	// Write PCM data
	err = encoder.Write(pcmBuffer)
	if err != nil {
		return err
	}

	// Close encoder properly
	if err := encoder.Close(); err != nil {
		return err
	}

	return nil
}

// GetEmbedSize calculates the available space for embedding data in PCM.
func GetEmbedSize(pcmData []byte) int {
	return len(pcmData) / 8
}

// EmbedToLSB embeds a message into the least significant bit (LSB) of PCM data.
func EmbedToLSB(pcmData []byte, message []byte) ([]byte, error) {
	if len(message) > GetEmbedSize(pcmData) {
		return nil, errors.New("message too large to embed in PCM data")
	}

	encodedPCM := make([]byte, len(pcmData))
	copy(encodedPCM, pcmData)

	// Embed message bit by bit
	for i := 0; i < len(message)*8; i++ {
		byteIndex := i / 8
		bitIndex := i % 8
		bitValue := (message[byteIndex] >> bitIndex) & 0x01

		encodedPCM[i] = (encodedPCM[i] & 0xFE) | bitValue
	}

	return encodedPCM, nil
}

// ExtractGDPFromLSB extracts a GDP file from the LSB of PCM data.
func ExtractGDPFromLSB(pcmData []byte) ([]byte, error) {
	// We set the minimum size as 1024 because there is no reaseon to deal with edge cases.
	dummyHeaderSize := 1024

	if len(pcmData) < dummyHeaderSize {
		return nil, errors.New("not enough data to contain a valid GDP file")
	}

	gdpHeaderBytes := make([]byte, dummyHeaderSize)
	for i := 0; i < dummyHeaderSize; i++ {
		for bit := 0; bit < 8; bit++ {
			gdpHeaderBytes[i] |= (pcmData[i*8+bit] & 0x01) << bit
		}
	}

	// Extract full GDP file size
	_, nonceSize, _, ciphertextSize, _, err := ParseGDPFile(gdpHeaderBytes, true)
	if err != nil {
		return nil, err
	}

	totalGDPSize := 13 + int(nonceSize) + int(ciphertextSize)

	// Ensure the PCM data contains enough bits
	if totalGDPSize*8 > len(pcmData) {
		return nil, errors.New("not enough PCM data to extract the full GDP file")
	}

	// Extract the full GDP file from LSB
	gdpFile := make([]byte, totalGDPSize)
	for i := 0; i < totalGDPSize; i++ {
		for bit := 0; bit < 8; bit++ {
			gdpFile[i] |= (pcmData[i*8+bit] & 0x01) << bit
		}
	}

	return gdpFile, nil
}
