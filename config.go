package instadm

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
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

	BotConfig BotConfig
}

// Account store personnal credentials
type Account struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// AutoDM store messaging module configuration
type AutoDM struct {
	Activated bool `yaml:"activated"`
	// List of all availlables message templates
	DmTemplates []string `yaml:"dm_templates"`
	// Greeting module add a customized DM header with recipient username
	Greeting struct {
		Activated bool `yaml:"activated"`
		// Add a string before the username
		Template string `yaml:"template"`
	} `yaml:"greeting"`
}

// BotConfig struct store all bot and ig related configuration .
// These parameters are readed from Yaml config files.
type BotConfig struct {
	// User credentials
	UserAccount Account `yaml:"account"`
	// Automatic messages sending module
	DmModule AutoDM `yaml:"auto_dm"`
	// Quotas
	Quotas QuotaManager `yaml:"quotas"`
	// Interracted users blacklist
	Blacklist bool `yaml:"blacklist_interacted_users"`
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
		BotConfig:           readBotConfigYaml(),
	}
}

// Read config yml file
func readBotConfigYaml() BotConfig {
	var c BotConfig
	file, err := ioutil.ReadFile("./config/config.yaml")
	if err != nil {
		logrus.Fatalf("Error opening config file: %s", err)
	}

	err = yaml.Unmarshal(file, &c)
	if err != nil {
		logrus.Fatalf("Error during unmarshal config file: %s\n", err)
	}

	c.Quotas.InitializeQuotaManager()
	logrus.Debugf("config.yaml: %+v\n\n", c)
	return c
}
