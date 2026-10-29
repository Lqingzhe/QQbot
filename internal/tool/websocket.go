package until

import (
	"qqbot/internal/model"

	"github.com/gorilla/websocket"
)

func MakeTextWebsocketMessage(message []byte) model.WebsocketMessage {
	return model.WebsocketMessage{
		MessageType: websocket.TextMessage,
		Message:     message,
	}
}
func MakeBinaryWebsocketMessage(message []byte) model.WebsocketMessage {
	return model.WebsocketMessage{
		MessageType: websocket.BinaryMessage,
		Message:     message,
	}
}
