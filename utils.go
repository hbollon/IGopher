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
	time.Sleep(randomMillisecondDuration(sleepMin, sleepMax))
}

// Sleep random time between custom values
func randomSleepCustom(min, max float64) {
	time.Sleep(randomMillisecondDuration(min, max))
}

// Generate time duration between two limits
func randomMillisecondDuration(min, max float64) time.Duration {
	// Convert arguments (in seconds) to milliseconds
	min *= 1000
	max *= 1000
	return time.Duration(min+rand.Float64()*(max-min)) * time.Millisecond
}

// Closes all selenium stuff and call logrus fatal with error printing
func (s *Selenium) fatal(msg string, err error) {
	s.CloseSelenium()
	logrus.Fatal(msg, err)
}
