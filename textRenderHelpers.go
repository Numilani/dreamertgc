package main

import "github.com/charmbracelet/lipgloss"

func RenderSentChat(name string, msg string) string {
	x := lipgloss.NewStyle().Bold(true).Render(name + ": ")
	return x + msg
}

func RenderReceivedChat(name string, msg string) string {
	x := lipgloss.NewStyle().Bold(true).Render(name + ": ")
	return lipgloss.NewStyle().Align(lipgloss.Right).Render(x + msg)
}

func RenderWorldMessage(msg string) string {
	return lipgloss.NewStyle().Align(lipgloss.Center).Render(msg)
}
