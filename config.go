package igopher

import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	// Blank import to avoid deletion by linter
	// Used for struct fieldl validate metadata
	_ "github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// BotConfigYaml is the raw representation of the yaml bot config file
type BotConfigYaml struct {
	Account   AccountYaml   `yaml:"account"`
	SrcUsers  SrcUsersYaml  `yaml:"users_src"`
	AutoDm    AutoDmYaml    `yaml:"auto_dm"`
	Quotas    QuotasYaml    `yaml:"quotas"`
	Schedule  ScheduleYaml  `yaml:"schedule"`
	Blacklist BlacklistYaml `yaml:"blacklist"`
}

// AccountYaml is the yaml account configuration representation
type AccountYaml struct {
	Username string `yaml:"username" validate:"required,min=1,max=30"`
	Password string `yaml:"password" validate:"required"`
}

// SrcUsersYaml is the yaml user scrapping configuration representation
type SrcUsersYaml struct {
	Accounts []string `yaml:"src_accounts"`
	Quantity int      `yaml:"fetch_quantity" validate:"numeric"`
}

// AutoDmYaml is the yaml autodm module configuration representation
type AutoDmYaml struct {
	DmTemplates []string     `yaml:"dm_templates"`
	Greeting    GreetingYaml `yaml:"greeting"`
	Activated   bool         `yaml:"activated"`
}

// GreetingYaml is the yaml dm greeting configuration representation
type GreetingYaml struct {
	Template  string `yaml:"template"`
	Activated bool   `yaml:"activated"`
}

// QuotasYaml is the yaml quotas module configuration representation
type QuotasYaml struct {
	DmDay     int  `yaml:"dm_per_day" validate:"numeric"`
	DmHour    int  `yaml:"dm_per_hour" validate:"numeric"`
	Activated bool `yaml:"activated"`
}

// ScheduleYaml is the yaml scheduler module configuration representation
type ScheduleYaml struct {
	BeginAt   string `yaml:"begin_at" validate:"contains=:"`
	EndAt     string `yaml:"end_at" validate:"contains=:"`
	Activated bool   `yaml:"activated"`
}

// BlacklistYaml is the yaml blacklist module configuration representation
type BlacklistYaml struct {
	Activated bool `yaml:"activated"`
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

// ScrapperConfig store scrapper configuration for user fetching
// It also store fetched usernames
type ScrapperConfig struct {
	SrcAccounts     []string `yaml:"src_accounts"`
	FetchedAccounts []string
	Quantity        int `yaml:"fetch_quantity"`
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
	// Scrapper
	Scrapper ScrapperConfig `yaml:"users_src"`
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

// Read config yml file and initialize it for use with bot
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

// ImportConfig read config.yaml, parse it in BotCvalidate:"numeric"onfigYaml instance and finally return it
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
		Account: AccountYaml{
			Username: "",
			Password: "",
		},
		SrcUsers: SrcUsersYaml{
			Accounts: []string{""},
			Quantity: 500,
		},
		AutoDm: AutoDmYaml{
			DmTemplates: []string{"Hey ! What's up?"},
			Greeting: GreetingYaml{
				Template:  "Hello",
				Activated: false,
			},
			Activated: true,
		},
		Quotas: QuotasYaml{
			DmDay:     50,
			DmHour:    5,
			Activated: true,
		},
		Schedule: ScheduleYaml{
			BeginAt:   "8:00",
			EndAt:     "18:00",
			Activated: true,
		},
		Blacklist: BlacklistYaml{
			Activated: true,
		},
	}
}
