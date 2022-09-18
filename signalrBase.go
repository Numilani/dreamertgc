package main

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-kit/log"
	"github.com/philippseith/signalr"
	"strconv"
	"time"
)

// ServerEventReceiver
//
// Hub handles any incoming messages from the server, and dispatches all RPC calls accordingly.
//
// UiUpdateChannel blocks the Listen() thread until a UI update is required by the received Data.
type ServerEventReceiver struct {
	Hub             signalr.Hub
	UiUpdateChannel chan ServerDataChunk
}

type ServerDataChunk struct {
	CallerName string
	Data       any
}

// RunSignalRClient runs within the BubbleTea TUI application context.
// It attempts to establish an HTTP connection to the Oneiros server.
//
// Once connection is established, the SignalR websocket connection is created and
// the function yields a ServerConnectionEstablishedMsg.
//
// If the connection is not established within 5 seconds, the function yields an ErrMsg
// and may be retried upon further prompts.
func RunSignalRClient(receiver *ServerEventReceiver) tea.Cmd {
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
			return ErrMsg{}
		}
		client.Start()

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		connectedSignal := client.WaitForState(ctx, signalr.ClientConnected)
		select {
		// I think selects pick whichever one completes first,
		// so if the connection happens before the context times out, it continues
		// otherwise it sends the error signal.
		case <-connectedSignal:
			return ServerConnectionEstablishedMsg{}
		case <-ctx.Done():
			return ErrMsg{}
		}

	}
}

// Listen runs within the BubbleTea TUI application context.
//
// Its thread is created after a successful ServerConnectionEstablishedMsg is received.
//
// Listen blocks its thread until a ServerDataChunk is received from the receiver
// established in AppMainModel. The chunk is then processed according to caller and data type,
// and the UI is updated accordingly.
func (scr *AppMainModel) Listen(ch chan ServerDataChunk) tea.Cmd { // TODO: This should probably be moved to appUi.go to be with other AMM functions
	return func() tea.Msg {
		var chunk = <-ch

		switch chunk.Data.(type) {

		case CharacterStatusData:
			scr.infoPane.Contents = chunk.Data.(CharacterStatusData).Name + "\n    HP: " + strconv.Itoa(chunk.Data.(CharacterStatusData).Hp)
			client.Invoke("Confirmed")
		}

		return ServerDataReceivedMsg{}
	}
}
