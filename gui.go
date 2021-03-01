package igopher

import (
	"fmt"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"github.com/sirupsen/logrus"
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

	// Add a listener on Astilectron crash event for selenium cleaning
	a.On(astilectron.EventNameAppCrash, func(e astilectron.Event) (deleteListener bool) {
		logrus.Error("Electron app has crashed!")
		BotStruct.SeleniumStruct.CloseSelenium()
		return
	})

	// Add a listener on Astilectron close event for selenium cleaning
	a.On(astilectron.EventNameAppClose, func(e astilectron.Event) (deleteListener bool) {
		logrus.Debug("Electron app was closed")
		BotStruct.SeleniumStruct.CloseSelenium()
		return
	})

	// Start
	if err = a.Start(); err != nil {
		log.Fatal(fmt.Errorf("main: starting astilectron failed: %w", err))
	}

	// New window
	if w, err = a.NewWindow("./gui/dm_automation.html", &astilectron.WindowOptions{
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

	// Open dev tools panel if flag is set
	if *flags.DevToolsFlag {
		w.OpenDevTools()
	}

	// Blocking pattern
	a.Wait()
}
