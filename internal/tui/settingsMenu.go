package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hbollon/igopher"
	log "github.com/sirupsen/logrus"
)

func (m model) UpdateSettingsMenu(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case ctrlC:
			return m, tea.Quit

		case ctrlB:
			m.screen = mainMenu

		case up, "k":
			if m.configScreen.cursor > 0 {
				m.configScreen.cursor--
			}

		case down, "j":
			if m.configScreen.cursor < len(m.configScreen.choices)-1 {
				m.configScreen.cursor++
			}

		case enter:
			switch m.configScreen.cursor {
			case 0:
				m.settingsInputsScreen = getAccountSettings()
				m.screen = settingsInputsScreen
				m.settingsChoice = accountSettings
			case 1:
				m.settingsInputsScreen = getUsersScrappingSettings()
				m.screen = settingsInputsScreen
				m.settingsChoice = scrappingSettings
			case 2:
				m.genericMenuScreen = menu{choices: []string{"Enable/Disable Module", "Configuration"}}
				m.screen = genericMenu
				m.settingsChoice = autodmSettingsMenu
			case 3:
				m.genericMenuScreen = menu{choices: []string{"Enable/Disable Module", "Configuration"}}
				m.screen = genericMenu
				m.settingsChoice = autodmGreetingMenu
			case 4:
				m.genericMenuScreen = menu{choices: []string{"Enable/Disable Module", "Configuration"}}
				m.screen = genericMenu
				m.settingsChoice = quotasSettingsMenu
			case 5:
				m.genericMenuScreen = menu{choices: []string{"Enable/Disable Module", "Configuration"}}
				m.screen = genericMenu
				m.settingsChoice = scheduleSettingsMenu
			case 6:
				m.screen = settingsBoolScreen
				m.settingsChoice = blacklistEnablingSettings
			case 7:
				m.screen = settingsProxyScreen
			case 8:
				igopher.ExportConfig(config)
				m.screen = mainMenu
			default:
				log.Warn("Invalid input!")
			}
		}
	}
	return m, nil
}

func (m model) ViewSettingsMenu() string {
	s := fmt.Sprintf("\nWhat would you like to %s?\n\n", keyword("tweak"))

	for i, choice := range m.configScreen.choices {
		cursor := " "
		if m.configScreen.cursor == i {
			cursor = cursorColor(">")
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += subtle("\nup/down: select") + dot + subtle("enter: choose") + dot + subtle("ctrl+b: save & back") + dot + subtle("ctrl+c: quit")
	return s
}
