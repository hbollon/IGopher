package igopher

import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// BotConfigYaml is the raw representation of the yaml bot config file
type BotConfigYaml struct {
	Account struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"account"`
	SrcUsers struct {
		Accounts []string `yaml:"src_accounts"`
		Quantity int      `yaml:"fetch_quantity"`
	} `yaml:"users_src"`
	AutoDm struct {
		DmTemplates []string `yaml:"dm_templates"`
		Greeting    struct {
			Template  string `yaml:"template"`
			Activated bool   `yaml:"activated"`
		} `yaml:"greeting"`
		Activated bool `yaml:"activated"`
	} `yaml:"auto_dm"`
	Quotas struct {
		DmDay     int  `yaml:"dm_per_day"`
		DmHour    int  `yaml:"dm_per_hour"`
		Activated bool `yaml:"activated"`
	} `yaml:"quotas"`
	Schedule struct {
		BeginAt   string `yaml:"begin_at"`
		EndAt     string `yaml:"end_at"`
		Activated bool   `yaml:"activated"`
	} `yaml:"schedule"`
	Blacklist struct {
		Activated bool `yaml:"activated"`
	} `yaml:"blacklist"`
}

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
	Greeting GreetingConfig `yaml:"greeting"`
}

// GreetingConfig store greeting configuration for AutoDM module
type GreetingConfig struct {
	Activated bool `yaml:"activated"`
	// Add a string before the username
	Template string `yaml:"template"`
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
	// Scheduler
	Scheduler SchedulerManager `yaml:"schedule"`
	// Interracted users blacklist
	Blacklist BlacklistManager `yaml:"blacklist"`
}

/* Yaml custom parser */

// CustomTime is a custom time.Time used to set a custom yaml unmarshal rule
type CustomTime struct {
	time.Time
}

func (t *CustomTime) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var buf string
	err := unmarshal(&buf)
	if err != nil {
		return nil
	}

	tt, err := time.Parse("15:04", strings.TrimSpace(buf))
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
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
	err = c.Blacklist.InitializeBlacklist()
	if err != nil {
		logrus.Errorf("Failed to initialize blacklist: %v", err)
	}
	logrus.Debugf("config.yaml: %+v\n\n", c)
	return c
}

// ImportConfig read config.yaml and parse it in BotConfigYaml instance
func ImportConfig() BotConfigYaml {
	var c BotConfigYaml
	file, err := ioutil.ReadFile("./config/config.yaml")
	if err != nil {
		log.Fatalf("Error opening config file: %s", err)
	}

	err = yaml.Unmarshal(file, &c)
	if err != nil {
		log.Fatalf("Error during unmarshal config file: %s\n", err)
	}

	return c
}

// ExportConfig export BotConfigYaml instance to config.yaml config file
func ExportConfig(c BotConfigYaml) {
	out, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatalf("Error during marshal config file: %s\n", err)
	}

	err = ioutil.WriteFile("./config/config.yaml", out, os.ModePerm)
	if err != nil {
		log.Fatalf("Error during config file writing: %s\n", err)
	}
}

// ResetBotConfig return default bot configuration instance
func ResetBotConfig() BotConfigYaml {
	return BotConfigYaml{
		Account: struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		}{
			Username: "",
			Password: "",
		},
		SrcUsers: struct {
			Accounts []string `yaml:"src_accounts"`
			Quantity int      `yaml:"fetch_quantity"`
		}{
			Accounts: []string{""},
			Quantity: 500,
		},
		AutoDm: struct {
			DmTemplates []string `yaml:"dm_templates"`
			Greeting    struct {
				Template  string `yaml:"template"`
				Activated bool   `yaml:"activated"`
			} `yaml:"greeting"`
			Activated bool `yaml:"activated"`
		}{
			DmTemplates: []string{"Hey ! What's up?"},
			Greeting: struct {
				Template  string `yaml:"template"`
				Activated bool   `yaml:"activated"`
			}{
				Template:  "Hello",
				Activated: false,
			},
			Activated: true,
		},
		Quotas: struct {
			DmDay     int  `yaml:"dm_per_day"`
			DmHour    int  `yaml:"dm_per_hour"`
			Activated bool `yaml:"activated"`
		}{
			DmDay:     50,
			DmHour:    5,
			Activated: true,
		},
		Schedule: struct {
			BeginAt   string `yaml:"begin_at"`
			EndAt     string `yaml:"end_at"`
			Activated bool   `yaml:"activated"`
		}{
			BeginAt:   "8:00",
			EndAt:     "18:00",
			Activated: true,
		},
		Blacklist: struct {
			Activated bool `yaml:"activated"`
		}{
			Activated: true,
		},
	}
}
