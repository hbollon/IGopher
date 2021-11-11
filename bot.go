package igopher

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/hbollon/igopher/internal/process"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// BotStruct is the main struct instance used by this bot
var BotStruct IGopher

// Flags declarations
var Flags = struct {
	// LogLevelFlag set loglevel threshold
	// If undefined or wrong set it to INFO level
	LogLevelFlag *string

	// ForceDlFlag force re-download of all dependencies
	ForceDlFlag *bool

	// DebugFlag set selenium debug mode and display its logging to stderr
	DebugFlag *bool

	// DevToolsFlag launch Electron gui with devtools openned
	DevToolsFlag *bool

	// IgnoreDependenciesFlag disable dependencies manager on startup
	IgnoreDependenciesFlag *bool

	// BackgroundFlag IGopher as background task with actual configuration and ignore TUI
	BackgroundFlag *bool

	// HeadlessFlag execute Selenium webdriver in headless mode
	HeadlessFlag *bool

	// PortFlag specifie custom communication port for Selenium and web drivers
	PortFlag *int
}{
	LogLevelFlag:           flag.String("loglevel", "info", "Log level threshold"),
	ForceDlFlag:            flag.Bool("force-download", false, "Force redownload of all dependencies even if exists"),
	DebugFlag:              flag.Bool("debug", false, "Display debug and selenium output"),
	DevToolsFlag:           flag.Bool("dev-tools", false, "Launch Electron gui with dev tools openned"),
	IgnoreDependenciesFlag: flag.Bool("ignore-dependencies", false, "Skip dependencies management"),
	HeadlessFlag:           flag.Bool("headless", false, "Run WebDriver with frame buffer"),
	PortFlag:               flag.Int("port", 8080, "Specify custom communication port"),
}

// errStopBot is used to trigger bot stopping from some function
var errStopBot = errors.New("Bot stop process triggered")

func initClientConfig() *ClientConfig {
	clientConfig := CreateClientConfig()
	clientConfig.LogLevel, _ = log.ParseLevel(*Flags.LogLevelFlag)
	clientConfig.ForceDependenciesDl = *Flags.ForceDlFlag
	clientConfig.Debug = *Flags.DebugFlag
	clientConfig.DevTools = *Flags.DevToolsFlag
	clientConfig.IgnoreDependencies = *Flags.IgnoreDependenciesFlag
	clientConfig.Headless = *Flags.HeadlessFlag

	if *Flags.PortFlag > math.MaxUint16 || *Flags.PortFlag < 8080 {
		log.Warnf("Invalid port argument '%d'. Use default 8080.", *Flags.PortFlag)
	} else {
		clientConfig.Port = uint16(*Flags.PortFlag)
	}

	return clientConfig
}

// LaunchBotTui start dm bot on main goroutine
func LaunchBotTui() {
	// Initialize client configuration
	var err error
	clientConfig := initClientConfig()
	BotStruct, err = ReadBotConfigYaml()
	if err != nil {
		logrus.Warn(err)
	}

	// Download dependencies
	if !clientConfig.IgnoreDependencies {
		DownloadDependencies(true, true, clientConfig.ForceDependenciesDl)
	}

	// Initialize Selenium and WebDriver and defer their closing
	BotStruct.SeleniumStruct.InitializeSelenium(clientConfig)
	BotStruct.SeleniumStruct.InitChromeWebDriver()
	defer BotStruct.SeleniumStruct.CloseSelenium()

	process.DumpProcessPidToFile()

	rand.Seed(time.Now().Unix())
	if err = BotStruct.Scheduler.CheckTime(); err == nil {
		BotStruct.ConnectToInstagram()
		for {
			var users []string
			users, err = BotStruct.FetchUsersFromUserFollowers()
			if err != nil {
				BotStruct.SeleniumStruct.Fatal("Failed users fetching: ", err)
			}
			for _, username := range users {
				var res bool
				res, err = BotStruct.SendMessage(username, BotStruct.DmModule.DmTemplates[rand.Intn(len(BotStruct.DmModule.DmTemplates))])
				if !res || err != nil {
					log.Errorf("Error during message sending: %v", err)
				}
			}
		}
	} else {
		BotStruct.SeleniumStruct.Fatal("Error on bot launch: ", err)
	}
}

func checkBotChannels() bool {
	select {
	case <-BotStruct.hotReloadCallback:
		if err := BotStruct.HotReload(); err != nil {
			logrus.Errorf("Bot hot reload failed: %v", err)
			BotStruct.hotReloadCallback <- false
		} else {
			logrus.Info("Bot hot reload successfully.")
			BotStruct.hotReloadCallback <- true
		}
		break
	case <-BotStruct.reloadCallback:
		logrus.Info("Bot reload successfully.")
		break
	case <-BotStruct.exitCh:
		logrus.Info("Bot process successfully stopped.")
		return true
	default:
		break
	}

	return false
}

// Initialize client and bot configs, download dependencies,
// launch Selenium instance and finally run dm bot routine
func launchBot(ctx context.Context) {
	// Initialize client configuration
	var err error
	clientConfig := initClientConfig()
	BotStruct, err = ReadBotConfigYaml()
	if err != nil {
		logrus.Warn(err)
	}
	BotStruct.running = true

	// Download dependencies
	if !clientConfig.IgnoreDependencies {
		DownloadDependencies(true, true, clientConfig.ForceDependenciesDl)
	}

	// Initialize Selenium and WebDriver and defer their closing
	BotStruct.SeleniumStruct.InitializeSelenium(clientConfig)
	BotStruct.SeleniumStruct.InitChromeWebDriver()
	defer BotStruct.SeleniumStruct.CloseSelenium()
	defer BotStruct.SeleniumStruct.Proxy.StopForwarderProxy()

	// Creation of needed communication channels and deferring their closing
	exitedCh = make(chan bool)
	defer close(exitedCh)
	hotReloadCh = make(chan bool)
	defer close(hotReloadCh)
	reloadCh = make(chan bool)
	defer close(reloadCh)

	BotStruct.infoCh = make(chan string)
	defer close(BotStruct.infoCh)
	BotStruct.errCh = make(chan string)
	defer close(BotStruct.errCh)
	BotStruct.crashCh = make(chan error)
	defer close(BotStruct.crashCh)
	BotStruct.exitCh = make(chan bool)
	defer close(BotStruct.exitCh)
	BotStruct.reloadCallback = make(chan bool)
	defer close(BotStruct.reloadCallback)
	BotStruct.hotReloadCallback = make(chan bool)
	defer close(BotStruct.hotReloadCallback)

	process.DumpProcessPidToFile()

	// Start bot routine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("Unknown error: %v", r)
				SendMessageToElectron(
					MessageOut{
						Status:  ERROR,
						Msg:     "bot crash",
						Payload: fmt.Errorf("Unknown error: %v", r),
					},
				)
				BotStruct.running = false
			}
		}()
		rand.Seed(time.Now().Unix())
		if err = BotStruct.Scheduler.CheckTime(); err == nil {
			if exit := checkBotChannels(); exit {
				return
			}
			BotStruct.ConnectToInstagram()
			for {
				var users []string
				if exit := checkBotChannels(); exit {
					return
				}
				users, err = BotStruct.FetchUsersFromUserFollowers()
				if err != nil {
					BotStruct.crashCh <- fmt.Errorf("Failed users fetching: %v. Check logs tab for more details", err)
					return
				}
				for _, username := range users {
					if exit := checkBotChannels(); exit {
						return
					}
					var res bool
					res, err = BotStruct.SendMessage(username, BotStruct.DmModule.DmTemplates[rand.Intn(len(BotStruct.DmModule.DmTemplates))])
					if !res || err != nil {
						BotStruct.errCh <- fmt.Sprintf("Error during message sending: %v", err)
						log.Errorf("Error during message sending: %v", err)
					}
				}
			}
		} else {
			if err == errStopBot {
				return
			}
			BotStruct.crashCh <- err
			BotStruct.SeleniumStruct.Fatal("Error on bot launch: ", err)
		}
	}()
	var msg string
	for {
		select {
		case msg = <-BotStruct.infoCh:
			log.Infof("infoCh: %s", msg)
			break
		case msg = <-BotStruct.errCh:
			log.Errorf("errCh: %s", msg)
			break
		case err := <-BotStruct.crashCh:
			log.Errorf("crashCh: %v", err)
			SendMessageToElectron(
				MessageOut{
					Status:  ERROR,
					Msg:     "bot crash",
					Payload: err.Error(),
				},
			)
			BotStruct.running = false
			return
		case <-hotReloadCh:
			BotStruct.hotReloadCallback <- true
			if <-BotStruct.hotReloadCallback {
				hotReloadCh <- true
			} else {
				hotReloadCh <- false
			}
			break
		case <-reloadCh:
			BotStruct.reloadCallback <- true
			return
		case <-ctx.Done():
			BotStruct.exitCh <- true
			exitedCh <- true
			BotStruct.running = false
			return
		default:
			break
		}

		if ws, err := BotStruct.SeleniumStruct.WebDriver.WindowHandles(); len(ws) == 0 || err != nil {
			BotStruct.SeleniumStruct.CleanUp()
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
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
