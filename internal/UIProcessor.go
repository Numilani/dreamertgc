package internal

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"golang.org/x/term"
	"os"
	"strings"
)

func (scr *AppModel) Init() tea.Cmd {
	return RunSignalRClient(scr)
}

func (scr *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case ServerConnectionEstablishedMsg:
		scr.AltWindow.Contents = append(scr.AltWindow.Contents, "Connected!")
		scr.AltWindow.IsFocused = !scr.AltWindow.IsFocused
		tea.ExitAltScreen()
		scr.PrimaryPane.IsFocused = true
		scr.PrimaryPane.ChatInput.Blink()
		scr.PrimaryPane.ChatInput.Focus()
		return scr, scr.ProcessCall(scr.Rcv.UiUpdateChannel) // needs to kick off some sort of listener for incoming signalR invokes

	case ServerDataReceivedMsg:
		return scr, scr.ProcessCall(scr.Rcv.UiUpdateChannel)

	case ErrMsg:
		switch msg.ErrType {

		case ServerConnectionTimeout:
			scr.State.ErrorState = ServerConnectionTimeout
			scr.AltWindow.Contents = append(scr.AltWindow.Contents, "Server Connection Failed. Retry? (Y/N)")
			return scr, nil

		case FatalError:
			return scr, tea.Quit

		}

	case tea.KeyMsg:
		// altscreen keystroke handlers
		if scr.AltWindow.IsFocused {
			if scr.State.ErrorState == ServerConnectionTimeout {
				switch msg.Type {
				case tea.KeyRunes:
					switch string(msg.Runes) {
					case "y":
						scr.State.ErrorState = NoError
						return scr, RunSignalRClient(scr)
					case "n":
						return scr, tea.Quit
					}
				}
			}
		}

		// chat pane keystroke handlers
		if scr.PrimaryPane.IsFocused {
			switch msg.Type {

			case tea.KeyEnter:
				if scr.PrimaryPane.ChatInput.Focused() && len(scr.PrimaryPane.ChatInput.Value()) > 0 { // If there's something typed, handle it
					if string(scr.PrimaryPane.ChatInput.Value()[0]) == "/" { // if it's a command, send to cmd handler
						scr.ProcessCommand(strings.Split(scr.PrimaryPane.ChatInput.Value(), " "))
					} else if string(scr.PrimaryPane.ChatInput.Value()[0]) != "/" {
						if viper.GetString("sessionToken") == "" { // no chatting if not logged in!
							scr.SecondaryPane.Contents = append(scr.SecondaryPane.Contents, "You can't chat before you log in!")
						} else {
							//scr.ProcessChat()
						}
					}
					scr.PrimaryPane.ChatInput.Reset()
				}
			}

			scr.PrimaryPane.ChatInput, _ = scr.PrimaryPane.ChatInput.Update(msg)
		}

		// system pane keystroke handlers
		if scr.SecondaryPane.IsFocused {
			switch msg.Type {

			}
		}

		// character pane keystroke handlers
		if scr.InfoPane.IsFocused {
			switch msg.Type {

			}
		}

		// universal keystroke handlers
		switch msg.Type {
		case tea.KeyF2: // focus chat pane
			scr.PrimaryPane.IsFocused = true
			scr.PrimaryPane.ChatInput.Focus()

			scr.SecondaryPane.IsFocused = false
			scr.InfoPane.IsFocused = false
			scr.AltWindow.IsFocused = false
		case tea.KeyF3: // focus system pane
			scr.SecondaryPane.IsFocused = true

			scr.PrimaryPane.IsFocused = false
			scr.InfoPane.IsFocused = false
			scr.AltWindow.IsFocused = false
		case tea.KeyF4: // focus character pane
			scr.InfoPane.IsFocused = true

			scr.PrimaryPane.IsFocused = false
			scr.SecondaryPane.IsFocused = false
			scr.AltWindow.IsFocused = false
		case tea.KeyF5: // focus debug screen
			scr.AltWindow.IsFocused = true

			scr.PrimaryPane.IsFocused = false
			scr.SecondaryPane.IsFocused = false
			scr.InfoPane.IsFocused = false

		case tea.KeyCtrlD: // dump current State to debug screen
			scr.AltWindow.Contents = append(scr.AltWindow.Contents, fmt.Sprintf("chatFocused: %v \nchatContents: %v", scr.PrimaryPane.ChatInput.Focused(), scr.PrimaryPane.ChatInput.Value()))

		case tea.KeyCtrlQ: // exit program
			return scr, tea.Quit
		}
	}
	return scr, nil
}

func (scr *AppModel) View() string {
	if scr.AltWindow.IsFocused {
		return RenderAltView(scr)
	} else {
		return RenderMainView(scr)
	}
}

func RenderMainView(scr *AppModel) string {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))

	rightStack := lipgloss.JoinVertical(lipgloss.Right, scr.PrimaryPane.RenderChatPane(w, h), scr.SecondaryPane.RenderCommandPane(w, h))
	mainApp := lipgloss.JoinHorizontal(lipgloss.Top, scr.InfoPane.RenderCharacterPane(w, h), rightStack)

	return mainApp + "\n" + scr.StatusBar.RenderStatusBar(w)
}

type AltWindow struct {
	IsFocused bool
	Contents  []string
}

func RenderAltView(scr *AppModel) string {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))

	mainStyle := lipgloss.NewStyle().
		Width(w-2).Height(int(h-3)).Border(lipgloss.DoubleBorder(), true)

	mainApp := mainStyle.Render(strings.Join(scr.AltWindow.Contents, "\n"))

	return mainApp + "\n" + scr.StatusBar.RenderStatusBar(w)
}

func (scr *AppModel) ProcessCommand(cmdData []string) {
	scr.AltWindow.Contents = append(scr.AltWindow.Contents, fmt.Sprintf("Sent command: %v", scr.PrimaryPane.ChatInput.Value()))

	command := cmdData[0][1:]

	commandsPermittedWhileUnverified := []string{"login", "tokentest"}
	if !slices.Contains(commandsPermittedWhileUnverified, command) && viper.GetString("sessionToken") == "" {
		scr.SecondaryPane.Contents = append(scr.SecondaryPane.Contents, "You can't do that before logging in!\n/login <user> <pass>")
		scr.PrimaryPane.ChatInput.Reset()
		return
	}

	switch command {

	case "clearcon":
		scr.SecondaryPane.Contents = []string{}

	case "clearchat":
		scr.PrimaryPane.Contents = []string{}

	case "quit":
		tea.Quit()

	case "charsay":
		if len(cmdData) < 3 {
			scr.SecondaryPane.Contents = append(scr.SecondaryPane.Contents, "Missing params: /charsay <char name> <message>")
			return
		}
		scr.ConnectionClient.Invoke("CharSay", cmdData[1], cmdData[2])

	case "login":
		if len(cmdData) < 3 {
			scr.SecondaryPane.Contents = append(scr.SecondaryPane.Contents, "Missing params: /login <username> <password>")
			return
		}
		scr.ConnectionClient.Invoke("GetLoginToken", cmdData[1], cmdData[2])

	case "tokentest":
		if len(cmdData) < 2 {
			scr.SecondaryPane.Contents = append(scr.SecondaryPane.Contents, "Missing params: /tokentest <token>")
			return
		}
		scr.ConnectionClient.Invoke("LoginWithToken", cmdData[1])
	}

}
