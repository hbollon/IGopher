package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
)

func (m model) UpdateGenericMenu(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case ctrlC:
			return m, tea.Quit

		case ctrlB:
			m.screen = settingsMenu

		case up, "k":
			if m.genericMenuScreen.cursor > 0 {
				m.genericMenuScreen.cursor--
			}

		case down, "j":
			if m.genericMenuScreen.cursor < len(m.genericMenuScreen.choices)-1 {
				m.genericMenuScreen.cursor++
			}

		case enter:
			switch m.genericMenuScreen.cursor {
			case 0:
				switch m.settingsChoice {
				case autodmSettingsMenu:
					m.settingsChoice = autodmEnablingSettings
				case autodmGreetingMenu:
					m.settingsChoice = autodmGreetingEnablingSettings
				case quotasSettingsMenu:
					m.settingsChoice = quotasEnablingSettings
				case scheduleSettingsMenu:
					m.settingsChoice = scheduleEnablingSettings
				default:
					log.Warn("Invalid input!")
				}
				m.screen = settingsBoolScreen
			case 1:
				switch m.settingsChoice {
				case autodmSettingsMenu:
					m.settingsInputsScreen = getAutoDmSettings()
					m.settingsChoice = autodmSettings
				case autodmGreetingMenu:
					m.settingsInputsScreen = getAutoDmGreetingSettings()
					m.settingsChoice = autodmGreetingSettings
				case quotasSettingsMenu:
					m.settingsInputsScreen = getQuotasSettings()
					m.settingsChoice = quotasSettings
				case scheduleSettingsMenu:
					m.settingsInputsScreen = getSchedulerSettings()
					m.settingsChoice = scheduleSettings
				default:
					log.Warn("Invalid input!")
				}
				m.screen = settingsInputsScreen
			default:
				log.Warn("Invalid input!")
			}
		}
	}
	return m, nil
}

func (m model) ViewGenericMenu() string {
	s := "\n\n"
	for i, choice := range m.genericMenuScreen.choices {
		cursor := " "
		if m.genericMenuScreen.cursor == i {
			cursor = cursorColor(">")
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += subtle("\nup/down: select") + dot + subtle("enter: choose") + dot + subtle("ctrl+b: back") + dot + subtle("ctrl+c: quit")
	return s
}
