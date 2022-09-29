package main

import (
	"context"
	"encoding/json"
	"fmt"
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
// established in AppModel. The chunk is then processed according to caller and data type,
// and the UI is updated accordingly.
func (scr *AppModel) Listen(ch chan ServerDataChunk) tea.Cmd {
	return func() tea.Msg {
		var chunk = <-ch

		switch chunk.CallerName {

		case "ReceivePlayerStats":
			var stats CharacterStatusData
			err := json.Unmarshal([]byte(chunk.Data.(string)), &stats)
			if err != nil {
				fmt.Println(err)
			}
			scr.infoPane.Contents = stats.Name + "\n    HP: " + strconv.Itoa(stats.Hp)

		case "ReceiveLoginToken":
			if chunk.Data.(string) == "invalid_credentials" {
				scr.secondaryPane.Contents = append(scr.secondaryPane.Contents, fmt.Sprintf("Login rejected: %v", chunk.Data.(string)))
				break
			}
			scr.secondaryPane.Contents = append(scr.secondaryPane.Contents, fmt.Sprintf("Login token received: %v", chunk.Data.(string)))
			scr.state.sessionToken = chunk.Data.(string)
			client.Invoke("LoginWithToken", scr.state.sessionToken)

		case "ReceiveSessionToken":
			if chunk.Data.(int) == -1 {
				scr.secondaryPane.Contents = append(scr.secondaryPane.Contents, fmt.Sprintf("User token rejected: %v", chunk.Data.(int)))
				break
			}
			scr.secondaryPane.Contents = append(scr.secondaryPane.Contents, fmt.Sprintf("Logged in, session token is: %v", chunk.Data.(int)))
		}

		return ServerDataReceivedMsg{}
	}
}

func (c *ServerEventReceiver) Receive(caller string, data any) {
	c.UiUpdateChannel <- ServerDataChunk{
		CallerName: caller,
		Data:       data,
	}
}
