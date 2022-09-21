package main

import tea "github.com/charmbracelet/bubbletea"

func (scr *AppMainModel) ProcessCommand(cmdData []string) {
	command := cmdData[0][1:]

	switch command {

	case "clearcon":
		scr.secondaryPane.Contents = []string{}

	case "clearchat":
		scr.primaryPane.Contents = []string{}

	case "quit":
		tea.Quit()

	case "login":
		if len(cmdData) < 3 {
			scr.secondaryPane.Contents = append(scr.secondaryPane.Contents, "Missing params: /login <username> <password>")
			return
		}
		client.Invoke("GetLoginToken", cmdData[1], cmdData[2])

	case "tokentest":
		if len(cmdData) < 2 {
			scr.secondaryPane.Contents = append(scr.secondaryPane.Contents, "Missing params: /tokentest <token>")
			return
		}
		client.Invoke("LoginWithToken", cmdData[1])
	}

}
