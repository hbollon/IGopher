package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-playground/validator/v10"
	"github.com/hbollon/igopher"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
	log "github.com/sirupsen/logrus"
)

type screen uint16
type settingsScreen uint16

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
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

	execBot        = false
	errorMessage   string
	settingsChoice settingsScreen = 0
	config         igopher.BotConfigYaml

	validate = validator.New()
)

type model struct {
	screen                  screen
	homeScreen              menu
	configScreen            menu
	configResetScreen       menu
	genericMenuScreen       menu
	settingsInputsScreen    inputs
	settingsTrueFalseScreen menu
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
	homeScreen:              menu{choices: []string{"🚀 - Launch!", "⚙️  - Configure", "🗒  - Reset settings", "🚪 - Exit"}},
	configScreen:            menu{choices: []string{"Account", "Users scraping", "AutoDM", "Greeting", "Quotas", "Schedule", "Blacklist", "Save & exit"}},
	configResetScreen:       menu{choices: []string{"Yes", "No"}},
	settingsTrueFalseScreen: menu{choices: []string{"True", "False"}},
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
	switch m.screen {
	case mainMenu:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit

			case "up", "k":
				if m.homeScreen.cursor > 0 {
					m.homeScreen.cursor--
				}

			case "down", "j":
				if m.homeScreen.cursor < len(m.homeScreen.choices)-1 {
					m.homeScreen.cursor++
				}

			case "enter":
				switch m.homeScreen.cursor {
				case 0:
					execBot = true
					return m, tea.Quit
				case 1:
					config = igopher.ImportConfig()
					m.screen = settingsMenu
					break
				case 2:
					m.screen = settingsResetMenu
					break
				case 3:
					return m, tea.Quit
				default:
					log.Warn("Invalid input!")
					break
				}
			}
		}
		break

	case settingsMenu:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit

			case "ctrl+b":
				m.screen = mainMenu
				break

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
					settingsChoice = accountSettings
					break
				case 1:
					m.settingsInputsScreen = getUsersScrappingSettings()
					m.screen = settingsInputsScreen
					settingsChoice = scrappingSettings
					break
				case 2:
					m.genericMenuScreen = menu{choices: []string{"Enable/Disable Module", "Configuration"}}
					m.screen = genericMenu
					settingsChoice = autodmSettingsMenu
					break
				case 3:
					m.genericMenuScreen = menu{choices: []string{"Enable/Disable Module", "Configuration"}}
					m.screen = genericMenu
					settingsChoice = autodmGreetingMenu
					break
				case 4:
					m.genericMenuScreen = menu{choices: []string{"Enable/Disable Module", "Configuration"}}
					m.screen = genericMenu
					settingsChoice = quotasSettingsMenu
					break
				case 5:
					m.genericMenuScreen = menu{choices: []string{"Enable/Disable Module", "Configuration"}}
					m.screen = genericMenu
					settingsChoice = scheduleSettingsMenu
					break
				case 6:
					m.screen = settingsBoolScreen
					settingsChoice = blacklistEnablingSettings
					break
				case 7:
					igopher.ExportConfig(config)
					m.screen = mainMenu
					break
				default:
					log.Warn("Invalid input!")
					break
				}
			}
		}
		break

	case settingsResetMenu:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit

			case "ctrl+b":
				m.screen = mainMenu
				break

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
					break
				case 1:
					m.screen = mainMenu
					break
				default:
					log.Warn("Invalid input!")
					break
				}
			}
		}
		break

	case genericMenu:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit

			case "ctrl+b":
				m.screen = settingsMenu
				break

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
					switch settingsChoice {
					case autodmSettingsMenu:
						settingsChoice = autodmEnablingSettings
						break
					case autodmGreetingMenu:
						settingsChoice = autodmGreetingEnablingSettings
						break
					case quotasSettingsMenu:
						settingsChoice = quotasEnablingSettings
						break
					case scheduleSettingsMenu:
						settingsChoice = scheduleEnablingSettings
						break
					default:
						log.Warn("Invalid input!")
						break
					}
					m.screen = settingsBoolScreen
					break
				case 1:
					switch settingsChoice {
					case autodmSettingsMenu:
						m.settingsInputsScreen = getAutoDmSettings()
						settingsChoice = autodmSettings
						break
					case autodmGreetingMenu:
						m.settingsInputsScreen = getAutoDmGreetingSettings()
						settingsChoice = autodmGreetingSettings
						break
					case quotasSettingsMenu:
						m.settingsInputsScreen = getQuotasSettings()
						settingsChoice = quotasSettings
						break
					case scheduleSettingsMenu:
						m.settingsInputsScreen = getSchedulerSettings()
						settingsChoice = scheduleSettings
						break
					default:
						log.Warn("Invalid input!")
						break
					}
					m.screen = settingsInputsScreen
					break
				default:
					log.Warn("Invalid input!")
					break
				}
			}
		}
		break

	case settingsInputsScreen:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit

			case "ctrl+b":
				m.screen = settingsMenu
				break

			case "enter":
				if m.settingsInputsScreen.index == len(m.settingsInputsScreen.input) {
					switch settingsChoice {
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
						break
					case scrappingSettings:
						val, err := strconv.Atoi(m.settingsInputsScreen.input[1].Value())
						if err == nil {
							scr := igopher.SrcUsersYaml{
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
						break
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
						break
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
						break
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
						break
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
						break
					default:
						log.Error("Unexpected settings screen value!\n\n")
						break
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
				break

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
					switch settingsChoice {
					case autodmEnablingSettings:
						config.AutoDm.Activated = true
						break

					case autodmGreetingEnablingSettings:
						config.AutoDm.Greeting.Activated = true
						break

					case quotasEnablingSettings:
						config.Quotas.Activated = true
						break

					case scheduleEnablingSettings:
						config.Schedule.Activated = true
						break

					case blacklistEnablingSettings:
						config.Blacklist.Activated = true
						break

					default:
						log.Error("Unexpected settings screen value!")
						break
					}
					m.screen = settingsMenu
					break
				case 1:
					switch settingsChoice {
					case autodmEnablingSettings:
						config.AutoDm.Activated = false
						break

					case autodmGreetingEnablingSettings:
						config.AutoDm.Greeting.Activated = false
						break

					case quotasEnablingSettings:
						config.Quotas.Activated = false
						break

					case scheduleEnablingSettings:
						config.Schedule.Activated = false
						break

					case blacklistEnablingSettings:
						config.Blacklist.Activated = false
						break

					default:
						log.Error("Unexpected settings screen value!")
						break
					}
					m.screen = settingsMenu
					break
				default:
					log.Warn("Invalid input!")
					m.screen = settingsMenu
					break
				}
			}
		}
		break
	}

	return m, nil
}

func (m model) View() string {
	var s string
	switch m.screen {
	case mainMenu:
		s = fmt.Sprintf("\n🦄 Welcome to %s, the (soon) most powerful and versatile %s bot!\n\n", keyword("IGopher"), keyword("Instagram"))

		for i, choice := range m.homeScreen.choices {
			cursor := " "
			if m.homeScreen.cursor == i {
				cursor = cursorColor(">")
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}

		s += subtle("\nup/down: select") + dot + subtle("enter: choose") + dot + subtle("ctrl+c: quit")
		break

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
		break

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
		break

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
		break

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
		break

	case settingsBoolScreen:
		switch settingsChoice {
		case autodmEnablingSettings:
			s = fmt.Sprintf("\nDo you want to enable %s module? (Default: %s)\n\n", keyword("AutoDM"), keyword("true"))
			break

		case autodmGreetingEnablingSettings:
			s = fmt.Sprintf("\nDo you want to enable %s sub-module with %s? (Default: %s)\n\n", keyword("Greeting"), keyword("AutoDm"), keyword("true"))
			break

		case quotasEnablingSettings:
			s = fmt.Sprintf("\nDo you want to enable %s module? (Default: %s)\n\n", keyword("Quotas"), keyword("true"))
			break

		case scheduleEnablingSettings:
			s = fmt.Sprintf("\nDo you want to enable %s module? (Default: %s)\n\n", keyword("Scheduler"), keyword("true"))
			break

		case blacklistEnablingSettings:
			s = fmt.Sprintf("\nDo you want to enable %s module? (Default: %s)\n\n", keyword("User Blacklist"), keyword("true"))
			break

		default:
			log.Error("Unexpected settings screen value!")
			s = ""
			break
		}

		for i, choice := range m.settingsTrueFalseScreen.choices {
			cursor := " "
			if m.settingsTrueFalseScreen.cursor == i {
				cursor = cursorColor(">")
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}
		s += subtle("\nup/down: select") + dot + subtle("enter: choose") + dot + subtle("ctrl+b: back") + dot + subtle("ctrl+c: quit")
		break
	}

	return s
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

// Actions

func launchBot() {
	// Initialize client configuration
	clientConfig := initClientConfig()

	// Download dependencies
	if !clientConfig.IgnoreDependencies {
		igopher.DownloadDependencies(true, true, clientConfig.ForceDependenciesDl)
	}

	// Initialize Selenium and WebDriver and defer their closing
	SeleniumStruct.InitializeSelenium(clientConfig)
	SeleniumStruct.InitChromeWebDriver()
	defer SeleniumStruct.CloseSelenium()

	if err := SeleniumStruct.Config.BotConfig.Scheduler.CheckTime(); err == nil {
		SeleniumStruct.ConnectToInstagram()
		res, err := SeleniumStruct.SendMessage("_motivation.business", "Test message ! :)")
		if res != true || err != nil {
			log.Errorf("Error during message sending: %v", err)
		}
	} else {
		SeleniumStruct.Fatal("Error on bot launch: ", err)
	}
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

// Convert a colorful.Color to a hexidecimal format compatible with termenv.
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
