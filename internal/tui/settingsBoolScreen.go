package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
)

func (m model) UpdateSettingsBoolMenu(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case ctrlC:
			return m, tea.Quit

		case ctrlB:
			m.screen = settingsMenu

		case up, "k":
			if m.settingsTrueFalseScreen.cursor > 0 {
				m.settingsTrueFalseScreen.cursor--
			}

		case down, "j":
			if m.settingsTrueFalseScreen.cursor < len(m.settingsTrueFalseScreen.choices)-1 {
				m.settingsTrueFalseScreen.cursor++
			}

		case enter:
			switch m.settingsTrueFalseScreen.cursor {
			case 0:
				switch m.settingsChoice {
				case autodmEnablingSettings:
					config.AutoDm.Activated = true

				case autodmGreetingEnablingSettings:
					config.AutoDm.Greeting.Activated = true

				case quotasEnablingSettings:
					config.Quotas.Activated = true

				case scheduleEnablingSettings:
					config.Schedule.Activated = true

				case blacklistEnablingSettings:
					config.Blacklist.Activated = true

				default:
					log.Error("Unexpected settings screen value!")
				}
				m.screen = settingsMenu
			case 1:
				switch m.settingsChoice {
				case autodmEnablingSettings:
					config.AutoDm.Activated = false

				case autodmGreetingEnablingSettings:
					config.AutoDm.Greeting.Activated = false

				case quotasEnablingSettings:
					config.Quotas.Activated = false

				case scheduleEnablingSettings:
					config.Schedule.Activated = false

				case blacklistEnablingSettings:
					config.Blacklist.Activated = false

				default:
					log.Error("Unexpected settings screen value!")
				}
				m.screen = settingsMenu
			default:
				log.Warn("Invalid input!")
				m.screen = settingsMenu
			}
		}
	}
	return m, nil
}

func (m model) ViewSettingsBoolMenu() string {
	var s string
	switch m.settingsChoice {
	case autodmEnablingSettings:
		s = fmt.Sprintf("\nDo you want to enable %s module? (Default: %s)\n\n", keyword("AutoDM"), keyword("true"))

	case autodmGreetingEnablingSettings:
		s = fmt.Sprintf("\nDo you want to enable %s sub-module with %s? (Default: %s)\n\n",
			keyword("Greeting"), keyword("AutoDm"), keyword("true"))

	case quotasEnablingSettings:
		s = fmt.Sprintf("\nDo you want to enable %s module? (Default: %s)\n\n", keyword("Quotas"), keyword("true"))

	case scheduleEnablingSettings:
		s = fmt.Sprintf("\nDo you want to enable %s module? (Default: %s)\n\n", keyword("Scheduler"), keyword("true"))

	case blacklistEnablingSettings:
		s = fmt.Sprintf("\nDo you want to enable %s module? (Default: %s)\n\n", keyword("User Blacklist"), keyword("true"))

	default:
		log.Error("Unexpected settings screen value!")
		s = ""
	}

	for i, choice := range m.settingsTrueFalseScreen.choices {
		cursor := " "
		if m.settingsTrueFalseScreen.cursor == i {
			cursor = cursorColor(">")
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	s += subtle("\nup/down: select") + dot + subtle("enter: choose") + dot + subtle("ctrl+b: back") + dot + subtle("ctrl+c: quit")

	return s
}
