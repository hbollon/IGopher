//go:build dev

package main

import (
	"flag"

	"github.com/hbollon/igopher"
	"github.com/hbollon/igopher/internal/process"
	tui "github.com/hbollon/igopher/internal/tui"
	"github.com/sirupsen/logrus"
)

func init() {
	igopher.Flags.BackgroundFlag = flag.Bool("background-task", false,
		"Run IGopher as background task with actual configuration (configure it normally and after re-run IGopher with this flag)")
}

func main() {
	flag.Parse()
	igopher.InitLogger()

	// Initialize environment
	igopher.CheckEnvironment()

	alreadyRunning, _ := process.CheckIfAlreadyRunning()
	if *igopher.Flags.BackgroundFlag {
		if alreadyRunning {
			logrus.Error("IGopher is already running! Kill it or close it through TUI interface and retry.")
			return
		}
		logrus.Debug("Successfully dump pid to tmp file!")
		igopher.LaunchBotTui()
	} else {
		// Clear terminal session
		igopher.ClearTerminal()

		// Launch TUI
		execBot := tui.InitTui(alreadyRunning)

		// Launch bot if option selected
		if execBot {
			igopher.LaunchBotTui()
		}
	}
}
