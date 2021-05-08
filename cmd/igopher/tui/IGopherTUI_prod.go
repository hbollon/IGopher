// +build !dev

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/hbollon/igopher"
	tui "github.com/hbollon/igopher/internal/tui"
)

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
	changeWorkingDir()

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
