package main

import (
	"github.com/hbollon/igopher"
	tui "github.com/hbollon/igopher/internal/tui"
)

func main() {
	// Change terminal current working directory with executable location one
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	} else {
		err = os.Chdir(dir)
		if err != nil {
			log.Fatal(err)
		}
	}

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
