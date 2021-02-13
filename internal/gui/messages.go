package gui

import (
	"encoding/json"

	"github.com/asticode/go-astilectron"
	"github.com/hbollon/igopher"
	"github.com/sirupsen/logrus"
)

// MessageOut represents a message going out
type MessageOut struct {
	Msg     string      `json:"msg"`
	Payload interface{} `json:"payload,omitempty"`
}

// MessageIn represents a message going in
type MessageIn struct {
	Msg     string          `json:"msg"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

func handleMessages(w *astilectron.Window) {
	w.OnMessage(func(m *astilectron.EventMessage) interface{} {
		// Unmarshal
		var i MessageIn
		var err error
		if err = m.Unmarshal(&i); err != nil {
			logrus.Errorf("Unmarshaling message %+v failed: %v", *m, err)
			return MessageOut{Msg: "Error during message reception"}
		}

		// Process message
		switch i.Msg {
		case "resetGlobalDefaultSettings":
			config := igopher.ResetBotConfig()
			igopher.ExportConfig(config)
			return MessageOut{Msg: "Global settings successfully reset!"}
		default:
			logrus.Error("Unexpected message received.")
			return MessageOut{Msg: "Error: Unexpected message received."}
		}
	})
}
