// +build dev

package main

import (
	"flag"
	"github.com/hbollon/igopher"
	tui "github.com/hbollon/igopher/internal/tui"
func init() {
	igopher.Flags.BackgroundFlag = flag.Bool("background-task", false,
		"Run IGopher as background task with actual configuration (configure it normally and after re-run IGopher with this flag)")
}

)

func main() {
	flag.Parse()
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
