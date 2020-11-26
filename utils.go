package main

import (
	"math/rand"
	"time"
)

const (
	sleepMin = 3.0
	sleepMax = 5.0
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func randomSleep() {
	time.Sleep(randomSecondsDuration(sleepMin, sleepMax))
}

func randomSecondsDuration(min, max float64) time.Duration {
	return time.Duration(min+rand.Float64()*(max-min)) * time.Second
}
