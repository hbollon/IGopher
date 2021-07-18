package tui

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hbollon/igopher"
	log "github.com/sirupsen/logrus"
)

const invalidInputMsg = "Invalid input, please check all fields.\n\n"

func (m model) UpdateSettingsInputsMenu(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case ctrlC:
			return m, tea.Quit

		case ctrlB:
			errorMessage = ""
			m.screen = settingsMenu

		case enter:
			if m.settingsInputsScreen.index == len(m.settingsInputsScreen.input) {
				switch m.settingsChoice {
				case accountSettings:
					acc := igopher.AccountYaml{
						Username: m.settingsInputsScreen.input[0].Value(),
						Password: m.settingsInputsScreen.input[1].Value(),
					}
					err := validate.Struct(acc)
					if err != nil {
						errorMessage = invalidInputMsg
						break
					} else {
						config.Account = acc
						errorMessage = ""
						m.screen = settingsMenu
					}
				case scrappingSettings:
					val, err := strconv.Atoi(m.settingsInputsScreen.input[1].Value())
					if err == nil {
						scr := igopher.ScrapperYaml{
							Accounts: strings.Split(m.settingsInputsScreen.input[0].Value(), ";"),
							Quantity: val,
						}
						err := validate.Struct(scr)
						if err != nil {
							errorMessage = invalidInputMsg
							break
						} else {
							config.SrcUsers = scr
							errorMessage = ""
							m.screen = settingsMenu
						}
					} else {
						errorMessage = "Invalid quantity field, value must be numeric.\n\n"
					}
				case autodmSettings:
					dm := igopher.AutoDmYaml{
						DmTemplates: strings.Split(m.settingsInputsScreen.input[0].Value(), ";"),
					}
					err := validate.Struct(dm)
					if err != nil {
						errorMessage = invalidInputMsg
						break
					} else {
						config.AutoDm.DmTemplates = dm.DmTemplates
						errorMessage = ""
						m.screen = settingsMenu
					}
				case autodmGreetingSettings:
					gre := igopher.GreetingYaml{
						Template: m.settingsInputsScreen.input[0].Value(),
					}
					err := validate.Struct(gre)
					if err != nil {
						errorMessage = invalidInputMsg
						break
					} else {
						config.AutoDm.Greeting.Template = gre.Template
						errorMessage = ""
						m.screen = settingsMenu
					}
				case quotasSettings:
					dmDay, err := strconv.Atoi(m.settingsInputsScreen.input[0].Value())
					dmHour, err2 := strconv.Atoi(m.settingsInputsScreen.input[1].Value())
					if err == nil && err2 == nil {
						quo := igopher.QuotasYaml{
							DmDay:  dmDay,
							DmHour: dmHour,
						}
						err := validate.Struct(quo)
						if err != nil {
							errorMessage = invalidInputMsg
							break
						} else {
							config.Quotas.DmDay = quo.DmDay
							config.Quotas.DmHour = quo.DmHour
							errorMessage = ""
							m.screen = settingsMenu
						}
					} else {
						errorMessage = invalidInputMsg
					}
				case scheduleSettings:
					sche := igopher.ScheduleYaml{
						BeginAt: m.settingsInputsScreen.input[0].Value(),
						EndAt:   m.settingsInputsScreen.input[1].Value(),
					}
					err := validate.Struct(sche)
					if err != nil {
						errorMessage = invalidInputMsg
						break
					} else {
						config.Schedule.BeginAt = sche.BeginAt
						config.Schedule.EndAt = sche.EndAt
						errorMessage = ""
						m.screen = settingsMenu
					}
				default:
					log.Error("Unexpected settings screen value!\n\n")
				}
				break
			}

		// Cycle between inputs
		case "tab", shiftTab, up, down:
			s := msg.String()

			// Cycle indexes
			if s == up || s == shiftTab {
				m.settingsInputsScreen.index--
			} else {
				m.settingsInputsScreen.index++
			}

			if m.settingsInputsScreen.index > len(m.settingsInputsScreen.input) {
				m.settingsInputsScreen.index = 0
			} else if m.settingsInputsScreen.index < 0 {
				m.settingsInputsScreen.index = len(m.settingsInputsScreen.input)
			}

			for i := 0; i < len(m.settingsInputsScreen.input); i++ {
				if i == m.settingsInputsScreen.index {
					// Set focused state
					m.settingsInputsScreen.input[i].Focus()
					m.settingsInputsScreen.input[i].Prompt = focusedPrompt
					m.settingsInputsScreen.input[i].TextColor = focusedTextColor
					continue
				}
				// Remove focused state
				m.settingsInputsScreen.input[i].Blur()
				m.settingsInputsScreen.input[i].Prompt = blurredPrompt
				m.settingsInputsScreen.input[i].TextColor = ""
			}

			if m.settingsInputsScreen.index == len(m.settingsInputsScreen.input) {
				m.settingsInputsScreen.submitButton = focusedSubmitButton
			} else {
				m.settingsInputsScreen.submitButton = blurredSubmitButton
			}

			return m, nil
		}
	}
	// Handle character input and blinks
	m, cmd := updateInputs(msg, m)
	return m, cmd
}

func (m model) ViewSettingsInputsMenu() string {
	s := m.settingsInputsScreen.title
	s += errorColor(errorMessage)
	for i := 0; i < len(m.settingsInputsScreen.input); i++ {
		s += m.settingsInputsScreen.input[i].View()
		if i < len(m.settingsInputsScreen.input)-1 {
			s += "\n"
		}
	}
	s += "\n\n" + m.settingsInputsScreen.submitButton + "\n"
	s += subtle("\nup/down: select") + dot + subtle("enter: choose") + dot + subtle("ctrl+b: back") + dot + subtle("ctrl+c: quit")
	return s
}
