package gui

import (
	"github.com/asticode/go-astilectron"
	"github.com/hbollon/igopher"
)

func handleMessages(w *astilectron.Window) {
	w.OnMessage(func(m *astilectron.EventMessage) interface{} {
		// Unmarshal
		var s string
		m.Unmarshal(&s)

		// Process message
		if s == "resetGlobalDefaultSettings" {
			config := igopher.ResetBotConfig()
			igopher.ExportConfig(config)
			return "Global settings successfully reset!"
		}
		return nil
	})
}
