package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	signalr "github.com/philippseith/signalr"
	"os"
)

var client signalr.Client
var appState = AppMainModel{
	stage:         Starting,
	errorState:    NoError,
	rcv:           ServerEventReceiver{UiUpdateChannel: make(chan ServerDataChunk)},
	altWindow:     AltWindow{IsFocused: true, Contents: []string{"Connecting to server..."}},
	infoPane:      CharacterPane{Contents: []string{"Log in to view character Data"}},
	primaryPane:   ChatPane{Contents: []string{}, ChatInput: textinput.New()},
	secondaryPane: SystemPane{Contents: []string{}},
	statusBar:     StatusBar{LeftBlurb: "NUMI'S TEST CLIENT", RightBlurb: "v0.01a", MiddleString: "L: Loading..."},
}

func main() {
	programUi := tea.NewProgram(&appState, tea.WithAltScreen())
	if err := programUi.Start(); err != nil {
		fmt.Printf("An error occurred: %v", err)
		os.Exit(1)
	}
}
