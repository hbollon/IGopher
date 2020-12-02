package main

import (
	"flag"
	"os"

	"github.com/hbollon/go-instadm"
	log "github.com/sirupsen/logrus"
)

// SeleniumStruct is the main Selenium instance used by this bot
var SeleniumStruct instadm.Selenium

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
}{
	LogLevelFlag:           flag.String("loglevel", "info", "Log level threasold"),
	ForceDlFlag:            flag.Bool("force-download", false, "Force redownload of all dependencies even if exists"),
	DebugFlag:              flag.Bool("debug", false, "Display debug and selenium output"),
	IgnoreDependenciesFlag: flag.Bool("ignore-dependencies", false, "Skip dependencies management"),
	HeadlessFlag:           flag.Bool("headless", false, "Run WebDriver with frame buffer"),
}

func init() {
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

func initClientConfig() *instadm.ClientConfig {
	clientConfig := instadm.CreateClientConfig()
	clientConfig.LogLevel, _ = log.ParseLevel(*flags.LogLevelFlag)
	clientConfig.ForceDependenciesDl = *flags.ForceDlFlag
	clientConfig.Debug = *flags.DebugFlag
	clientConfig.IgnoreDependencies = *flags.IgnoreDependenciesFlag
	clientConfig.Headless = *flags.HeadlessFlag
	return clientConfig
}

func main() {
	// Initialize client configuration
	clientConfig := initClientConfig()

	// Download dependencies
	if !clientConfig.IgnoreDependencies {
		instadm.DownloadDependencies(true, true, clientConfig.ForceDependenciesDl)
	}

	// Initialize Selenium and WebDriver and defer their closing
	SeleniumStruct.InitializeSelenium(clientConfig)
	SeleniumStruct.InitChromeWebDriver()
	defer SeleniumStruct.CloseSelenium()

	SeleniumStruct.ConnectToInstagram()
	res, err := SeleniumStruct.SendMessage("_motivation.business", "Test message ! :)")
	if res == true && err == nil {
		log.Info("Message successfuly sent!")
	} else {
		log.Errorf("Error during message sending: %v", err)
	}
}
