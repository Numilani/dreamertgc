package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	signalr "github.com/philippseith/signalr"
	viper "github.com/spf13/viper"
	"os"
)

var client signalr.Client
var state = AppState{
	stage:      Starting,
	errorState: NoError,
}
var Application = AppModel{
	state:         state,
	rcv:           ServerEventReceiver{UiUpdateChannel: make(chan ServerDataChunk)},
	altWindow:     AltWindow{IsFocused: true, Contents: []string{"Connecting to server..."}},
	infoPane:      CharacterPane{Contents: "Log in to view character Data"},
	primaryPane:   ChatPane{Contents: []string{}, ChatInput: textinput.New()},
	secondaryPane: SystemPane{Contents: []string{}},
	statusBar:     StatusBar{LeftBlurb: "NUMI'S TEST CLIENT", RightBlurb: "v0.01a", MiddleString: "L: Loading..."},
}

func main() {
	// set up config file
	viper.SetConfigName("settings")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Couldn't read config!")
	}

	programUi := tea.NewProgram(&Application, tea.WithAltScreen())
	if err := programUi.Start(); err != nil {
		fmt.Printf("An error occurred: %v", err)
		os.Exit(1)
	}
}
