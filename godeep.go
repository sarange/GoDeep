package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"log"
	"encoding/hex"

	"golang.org/x/crypto/pbkdf2"
	"github.com/spf13/cobra"

	"godeep/utils"
)

func main() {
	// Define the root command
	var rootCmd = &cobra.Command{
		Use:   "godeep",
		Short: "A tool to embed or extract data from a WAV file",
		Run: func(cmd *cobra.Command, args []string) {
			// This can be used for general functionality if needed
			fmt.Println("You must specify a mode: embed or extract.")
		},
	}

	// Define flags for the root command (applies to the whole program)
	var inputFile string
	var outputFile string
	var container string
	var password string
	var verbose bool
	var noEncryption bool

	// Root command flags
	rootCmd.PersistentFlags().StringVarP(&inputFile, "input", "i", "", "Path to the input WAV file")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "Path to the output WAV file or extracted file")
	rootCmd.PersistentFlags().StringVarP(&container, "container", "c", "", "WAV Container to embed the data within")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Encryption password (required unless --noencryption is used)")
	rootCmd.PersistentFlags().BoolVarP(&noEncryption, "noencryption", "", false, "Disable encryption")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Define the "embed" command
	var embedCmd = &cobra.Command{
		Use:   "embed",
		Short: "Embed data into a WAV file",
		Run: func(cmd *cobra.Command, args []string) {
			// Validate Required Inputs for "embed" command
			if inputFile == "" {
				fmt.Println("Error: Input file path is required for embedding.")
				cmd.Usage()
				os.Exit(1)
			}

			if outputFile == "" {
				fmt.Println("Error: Output file path is required for embedding.")
				cmd.Usage()
				os.Exit(1)
			}

			if container == "" {
				fmt.Println("Error: Container WAV file is required for embedding.")
				cmd.Usage()
				os.Exit(1)
			}

			// Check Password Requirement for embedding (when encryption is not disabled)
			if password == "" && !noEncryption {
				fmt.Println("Error: Password is required for encryption when --noencryption is not used.")
				cmd.Usage()
				os.Exit(1)
			}

			// If validation passed, print out the parameters and proceed with the embed logic
			fmt.Printf("Embedding data from '%s' into '%s' (container: '%s', encryption: %v)\n",
				inputFile, outputFile, container, !noEncryption)
			
			// Derive Key if Encryption is Enabled
			var key []byte
			var encryption bool
			if !noEncryption {
				key = pbkdf2.Key([]byte(password), []byte("GoDeepSalt"), 100000, 32, sha256.New)
				encryption = true
			} else {
				encryption = false
			}

			if verbose {
				// Verbose Output for Encryption/Compression
				if encryption {
					fmt.Println("[DEBUG] Encryption enabled. Deriving key...")
				} else {
					fmt.Println("[DEBUG] Encryption disabled.")
				}
			}

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
				ciphertext, nonce, err = utils.CompressAndEncrypt(file, key)
				if err != nil {
					fmt.Println("Encryption failed:", err)
					os.Exit(1)
				}
			} else {
				// Compress the input file to GDP without encryption
				ciphertext, err = utils.CompressXZ(file)
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
			gdpFile, err := utils.MakeGDPFile(encryption, nonce, ciphertext)
			if err != nil {
				fmt.Println("Error creating GDP file:", err)
				os.Exit(1)
			}

			// Verbose output for GDP file size
			if verbose {
				fmt.Printf("[DEBUG] GDP file size: %d bytes\n", len(gdpFile))
			}

			// Read container WAV file
			containerData, metadata, err := utils.WAVToPCM(container)
			if err != nil {
				fmt.Println("Error reading container file:", err)
				os.Exit(1)
			}

			// Verbose output for container WAV size
			if verbose {
				fmt.Printf("[DEBUG] Container WAV file size: %d bytes\n", len(containerData))
			}

			// Embed GDP file into container WAV file using LSB encoding
			embeddedWAV, err := utils.EmbedToLSB(containerData, gdpFile)
			if err != nil {
				fmt.Println("Error embedding GDP into WAV:", err)
				os.Exit(1)
			}

			// Write the embedded WAV data to output file
			err = utils.PCMToWAV(outputFile, embeddedWAV, *metadata)
			if err != nil {
				fmt.Println("Error writing to output file:", err)
				os.Exit(1)
			}

			// Print Success
			fmt.Println("Success: Embedded data successfully written to output file")
		},
	}

	// Define the "extract" command
	var extractCmd = &cobra.Command{
		Use:   "extract",
		Short: "Extract data from a WAV file",
		Run: func(cmd *cobra.Command, args []string) {
			// Validate Required Inputs for "extract" command
			if container == "" {
				fmt.Println("Error: Container WAV file is required for embedding.")
				cmd.Usage()
				os.Exit(1)
			}

			if outputFile == "" {
				fmt.Println("Error: Output file path is required for extraction.")
				cmd.Usage()
				os.Exit(1)
			}

			// Check Password Requirement for extraction (when encryption is not disabled)
			if password == "" && !noEncryption {
				fmt.Println("Error: Password is required for decryption when --noencryption is not used.")
				cmd.Usage()
				os.Exit(1)
			}

			// If validation passed, print out the parameters and proceed with the extract logic
			fmt.Printf("Extracting data from '%s' to '%s' (container: '%s', encryption: %v)\n",
				inputFile, outputFile, container, !noEncryption)
			
			// Derive Key if Encryption is Enabled
			var key []byte
			var encryption bool
			if !noEncryption {
				key = pbkdf2.Key([]byte(password), []byte("GoDeepSalt"), 100000, 32, sha256.New)
				encryption = true
			} else {
				encryption = false
			}

			if verbose {
				// Verbose Output for Encryption/Compression
				if encryption {
					fmt.Println("[DEBUG] Encryption enabled. Deriving key...")
				} else {
					fmt.Println("[DEBUG] Encryption disabled.")
				}
			}

			// Read container WAV to PCM file
			if verbose {
				fmt.Println("[DEBUG] Reading container WAV to PCM file")
			}
			containerData, _, err := utils.WAVToPCM(container)
			if err != nil {
				fmt.Println("Error reading container file:", err)
				os.Exit(1)
			}

			// Verbose output for container WAV size
			if verbose {
				fmt.Printf("[DEBUG] Container WAV file size: %d bytes\n", len(containerData))
			}

			// Extract GDP file from LSB of container WAV file
			gdpFile, err := utils.ExtractGDPFromLSB(containerData)
			if err != nil {
				fmt.Println("Error extracting GDP file from container:", err)
				os.Exit(1)
			}

			// Verbose output for GDP file size
			if verbose {
				fmt.Printf("[DEBUG] Extracted GDP file size: %d bytes\n", len(gdpFile))
			}

			// Parse the GDP file to get encryption flag, nonce, and ciphertext
			encryption, _, nonce, _, ciphertext, err := utils.ParseGDPFile(gdpFile, false)
			if err != nil {
				fmt.Println("Error parsing GDP file:", err)
				os.Exit(1)
			}

			var plaintext []byte
			if encryption {
				plaintext, err = utils.DecryptAndDecompress(ciphertext, key, nonce)
			} else {
				plaintext, err = utils.DecompressXZ(ciphertext)
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

		},
	}

	// Add the embed and extract commands to the root command
	rootCmd.AddCommand(embedCmd, extractCmd)

	// Add bash completion command
	var completionCmd = &cobra.Command{
		Use:   "completion",
		Short: "Generate bash completion script",
		Run: func(cmd *cobra.Command, args []string) {
			// Output the bash completion script
			if err := rootCmd.GenBashCompletion(os.Stdout); err != nil {
				log.Fatal(err)
			}
		},
	}

	// Add the completion command to the root command
	rootCmd.AddCommand(completionCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
	
}
