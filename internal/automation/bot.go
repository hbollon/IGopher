package automation

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/hbollon/igopher/internal/actions"
	"github.com/hbollon/igopher/internal/config/flags"
	confdata "github.com/hbollon/igopher/internal/config/types"
	dep "github.com/hbollon/igopher/internal/dependency"
	"github.com/hbollon/igopher/internal/engine"
	"github.com/hbollon/igopher/internal/gui/comm"
	"github.com/hbollon/igopher/internal/gui/datatypes"
	"github.com/hbollon/igopher/internal/process"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var (
	// BotStruct is the main struct instance used by this bot
	BotStruct                       confdata.IGopher
	ReloadCh, HotReloadCh, ExitedCh chan bool
)

// ErrStopBot is used to trigger bot stopping from some function
var ErrStopBot = errors.New("Bot stop process triggered")

func initClientConfig() *engine.ClientConfig {
	clientConfig := engine.CreateClientConfig()
	clientConfig.LogLevel, _ = log.ParseLevel(*flags.Flags.LogLevelFlag)
	clientConfig.ForceDependenciesDl = *flags.Flags.ForceDlFlag
	clientConfig.Debug = *flags.Flags.DebugFlag
	clientConfig.DevTools = *flags.Flags.DevToolsFlag
	clientConfig.IgnoreDependencies = *flags.Flags.IgnoreDependenciesFlag
	clientConfig.Headless = *flags.Flags.HeadlessFlag

	if *flags.Flags.PortFlag > math.MaxUint16 || *flags.Flags.PortFlag < 8080 {
		log.Warnf("Invalid port argument '%d'. Use default 8080.", *flags.Flags.PortFlag)
	} else {
		clientConfig.Port = uint16(*flags.Flags.PortFlag)
	}

	return clientConfig
}

// LaunchBotTui start dm bot on main goroutine
func LaunchBotTui() {
	// Initialize client configuration
	var err error
	clientConfig := initClientConfig()
	BotStruct, err = confdata.ReadBotConfigYaml()
	if err != nil {
		logrus.Warn(err)
	}

	// Download dependencies
	if !clientConfig.IgnoreDependencies {
		dep.DownloadDependencies(true, false, clientConfig.ForceDependenciesDl)
	}

	// Initialize Selenium and WebDriver and defer their closing
	BotStruct.SeleniumStruct.InitializeSelenium(clientConfig)
	BotStruct.SeleniumStruct.InitChromeWebDriver()
	defer BotStruct.SeleniumStruct.CloseSelenium()

	process.DumpProcessPidToFile()

	rand.Seed(time.Now().Unix())
	if err = BotStruct.Scheduler.CheckTime(); err == nil {
		actions.ConnectToInstagram(&BotStruct)
		for {
			var users []string
			users, err = actions.FetchUsersFromUserFollowers(&BotStruct)
			if err != nil {
				BotStruct.SeleniumStruct.Fatal("Failed users fetching: ", err)
			}
			for _, username := range users {
				var res bool
				res, err = actions.SendMessage(&BotStruct, username, BotStruct.DmModule.DmTemplates[rand.Intn(len(BotStruct.DmModule.DmTemplates))])
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
	case <-BotStruct.HotReloadCallback:
		if err := BotStruct.HotReload(); err != nil {
			logrus.Errorf("Bot hot reload failed: %v", err)
			BotStruct.HotReloadCallback <- false
		} else {
			logrus.Info("Bot hot reload successfully.")
			BotStruct.HotReloadCallback <- true
		}
		break
	case <-BotStruct.ReloadCallback:
		logrus.Info("Bot reload successfully.")
		break
	case <-BotStruct.ExitCh:
		logrus.Info("Bot process successfully stopped.")
		return true
	default:
		break
	}

	return false
}

// Initialize client and bot configs, download dependencies,
// launch Selenium instance and finally run dm bot routine
func LaunchBot(ctx context.Context) {
	// Initialize client configuration
	var err error
	clientConfig := initClientConfig()
	BotStruct, err = confdata.ReadBotConfigYaml()
	if err != nil {
		logrus.Warn(err)
	}
	BotStruct.Running = true

	// Download dependencies
	if !clientConfig.IgnoreDependencies {
		dep.DownloadDependencies(true, false, clientConfig.ForceDependenciesDl)
	}

	// Initialize Selenium and WebDriver and defer their closing
	BotStruct.SeleniumStruct.InitializeSelenium(clientConfig)
	BotStruct.SeleniumStruct.InitChromeWebDriver()
	defer BotStruct.SeleniumStruct.CloseSelenium()
	defer BotStruct.SeleniumStruct.Proxy.StopForwarderProxy()

	// Creation of needed communication channels and deferring their closing
	ExitedCh = make(chan bool)
	defer close(ExitedCh)
	HotReloadCh = make(chan bool)
	defer close(HotReloadCh)
	ReloadCh = make(chan bool)
	defer close(ReloadCh)

	BotStruct.InfoCh = make(chan string)
	defer close(BotStruct.InfoCh)
	BotStruct.ErrCh = make(chan string)
	defer close(BotStruct.ErrCh)
	BotStruct.CrashCh = make(chan error)
	defer close(BotStruct.CrashCh)
	BotStruct.ExitCh = make(chan bool)
	defer close(BotStruct.ExitCh)
	BotStruct.ReloadCallback = make(chan bool)
	defer close(BotStruct.ReloadCallback)
	BotStruct.HotReloadCallback = make(chan bool)
	defer close(BotStruct.HotReloadCallback)

	process.DumpProcessPidToFile()

	// Start bot routine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("Unknown error: %v", r)
				comm.SendMessageToElectron(
					datatypes.MessageOut{
						Status:  datatypes.ERROR,
						Msg:     "bot crash",
						Payload: fmt.Errorf("Unknown error: %v", r),
					},
				)
				BotStruct.Running = false
			}
		}()
		rand.Seed(time.Now().Unix())
		if err = BotStruct.Scheduler.CheckTime(); err == nil {
			if exit := checkBotChannels(); exit {
				return
			}
			actions.ConnectToInstagram(&BotStruct)
			for {
				var users []string
				if exit := checkBotChannels(); exit {
					return
				}
				users, err = actions.FetchUsersFromUserFollowers(&BotStruct)
				if err != nil {
					BotStruct.CrashCh <- fmt.Errorf("Failed users fetching: %v. Check logs tab for more details", err)
					return
				}
				for _, username := range users {
					if exit := checkBotChannels(); exit {
						return
					}
					var res bool
					res, err = actions.SendMessage(&BotStruct, username, BotStruct.DmModule.DmTemplates[rand.Intn(len(BotStruct.DmModule.DmTemplates))])
					if !res || err != nil {
						BotStruct.ErrCh <- fmt.Sprintf("Error during message sending: %v", err)
						log.Errorf("Error during message sending: %v", err)
					}
				}
			}
		} else {
			if err == ErrStopBot {
				return
			}
			BotStruct.CrashCh <- err
			BotStruct.SeleniumStruct.Fatal("Error on bot launch: ", err)
		}
	}()
	var msg string
	for {
		select {
		case msg = <-BotStruct.InfoCh:
			log.Infof("infoCh: %s", msg)
			break
		case msg = <-BotStruct.ErrCh:
			log.Errorf("errCh: %s", msg)
			break
		case err := <-BotStruct.CrashCh:
			log.Errorf("crashCh: %v", err)
			comm.SendMessageToElectron(
				datatypes.MessageOut{
					Status:  datatypes.ERROR,
					Msg:     "bot crash",
					Payload: err.Error(),
				},
			)
			BotStruct.Running = false
			return
		case <-HotReloadCh:
			BotStruct.HotReloadCallback <- true
			if <-BotStruct.HotReloadCallback {
				HotReloadCh <- true
			} else {
				HotReloadCh <- false
			}
			break
		case <-ReloadCh:
			BotStruct.ReloadCallback <- true
			return
		case <-ctx.Done():
			BotStruct.ExitCh <- true
			ExitedCh <- true
			BotStruct.Running = false
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
