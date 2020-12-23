package main

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hbollon/igopher"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
	log "github.com/sirupsen/logrus"
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

	ramp = makeRamp("#B14FFF", "#00FFA3", progressBarWidth)

	execBot = false
)

type menu struct {
	choices []string
	cursor  int
}

type mcList struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

var homeScreen = menu{
	choices: []string{"🚀 - Launch!", "⚙️  - Configure", "🗒  - Edit raw config file", "🚪 - Exit"},
}

func (m menu) Init() tea.Cmd {
	return nil
}

func (m mcList) Init() tea.Cmd {
	return nil
}

func (m menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ":
			switch m.cursor {
			case 0:
				execBot = true
				return m, tea.Quit
			case 1:
				fmt.Println("1")
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

	return m, nil
}

func (m menu) View() string {
	s := fmt.Sprintf("\n🦄 Welcome to %s, the (soon) most powerful and versatile %s bot!\n\n", keyword("IGopher"), keyword("Instagram"))

	for i, choice := range m.choices {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = cursorColor(">") // cursor!
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += subtle("\nup/down: select") + dot + subtle("enter: choose") + dot + subtle("q: quit") // Send the UI for rendering

	return s
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
