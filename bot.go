package igopher

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"

	logRuntime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// BotStruct is the main struct instance used by this bot
var BotStruct IGopher

/// Flags
var flags = struct {
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

func init() {
	// Add formatter to logrus in order to display line and function with messages
	formatter := logRuntime.Formatter{ChildFormatter: &log.TextFormatter{
		FullTimestamp: false,
		ForceColors:   true,
	}}
	formatter.Line = true
	log.SetFormatter(&formatter)

	// Output to stderr and log file
	if _, err := os.Stat("./logs/"); os.IsNotExist(err) {
		os.Mkdir("./logs/", os.ModePerm)
	}
	logFile, err := os.OpenFile("./logs/igopher.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Can't initialize logger: %v", err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	if runtime.GOOS == "windows" {
		log.SetOutput(ansicolor.NewAnsiColorWriter(mw))
	} else {
		log.SetOutput(mw)
	}

	flag.Parse()
	level, err := log.ParseLevel(*flags.LogLevelFlag)
	if err == nil {
		log.SetLevel(level)
	} else {
		log.SetLevel(log.InfoLevel)
		log.Warnf("Invalid log level '%s', use default one.", *flags.LogLevelFlag)
	}
}

func initClientConfig() *ClientConfig {
	clientConfig := CreateClientConfig()
	clientConfig.LogLevel, _ = log.ParseLevel(*flags.LogLevelFlag)
	clientConfig.ForceDependenciesDl = *flags.ForceDlFlag
	clientConfig.Debug = *flags.DebugFlag
	clientConfig.DevTools = *flags.DevToolsFlag
	clientConfig.IgnoreDependencies = *flags.IgnoreDependenciesFlag
	clientConfig.Headless = *flags.HeadlessFlag

	if *flags.PortFlag > math.MaxUint16 || *flags.PortFlag < 8080 {
		log.Warnf("Invalid port argument '%d'. Use default 8080.", *flags.PortFlag)
	} else {
		clientConfig.Port = uint16(*flags.PortFlag)
	}

	return clientConfig
}

// LaunchGui initialize environment needed by IGopher and run it with his Gui
func LaunchGui() {
	// Initialize environment
	CheckEnvironment()

	// Launch GUI
	InitGui()
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

	rand.Seed(time.Now().Unix())
	if err = BotStruct.Scheduler.CheckTime(); err == nil {
		BotStruct.ConnectToInstagram()
		for {
			var users []string
			users, err = BotStruct.FetchUsersFromUserFollowers()
			if err != nil {
				log.Error(err)
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

// Initialize client and bot configs, download dependencies,
// launch Selenium instance and finally run dm bot routine
func launchDmBot(ctx context.Context) {
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

	// Start bot routine
	go func() {
		rand.Seed(time.Now().Unix())
		if err = BotStruct.Scheduler.CheckTime(); err == nil {
			BotStruct.ConnectToInstagram()
			for {
				users, err := BotStruct.FetchUsersFromUserFollowers()
				if err != nil {
					BotStruct.crashCh <- err
					BotStruct.SeleniumStruct.Fatal("Failed usersDm bot successfully stopped! fetching: ", err)
				}
				for _, username := range users {
					select {
					case <-BotStruct.hotReloadCallback:
						if err = BotStruct.HotReload(); err != nil {
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
						return
					default:
						break
					}
					res, err := BotStruct.SendMessage(username, BotStruct.DmModule.DmTemplates[rand.Intn(len(BotStruct.DmModule.DmTemplates))])
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
			fmt.Printf("infoCh: %s", msg)
			break
		case msg = <-BotStruct.errCh:
			fmt.Printf("errCh: %s", msg)
			break
		case err := <-BotStruct.crashCh:
			fmt.Printf("crashCh: %v", err)
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
			BotStruct.SeleniumStruct.CloseSelenium()
			exitedCh <- true
			BotStruct.running = false
			return
		default:
			break
		}
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
