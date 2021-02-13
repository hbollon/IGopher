package gui

import (
	"encoding/json"

	"github.com/asticode/go-astilectron"
	"github.com/go-playground/validator/v10"
	"github.com/hbollon/igopher"
	"github.com/sirupsen/logrus"
)

const (
	SUCCESS = "Success"
	FAIL    = "Fail"
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
	config := igopher.ImportConfig()
	validate := validator.New()
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
			config = igopher.ResetBotConfig()
			igopher.ExportConfig(config)
			return MessageOut{Msg: SUCCESS}

		case "igCredentialsForm":
			var credentialsConfig igopher.AccountYaml
			// Unmarshal payload
			if err = json.Unmarshal([]byte(i.Payload), &credentialsConfig); err != nil {
				logrus.Errorf("Failed to unmarshal message payload: %v", err)
				return MessageOut{Msg: FAIL}
			}

			err = validate.Struct(credentialsConfig)
			if err != nil {
				logrus.Warning("Validation issue on credentials form, abort.")
				return MessageOut{Msg: FAIL}
			}

			config.Account = credentialsConfig
			igopher.ExportConfig(config)
			return MessageOut{Msg: SUCCESS}
		default:
			logrus.Error("Unexpected message received.")
			return MessageOut{Msg: FAIL}
		}
	})
}
