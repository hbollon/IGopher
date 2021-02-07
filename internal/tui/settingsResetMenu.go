package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hbollon/igopher"
	log "github.com/sirupsen/logrus"
)

func (m model) UpdateSettingsResetMenu(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case ctrlC:
			return m, tea.Quit

		case ctrlB:
			m.screen = mainMenu

		case up, "k":
			if m.configResetScreen.cursor > 0 {
				m.configResetScreen.cursor--
			}

		case down, "j":
			if m.configResetScreen.cursor < len(m.configResetScreen.choices)-1 {
				m.configResetScreen.cursor++
			}

		case enter:
			switch m.configResetScreen.cursor {
			case 0:
				config = igopher.ResetBotConfig()
				igopher.ExportConfig(config)
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

func (m model) ViewSettingsResetMenu() string {
	s := fmt.Sprintf("\nAre you sure you want to %s the default %s? This operation cannot be undone!\n\n",
		keyword("reset"), keyword("settings"))

	for i, choice := range m.configResetScreen.choices {
		cursor := " "
		if m.configResetScreen.cursor == i {
			cursor = cursorColor(">")
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += subtle("\nup/down: select") + dot + subtle("enter: choose") + dot + subtle("ctrl+b: back") + dot + subtle("ctrl+c: quit")
	return s
}
