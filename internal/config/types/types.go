package types

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hbollon/igopher/internal/engine"
	"github.com/hbollon/igopher/internal/modules/blacklist"
	"github.com/hbollon/igopher/internal/modules/quotas"
	"github.com/hbollon/igopher/internal/modules/scheduler"
	"github.com/hbollon/igopher/internal/proxy"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// SplitStringSlice is a custom string slice type used to define a custom json unmarshal rule
type SplitStringSlice []string

// UnmarshalJSON custom rule for unmarshal string array from string by splitting it by ';'
func (strSlice *SplitStringSlice) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*strSlice = strings.Split(s, ";")
	return nil
}

// IGopher struct store all bot and ig related configuration and modules instances.
// Settings are readed from Yaml config files.
type IGopher struct {
	// SeleniumStruct contain all selenium stuff and config
	SeleniumStruct engine.Selenium `yaml:"webdriver"`
	// User credentials
	UserAccount Account `yaml:"account"`
	// Automatic messages sending module
	DmModule AutoDM `yaml:"auto_dm"`
	// Quotas
	Quotas quotas.QuotaManager `yaml:"quotas"`
	// Scrapper
	ScrapperManager ScrapperConfig `yaml:"scrapper"`
	// Scheduler
	Scheduler scheduler.SchedulerManager `yaml:"schedule"`
	// Interracted users blacklist
	Blacklist blacklist.BlacklistManager `yaml:"blacklist"`
	// Channels
	InfoCh            chan string `yaml:"-"`
	ErrCh             chan string `yaml:"-"`
	CrashCh           chan error  `yaml:"-"`
	ExitCh            chan bool   `yaml:"-"`
	HotReloadCallback chan bool   `yaml:"-"`
	ReloadCallback    chan bool   `yaml:"-"`
	// Running state
	Running bool `yaml:"-"`
}

// HotReload update bot config without stopping it
// Some settings cannot be updated this way like account credentials
func (bot *IGopher) HotReload() error {
	newConfig, err := ReadBotConfigYaml()
	if err != nil {
		return err
	}

	bot.DmModule = newConfig.DmModule
	bot.Quotas = newConfig.Quotas
	bot.ScrapperManager = newConfig.ScrapperManager
	bot.Scheduler = newConfig.Scheduler
	bot.Blacklist = newConfig.Blacklist
	return nil
}

// ReadBotConfigYaml read config yml file and initialize it for use with bot
func ReadBotConfigYaml() (IGopher, error) {
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
		return c, err
	}
	err = c.Blacklist.InitializeBlacklist()
	if err != nil {
		logrus.Errorf("Failed to initialize blacklist: %v", err)
		return c, err
	}
	return c, nil
}

// ScrapperConfig store scrapper configuration for user fetching
// It also store fetched usernames
type ScrapperConfig struct {
	SrcAccounts     []string `yaml:"src_accounts"`
	FetchedAccounts []string
	Quantity        int `yaml:"fetch_quantity" validate:"numeric"`
}

// Account store personnal credentials
type Account struct {
	Username string `json:"username" yaml:"username" validate:"required,min=1,max=30"`
	Password string `json:"password" yaml:"password" validate:"required,min=1"`
}

// AutoDM store messaging module configuration
type AutoDM struct {
	Activated bool `json:"dmActivated" yaml:"activated"`
	// List of all availlables message templates
	DmTemplates []string `json:"dmTemplates" yaml:"dm_templates" validate:"required"`
	// Greeting module add a customized DM header with recipient username
	Greeting GreetingConfig `yaml:"greeting"`
}

// GreetingConfig store greeting configuration for AutoDM module
type GreetingConfig struct {
	Activated bool `json:"greetingActivated" yaml:"activated"`
	// Add a string before the username
	Template string `json:"greetingTemplate" yaml:"template" validate:"required"`
}

/* Yaml */

// BotConfigYaml is the raw representation of the yaml bot config file
type BotConfigYaml struct {
	Account   AccountYaml   `json:"account" yaml:"account"`
	SrcUsers  ScrapperYaml  `json:"scrapper" yaml:"scrapper"`
	AutoDm    AutoDmYaml    `json:"auto_dm" yaml:"auto_dm"`
	Quotas    QuotasYaml    `json:"quotas" yaml:"quotas"`
	Schedule  ScheduleYaml  `json:"schedule" yaml:"schedule"`
	Blacklist BlacklistYaml `json:"blacklist" yaml:"blacklist"`
	Selenium  SeleniumYaml  `json:"webdriver" yaml:"webdriver"`
}

// AccountYaml is the yaml account configuration representation
type AccountYaml struct {
	Username string `json:"username" yaml:"username" validate:"required,min=1,max=30"`
	Password string `json:"password" yaml:"password" validate:"required"`
}

// ScrapperYaml is the yaml user scrapping configuration representation
type ScrapperYaml struct {
	Accounts SplitStringSlice `json:"srcUsers" yaml:"src_accounts" validate:"required"`
	Quantity int              `json:"scrappingQuantity,string" yaml:"fetch_quantity" validate:"numeric,min=1"`
}

// AutoDmYaml is the yaml autodm module configuration representation
type AutoDmYaml struct {
	DmTemplates SplitStringSlice `json:"dmTemplates" yaml:"dm_templates" validate:"required"`
	Greeting    GreetingYaml     `json:"greeting" yaml:"greeting"`
	Activated   bool             `json:"dmActivation,string" yaml:"activated"`
}

// GreetingYaml is the yaml dm greeting configuration representation
type GreetingYaml struct {
	Template  string `json:"greetingTemplate" yaml:"template"`
	Activated bool   `json:"greetingActivation,string" yaml:"activated"`
}

// QuotasYaml is the yaml quotas module configuration representation
type QuotasYaml struct {
	DmDay     int  `json:"dmDay,string" yaml:"dm_per_day" validate:"numeric,min=1"`
	DmHour    int  `json:"dmHour,string" yaml:"dm_per_hour" validate:"numeric,min=1"`
	Activated bool `json:"quotasActivation,string" yaml:"activated"`
}

// ScheduleYaml is the yaml scheduler module configuration representation
type ScheduleYaml struct {
	BeginAt   string `json:"beginAt" yaml:"begin_at" validate:"contains=:"`
	EndAt     string `json:"endAt" yaml:"end_at" validate:"contains=:"`
	Activated bool   `json:"scheduleActivation,string" yaml:"activated"`
}

// BlacklistYaml is the yaml blacklist module configuration representation
type BlacklistYaml struct {
	Activated bool `json:"blacklistActivation,string" yaml:"activated"`
}

// SeleniumYaml is the yaml selenium configuration representation
type SeleniumYaml struct {
	Proxy proxy.Proxy `json:"proxy" yaml:"proxy"`
}
