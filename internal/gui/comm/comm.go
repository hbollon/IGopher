package comm

import (
	"github.com/asticode/go-astilectron"
	"github.com/hbollon/igopher/internal/gui/datatypes"
)

var (
	Window *astilectron.Window
)

// IsElectronRunning checks if electron is running
func IsElectronRunning() bool {
	return Window != nil
}

// SendMessageToElectron will send a message to Electron Gui and execute a callback
// Callback function is optional
func SendMessageToElectron(msg datatypes.MessageOut, callbacks ...astilectron.CallbackMessage) {
	if IsElectronRunning() {
		Window.SendMessage(msg, callbacks...)
	}
}
