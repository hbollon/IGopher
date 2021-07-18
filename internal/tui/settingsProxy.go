package tui

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hbollon/igopher/internal/proxy"
)

func (m model) UpdateSettingsProxy(msg tea.Msg) (model, tea.Cmd) {
	menuLength := len(m.settingsProxy.states) + len(m.settingsProxy.inputs)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case ctrlC:
			return m, tea.Quit

		case ctrlB:
			errorMessage = ""
			m.screen = settingsMenu

		case enter:
			switch m.settingsProxy.index {
			case 4:
				m.settingsProxy.states["Authentication"] = !m.settingsProxy.states["Authentication"]

			case 5:
				m.settingsProxy.states["Enabled"] = !m.settingsProxy.states["Enabled"]

			case menuLength:
				if m.settingsProxy.index == menuLength {
					port, _ := strconv.Atoi(m.settingsProxy.inputs[1].Value())
					proxy := proxy.Proxy{
						RemoteIP:       m.settingsProxy.inputs[0].Value(),
						RemotePort:     port,
						RemoteUsername: m.settingsProxy.inputs[2].Value(),
						RemotePassword: m.settingsProxy.inputs[3].Value(),
						WithAuth:       m.settingsProxy.states["Authentication"],
						Enabled:        m.settingsProxy.states["Enabled"],
					}
					err := validate.Struct(proxy)
					if err != nil {
						errorMessage = invalidInputMsg
						break
					} else {
						config.Selenium.Proxy = proxy
						errorMessage = ""
						m.screen = settingsMenu
					}
				}
			}

		// Cycle between inputs
		case "tab", shiftTab, up, down:
			s := msg.String()

			// Cycle indexes
			if s == up || s == shiftTab {
				m.settingsProxy.index--
			} else {
				m.settingsProxy.index++
			}

			if m.settingsProxy.index > menuLength {
				m.settingsProxy.index = 0
			} else if m.settingsProxy.index < 0 {
				m.settingsProxy.index = menuLength
			}

			for i := 0; i < len(m.settingsProxy.inputs); i++ {
				if i == m.settingsProxy.index {
					// Set focused state
					m.settingsProxy.inputs[i].Focus()
					m.settingsProxy.inputs[i].Prompt = focusedPrompt
					m.settingsProxy.inputs[i].TextColor = focusedTextColor
					continue
				}
				// Remove focused state
				m.settingsProxy.inputs[i].Blur()
				m.settingsProxy.inputs[i].Prompt = blurredPrompt
				m.settingsProxy.inputs[i].TextColor = ""
			}

			if m.settingsProxy.index == menuLength {
				m.settingsProxy.submitButton = focusedSubmitButton
			} else {
				m.settingsProxy.submitButton = blurredSubmitButton
			}

			return m, nil
		}
	}

	// Handle character input and blinks
	m, cmd := updateInputsProxy(msg, m)
	return m, cmd
}

func updateInputsProxy(msg tea.Msg, m model) (model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	for i := 0; i < len(m.settingsProxy.inputs); i++ {
		m.settingsProxy.inputs[i], cmd = m.settingsProxy.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) ViewSettingsProxy() string {
	s := m.settingsProxy.title
	s += errorColor(errorMessage)
	for _, input := range m.settingsProxy.inputs {
		s += input.View()
		s += "\n\n"
	}

	if m.settingsProxy.index == 4 {
		s += fmt.Sprintf("%s : %s\n", focusColor("Authentication"), strconv.FormatBool(m.settingsProxy.states["Authentication"]))
	} else {
		s += fmt.Sprintf("%s : %s\n", "Authentication", strconv.FormatBool(m.settingsProxy.states["Authentication"]))
	}

	if m.settingsProxy.index == 5 {
		s += fmt.Sprintf("%s : %s\n", focusColor("Enabled"), strconv.FormatBool(m.settingsProxy.states["Enabled"]))
	} else {
		s += fmt.Sprintf("%s : %s\n", "Enabled", strconv.FormatBool(m.settingsProxy.states["Enabled"]))
	}

	s += "\n" + m.settingsProxy.submitButton + "\n"
	s += subtle("\nup/down: select") + dot + subtle("enter: choose/enable/disable") +
		dot + subtle("ctrl+b: back") + dot + subtle("ctrl+c: quit")
	return s
}
