package tui

import (
	"fmt"
	"strconv"

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

// Screen identifiers
const (
	mainMenu screen = iota
	settingsMenu
	settingsResetMenu
	genericMenu
	settingsInputsScreen
	settingsBoolScreen
	settingsProxyScreen
	stopRunningInstance
)

// Settings screen identifiers
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

// Keyboard key
const (
	ctrlC    = "ctrl+c"
	ctrlB    = "ctrl+b"
	enter    = "enter"
	up       = "up"
	down     = "down"
	shiftTab = "shift+tab"
)

var (
	term          = termenv.ColorProfile()
	keyword       = makeFgStyle("211")
	cursorColor   = makeFgStyle("14")
	subtle        = makeFgStyle("241")
	errorColor    = makeFgStyle("1")
	infoColor     = makeFgStyle("32")
	focusColor    = makeFgStyle("205")
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
	infoMessage  string
	errorMessage string
	config       igopher.BotConfigYaml

	validate = validator.New()
)

type model struct {
	instanceAlreadyRunning   bool
	screen                   screen
	settingsChoice           settingsScreen
	homeScreen               menu
	configScreen             menu
	configResetScreen        menu
	genericMenuScreen        menu
	stopRunningProcessScreen menu
	settingsInputsScreen     inputs
	settingsTrueFalseScreen  menu
	settingsProxy            proxyMenu

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

type proxyMenu struct {
	index        int
	title        string
	inputs       []textinput.Model
	states       map[string]bool
	submitButton string
}

var initialModel = model{
	screen:     0,
	homeScreen: menu{choices: []string{"🚀 - Launch!", "🔧 - Configure", "🧨 - Reset settings", "🚪 - Exit"}},
	configScreen: menu{choices: []string{"Account", "Users scraping", "AutoDM",
		"Greeting", "Quotas", "Schedule", "Blacklist", "Proxy", "Save & exit"}},
	configResetScreen:        menu{choices: []string{"Yes", "No"}},
	stopRunningProcessScreen: menu{choices: []string{"Yes", "No"}},
	settingsTrueFalseScreen:  menu{choices: []string{"True", "False"}},
	settingsProxy:            getProxySettings(),
}

// InitTui initialize and start a terminal user interface instance
// It take
// Return the bot execution state on tui exit
func InitTui(instanceRunning bool) bool {
	initialModel.instanceAlreadyRunning = instanceRunning
	initialModel.updateMenuItemsHomePage()
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
		title: fmt.Sprintf("\nPlease enter the list of %s you would like to use for %s (separated by ';') :\n\n",
			keyword("accounts"), keyword("users scraping")),
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
		title: fmt.Sprintf("\nPlease fill following %s with desired values for %s module configuration.\n\n",
			keyword("inputs"), keyword("Quotas")),
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
		title: fmt.Sprintf("\nPlease fill following %s with desired values for %s module configuration.\n\n",
			keyword("inputs"), keyword("Scheduler")),
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
		title: fmt.Sprintf("\nPlease fill following %s with desired greeting message template for %s sub-module configuration.\n\n",
			keyword("field"), keyword("Greeting")),
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
		title: fmt.Sprintf("\nPlease fill following %s with desired message templates (separated by ';') for %s configuration.\n\n",
			keyword("field"), keyword("AutoDm")),
		input: []textinput.Model{
			textinput.NewModel(),
		}, submitButton: blurredSubmitButton}

	inp.input[0].Placeholder = "Messages templates (separated by ';')"
	inp.input[0].Focus()
	inp.input[0].Prompt = focusedPrompt
	inp.input[0].TextColor = focusedTextColor

	return inp
}

func getProxySettings() proxyMenu {
	inp := proxyMenu{
		title: fmt.Sprintf("\nPlease enter your %s configuration."+
			" If you use a proxy with %s,"+
			" do not forget to activate the corresponding option below and add your connection %s.\n\n",
			keyword("proxy"), keyword("authentication"), keyword("credentials")),
		inputs: []textinput.Model{
			textinput.NewModel(),
			textinput.NewModel(),
			textinput.NewModel(),
			textinput.NewModel(),
		},
		states: map[string]bool{
			"Authentication": false,
			"Enabled":        false,
		},
		submitButton: blurredSubmitButton,
	}

	inp.inputs[0].Placeholder = "Remote proxy IP"
	inp.inputs[0].Focus()
	inp.inputs[0].Prompt = focusedPrompt
	inp.inputs[0].TextColor = focusedTextColor

	inp.inputs[1].Placeholder = "Remote proxy port"
	inp.inputs[1].Prompt = blurredPrompt
	inp.inputs[1].TextColor = focusedTextColor

	inp.inputs[2].Placeholder = "Remote proxy username (optional)"
	inp.inputs[2].Prompt = blurredPrompt
	inp.inputs[2].TextColor = focusedTextColor

	inp.inputs[3].Placeholder = "Remote proxy password (optional)"
	inp.inputs[3].Prompt = blurredPrompt
	inp.inputs[3].TextColor = focusedTextColor

	return inp
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		return m.UpdateSettingsMenu(msg)

	case settingsResetMenu:
		return m.UpdateSettingsResetMenu(msg)

	case genericMenu:
		return m.UpdateGenericMenu(msg)

	case settingsInputsScreen:
		return m.UpdateSettingsInputsMenu(msg)

	case settingsBoolScreen:
		return m.UpdateSettingsBoolMenu(msg)

	case settingsProxyScreen:
		return m.UpdateSettingsProxy(msg)

	case stopRunningInstance:
		return m.UpdateStopRunningProcess(msg)
	}

	return m, nil
}

func (m model) View() string {
	var s string
	switch m.screen {
	case mainMenu:
		s = m.ViewHomePage()

	case settingsMenu:
		s = m.ViewSettingsMenu()

	case settingsResetMenu:
		s = m.ViewSettingsResetMenu()

	case genericMenu:
		s = m.ViewGenericMenu()

	case settingsInputsScreen:
		s = m.ViewSettingsInputsMenu()

	case settingsBoolScreen:
		s = m.ViewSettingsBoolMenu()

	case settingsProxyScreen:
		s = m.ViewSettingsProxy()

	case stopRunningInstance:
		s = m.ViewStopRunningProcess()
	}

	return wordwrap.String(s, min(m.termWidth, maxLineWidth))
}

// Pass messages and models through to text input components. Only text inputs
// with Focus() set will respond, so it's safe to simply update all of them
// here without any further logic.
// Used to update generics input screens
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

func launchBot(m model) (model, tea.Cmd) {
	err := igopher.CheckConfigValidity()
	if err == nil {
		execBot = true
		return m, tea.Quit
	}
	errorMessage = err.Error() + "\n\n"
	return m, nil
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
