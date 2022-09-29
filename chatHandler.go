package main

func (scr *AppModel) ProcessChat() {
	scr.primaryPane.Contents = append(scr.primaryPane.Contents, RenderSentChat("You", scr.primaryPane.ChatInput.Value()))
}
