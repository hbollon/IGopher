package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hbollon/igopher/internal/process"
	log "github.com/sirupsen/logrus"
)

func (m model) UpdateStopRunningProcess(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case ctrlC:
			return m, tea.Quit

		case ctrlB:
			m.screen = mainMenu

		case up, "k":
			if m.stopRunningProcessScreen.cursor > 0 {
				m.stopRunningProcessScreen.cursor--
			}

		case down, "j":
			if m.stopRunningProcessScreen.cursor < len(m.stopRunningProcessScreen.choices)-1 {
				m.stopRunningProcessScreen.cursor++
			}

		case enter:
			switch m.stopRunningProcessScreen.cursor {
			case 0:
				if err := process.TerminateRunningInstance(); err != nil {
					errorMessage = "Failed to terminate running IGopher instance! If the problem persist try to manually kill it or restart your computer."
				} else {
					infoMessage = "IGopher running instance has been successfully killed!" +
						" You can now run it again or close this TUI and restart IGopher as background task using \"--background-task\" flag.\n\n"
					m.instanceAlreadyRunning = false
					m.updateMenuItemsHomePage()
				}
				m.screen = mainMenu
			case 1:
				m.screen = mainMenu
			default:
				log.Warn("Invalid input!")
			}
		}
	}
	return m, nil
}

func (m model) ViewStopRunningProcess() string {
	s := fmt.Sprintf("\nAn instance of %s is already running, do you want to end it and continue?\n\n",
		keyword("IGopher"))

	for i, choice := range m.stopRunningProcessScreen.choices {
		cursor := " "
		if m.stopRunningProcessScreen.cursor == i {
			cursor = cursorColor(">")
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += subtle("\nup/down: select") + dot + subtle("enter: choose") + dot + subtle("ctrl+b: back") + dot + subtle("ctrl+c: quit")
	return s
}
