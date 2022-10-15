package renderers

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

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

// CharacterPane handles the data for the leftmost pane of the main UI.
type CharacterPane struct {
	IsFocused bool
	Contents  string
}

func (ip *CharacterPane) RenderCharacterPane(w int, h int) string {
	style := lipgloss.NewStyle().
		Width(int(w/3)-2).Height(int(h-3)).Border(lipgloss.DoubleBorder(), true)

	return style.Render(ip.Contents)
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
