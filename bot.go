package igopher

import (
	"context"
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
	IgnoreDependenciesFlag: flag.Bool("ignore-dependencies", false, "Skip dependencies management"),
	HeadlessFlag:           flag.Bool("headless", false, "Run WebDriver with frame buffer"),
	PortFlag:               flag.Int("port", 8080, "Specify custom communication port"),
}

func init() {
	// Add formatter to logrus in order to display line and function with messages
	formatter := logRuntime.Formatter{ChildFormatter: &log.TextFormatter{
		FullTimestamp: false,
		ForceColors:   true,
	}}
	formatter.Line = true
	log.SetFormatter(&formatter)

	// Output to stderr
	if runtime.GOOS == "windows" {
		log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stderr))
	} else {
		log.SetOutput(os.Stderr)
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

// Initialize client and bot configs, download dependencies,
// launch Selenium instance and finally run dm bot routine
func launchDmBot(ctx context.Context, hotReloadCh, reloadCh chan bool) {
	// Initialize client configuration
	clientConfig := initClientConfig()
	BotStruct = ReadBotConfigYaml()

	// Download dependencies
	if !clientConfig.IgnoreDependencies {
		DownloadDependencies(true, true, clientConfig.ForceDependenciesDl)
	}

	// Initialize Selenium and WebDriver and defer their closing
	BotStruct.SeleniumStruct.InitializeSelenium(clientConfig)
	BotStruct.SeleniumStruct.InitChromeWebDriver()
	defer BotStruct.SeleniumStruct.CloseSelenium()

	// Creation of needed communication channels and deferring their closing
	infoCh := make(chan string)
	defer close(infoCh)
	errCh := make(chan string)
	defer close(errCh)
	crashCh := make(chan error)
	defer close(crashCh)
	exitCh := make(chan bool)
	defer close(exitCh)
	reloadCallback := make(chan bool)
	defer close(reloadCallback)
	hotReloadCallback := make(chan bool)
	defer close(hotReloadCallback)

	// Start bot routine
	go func() {
		for {
			rand.Seed(time.Now().Unix())
			if err := BotStruct.Scheduler.CheckTime(); err == nil {
				BotStruct.ConnectToInstagram()
				users, err := BotStruct.FetchUsersFromUserFollowers()
				if err != nil {
					crashCh <- err
					BotStruct.SeleniumStruct.Fatal("Failed users fetching: ", err)
				}
				for _, username := range users {
					select {
					case <-hotReloadCallback:
						fmt.Println("hotReloadCallback")
						break
					case <-reloadCallback:
						fmt.Println("reloadCallback")
						break
					case <-exitCh:
						logrus.Infof("Bot process successfully stopped.")
						return
					default:
						break
					}
					res, err := BotStruct.SendMessage(username, BotStruct.DmModule.DmTemplates[rand.Intn(len(BotStruct.DmModule.DmTemplates))])
					if !res || err != nil {
						errCh <- fmt.Sprintf("Error during message sending: %v", err)
						log.Errorf("Error during message sending: %v", err)
					}
				}
			} else {
				crashCh <- err
				BotStruct.SeleniumStruct.Fatal("Error on bot launch: ", err)
			}
		}
	}()
	var msg string
	for {
		select {
		case msg = <-infoCh:
			fmt.Printf("infoCh: %s", msg)
			break
		case msg = <-errCh:
			fmt.Printf("errCh: %s", msg)
			break
		case err := <-crashCh:
			fmt.Printf("crashCh: %v", err)
			return
		case <-hotReloadCh:
			hotReloadCallback <- true
			return
		case <-reloadCh:
			reloadCallback <- true
			return
		case <-ctx.Done():
			exitCh <- true
			return
		default:
			break
		}
	}
}
