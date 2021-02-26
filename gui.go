package igopher

import (
	"fmt"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	log "github.com/sirupsen/logrus"
)

var w *astilectron.Window

func InitGui() {
	// Create astilectron
	a, err := astilectron.New(log.StandardLogger(), astilectron.Options{
		AppName:           "IGopher",
		BaseDirectoryPath: "./lib/electron",
	})
	if err != nil {
		log.Fatal(fmt.Errorf("main: creating astilectron failed: %w", err))
	}
	defer a.Close()

	// Handle signals
	a.HandleSignals()

	// Start
	if err = a.Start(); err != nil {
		log.Fatal(fmt.Errorf("main: starting astilectron failed: %w", err))
	}

	// New window
	if w, err = a.NewWindow("./gui/dashboard.html", &astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
		Width:  astikit.IntPtr(1400),
		Height: astikit.IntPtr(1000),
	}); err != nil {
		log.Fatal(fmt.Errorf("main: new window failed: %w", err))
	}

	// Create windows
	if err = w.Create(); err != nil {
		log.Fatal(fmt.Errorf("main: creating window failed: %w", err))
	}
	handleMessages()
	w.OpenDevTools()

	// Blocking pattern
	a.Wait()
}
