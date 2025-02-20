package main

import (
	"fmt"
	"os"
	"log"

	"github.com/spf13/cobra"

	"github.com/sarange/godeep/utils"
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
			fmt.Printf("Embedding data from '%s' into '%s' (container: '%s', encryption: %v)\n", inputFile, outputFile, container, !noEncryption)

			var key []byte
			if !noEncryption {
				key = utils.DeriveKey(password)
				if verbose {
					fmt.Println("[DEBUG] Encryption disabled.")
				}
			} else if verbose {
				fmt.Println("[DEBUG] Encryption enabled. Deriving key...")
			}
		
			err := utils.Embed(inputFile, outputFile, container, key, !noEncryption, verbose)

			if err != nil {
				fmt.Println("Error embeding:", err)
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

			var key []byte
			if !noEncryption {
				key = utils.DeriveKey(password)
				if verbose {
					fmt.Println("[DEBUG] Encryption disabled.")
				}
			} else if verbose {
				fmt.Println("[DEBUG] Encryption enabled. Deriving key...")
			}

			err := utils.Extract(container, outputFile, key, !noEncryption, verbose)
			if err != nil {
				fmt.Println("Error extracting:", err)
				os.Exit(1)
			}

		},
	}

	// Define the "extract" command
	var guiCmd = &cobra.Command{
		Use:   "gui",
		Short: "Spawn a GUI interface",
		Run: func(cmd *cobra.Command, args []string) {
			utils.SpawnGui()
		},
	}

	// Add the embed and extract commands to the root command
	rootCmd.AddCommand(embedCmd, extractCmd, guiCmd)

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
