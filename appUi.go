package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"os"
	"strings"
)

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

type AppConfiguration struct {
	activeUsername string
	loginToken string
}

type AppState struct {
	stage          AppStage
	errorState     ErrorState
	cfg AppConfiguration
	sessionToken   string
}

type AppModel struct {
	state         AppState
	rcv           ServerEventReceiver
	altWindow     AltWindow
	infoPane      CharacterPane
	primaryPane   ChatPane
	secondaryPane SystemPane
	statusBar     StatusBar
}

type AltWindow struct {
	IsFocused bool
	Contents  []string
}

func (scr *AppModel) Init() tea.Cmd {
	return RunSignalRClient(&scr.rcv)
}

func (scr *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case ServerConnectionEstablishedMsg:
		scr.altWindow.Contents = append(scr.altWindow.Contents, "Connected!")
		scr.altWindow.IsFocused = !scr.altWindow.IsFocused
		tea.ExitAltScreen()
		scr.primaryPane.IsFocused = true
		scr.primaryPane.ChatInput.Blink()
		scr.primaryPane.ChatInput.Focus()
		return scr, scr.Listen(scr.rcv.UiUpdateChannel) // needs to kick off some sort of listener for incoming signalR invokes

	case ServerDataReceivedMsg:
		return scr, scr.Listen(scr.rcv.UiUpdateChannel)

	case ErrMsg:
		switch msg.ErrType {

		case ServerConnectionTimeout:
			scr.state.errorState = ServerConnectionTimeout
			scr.altWindow.Contents = append(scr.altWindow.Contents, "Server Connection Failed. Retry? (Y/N)")
			return scr, nil

		case FatalError:
			return scr, tea.Quit

		}

	case tea.KeyMsg:
		// altscreen keystroke handlers
		if scr.altWindow.IsFocused {
			if scr.state.errorState == ServerConnectionTimeout {
				switch msg.Type {
				case tea.KeyRunes:
					switch string(msg.Runes) {
					case "y":
						scr.state.errorState = NoError
						return scr, RunSignalRClient(&scr.rcv)
					case "n":
						return scr, tea.Quit
					}
				}
			}
		}

		// chat pane keystroke handlers
		if scr.primaryPane.IsFocused {
			switch msg.Type {

			case tea.KeyEnter:
				if scr.primaryPane.ChatInput.Focused() && len(scr.primaryPane.ChatInput.Value()) > 0 { // If there's something typed, handle it
					if string(scr.primaryPane.ChatInput.Value()[0]) == "/" { // if it's a command, send to cmd handler
						scr.ProcessCommand(strings.Split(scr.primaryPane.ChatInput.Value(), " "))
					} else if string(scr.primaryPane.ChatInput.Value()[0]) != "/" {
						if scr.state.sessionToken == "" { // no chatting if not logged in!
							scr.secondaryPane.Contents = append(scr.secondaryPane.Contents, "You can't chat before you log in!")
						} else {
							scr.ProcessChat()
						}
					}
					scr.primaryPane.ChatInput.Reset()
				}
			}

			scr.primaryPane.ChatInput, _ = scr.primaryPane.ChatInput.Update(msg)
		}

		// system pane keystroke handlers
		if scr.secondaryPane.IsFocused {
			switch msg.Type {

			}
		}

		// character pane keystroke handlers
		if scr.infoPane.IsFocused {
			switch msg.Type {

			}
		}

		// universal keystroke handlers
		switch msg.Type {
		case tea.KeyF2: // focus chat pane
			scr.primaryPane.IsFocused = true
			scr.primaryPane.ChatInput.Focus()

			scr.secondaryPane.IsFocused = false
			scr.infoPane.IsFocused = false
			scr.altWindow.IsFocused = false
		case tea.KeyF3: // focus system pane
			scr.secondaryPane.IsFocused = true

			scr.primaryPane.IsFocused = false
			scr.infoPane.IsFocused = false
			scr.altWindow.IsFocused = false
		case tea.KeyF4: // focus character pane
			scr.infoPane.IsFocused = true

			scr.primaryPane.IsFocused = false
			scr.secondaryPane.IsFocused = false
			scr.altWindow.IsFocused = false
		case tea.KeyF5: // focus debug screen
			scr.altWindow.IsFocused = true

			scr.primaryPane.IsFocused = false
			scr.secondaryPane.IsFocused = false
			scr.infoPane.IsFocused = false

		case tea.KeyCtrlD: // dump current state to debug screen
			scr.altWindow.Contents = append(scr.altWindow.Contents, fmt.Sprintf("chatFocused: %v \nchatContents: %v", scr.primaryPane.ChatInput.Focused(), scr.primaryPane.ChatInput.Value()))

		case tea.KeyCtrlQ: // exit program
			return scr, tea.Quit
		}
	}
	return scr, nil
}

func (scr *AppModel) View() string {
	if scr.altWindow.IsFocused {
		return RenderAltView(scr)
	} else {
		return RenderMainView(scr)
	}
}

func RenderMainView(scr *AppModel) string {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))

	rightStack := lipgloss.JoinVertical(lipgloss.Right, scr.primaryPane.RenderChatPane(w, h), scr.secondaryPane.RenderCommandPane(w, h))
	mainApp := lipgloss.JoinHorizontal(lipgloss.Top, scr.infoPane.RenderInfoPane(w, h), rightStack)

	return mainApp + "\n" + scr.statusBar.RenderStatusBar(w)
}

func RenderAltView(scr *AppModel) string {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))

	mainStyle := lipgloss.NewStyle().
		Width(w-2).Height(int(h-3)).Border(lipgloss.DoubleBorder(), true)

	mainApp := mainStyle.Render(strings.Join(scr.altWindow.Contents, "\n"))

	return mainApp + "\n" + scr.statusBar.RenderStatusBar(w)
}

type ChatPane struct {
	IsFocused bool
	Contents  []string
	ChatInput textinput.Model
}

func (pp *ChatPane) RenderChatPane(w int, h int) string {
	style := lipgloss.NewStyle().
		Width(int((w/3)*2)-2).Height(int(((2*h)/3)-2)).Border(lipgloss.DoubleBorder(), true)

	chatHistory := viewport.New(int((w/3)*2)-1, int(((2*h)/3)-3))
	chatHistory.SetContent(strings.Join(pp.Contents, "\n"))

	pp.ChatInput.TextStyle = lipgloss.NewStyle().Background(lipgloss.Color("#AFAFAF")).Foreground(lipgloss.Color("#000000"))
	//pp.ChatInput.BackgroundStyle = lipgloss.NewStyle().Background(lipgloss.Color("#AFAFAF")).Foreground(lipgloss.Color("#000000"))
	//pp.ChatInput.PlaceholderStyle = lipgloss.NewStyle().Background(lipgloss.Color("#AFAFAF")).Foreground(lipgloss.Color("#000000"))
	pp.ChatInput.Width = int((w/3)*2) - 5
	pp.ChatInput.CharLimit = 255
	//pp.ChatInput.Placeholder = "Chat Goes Here..."

	chatHistory.GotoBottom()
	return style.Render(chatHistory.View() + "\n" + pp.ChatInput.View())
}

type SystemPane struct {
	IsFocused    bool
	Contents     []string
	commandInput textinput.Model
}

func (sp *SystemPane) RenderCommandPane(w int, h int) string {
	style := lipgloss.NewStyle().
		Width(int((w/3)*2)-2).Height(int((h/3)-2)).Border(lipgloss.DoubleBorder(), true)

	vp := viewport.New(int((w/3)*2)-2, int((h/3)-2))
	vp.SetContent(strings.Join(sp.Contents, "\n"))

	vp.GotoBottom()
	return style.Render(vp.View())
}

type StatusBar struct {
	LeftBlurb    string
	MiddleString string
	RightBlurb   string
}

func (sb *StatusBar) RenderStatusBar(w int) string {
	var statusBarLeftChunk = lipgloss.NewStyle().Background(lipgloss.Color("#FF5F87")).Foreground(lipgloss.Color("#FFFDF5")).Align(lipgloss.Left)
	var statusBarRightChunk = lipgloss.NewStyle().Background(lipgloss.Color("#FF5F87")).Foreground(lipgloss.Color("#FFFDF5")).Align(lipgloss.Right)

	left := statusBarLeftChunk.Render(sb.LeftBlurb)
	right := statusBarRightChunk.Render(sb.RightBlurb)

	var middleText = lipgloss.NewStyle().Background(lipgloss.Color("#AFAFAF")).Foreground(lipgloss.Color("#000000")).Align(lipgloss.Center).Width(w - lipgloss.Width(left) - lipgloss.Width(right))
	middle := middleText.Render(sb.MiddleString)

	statusBarPrerender := lipgloss.JoinHorizontal(lipgloss.Top, left, middle, right)
	return statusBarPrerender
}
