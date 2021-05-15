package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hbollon/igopher"
	log "github.com/sirupsen/logrus"
)

func (m model) UpdateHomePage(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case ctrlC:
			return m, tea.Quit

		case up, "k":
			if m.homeScreen.cursor > 0 {
				m.homeScreen.cursor--
			}

		case down, "j":
			if m.homeScreen.cursor < len(m.homeScreen.choices)-1 {
				m.homeScreen.cursor++
			}

		case enter:
			errorMessage = ""
			infoMessage = ""
			switch m.homeScreen.cursor {
			case 0:
				return launchBot(m)
			case 1:
				config = igopher.ImportConfig()
				m.screen = settingsMenu
			case 2:
				m.screen = settingsResetMenu
			case 3:
				if m.instanceAlreadyRunning {
					m.screen = stopRunningInstance
				} else {
					return m, tea.Quit
				}
			case 4:
				if m.instanceAlreadyRunning {
					return m, tea.Quit
				}
				log.Warn("Invalid input!")
			default:
				log.Warn("Invalid input!")
			}
		}
	}
	return m, nil
}

func (m model) ViewHomePage() string {
	s := fmt.Sprintf("\n🦄 Welcome to %s, the (soon) most powerful and versatile %s bot!\n\n", keyword("IGopher"), keyword("Instagram"))
	if errorMessage != "" {
		s += errorColor(errorMessage)
	} else {
		s += infoColor(infoMessage)
	}

	for i, choice := range m.homeScreen.choices {
		cursor := " "
		if m.homeScreen.cursor == i {
			cursor = cursorColor(">")
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += subtle("\nup/down: select") + dot + subtle("enter: choose") + dot + subtle("ctrl+c: quit")
	return s
}

func (m *model) updateMenuItemsHomePage() {
	if m.instanceAlreadyRunning {
		m.homeScreen.choices = []string{"🚀 - Launch!", "🔧 - Configure", "🧨 - Reset settings", "☠️ - Stop running instance", "🚪 - Exit"}
	} else {
		m.homeScreen.choices = []string{"🚀 - Launch!", "🔧 - Configure", "🧨 - Reset settings", "🚪 - Exit"}
	}
}
