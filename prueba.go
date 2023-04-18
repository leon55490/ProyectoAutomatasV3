package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

var automaton *Automaton // Define the loaded automaton pointer
var loadedContent []byte

type Automaton struct {
	States       []string                     `json:"states"`
	Alphabet     []string                     `json:"alphabet"`
	Transitions  map[string]map[string]string `json:"transitions"`
	InitialState string                       `json:"initialState"`
	FinalStates  []string                     `json:"finalStates"`
}

func (a *Automaton) Run(input string) string {
	if a == nil {
		return "Error: no automaton loaded"
	}

	currentState := a.InitialState
	accepted := false
	for _, c := range input {
		transition, ok := a.Transitions[currentState][string(c)]
		if !ok {
			return "Rechazado"
		}
		currentState = transition
	}

	for _, finalState := range a.FinalStates {
		if currentState == finalState {
			accepted = true
			break
		}
	}

	if accepted {
		return "ACEPTADO"
	} else {
		return "Rechazado"
	}
}

func main() {
	// Create a new Fyne application
	myApp := app.New()

	// Create a new window
	myWindow := myApp.NewWindow("Automaton")

	// Create output widgets
	outputLabel := widget.NewLabel("Result:")
	outputText := widget.NewLabel("")

	// Create button widgets
	loadButton := widget.NewButton("Load Automaton", func() {
		fmt.Println("Muestra la opcion de cargar archivos")
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				defer reader.Close()
				bytes, err := ioutil.ReadAll(reader)
				if err != nil {
					dialog.ShowError(err, myWindow)
					return
				}
				err = loadAutomatonFromBytes(bytes)
				if err != nil {
					dialog.ShowError(err, myWindow)
					return
				}
				outputText.SetText("")
			}
		}, myWindow)
		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		fileDialog.Show()
	})

	runButton := widget.NewButton("Run Automaton", func() {
		if automaton == nil {
			dialog.ShowError(errors.New("no automaton loaded"), myWindow)
			return
		}

		result := automaton.Run(string(loadedContent))
		outputText.SetText(result)
	})

	// Create a container for output widgets
	outputContainer := container.NewVBox(outputLabel, outputText)

	// Create a container for the buttons
	buttonContainer := container.NewHBox(loadButton, runButton)

	// Create a container for the output container and the buttons
	content := container.NewVBox(outputContainer, buttonContainer)

	// Add the content container to the window
	myWindow.SetContent(content)

	// Show the window and start the app
	myWindow.ShowAndRun()
}

func loadAutomatonFromBytes(data []byte) error {
	fmt.Println("Lee y carga archivos")
	automaton = &Automaton{}
	err := json.Unmarshal(data, automaton)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}
	loadedContent = data
	return nil
}
