package igopher

import (
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	sleepMin = 3.0
	sleepMax = 5.0
)

var clear map[string]func() // Map storing clear funcs for different os

func init() {
	// Initialize random engine
	rand.Seed(time.Now().UTC().UnixNano())

	// Prepare terminal cleaning functions for all os
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["darwin"] = clear["linux"]
}

// Sleep random time between default sleepMin and sleepMax
func randomSleep() {
	time.Sleep(randomMillisecondDuration(sleepMin, sleepMax))
}

// Sleep random time between custom values
func randomSleepCustom(min, max float64) {
	time.Sleep(randomMillisecondDuration(min, max))
}

// Generate time duration (in milliseconds) between two limits (in seconds)
func randomMillisecondDuration(min, max float64) time.Duration {
	// Convert arguments (in seconds) to milliseconds
	min *= 1000
	max *= 1000
	return time.Duration(min+rand.Float64()*(max-min)) * time.Millisecond
}

// Fatal closes all selenium stuff and call logrus fatal with error printing
func (s *Selenium) Fatal(msg string, err error) {
	s.CleanUp()
	logrus.Fatal(msg, err)
}

// ClearTerminal clear current terminal session according to user OS
func ClearTerminal() {
	value, ok := clear[runtime.GOOS] // runtime.GOOS -> linux, windows, darwin etc.
	if ok {
		value()
	} else {
		logrus.Errorf("Can't clear terminal, os unsupported !")
	}
}
