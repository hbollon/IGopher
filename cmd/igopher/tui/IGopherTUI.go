package main

import (
	"flag"
	"math"
	"os"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hbollon/igopher"
	log "github.com/sirupsen/logrus"
)

// SeleniumStruct is the main Selenium instance used by this bot
var SeleniumStruct igopher.Selenium

/// Flags
var flags = struct {
	// LogLevelFlag set loglevel threshold
	// If undefined or wrong set it to INFO level
	LogLevelFlag *string

	// ForceDlFlag force re-download of all dependencies
	ForceDlFlag *bool

	// DebugFlag set selenium debug mode and display its logging to stderr
	DebugFlag *bool

	// IgnoreDependenciesFlag disable dependencies manager on startup
	IgnoreDependenciesFlag *bool

	// HeadlessFlag execute Selenium webdriver in headless mode
	HeadlessFlag *bool

	// PortFlag specifie custom communication port for Selenium and web drivers
	PortFlag *int
}{
	LogLevelFlag:           flag.String("loglevel", "info", "Log level threasold"),
	ForceDlFlag:            flag.Bool("force-download", false, "Force redownload of all dependencies even if exists"),
	DebugFlag:              flag.Bool("debug", false, "Display debug and selenium output"),
	IgnoreDependenciesFlag: flag.Bool("ignore-dependencies", false, "Skip dependencies management"),
	HeadlessFlag:           flag.Bool("headless", false, "Run WebDriver with frame buffer"),
	PortFlag:               flag.Int("port", 8080, "Specify custom communication port"),
}

func init() {
	// Add formatter to logrus in order to display line and function with messages
	formatter := runtime.Formatter{ChildFormatter: &log.TextFormatter{
		FullTimestamp: true,
	}}
	formatter.Line = true
	log.SetFormatter(&formatter)

	// Output to stderr
	log.SetOutput(os.Stderr)

	flag.Parse()
	level, err := log.ParseLevel(*flags.LogLevelFlag)
	if err == nil {
		log.SetLevel(level)
	} else {
		log.SetLevel(log.InfoLevel)
		log.Warnf("Invalid log level '%s', use default one.", *flags.LogLevelFlag)
	}
}

func initClientConfig() *igopher.ClientConfig {
	clientConfig := igopher.CreateClientConfig()
	clientConfig.LogLevel, _ = log.ParseLevel(*flags.LogLevelFlag)
	clientConfig.ForceDependenciesDl = *flags.ForceDlFlag
	clientConfig.Debug = *flags.DebugFlag
	clientConfig.IgnoreDependencies = *flags.IgnoreDependenciesFlag
	clientConfig.Headless = *flags.HeadlessFlag

	if *flags.PortFlag > math.MaxUint16 || *flags.PortFlag < 8080 {
		log.Warnf("Invalid port argument '%d'. Use default 8080.", *flags.PortFlag)
	} else {
		clientConfig.Port = uint16(*flags.PortFlag)
	}

	return clientConfig
}

func main() {
	// Launch TUI
	p := tea.NewProgram(homeScreen)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}

	// Lauch bot if option selected
	if execBot {
		launchBot()
	}
}
