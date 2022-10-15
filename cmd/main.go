package main

import (
	. "dreamer_tgc/internal"
	. "dreamer_tgc/internal/renderers"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	signalr "github.com/philippseith/signalr"
	viper "github.com/spf13/viper"
	"os"
)

var client signalr.Client
var state = AppState{
	Stage:      Starting,
	ErrorState: NoError,
}
var Application = AppModel{
	State:            state,
	ConnectionClient: client,
	Rcv:              ServerEventReceiver{UiUpdateChannel: make(chan ServerDataChunk)},
	AltWindow:        AltWindow{IsFocused: true, Contents: []string{"Connecting to server..."}},
	InfoPane:         CharacterPane{Contents: "Log in to view character Data"},
	PrimaryPane:      ChatPane{Contents: []string{}, ChatInput: textinput.New()},
	SecondaryPane:    SystemPane{Contents: []string{}},
	StatusBar:        StatusBar{LeftBlurb: "NUMI'S TEST CLIENT", RightBlurb: "v0.01a", MiddleString: "L: Loading..."},
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
