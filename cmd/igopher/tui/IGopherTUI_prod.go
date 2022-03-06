//go:build !dev

package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/hbollon/igopher/internal/automation"
	"github.com/hbollon/igopher/internal/config"
	"github.com/hbollon/igopher/internal/config/flags"
	"github.com/hbollon/igopher/internal/logger"
	"github.com/hbollon/igopher/internal/process"
	tui "github.com/hbollon/igopher/internal/tui"
	"github.com/hbollon/igopher/internal/utils"
	"github.com/sirupsen/logrus"
)

func init() {
	flags.Flags.BackgroundFlag = flag.Bool("background-task", false,
		"Run IGopher as background task with actual configuration (configure it normally and after re-run IGopher with this flag)")
}

// Change the current working directory by executable location one
func changeWorkingDir() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	changeWorkingDir()
	logger.InitLogger()

	// Initialize environment
	config.CheckEnvironment()

	alreadyRunning, _ := process.CheckIfAlreadyRunning()
	if *flags.Flags.BackgroundFlag {
		if alreadyRunning {
			logrus.Error("IGopher is already running! Kill it or close it through TUI interface and retry.")
			return
		}
		logrus.Debug("Successfully dump pid to tmp file!")
		automation.LaunchBotTui()
	} else {
		// Clear terminal session
		utils.ClearTerminal()

		// Launch TUI
		execBot := tui.InitTui(alreadyRunning)

		// Launch bot if option selected
		if execBot {
			automation.LaunchBotTui()
		}
	}
}
