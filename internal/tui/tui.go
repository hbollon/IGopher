package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-playground/validator/v10"
	"github.com/hbollon/igopher"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/termenv"
	log "github.com/sirupsen/logrus"
)

type screen uint16
type settingsScreen uint16

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	maxLineWidth      = 80
)

const (
	mainMenu screen = iota
	settingsMenu
	settingsResetMenu
	genericMenu
	settingsInputsScreen
	settingsBoolScreen
)

const (
	accountSettings settingsScreen = iota
	scrappingSettings
	autodmSettingsMenu
	autodmEnablingSettings
	autodmSettings
	autodmGreetingMenu
	autodmGreetingEnablingSettings
	autodmGreetingSettings
	quotasSettingsMenu
	quotasEnablingSettings
	quotasSettings
	scheduleSettingsMenu
	scheduleEnablingSettings
	scheduleSettings
	blacklistEnablingSettings
)

var (
	term          = termenv.ColorProfile()
	keyword       = makeFgStyle("211")
	cursorColor   = makeFgStyle("14")
	subtle        = makeFgStyle("241")
	errorColor    = makeFgStyle("1")
	progressEmpty = subtle(progressEmptyChar)
	dot           = colorFg(" • ", "236")

	color               = termenv.ColorProfile().Color
	focusedTextColor    = "205"
	focusedPrompt       = termenv.String("> ").Foreground(color("205")).String()
	blurredPrompt       = "> "
	focusedSubmitButton = "[ " + termenv.String("Submit").Foreground(color("205")).String() + " ]"
	blurredSubmitButton = "[ Submit ]"

	ramp = makeRamp("#B14FFF", "#00FFA3", progressBarWidth)

	execBot      = false
	errorMessage string
	config       igopher.BotConfigYaml

	validate = validator.New()
)

type model struct {
	screen                  screen
	settingsChoice          settingsScreen
	homeScreen              menu
	configScreen            menu
	configResetScreen       menu
	genericMenuScreen       menu
	settingsInputsScreen    inputs
	settingsTrueFalseScreen menu

	termWidth  int
	termHeight int
}

type menu struct {
	choices []string
	cursor  int
}

type mcList struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

type inputs struct {
	index        int
	title        string
	input        []textinput.Model
	submitButton string
}

var initialModel = model{
	screen:                  0,
	homeScreen:              menu{choices: []string{"🚀 - Launch!", "🔧 - Configure", "🧨 - Reset settings", "🚪 - Exit"}},
	configScreen:            menu{choices: []string{"Account", "Users scraping", "AutoDM", "Greeting", "Quotas", "Schedule", "Blacklist", "Save & exit"}},
	configResetScreen:       menu{choices: []string{"Yes", "No"}},
	settingsTrueFalseScreen: menu{choices: []string{"True", "False"}},
}

// InitTui initialize and start a terminal user interface instance
// Return the bot execution state on tui exit
func InitTui() bool {
	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}

	return execBot
}

func getAccountSettings() inputs {
	inp := inputs{
		title: fmt.Sprintf("\nPlease enter your %s credentials:\n\n", keyword("account")),
		input: []textinput.Model{
			textinput.NewModel(),
			textinput.NewModel(),
		}, submitButton: blurredSubmitButton}

	inp.input[0].Placeholder = "Username"
	inp.input[0].Focus()
	inp.input[0].Prompt = focusedPrompt
	inp.input[0].TextColor = focusedTextColor

	inp.input[1].Placeholder = "Password"
	inp.input[1].Prompt = blurredPrompt
	inp.input[1].EchoMode = textinput.EchoPassword
	inp.input[1].EchoCharacter = '•'

	return inp
}

func getUsersScrappingSettings() inputs {
	inp := inputs{
		title: fmt.Sprintf("\nPlease enter the list of %s you would like to use for %s (separated by a comma) :\n\n", keyword("accounts"), keyword("users scraping")),
		input: []textinput.Model{
			textinput.NewModel(),
			textinput.NewModel(),
		}, submitButton: blurredSubmitButton}

	inp.input[0].Placeholder = "Usernames (comma separated)"
	inp.input[0].Focus()
	inp.input[0].Prompt = focusedPrompt
	inp.input[0].TextColor = focusedTextColor

	inp.input[1].Placeholder = "Fetch quantity (default: 500)"
	inp.input[1].Prompt = blurredPrompt
	inp.input[1].TextColor = focusedTextColor

	return inp
}

func getQuotasSettings() inputs {
	inp := inputs{
		title: fmt.Sprintf("\nPlease fill following %s with desired values for %s module configuration.\n\n", keyword("inputs"), keyword("Quotas")),
		input: []textinput.Model{
			textinput.NewModel(),
			textinput.NewModel(),
		}, submitButton: blurredSubmitButton}

	inp.input[0].Placeholder = "Maximum dm quantity per day (default: 50)"
	inp.input[0].Focus()
	inp.input[0].Prompt = focusedPrompt
	inp.input[0].TextColor = focusedTextColor

	inp.input[1].Placeholder = "Maximum dm quantity per hour (default: 5)"
	inp.input[1].Prompt = blurredPrompt
	inp.input[1].TextColor = focusedTextColor

	return inp
}

func getSchedulerSettings() inputs {
	inp := inputs{
		title: fmt.Sprintf("\nPlease fill following %s with desired values for %s module configuration.\n\n", keyword("inputs"), keyword("Scheduler")),
		input: []textinput.Model{
			textinput.NewModel(),
			textinput.NewModel(),
		}, submitButton: blurredSubmitButton}

	inp.input[0].Placeholder = "Starting time (default: 08:00)"
	inp.input[0].Focus()
	inp.input[0].Prompt = focusedPrompt
	inp.input[0].TextColor = focusedTextColor

	inp.input[1].Placeholder = "Ending time (default: 18:00)"
	inp.input[1].Prompt = blurredPrompt
	inp.input[1].TextColor = focusedTextColor

	return inp
}

func getAutoDmGreetingSettings() inputs {
	inp := inputs{
		title: fmt.Sprintf("\nPlease fill following %s with desired greeting message template for %s sub-module configuration.\n\n", keyword("field"), keyword("Greeting")),
		input: []textinput.Model{
			textinput.NewModel(),
		}, submitButton: blurredSubmitButton}

	inp.input[0].Placeholder = "Greeting message (default: \"Hello\", will produce -> \"Hello <username>,\")"
	inp.input[0].Focus()
	inp.input[0].Prompt = focusedPrompt
	inp.input[0].TextColor = focusedTextColor

	return inp
}

func getAutoDmSettings() inputs {
	inp := inputs{
		title: fmt.Sprintf("\nPlease fill following %s with desired message templates (separated by ';') for %s configuration.\n\n", keyword("field"), keyword("AutoDm")),
		input: []textinput.Model{
			textinput.NewModel(),
		}, submitButton: blurredSubmitButton}

	inp.input[0].Placeholder = "Messages templates (separated by ';')"
	inp.input[0].Focus()
	inp.input[0].Prompt = focusedPrompt
	inp.input[0].TextColor = focusedTextColor

	return inp
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	// Get terminal windows dimensions from msg and update model
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
	}

	switch m.screen {
	case mainMenu:
		return m.UpdateHomePage(msg)

	case settingsMenu:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit

			case "ctrl+b":
				m.screen = mainMenu

			case "up", "k":
				if m.configScreen.cursor > 0 {
					m.configScreen.cursor--
				}

			case "down", "j":
				if m.configScreen.cursor < len(m.configScreen.choices)-1 {
					m.configScreen.cursor++
				}

			case "enter":
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
					igopher.ExportConfig(config)
					m.screen = mainMenu
				default:
					log.Warn("Invalid input!")
				}
			}
		}

	case settingsResetMenu:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit

			case "ctrl+b":
				m.screen = mainMenu

			case "up", "k":
				if m.configResetScreen.cursor > 0 {
					m.configResetScreen.cursor--
				}

			case "down", "j":
				if m.configResetScreen.cursor < len(m.configResetScreen.choices)-1 {
					m.configResetScreen.cursor++
				}

			case "enter":
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

	case genericMenu:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit

			case "ctrl+b":
				m.screen = settingsMenu

			case "up", "k":
				if m.genericMenuScreen.cursor > 0 {
					m.genericMenuScreen.cursor--
				}

			case "down", "j":
				if m.genericMenuScreen.cursor < len(m.genericMenuScreen.choices)-1 {
					m.genericMenuScreen.cursor++
				}

			case "enter":
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

	case settingsInputsScreen:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit

			case "ctrl+b":
				errorMessage = ""
				m.screen = settingsMenu

			case "enter":
				if m.settingsInputsScreen.index == len(m.settingsInputsScreen.input) {
					switch m.settingsChoice {
					case accountSettings:
						acc := igopher.AccountYaml{
							Username: m.settingsInputsScreen.input[0].Value(),
							Password: m.settingsInputsScreen.input[1].Value(),
						}
						err := validate.Struct(acc)
						if err != nil {
							errorMessage = "Invalid input, please check all fields.\n\n"
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
								Accounts: strings.Split(m.settingsInputsScreen.input[0].Value(), ","),
								Quantity: val,
							}
							err := validate.Struct(scr)
							if err != nil {
								errorMessage = "Invalid input, please check all fields.\n\n"
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
							errorMessage = "Invalid input, please check all fields.\n\n"
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
							errorMessage = "Invalid input, please check all fields.\n\n"
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
								errorMessage = "Invalid input, please check all fields.\n\n"
								break
							} else {
								config.Quotas.DmDay = quo.DmDay
								config.Quotas.DmHour = quo.DmHour
								errorMessage = ""
								m.screen = settingsMenu
							}
						} else {
							errorMessage = "Invalid input, please check all fields.\n\n"
						}
					case scheduleSettings:
						sche := igopher.ScheduleYaml{
							BeginAt: m.settingsInputsScreen.input[0].Value(),
							EndAt:   m.settingsInputsScreen.input[1].Value(),
						}
						err := validate.Struct(sche)
						if err != nil {
							errorMessage = "Invalid input, please check all fields.\n\n"
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
			case "tab", "shift+tab", "up", "down":
				s := msg.String()

				// Cycle indexes
				if s == "up" || s == "shift+tab" {
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
		m, cmd = updateInputs(msg, m)
		return m, cmd

	case settingsBoolScreen:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit

			case "ctrl+b":
				m.screen = settingsMenu

			case "up", "k":
				if m.settingsTrueFalseScreen.cursor > 0 {
					m.settingsTrueFalseScreen.cursor--
				}

			case "down", "j":
				if m.settingsTrueFalseScreen.cursor < len(m.settingsTrueFalseScreen.choices)-1 {
					m.settingsTrueFalseScreen.cursor++
				}

			case "enter":
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
	}

	return m, nil
}

func (m model) View() string {
	var s string
	switch m.screen {
	case mainMenu:
		s = fmt.Sprintf("\n🦄 Welcome to %s, the (soon) most powerful and versatile %s bot!\n\n", keyword("IGopher"), keyword("Instagram"))
		s += errorColor(errorMessage)

		for i, choice := range m.homeScreen.choices {
			cursor := " "
			if m.homeScreen.cursor == i {
				cursor = cursorColor(">")
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}

		s += subtle("\nup/down: select") + dot + subtle("enter: choose") + dot + subtle("ctrl+c: quit")

	case settingsMenu:
		s = fmt.Sprintf("\nWhat would you like to %s?\n\n", keyword("tweak"))

		for i, choice := range m.configScreen.choices {
			cursor := " "
			if m.configScreen.cursor == i {
				cursor = cursorColor(">")
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}

		s += subtle("\nup/down: select") + dot + subtle("enter: choose") + dot + subtle("ctrl+b: save & back") + dot + subtle("ctrl+c: quit")

	case settingsResetMenu:
		s = fmt.Sprintf("\nAre you sure you want to %s the default %s? This operation cannot be undone!\n\n", keyword("reset"), keyword("settings"))

		for i, choice := range m.configResetScreen.choices {
			cursor := " "
			if m.configResetScreen.cursor == i {
				cursor = cursorColor(">")
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}

		s += subtle("\nup/down: select") + dot + subtle("enter: choose") + dot + subtle("ctrl+b: back") + dot + subtle("ctrl+c: quit")

	case genericMenu:
		s += "\n\n"

		for i, choice := range m.genericMenuScreen.choices {
			cursor := " "
			if m.genericMenuScreen.cursor == i {
				cursor = cursorColor(">")
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}

		s += subtle("\nup/down: select") + dot + subtle("enter: choose") + dot + subtle("ctrl+b: back") + dot + subtle("ctrl+c: quit")

	case settingsInputsScreen:
		s = m.settingsInputsScreen.title
		s += errorColor(errorMessage)
		for i := 0; i < len(m.settingsInputsScreen.input); i++ {
			s += m.settingsInputsScreen.input[i].View()
			if i < len(m.settingsInputsScreen.input)-1 {
				s += "\n"
			}
		}
		s += "\n\n" + m.settingsInputsScreen.submitButton + "\n"
		s += subtle("\nup/down: select") + dot + subtle("enter: choose") + dot + subtle("ctrl+b: back") + dot + subtle("ctrl+c: quit")

	case settingsBoolScreen:
		switch m.settingsChoice {
		case autodmEnablingSettings:
			s = fmt.Sprintf("\nDo you want to enable %s module? (Default: %s)\n\n", keyword("AutoDM"), keyword("true"))

		case autodmGreetingEnablingSettings:
			s = fmt.Sprintf("\nDo you want to enable %s sub-module with %s? (Default: %s)\n\n", keyword("Greeting"), keyword("AutoDm"), keyword("true"))

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
	}

	return wordwrap.String(s, min(m.termWidth, maxLineWidth))
}

// Pass messages and models through to text input components. Only text inputs
// with Focus() set will respond, so it's safe to simply update all of them
// here without any further logic.
func updateInputs(msg tea.Msg, m model) (model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	for i := 0; i < len(m.settingsInputsScreen.input); i++ {
		m.settingsInputsScreen.input[i], cmd = m.settingsInputsScreen.input[i].Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// Utils

// Color a string's foreground with the given value.
func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

// Return a function that will colorize the foreground of a given string.
func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}

// Color a string's foreground and background with the given value.
func makeFgBgStyle(fg, bg string) func(string) string {
	return termenv.Style{}.
		Foreground(term.Color(fg)).
		Background(term.Color(bg)).
		Styled
}

// Generate a blend of colors.
func makeRamp(colorA, colorB string, steps float64) (s []string) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, colorToHex(c))
	}
	return
}

// Convert a colorful.Color to a hexadecimal format compatible with termenv.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and
// 1.
func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
