package main

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
)

// ConnectToInstagram get ig login webpage and connect user account
func (s *Selenium) ConnectToInstagram() {
	log.Info("Connecting to Instagram account...")
	// Access Instagram url
	if err := s.WebDriver.Get("https://instagram.com/?hl=en"); err != nil {
		log.Fatal(err)
	}
	randomSleep()
	// Accept cookies if requested
	if elem, err := s.GetElement("//button[text()='Accept']", "xpath"); err == nil {
		elem.Click()
		log.Debug("Cookies validation done!")
	} else {
		log.Info("Cookies validation button not found, skipping.")
	}
	randomSleep()
	// Access to login screen if needed
	if elem, err := s.GetElement("//button[text()='Log In']", "xpath"); err == nil {
		elem.Click()
		log.Debug("Log in screen access done!")
	} else {
		log.Info("Login button not found, skipping.")
	}
	randomSleep()
	// Inject username and password to input fields and log in
	if elem, err := s.GetElement("username", "name"); err == nil {
		elem.SendKeys("")
		log.Debug("Username injection done!")
	} else {
		log.Fatal(err)
	}
	if elem, err := s.GetElement("password", "name"); err == nil {
		elem.SendKeys("")
		log.Debug("Password injection done!")
	} else {
		log.Fatal(err)
	}
	if elem, err := s.GetElement("//button/*[text()='Log In']", "xpath"); err == nil {
		elem.Click()
		log.Debug("Sent login request")
	} else {
		log.Fatal(err)
	}
	randomSleep()
	time.Sleep(10 * time.Second)
	// Check if login was successful
	if s.IsElementPresent(selenium.ByXPATH, "//*[@aria-label='Home'] | //button[text()='Save Info'] | //button[text()='Not Now']") {
		log.Info("Login Successful!")
	} else {
		log.Fatal("Login Failed: Incorrect credentials")
	}
}
