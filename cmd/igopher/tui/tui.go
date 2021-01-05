package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hbollon/igopher"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
)

var (
	term          = termenv.ColorProfile()
	keyword       = makeFgStyle("211")
	cursorColor   = makeFgStyle("14")
	subtle        = makeFgStyle("241")
	progressEmpty = subtle(progressEmptyChar)
	dot           = colorFg(" • ", "236")

	color               = termenv.ColorProfile().Color
	focusedTextColor    = "205"
	focusedPrompt       = termenv.String("> ").Foreground(color("205")).String()
	blurredPrompt       = "> "
	focusedSubmitButton = "[ " + termenv.String("Submit").Foreground(color("205")).String() + " ]"
	blurredSubmitButton = "[ " + termenv.String("Submit").Foreground(color("240")).String() + " ]"

	ramp = makeRamp("#B14FFF", "#00FFA3", progressBarWidth)

	execBot = false
)

type statusMsg int

type model struct {
	screen               int
	homeScreen           menu
	configScreen         menu
	settingsInputsScreen inputs
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
	screen:       0,
	homeScreen:   menu{choices: []string{"🚀 - Launch!", "⚙️  - Configure", "🗒  - Reset settings", "🚪 - Exit"}},
	configScreen: menu{choices: []string{"Account", "Users scrapping", "AutoDM", "Quotas", "Schedule", "Blacklist", "Save & exit"}},
}

func getAccountSettings() inputs {
	return inputs{
		title: fmt.Sprintf("\nPlease enter your %s credentials:\n\n", keyword("account")),
		input: []textinput.Model{
			textinput.Model{Placeholder: "Username", Prompt: focusedPrompt, TextColor: focusedTextColor},
			textinput.Model{Placeholder: "Password", Prompt: blurredPrompt, EchoMode: textinput.EchoPassword, EchoCharacter: '*'},
		}, submitButton: blurredSubmitButton}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.screen {
	case 0:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit

			case "up", "k":
				if m.homeScreen.cursor > 0 {
					m.homeScreen.cursor--
				}

			case "down", "j":
				if m.homeScreen.cursor < len(m.homeScreen.choices)-1 {
					m.homeScreen.cursor++
				}

			case "enter", " ":
				switch m.homeScreen.cursor {
				case 0:
					execBot = true
					return m, tea.Quit
				case 1:
					m.screen = 1
					break
				case 2:
					fmt.Println("2")
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
	case 1:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit

			case "up", "k":
				if m.configScreen.cursor > 0 {
					m.configScreen.cursor--
				}

			case "down", "j":
				if m.configScreen.cursor < len(m.configScreen.choices)-1 {
					m.configScreen.cursor++
				}

			case "enter", " ":
				switch m.configScreen.cursor {
				case 0:
					m.settingsInputsScreen = getAccountSettings()
					m.screen = 2
					break
				case 1:
					break
				case 2:
					break
				case 3:
					break
				case 4:
					break
				case 5:
					break
				case 6:
					m.screen = 0
					break
				default:
					log.Warn("Invalid input!")
					break
				}
			}
		}
		break
	case 2:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit

			case "backspace":
				m.screen--
				break

			// Cycle between inputs
			case "tab", "shift+tab", "enter", "up", "down":
				s := msg.String()

				if s == "enter" && m.settingsInputsScreen.index == len(m.settingsInputsScreen.input) {
					m.screen--
					break
				}

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

				for i := 0; i <= len(m.settingsInputsScreen.input)-1; i++ {
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
			}
		}
		break
	}

	return m, nil
}

func (m model) View() string {
	var s string
	switch m.screen {
	case 0:
		s = fmt.Sprintf("\n🦄 Welcome to %s, the (soon) most powerful and versatile %s bot!\n\n", keyword("IGopher"), keyword("Instagram"))

		for i, choice := range m.homeScreen.choices {
			cursor := " "
			if m.homeScreen.cursor == i {
				cursor = cursorColor(">")
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}

		s += subtle("\nup/down: select") + dot + subtle("enter: choose") + dot + subtle("q: quit")
		break

	case 1:
		s = fmt.Sprintf("\nWhat would you like to %s?\n\n", keyword("tweak"))

		for i, choice := range m.configScreen.choices {
			cursor := " "
			if m.configScreen.cursor == i {
				cursor = cursorColor(">")
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}

		s += subtle("\nup/down: select") + dot + subtle("enter: choose") + dot + subtle("q: quit")
		break

	case 2:
		s = m.settingsInputsScreen.title
		for i := 0; i < len(m.settingsInputsScreen.input); i++ {
			s += m.settingsInputsScreen.input[i].View()
			if i < len(m.settingsInputsScreen.input)-1 {
				s += "\n"
			}
		}
		s += "\n\n" + m.settingsInputsScreen.submitButton + "\n"
		s += subtle("\nup/down: select") + dot + subtle("enter: choose") + dot + subtle("q: quit")
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

	for _, i := range m.settingsInputsScreen.input {
		i, cmd = i.Update(msg)
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

type configYaml struct {
	account struct {
		username string `yaml:"username"`
		password string `yaml:"password"`
	} `yaml:"account"`
	srcUsers struct {
		accounts []string `yaml:"src_accounts"`
		quantity int      `yaml:"fetch_quantity"`
	} `yaml:"users_src"`
	autoDm struct {
		dmTemplates []string `yaml:"dm_templates"`
		greeting    struct {
			template  string `yaml:"template"`
			activated bool   `yaml:"activated"`
		} `yaml:"greeting"`
		activated bool `yaml:"activated"`
	} `yaml:"auto_dm"`
	quotas struct {
		dmDay     int  `yaml:"dm_per_day"`
		dmHour    int  `yaml:"dm_per_hour"`
		activated bool `yaml:"activated"`
	} `yaml:"quotas"`
	schedule struct {
		dmDay     int  `yaml:"begin_at"`
		dmHour    int  `yaml:"end_at"`
		activated bool `yaml:"activated"`
	} `yaml:"schedule"`
	blacklist struct {
		activated bool `yaml:"activated"`
	} `yaml:"blacklist"`
}

func importConfig() configYaml {
	var c configYaml
	file, err := ioutil.ReadFile("./config/config.yaml")
	if err != nil {
		log.Fatalf("Error opening config file: %s", err)
	}

	err = yaml.Unmarshal(file, &c)
	if err != nil {
		log.Fatalf("Error during unmarshal config file: %s\n", err)
	}

	return c
}

func exportConfig(c configYaml) {
	out, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatalf("Error during marshal config file: %s\n", err)
	}

	err = ioutil.WriteFile("./config/config.yaml", out, os.ModePerm)
	if err != nil {
		log.Fatalf("Error during config file writing: %s\n", err)
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
