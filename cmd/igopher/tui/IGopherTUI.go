package main

import (
	"flag"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"

	logRuntime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/hbollon/igopher"
	tui "github.com/hbollon/igopher/internal/tui"
	"github.com/shiena/ansicolor"
	log "github.com/sirupsen/logrus"
)

// BotStruct is the main struct instance used by this bot
var BotStruct igopher.IGopher

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
	LogLevelFlag:           flag.String("loglevel", "info", "Log level threshold"),
	ForceDlFlag:            flag.Bool("force-download", false, "Force redownload of all dependencies even if exists"),
	DebugFlag:              flag.Bool("debug", false, "Display debug and selenium output"),
	IgnoreDependenciesFlag: flag.Bool("ignore-dependencies", false, "Skip dependencies management"),
	HeadlessFlag:           flag.Bool("headless", false, "Run WebDriver with frame buffer"),
	PortFlag:               flag.Int("port", 8080, "Specify custom communication port"),
}

func init() {
	// Add formatter to logrus in order to display line and function with messages
	formatter := logRuntime.Formatter{ChildFormatter: &log.TextFormatter{
		FullTimestamp: false,
		ForceColors:   true,
	}}
	formatter.Line = true
	log.SetFormatter(&formatter)

	// Output to stderr
	if runtime.GOOS == "windows" {
		log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stderr))
	} else {
		log.SetOutput(os.Stderr)
	}

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
	// Initialize environment
	igopher.CheckEnvironment()

	// Clear terminal session
	igopher.ClearTerminal()

	// Launch TUI
	execBot := tui.InitTui()

	// Lauch bot if option selected
	if execBot {
		launchBot()
	}
}

func launchBot() {
	// Initialize client configuration
	clientConfig := initClientConfig()
	BotStruct = igopher.ReadBotConfigYaml()

	// Download dependencies
	if !clientConfig.IgnoreDependencies {
		igopher.DownloadDependencies(true, true, clientConfig.ForceDependenciesDl)
	}

	// Initialize Selenium and WebDriver and defer their closing
	BotStruct.SeleniumStruct.InitializeSelenium(clientConfig)
	BotStruct.SeleniumStruct.InitChromeWebDriver()
	defer BotStruct.SeleniumStruct.CloseSelenium()

	rand.Seed(time.Now().Unix())
	if err := BotStruct.Scheduler.CheckTime(); err == nil {
		BotStruct.ConnectToInstagram()
		users, err := BotStruct.FetchUsersFromUserFollowers()
		if err != nil {
			log.Error(err)
		}
		for _, username := range users {
			res, err := BotStruct.SendMessage(username, BotStruct.DmModule.DmTemplates[rand.Intn(len(BotStruct.DmModule.DmTemplates))])
			if !res || err != nil {
				log.Errorf("Error during message sending: %v", err)
			}
		}
	} else {
		BotStruct.SeleniumStruct.Fatal("Error on bot launch: ", err)
	}
}
