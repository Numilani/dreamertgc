package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"hello_go/models"
	"os"
)

type appmodel struct {
	infoPane      models.InfoPaneData
	primaryPane   models.PrimaryPaneData
	secondaryPane models.SecondaryPaneData
	statusBar     models.StatusBarData
	rcv           receiver
}

func (scr *appmodel) Init() tea.Cmd {
	return runSignalRClient(&scr.rcv)
}

func (scr *appmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case models.ServerConnectionEstablishedMsg:
		return scr, scr.Listen(scr.rcv.updateChannel) // needs to kick off some sort of listener for incoming signalR invokes

	case models.ServerDataReceivedMsg:
		return scr, scr.Listen(scr.rcv.updateChannel)

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl-c", "q":
			return scr, tea.Quit

		}
	}
	return scr, nil
}

func (scr *appmodel) View() string {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))

	doc := ""

	leftModule := lipgloss.NewStyle().
		Width(int(w/3)-2).Height(int(h-3)).Border(lipgloss.DoubleBorder(), true)

	topModule := lipgloss.NewStyle().
		Width(int((w/3)*2)-1).Height(int((h/2)-2)).Border(lipgloss.DoubleBorder(), true)

	bottomModule := topModule.Copy()

	//statusBar := lipgloss.NewStyle().
	//	Background(lipgloss.Color("#fafafa")).
	//	Width(w).Height(1)

	var statusBarLeftChunk = lipgloss.NewStyle().Background(lipgloss.Color("#FF5F87")).Foreground(lipgloss.Color("#FFFDF5")).Align(lipgloss.Left)
	var statusBarRightChunk = lipgloss.NewStyle().Background(lipgloss.Color("#FF5F87")).Foreground(lipgloss.Color("#FFFDF5")).Align(lipgloss.Right)

	left := statusBarLeftChunk.Render(scr.statusBar.LeftBlurb)
	right := statusBarRightChunk.Render(scr.statusBar.RightBlurb)

	var middleText = lipgloss.NewStyle().Background(lipgloss.Color("#AFAFAF")).Foreground(lipgloss.Color("#000000")).Align(lipgloss.Center).Width(w - lipgloss.Width(left) - lipgloss.Width(right))
	middle := middleText.Render(scr.statusBar.MiddleString)

	statusBarPrerender := lipgloss.JoinHorizontal(lipgloss.Top, left, middle, right)

	rightStack := lipgloss.JoinVertical(lipgloss.Right, topModule.Render(scr.primaryPane.Contents), bottomModule.Render(scr.secondaryPane.Contents))
	mainApp := lipgloss.JoinHorizontal(lipgloss.Top, leftModule.Render(scr.infoPane.Contents), rightStack)

	doc += mainApp + "\n" + statusBarPrerender
	return doc
}
