package main

import (
	log "github.com/sirupsen/logrus"
)

// ConnectToInstagram get ig login webpage and connect user account
func (s *Selenium) ConnectToInstagram() {
	log.Info("Connecting to Instagram account...")
	if err := s.WebDriver.Get("https://instagram.com/?hl=en"); err != nil {
		log.Fatal(err)
	}
	randomSleep()
	if elem, err := s.GetElement("//button[text()='Accept']", "xpath"); err == nil {
		elem.Click()
		log.Debug("Cookies validation done!")
	} else {
		log.Info("Cookies validation button not found, skipping.")
	}
	randomSleep()
	if elem, err := s.GetElement("//button[text()='Log In']", "xpath"); err == nil {
		elem.Click()
		log.Debug("Log in screen access done!")
	} else {
		log.Info("Login button not found, skipping.")
	}
	randomSleep()
	if elem, err := s.GetElement("username", "name"); err == nil {
		elem.SendKeys("boursorama_parrainage__")
		log.Debug("Username injection done!")
	} else {
		log.Fatal(err)
	}
	if elem, err := s.GetElement("password", "name"); err == nil {
		elem.SendKeys("IAMHARDSTYLE")
		log.Debug("Password injection done!")
	} else {
		log.Fatal(err)
	}
}
