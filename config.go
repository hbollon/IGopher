package instadm

import (
	"github.com/sirupsen/logrus"
)

// ClientConfig struct centralize all client configuration and flags.
// Inizialized at program startup, not safe to modify this after.
type ClientConfig struct {
	// LogLevel set loglevel threshold
	// If undefined or wrong set it to INFO level
	LogLevel logrus.Level

	// ForceDependenciesDl force re-download of all dependencies
	ForceDependenciesDl bool

	// Debug set selenium debug mode and display its logging to stderr
	Debug bool

	// IgnoreDependencies disable dependencies manager on startup
	IgnoreDependencies bool

	// Headless execute Selenium webdriver in headless mode
	Headless bool

	// Port : communication port
	Port uint16
}

// CreateClientConfig create default ClientConfig instance and return a pointer on it
func CreateClientConfig() *ClientConfig {
	return &ClientConfig{
		LogLevel:            logrus.InfoLevel,
		ForceDependenciesDl: false,
		Debug:               false,
		IgnoreDependencies:  false,
		Headless:            false,
		Port:                8080,
	}
}
