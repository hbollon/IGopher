// +build !dev

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/hbollon/igopher"
	tui "github.com/hbollon/igopher/internal/tui"
)

func init() {
	igopher.Flags.BackgroundFlag = flag.Bool("background-task", false,
		"Run IGopher as background task with actual configuration (configure it normally and after re-run IGopher with this flag)")
}

// Change the current working directory by executable location one
func changeWorkingDir() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	} else {
		err = os.Chdir(dir)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	flag.Parse()
	changeWorkingDir()
	igopher.InitLogger()

	// Initialize environment
	igopher.CheckEnvironment()

	// Clear terminal session
	igopher.ClearTerminal()

	// Launch TUI
	execBot := tui.InitTui()

	// Launch bot if option selected
	if execBot {
		igopher.LaunchBotTui()
	}
}
