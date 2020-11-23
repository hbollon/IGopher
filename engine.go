package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
)

const (
	seleniumPath    = "./lib/selenium-server.jar"
	geckoDriverPath = "./lib/geckodriver"
	port            = 8080
)

var err error

// Selenium instance and opts
type Selenium struct {
	Instance  *selenium.Service
	Opts      []selenium.ServiceOption
	WebDriver selenium.WebDriver
}

// InitializeSelenium start a Selenium WebDriver server instance
// (if one is not already running).
func (s *Selenium) InitializeSelenium() {
	s.Opts = []selenium.ServiceOption{
		selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
		selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(os.Stderr),            // Output debug information to stderr.
	}

	selenium.SetDebug(true)
	s.Instance, err = selenium.NewSeleniumService(seleniumPath, port, s.Opts...)
	if err != nil {
		log.Fatal(err) // Fatal error, exit if webdriver can't be initialize.
	}
}

// InitFirefoxWebDriver init and launch web driver with Firefox
func (s *Selenium) InitFirefoxWebDriver() {
	caps := selenium.Capabilities{"browserName": "firefox"}
	s.WebDriver, err = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		log.Error(err)
	}
}

// CloseSelenium close webdriver and selenium instance
func (s *Selenium) CloseSelenium() {
	s.Instance.Stop()
	if s.WebDriver != nil {
		s.WebDriver.Quit()
	}
}
