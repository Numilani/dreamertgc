package main

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

// CharacterPane handles the data for the leftmost pane of the main UI.
type CharacterPane struct {
	IsFocused bool
	Contents  string
}

func (ip *CharacterPane) RenderInfoPane(w int, h int) string {
	style := lipgloss.NewStyle().
		Width(int(w/3)-2).Height(int(h-3)).Border(lipgloss.DoubleBorder(), true)

	return style.Render(ip.Contents)
}

// CharacterStatusData represents a response from the server containing the full status of a character.
type CharacterStatusData struct {
	Guid     string
	Name     string
	Hp       int
	Statuses []string
}

// ReceivePlayerStats is an RPC function, deserializing a JSONified CharacterStatusData and passing it to the UI update channel.
func (c *ServerEventReceiver) ReceivePlayerStats(msg string) {
	var stats CharacterStatusData
	err := json.Unmarshal([]byte(msg), &stats)
	if err != nil {
		fmt.Println(err)
	}
	c.UiUpdateChannel <- ServerDataChunk{
		CallerName: "ReceivePlayerStats",
		Data:       stats,
	}
}
