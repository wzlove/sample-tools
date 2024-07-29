package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"

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
