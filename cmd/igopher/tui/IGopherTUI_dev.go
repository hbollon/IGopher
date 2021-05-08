// +build dev

package main

import (
	"github.com/hbollon/igopher"
	tui "github.com/hbollon/igopher/internal/tui"
)

func main() {
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
