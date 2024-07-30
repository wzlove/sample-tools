package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fyne.io/fyne/v2/dialog"
	"os/exec"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Utility Tools")

	resultLabel := widget.NewLabel("Result will be shown here")
	copyButton := widget.NewButton("Copy Result", func() {
		content := resultLabel.Text
		myWindow.Clipboard().SetContent(content)
	})

	md5Form := createMD5Form(resultLabel)
	base64Form := createBase64Form(resultLabel)
	jsonForm := createJSONForm(resultLabel)
	protoForm := createProtoForm(resultLabel, myWindow)

	toolbar := container.NewVBox(
		widget.NewButton("MD5", func() {
			setContent(myWindow, md5Form, resultLabel, copyButton)
		}),
		widget.NewButton("Base64", func() {
			setContent(myWindow, base64Form, resultLabel, copyButton)
		}),
		widget.NewButton("JSON", func() {
			setContent(myWindow, jsonForm, resultLabel, copyButton)
		}),
		widget.NewButton("Proto", func() {
			setContent(myWindow, protoForm, resultLabel, copyButton)
		}),
	)

	leftSidebar := container.NewVBox(toolbar)
	initialContent := container.NewHSplit(leftSidebar, widget.NewLabel("Select a tool"))
	initialContent.SetOffset(0.3)

	myWindow.SetContent(initialContent)
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}

func setContent(window fyne.Window, form fyne.CanvasObject, resultLabel *widget.Label, copyButton *widget.Button) {
	// Clear previous result
	resultLabel.SetText("Result will be shown here")

	// Update the right content with the new form and result display
	rightContent := container.NewVBox(form, resultLabel, copyButton)

	// Recreate the split content to ensure it updates correctly
	leftSidebar := container.NewVBox(
		widget.NewButton("MD5", func() {
			setContent(window, createMD5Form(resultLabel), resultLabel, copyButton)
		}),
		widget.NewButton("Base64", func() {
			setContent(window, createBase64Form(resultLabel), resultLabel, copyButton)
		}),
		widget.NewButton("JSON", func() {
			setContent(window, createJSONForm(resultLabel), resultLabel, copyButton)
		}),
		widget.NewButton("Proto", func() {
			setContent(window, createProtoForm(resultLabel, window), resultLabel, copyButton)
		}),
	)

	content := container.NewHSplit(leftSidebar, rightContent)
	content.SetOffset(0.3)

	window.SetContent(content)
}

func createMD5Form(resultLabel *widget.Label) fyne.CanvasObject {
	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder("Enter text to hash")

	hashButton := widget.NewButton("Generate MD5", func() {
		hash := md5.Sum([]byte(inputEntry.Text))
		resultLabel.SetText(hex.EncodeToString(hash[:]))
	})

	return container.NewVBox(inputEntry, hashButton)
}

func createBase64Form(resultLabel *widget.Label) fyne.CanvasObject {
	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder("Enter text to encode/decode")

	encodeButton := widget.NewButton("Encode Base64", func() {
		resultLabel.SetText(base64.StdEncoding.EncodeToString([]byte(inputEntry.Text)))
	})

	decodeButton := widget.NewButton("Decode Base64", func() {
		decoded, err := base64.StdEncoding.DecodeString(inputEntry.Text)
		if err != nil {
			resultLabel.SetText("Error decoding Base64")
		} else {
			resultLabel.SetText(string(decoded))
		}
	})

	return container.NewVBox(inputEntry, encodeButton, decodeButton)
}

func createJSONForm(resultLabel *widget.Label) fyne.CanvasObject {
	inputEntry := widget.NewMultiLineEntry()
	inputEntry.SetPlaceHolder("Enter JSON to format")

	formatButton := widget.NewButton("Format JSON", func() {
		var formattedJSON map[string]interface{}
		err := json.Unmarshal([]byte(inputEntry.Text), &formattedJSON)
		if err != nil {
			resultLabel.SetText("Invalid JSON")
			return
		}
		formatted, err := json.MarshalIndent(formattedJSON, "", "  ")
		if err != nil {
			resultLabel.SetText("Error formatting JSON")
		} else {
			resultLabel.SetText(string(formatted))
		}
	})

	return container.NewVBox(inputEntry, formatButton)
}

func createProtoForm(resultLabel *widget.Label, window fyne.Window) fyne.CanvasObject {
	protoPathLabel := widget.NewLabel("No file selected")
	selectProtoButton := widget.NewButton("Select .proto file", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				resultLabel.SetText("Error: " + err.Error())
				return
			}
			if reader == nil {
				resultLabel.SetText("No file selected")
				return
			}
			protoPathLabel.SetText(reader.URI().Path())
		}, window)
	})

	outputDirLabel := widget.NewLabel("No directory selected")
	selectOutputDirButton := widget.NewButton("Select output directory", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				resultLabel.SetText("Error: " + err.Error())
				return
			}
			if list == nil {
				resultLabel.SetText("No directory selected")
				return
			}
			outputDirLabel.SetText(list.Path())
		}, window)
	})

	generateButton := widget.NewButton("Generate .go file", func() {
		protoPath := protoPathLabel.Text
		outputDir := outputDirLabel.Text
		if protoPath == "No file selected" || outputDir == "No directory selected" {
			resultLabel.SetText("Please select both .proto file and output directory")
			return
		}
		protoDir := filepath.Dir(protoPath)
		generateGoFromProto(protoPath, protoDir, outputDir, resultLabel)
	})

	return container.NewVBox(
		selectProtoButton,
		protoPathLabel,
		selectOutputDirButton,
		outputDirLabel,
		generateButton,
	)
}

func generateGoFromProto(protoPath string, protoDir string, outputDir string, resultLabel *widget.Label) {
	cmd := exec.Command("protoc", "--proto_path="+protoDir, "--go_out="+outputDir, protoPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		resultLabel.SetText("Error: " + err.Error() + "\n" + string(output))
		return
	}
	resultLabel.SetText("Successfully generated .go file from " + filepath.Base(protoPath) + " to " + outputDir)
}
