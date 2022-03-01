package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/hbollon/igopher/internal/config/types"
	"github.com/hbollon/igopher/internal/logger"
	"github.com/hbollon/igopher/internal/proxy"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	requiredDirectories = [...]string{"./lib", "./config"}
)

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

// ClearData remove all IGopher data sub-folder and their content.
// It will recreate the necessary environment at the end no matter if an error has occurred or not.
func ClearData() error {
	defer CheckEnvironment()
	defer logger.SetLoggerOutput()
	var err error
	dirs := []string{"./logs", "./config", "./data"}
	for _, dir := range dirs {
		err = os.RemoveAll(dir)
		if err != nil {
			return err
		}
	}
	return nil
}

// ImportConfig read config.yaml, parse it in BotConfigYaml instance and finally return it
func ImportConfig() types.BotConfigYaml {
	var c types.BotConfigYaml
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
func ExportConfig(c types.BotConfigYaml) {
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
func ResetBotConfig() types.BotConfigYaml {
	return types.BotConfigYaml{
		Account: types.AccountYaml{
			Username: "",
			Password: "",
		},
		SrcUsers: types.ScrapperYaml{
			Accounts: []string{""},
			Quantity: 500,
		},
		AutoDm: types.AutoDmYaml{
			DmTemplates: []string{"Hey ! What's up?"},
			Greeting: types.GreetingYaml{
				Template:  "Hello",
				Activated: false,
			},
			Activated: true,
		},
		Quotas: types.QuotasYaml{
			DmDay:     50,
			DmHour:    5,
			Activated: true,
		},
		Schedule: types.ScheduleYaml{
			BeginAt:   "08:00",
			EndAt:     "18:00",
			Activated: true,
		},
		Blacklist: types.BlacklistYaml{
			Activated: true,
		},
		Selenium: types.SeleniumYaml{
			Proxy: proxy.Proxy{
				RemoteIP:       "",
				RemotePort:     8080,
				RemoteUsername: "",
				RemotePassword: "",
				WithAuth:       false,
				Enabled:        false,
			},
		},
	}
}
