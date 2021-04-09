package main

/*
This package aims to allow developers to run the GUI without having to bundle using resources directly.
To execute it just run: go run ./cmd/igopher/gui

For release purpose use gui-bundler package
*/

import (
	"fmt"
	"path/filepath"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"github.com/hbollon/igopher"
	log "github.com/sirupsen/logrus"
)

const (
	AppName            = "IGopher"
	VersionAstilectron = "0.44.0"
	VersionElectron    = "11.1.0"
)

func main() {
	igopher.CheckEnvironment()

	var w *astilectron.Window
	// Create astilectron
	a, err := astilectron.New(log.StandardLogger(), astilectron.Options{
		AppName:            "IGopher",
		AppIconDarwinPath:  filepath.FromSlash("resources/favicon.icns"),
		AppIconDefaultPath: filepath.FromSlash("resources/favicon.png"),
		BaseDirectoryPath:  "./lib/electron",
		SingleInstance:     true,
		VersionAstilectron: VersionAstilectron,
		VersionElectron:    VersionElectron,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("main: creating astilectron failed: %w", err))
	}
	defer a.Close()

	// Handle signals
	a.HandleSignals()

	// Add a listener on Astilectron crash event for selenium cleaning
	a.On(astilectron.EventNameAppCrash, func(e astilectron.Event) (deleteListener bool) {
		log.Error("Electron app has crashed!")
		igopher.BotStruct.SeleniumStruct.CloseSelenium()
		return
	})

	// Add a listener on Astilectron close event for selenium cleaning
	a.On(astilectron.EventNameAppClose, func(e astilectron.Event) (deleteListener bool) {
		log.Debug("Electron app was closed")
		igopher.BotStruct.SeleniumStruct.CloseSelenium()
		return
	})

	// Start
	if err = a.Start(); err != nil {
		log.Fatal(fmt.Errorf("main: starting astilectron failed: %w", err))
	}

	// New window
	if w, err = a.NewWindow("./resources/static/app/dm_automation.html", &astilectron.WindowOptions{
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
	igopher.HandleMessages(w)

	// Open dev tools panel if flag is set
	// if *flags.DevToolsFlag {
	// 	w.OpenDevTools()
	// }

	// Blocking pattern
	a.Wait()
}