package main

import (
	"log"
	"os"

	"github.com/tebeka/selenium"
)

const (
	seleniumPath    = "lib/selenium-server.jar"
	geckoDriverPath = "lib/geckodriver"
	port            = 8080
)

// Selenium instance and opts
type Selenium struct {
	Instance selenium.Service
	Opts     []selenium.ServiceOption
}

// InitializeWebDriver start a Selenium WebDriver server instance
// (if one is not already running).
func (s *Selenium) InitializeWebDriver() {
	s.Opts = []selenium.ServiceOption{
		selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
		selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(os.Stderr),            // Output debug information to stderr.
	}

	selenium.SetDebug(true)
	service, err := selenium.NewSeleniumService(seleniumPath, port, s.Opts...)
	if err != nil {
		log.Fatal(err) // Fatal error, exit if webdriver can't be initialize.
	}
	defer service.Stop()
}
