package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/exp/slices"
)

func (scr *AppModel) ProcessCommand(cmdData []string) {
	scr.altWindow.Contents = append(scr.altWindow.Contents, fmt.Sprintf("Sent command: %v", scr.primaryPane.ChatInput.Value()))

	command := cmdData[0][1:]

	commandsPermittedWhileUnverified := []string{"login", "tokentest"}
	if !slices.Contains(commandsPermittedWhileUnverified, command) && scr.state.sessionToken == "" {
		scr.secondaryPane.Contents = append(scr.secondaryPane.Contents, "You can't do that before logging in!\n/login <user> <pass>")
		scr.primaryPane.ChatInput.Reset()
		return
	}

	switch command {

	case "clearcon":
		scr.secondaryPane.Contents = []string{}

	case "clearchat":
		scr.primaryPane.Contents = []string{}

	case "quit":
		tea.Quit()

	case "charsay":
		if len(cmdData) < 3 {
			scr.secondaryPane.Contents = append(scr.secondaryPane.Contents, "Missing params: /charsay <char name> <message>")
			return
		}
		client.Invoke("CharSay", cmdData[1], cmdData[2])

	case "login":
		if len(cmdData) < 3 {
			scr.secondaryPane.Contents = append(scr.secondaryPane.Contents, "Missing params: /login <username> <password>")
			return
		}
		client.Invoke("GetLoginToken", cmdData[1], cmdData[2])

	case "tokentest":
		if len(cmdData) < 3 {
			scr.secondaryPane.Contents = append(scr.secondaryPane.Contents, "Missing params: /tokentest <username> <token>")
			return
		}
		client.Invoke("LoginWithToken", cmdData[1], cmdData[2])
	}

}
