package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
)

// ConnectToInstagram get ig login webpage and connect user account
func (s *Selenium) ConnectToInstagram() {
	log.Info("Connecting to Instagram account...")
	// Access Instagram url
	if err := s.WebDriver.Get("https://instagram.com/?hl=en"); err != nil {
		s.fatal("Can't access to Instagram. ", err)
	}
	randomSleep()
	// Accept cookies if requested
	if find, err := s.WaitForElement("//button[text()='Accept']", "xpath", 10); err == nil && find {
		elem, _ := s.GetElement("//button[text()='Accept']", "xpath")
		elem.Click()
		log.Debug("Cookies validation done!")
	} else {
		log.Info("Cookies validation button not found, skipping.")
	}
	randomSleep()
	// Access to login screen if needed
	if find, err := s.WaitForElement("//button[text()='Log In']", "xpath", 10); err == nil && find {
		elem, _ := s.GetElement("//button[text()='Log In']", "xpath")
		elem.Click()
		log.Debug("Log in screen access done!")
	} else {
		log.Info("Login button not found, skipping.")
	}
	randomSleep()
	// Inject username and password to input fields and log in
	if find, err := s.WaitForElement("username", "name", 10); err == nil && find {
		elem, _ := s.GetElement("username", "name")
		elem.SendKeys("boursorama_parrainage__")
		log.Debug("Username injection done!")
	} else {
		s.fatal("Exception during username inject: ", err)
	}
	if find, err := s.WaitForElement("password", "name", 10); err == nil && find {
		elem, _ := s.GetElement("password", "name")
		elem.SendKeys("IAMHARDSTYLE74")
		log.Debug("Password injection done!")
	} else {
		s.fatal("Exception during password inject: ", err)
	}
	if find, err := s.WaitForElement("//button/*[text()='Log In']", "xpath", 10); err == nil && find {
		elem, _ := s.GetElement("//button/*[text()='Log In']", "xpath")
		elem.Click()
		log.Debug("Sent login request")
	} else {
		s.fatal("Log in button not found: ", err)
	}
	randomSleepCustom(10, 15)
	// Check if login was successful
	if s.IsElementPresent(selenium.ByXPATH, "//*[@aria-label='Home'] | //button[text()='Save Info'] | //button[text()='Not Now']") {
		log.Info("Login Successful!")
	} else {
		s.fatal("Login failed! ", err)
	}
}
