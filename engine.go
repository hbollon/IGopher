package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
)

const (
	seleniumPath    = "./lib/selenium-server.jar"
	geckoDriverPath = "./lib/geckodriver"
	port            = 8081
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
	var output *os.File
	if *debug {
		output = os.Stderr
	} else {
		output = nil
	}

	s.Opts = []selenium.ServiceOption{
		selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
		selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(output),               // Output debug information to stderr.
	}

	selenium.SetDebug(*debug)
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

/* Browser methods */

// IsElementPresent check if an element is present on the current webpage
func (s *Selenium) IsElementPresent(by, value string) bool {
	_, err := s.WebDriver.FindElement(by, value)
	if err != nil {
		log.Infof("Element not found by %s: %v", by, err)
		return false
	}
	return true
}

// GetElement wait for element and then return when it is available
func (s *Selenium) GetElement(elementTag, locator string) (selenium.WebElement, error) {
	locator = strings.ToUpper(locator)
	if locator == "ID" && s.IsElementPresent(selenium.ByID, elementTag) {
		return s.WebDriver.FindElement(selenium.ByID, elementTag)
	} else if locator == "NAME" && s.IsElementPresent(selenium.ByName, elementTag) {
		return s.WebDriver.FindElement(selenium.ByName, elementTag)
	} else if locator == "XPATH" && s.IsElementPresent(selenium.ByXPATH, elementTag) {
		return s.WebDriver.FindElement(selenium.ByXPATH, elementTag)
	} else if locator == "CSS" && s.IsElementPresent(selenium.ByCSSSelector, elementTag) {
		return s.WebDriver.FindElement(selenium.ByCSSSelector, elementTag)
	} else {
		log.Debugf("Incorrect locator '%s'", locator)
		return nil, errors.New("Incorrect locator")
	}
}
