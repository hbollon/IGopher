package main

/*
This package aims to allow developers to run the GUI without having to bundle using resources directly.
To execute it just run: go run ./cmd/igopher/gui

For release purpose use gui-bundler package
*/

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"github.com/hbollon/igopher"
	log "github.com/sirupsen/logrus"
)

const (
	AppName            = "IGopher"
	VersionAstilectron = "0.46.0"
	VersionElectron    = "11.1.0"
)

func main() {
	flag.Parse()
	igopher.InitLogger()
	igopher.CheckEnvironment()
	defer igopher.BotStruct.SeleniumStruct.CleanUp()

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

	// Start
	if err = a.Start(); err != nil {
		log.Fatal(fmt.Errorf("main: starting astilectron failed: %w", err))
	}

	// New window
	if w, err = a.NewWindow("./resources/static/vue-igopher/dist/app/index.html", &astilectron.WindowOptions{
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
	w.OpenDevTools()
	// }

	// Blocking pattern
	a.Wait()
}
