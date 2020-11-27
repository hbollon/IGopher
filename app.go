package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

// SeleniumStruct is the main Selenium instance used by this bot
var SeleniumStruct Selenium

/// Flags
// Get loglevel flag and set thresold to it
// If undefined or wrong set it to WARNING
var loglevel = flag.String("loglevel", "info", "Log level threasold")
var forceDl = flag.Bool("force-download", false, "Force redownload of all dependencies even if exists")
var debug = flag.Bool("debug", false, "Display debug and selenium output")
var ignoreDependencies = flag.Bool("ignore-dependencies", false, "Skip dependencies management")

func init() {
	// Output to stderr
	log.SetOutput(os.Stderr)

	flag.Parse()
	level, err := log.ParseLevel(*loglevel)
	if err == nil {
		log.SetLevel(level)
	} else {
		log.SetLevel(log.WarnLevel)
	}
}

func main() {
	// Download dependencies
	if !*ignoreDependencies {
		DownloadDependencies(true, true)
	}

	// Initialize Selenium and WebDriver and defer their closing
	SeleniumStruct.InitializeSelenium()
	SeleniumStruct.InitFirefoxWebDriver()
	defer SeleniumStruct.CloseSelenium()

	SeleniumStruct.ConnectToInstagram()
}
