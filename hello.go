package main

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-kit/log"
	signalr "github.com/philippseith/signalr"
	"hello_go/models"
	"os"
	"strconv"
	"time"
)

var client signalr.Client
var appState = appmodel{
	rcv:           receiver{updateChannel: make(chan ServerDataChunk)},
	infoPane:      models.InfoPaneData{Contents: "Loading character data..."},
	primaryPane:   models.PrimaryPaneData{Contents: "Loading world data..."},
	secondaryPane: models.SecondaryPaneData{Contents: "Loading details..."},
	statusBar:     models.StatusBarData{LeftBlurb: "NUMI'S TEST CLIENT", RightBlurb: "v0.01a", MiddleString: "L: Loading..."},
}

func main() {
	programUi := tea.NewProgram(&appState)
	if err := programUi.Start(); err != nil {
		fmt.Printf("An error occurred: %v", err)
		os.Exit(1)
	}
}

func runSignalRClient(receiver *receiver) tea.Cmd {
	return func() tea.Msg {
		var err error
		client, err = signalr.NewClient(context.Background(), nil,
			signalr.WithReceiver(receiver),
			signalr.WithConnector(func() (signalr.Connection, error) {
				creationCtx, _ := context.WithTimeout(context.Background(), 2*time.Second)
				return signalr.NewHTTPConnection(creationCtx, "https://localhost:7277/test")
			}),
			signalr.Logger(log.NewNopLogger(), false))
		//nil)
		if err != nil {
			return models.ErrMsg{}
		}
		client.Start()
		return models.ServerConnectionEstablishedMsg{}
	}
}

func (scr *appmodel) Listen(ch chan ServerDataChunk) tea.Cmd {
	return func() tea.Msg {
		var chunk = <-ch

		switch chunk.data.(type) {

		case models.CharacterStatusData:
			scr.infoPane.Contents = chunk.data.(models.CharacterStatusData).Name + "\n    HP: " + strconv.Itoa(chunk.data.(models.CharacterStatusData).Hp)
			client.Invoke("Confirmed")
		}

		return models.ServerDataReceivedMsg{}
	}
}
