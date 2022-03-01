package utils

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

// RandomSleep sleep random time between default sleepMin and sleepMax
func RandomSleep() {
	time.Sleep(RandomMillisecondDuration(sleepMin, sleepMax))
}

// RandomSleepCustom sleep random time between custom values
func RandomSleepCustom(min, max float64) {
	time.Sleep(RandomMillisecondDuration(min, max))
}

// RandomMillisecondDuration generate time duration (in milliseconds) between two limits (in seconds)
func RandomMillisecondDuration(min, max float64) time.Duration {
	// Convert arguments (in seconds) to milliseconds
	min *= 1000
	max *= 1000
	return time.Duration(min+rand.Float64()*(max-min)) * time.Millisecond
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
