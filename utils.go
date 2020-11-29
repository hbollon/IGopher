package main

import (
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	sleepMin = 3.0
	sleepMax = 5.0
)

// Initialize random engine
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Sleep random time between default sleepMin and sleepMax
func randomSleep() {
	time.Sleep(randomSecondsDuration(sleepMin, sleepMax))
}

// Sleep random time between custom values
func randomSleepCustom(min, max float64) {
	time.Sleep(randomSecondsDuration(min, max))
}

// Generate time duration between two limits
func randomSecondsDuration(min, max float64) time.Duration {
	return time.Duration(min+rand.Float64()*(max-min)) * time.Second
}

// Closes all selenium stuff and call logrus fatal with error printing
func (s *Selenium) fatal(msg string, err error) {
	s.CloseSelenium()
	logrus.Fatal(msg, err)
}
