package utils

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// SpawnGui initializes and runs the GoDeep GUI
func SpawnGui() {
	// Initialize Fyne app
	myApp := app.New()
	myApp.Settings().SetTheme(theme.DarkTheme())
	win := myApp.NewWindow("GoDeep - Audio Steganography")
	win.Resize(fyne.NewSize(600, 450))

	// Title
	title := widget.NewLabelWithStyle("GoDeep - Hide Data in Audio", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Mode Selection
	modeLabel := widget.NewLabel("Select Mode:")
	modeSelect := widget.NewSelect([]string{"Embed", "Extract"}, nil)
	modeSelect.SetSelected("Embed") // Default selection

	// Container File
	containerLabel := widget.NewLabel("WAV Container File:")
	containerEntry := widget.NewEntry()
	containerEntry.SetPlaceHolder("Select WAV container file")
	containerButton := widget.NewButton("Browse", func() {
		dialog.ShowFileOpen(func(f fyne.URIReadCloser, err error) {
			if f != nil {
				containerEntry.SetText(f.URI().Path())
			}
		}, win)
	})
	containerBox := container.NewBorder(nil, nil, nil, containerButton, containerEntry)

	// Input File Selection
	inputLabel := widget.NewLabel("File to embed:")
	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder("Select input file")
	inputButton := widget.NewButton("Browse", func() {
		dialog.ShowFileOpen(func(f fyne.URIReadCloser, err error) {
			if f != nil {
				inputEntry.SetText(f.URI().Path())
			}
		}, win)
	})
	inputBox := container.NewBorder(nil, nil, nil, inputButton, inputEntry)

	// Output File Selection
	outputLabel := widget.NewLabel("Output File:")
	outputEntry := widget.NewEntry()
	outputEntry.SetPlaceHolder("Select output file")
	outputButton := widget.NewButton("Browse", func() {
		dialog.ShowFileSave(func(f fyne.URIWriteCloser, err error) {
			if f != nil {
				outputEntry.SetText(f.URI().Path())
			}
		}, win)
	})
	outputBox := container.NewBorder(nil, nil, nil, outputButton, outputEntry)

	// Password Entry (Optional)
	passwordLabel := widget.NewLabel("Encryption Password (Optional):")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Enter password")

	// Progress Bar & Status
	progress := widget.NewProgressBar()
	statusLabel := widget.NewLabel("Ready")

	// Dynamically show/hide input file field based on mode
	modeSelect.OnChanged = func(selected string) {
		if selected == "Embed" {
			inputLabel.Show()
			inputBox.Show()
		} else {
			inputLabel.Hide()
			inputBox.Hide()
		}
	}
	modeSelect.OnChanged("Embed")

	// Run Button - Calls Embed or Extract function directly
	runButton := widget.NewButtonWithIcon("Run", theme.ConfirmIcon(), func() {
		mode := modeSelect.Selected
		input := inputEntry.Text
		output := outputEntry.Text
		container := containerEntry.Text
		password := passwordEntry.Text

		// Validate required inputs
		if container == "" || output == "" || mode == "" {
			dialog.ShowError(fmt.Errorf("Fill in required fields"), win)
			return
		}

		if mode == "Embed" && input == "" {
			dialog.ShowError(fmt.Errorf("Input file is required for embedding"), win)
			return
		}

		progress.SetValue(0.1) // Initial progress
		statusLabel.SetText("Processing...")

		// Derive encryption key if needed
		var key []byte
		encryption := password != ""
		if encryption {
			key = DeriveKey(password)
		}

		// Run the function in a goroutine to prevent UI freezing
		go func() {
			var err error

			if mode == "Embed" {
				err = Embed(input, output, container, key, encryption, true)
			} else {
				err = Extract(container, output, key, encryption, true)
			}

			if err != nil {
				statusLabel.SetText(fmt.Sprintf("Error: %s", err))
				progress.SetValue(0)
			} else {
				progress.SetValue(1)
				statusLabel.SetText("Success!")
			}
		}()
	})

	// Reset Button
	resetButton := widget.NewButtonWithIcon("Reset", theme.CancelIcon(), func() {
		inputEntry.SetText("")
		outputEntry.SetText("")
		containerEntry.SetText("")
		passwordEntry.SetText("")
		progress.SetValue(0)
		statusLabel.SetText("Ready")
	})

	// Form Layout
	form := container.NewVBox(
		title,
		widget.NewSeparator(),

		container.NewGridWithColumns(2, modeLabel, modeSelect),

		containerLabel, containerBox,
		inputLabel, inputBox,
		outputLabel, outputBox,

		passwordLabel, passwordEntry,

		widget.NewSeparator(),
		progress,
		container.NewGridWithColumns(2, runButton, resetButton),
		statusLabel,
	)

	// Set Content & Run App
	win.SetContent(form)
	win.ShowAndRun()
}
