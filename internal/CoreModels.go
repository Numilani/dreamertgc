package internal

import (
	. "dreamer_tgc/internal/renderers"
	"github.com/philippseith/signalr"
)

// CharacterStatusData represents a response from the server containing the full status of a character.
type CharacterStatusData struct {
	Guid     string
	Name     string
	Hp       int
	Statuses []string
}

type AppStage int

const (
	Starting AppStage = iota
	ConnectingToServer
	LoadingInitialData
	Ready
)

type ErrorState int

const (
	ServerConnectionTimeout ErrorState = iota
	UnknownError
	FatalError
	NoError
)

type AppState struct {
	Stage      AppStage
	ErrorState ErrorState
}

type AppModel struct {
	State            AppState
	ConnectionClient signalr.Client
	Rcv              ServerEventReceiver
	AltWindow        AltWindow
	InfoPane         CharacterPane
	PrimaryPane      ChatPane
	SecondaryPane    SystemPane
	StatusBar        StatusBar
}

type ErrMsg struct {
	ErrType ErrorState
}

type ServerConnectionEstablishedMsg struct{}

type ServerDataReceivedMsg struct{}
