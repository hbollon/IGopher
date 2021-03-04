package main

import (
	"path/filepath"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"github.com/hbollon/igopher"
	"github.com/sirupsen/logrus"
)

const (
	AppName            = "IGopher"
	VersionAstilectron = "0.44.0"
	VersionElectron    = "11.1.0"
)

func main() {
	igopher.CheckEnvironment()

	if err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  filepath.FromSlash("resources/favicon.icns"),
			AppIconDefaultPath: filepath.FromSlash("resources/favicon.png"),
			SingleInstance:     true,
			VersionAstilectron: VersionAstilectron,
			VersionElectron:    VersionElectron,
		},
		Debug:       false,
		Logger:      logrus.StandardLogger(),
		MenuOptions: []*astilectron.MenuItemOptions{},
		OnWait: func(a *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			// Add message handler
			igopher.HandleMessages(ws[0])

			// Add a listener on Astilectron crash event for selenium cleaning
			a.On(astilectron.EventNameAppCrash, func(e astilectron.Event) (deleteListener bool) {
				logrus.Error("Electron app has crashed!")
				igopher.BotStruct.SeleniumStruct.CloseSelenium()
				return
			})

			// Add a listener on Astilectron close event for selenium cleaning
			a.On(astilectron.EventNameAppClose, func(e astilectron.Event) (deleteListener bool) {
				logrus.Debug("Electron app was closed")
				igopher.BotStruct.SeleniumStruct.CloseSelenium()
				return
			})

			return nil
		},
		RestoreAssets: RestoreAssets,
		ResourcesPath: "resources/static",
		Windows: []*bootstrap.Window{{
			Homepage: "dm_automation.html",
			Options: &astilectron.WindowOptions{
				BackgroundColor: astikit.StrPtr("#333"),
				Center:          astikit.BoolPtr(true),
				Width:           astikit.IntPtr(1400),
				Height:          astikit.IntPtr(1000),
			},
		}},
	}); err != nil {
		logrus.Fatalf("running bootstrap failed: %v", err)
	}
}
