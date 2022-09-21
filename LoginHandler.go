package main

func (c *ServerEventReceiver) ReceiveLoginToken(msg string) {
	c.UiUpdateChannel <- ServerDataChunk{
		CallerName: "ReceiveLoginToken",
		Data:       msg,
	}
}
