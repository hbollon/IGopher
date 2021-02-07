package igopher

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	requiredDirectories = [...]string{"./lib", "./config"}
)

// IGopher struct store all bot and ig related configuration and modules instances.
// Settings are readed from Yaml config files.
type IGopher struct {
	// SeleniumStruct contain all selenium stuff and config
	SeleniumStruct Selenium
	// User credentials
	UserAccount Account `yaml:"account"`
	// Automatic messages sending module
	DmModule AutoDM `yaml:"auto_dm"`
	// Quotas
	Quotas QuotaManager `yaml:"quotas"`
	// Scrapper
	ScrapperManager ScrapperConfig `yaml:"scrapper"`
	// Scheduler
	Scheduler SchedulerManager `yaml:"schedule"`
	// Interracted users blacklist
	Blacklist BlacklistManager `yaml:"blacklist"`
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
}

// Account store personnal credentials
type Account struct {
	Username string `yaml:"username" validate:"required,min=1,max=30"`
	Password string `yaml:"password" validate:"required,min=1"`
}

// AutoDM store messaging module configuration
type AutoDM struct {
	Activated bool `yaml:"activated"`
	// List of all availlables message templates
	DmTemplates []string `yaml:"dm_templates" validate:"required"`
	// Greeting module add a customized DM header with recipient username
	Greeting GreetingConfig `yaml:"greeting"`
}

// GreetingConfig store greeting configuration for AutoDM module
type GreetingConfig struct {
	Activated bool `yaml:"activated"`
	// Add a string before the username
	Template string `yaml:"template"`
}

/* Yaml */

// BotConfigYaml is the raw representation of the yaml bot config file
type BotConfigYaml struct {
	Account   AccountYaml   `yaml:"account"`
	SrcUsers  ScrapperYaml  `yaml:"scrapper"`
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

// ScrapperYaml is the yaml user scrapping configuration representation
type ScrapperYaml struct {
	Accounts []string `yaml:"src_accounts" validate:"required"`
	Quantity int      `yaml:"fetch_quantity" validate:"numeric,min=1"`
}

// AutoDmYaml is the yaml autodm module configuration representation
type AutoDmYaml struct {
	DmTemplates []string     `yaml:"dm_templates" validate:"required"`
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

// CheckEnvironment check existence of sub-directories and files required
// for the operation of the program and creates them otherwise
func CheckEnvironment() {
	// Check and create directories
	for _, dir := range requiredDirectories {
		dir = filepath.FromSlash(dir)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err = os.Mkdir(dir, 0755); err != nil {
				log.Fatalf("Error during creation of '%s' sub-directory,"+
					" check root directory permissions or try to create it manually\nMkdir error:\n%v", dir, err)
			}
		}
	}

	// Check config.yaml existence
	if _, err := os.Stat(filepath.FromSlash("./config/config.yaml")); os.IsNotExist(err) {
		ExportConfig(ResetBotConfig())
	}
}

// CheckConfigValidity check bot config validity
func CheckConfigValidity() error {
	config := ImportConfig()
	validate := validator.New()
	if err := validate.Struct(config.Account); err != nil {
		return errors.New("Invalid credentials format! Please check your settings")
	}
	if err := validate.Struct(config.SrcUsers); err != nil {
		return errors.New("Invalid scrapper configuration! Please check your settings")
	}
	if err := validate.Struct(config.AutoDm); err != nil {
		return errors.New("Invalid autodm module configuration! Please check your settings")
	}

	return nil
}

// ReadBotConfigYaml read config yml file and initialize it for use with bot
func ReadBotConfigYaml() IGopher {
	var c IGopher
	file, err := ioutil.ReadFile(filepath.FromSlash("./config/config.yaml"))
	if err != nil {
		logrus.Fatalf("Error opening config file: %s", err)
	}

	err = yaml.Unmarshal(file, &c)
	if err != nil {
		logrus.Fatalf("Error during unmarshal config file: %s\n", err)
	}

	c.Quotas.InitializeQuotaManager()
	err = c.Scheduler.InitializeScheduler()
	if err != nil {
		logrus.Errorf("Failed to initialize scheduler: %v", err)
	}
	err = c.Blacklist.InitializeBlacklist()
	if err != nil {
		logrus.Errorf("Failed to initialize blacklist: %v", err)
	}
	logrus.Debugf("config.yaml: %+v\n\n", c)
	return c
}

// ImportConfig read config.yaml, parse it in BotConfigYaml instance and finally return it
func ImportConfig() BotConfigYaml {
	var c BotConfigYaml
	file, err := ioutil.ReadFile(filepath.FromSlash("./config/config.yaml"))
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

	err = ioutil.WriteFile(filepath.FromSlash("./config/config.yaml"), out, os.ModePerm)
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
		SrcUsers: ScrapperYaml{
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
