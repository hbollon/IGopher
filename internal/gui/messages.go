package gui

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/asticode/go-astilectron"
	"github.com/go-playground/validator/v10"
	"github.com/hbollon/igopher/internal/automation"
	bot "github.com/hbollon/igopher/internal/automation"
	conf "github.com/hbollon/igopher/internal/config"
	confdata "github.com/hbollon/igopher/internal/config/types"
	"github.com/hbollon/igopher/internal/gui/comm"
	"github.com/hbollon/igopher/internal/gui/datatypes"
	"github.com/hbollon/igopher/internal/logger"
	"github.com/hbollon/igopher/internal/proxy"
	"github.com/sirupsen/logrus"
)

var (
	config   confdata.BotConfigYaml
	validate = validator.New()
	ctx      context.Context
	cancel   context.CancelFunc
)

// CallbackMap is a map of callback functions for each message
var CallbackMap = map[string]func(m *datatypes.MessageIn) datatypes.MessageOut{
	"resetGlobalDefaultSettings":  resetGlobalSettingsCallback,
	"clearAllData":                clearDataCallback,
	"igCredentialsForm":           credentialsFormCallback,
	"quotasForm":                  quotasFormCallback,
	"schedulerForm":               schedulerCallback,
	"blacklistForm":               blacklistFormCallback,
	"dmSettingsForm":              dmBotFormCallback,
	"dmUserScrappingSettingsForm": dmScrapperFormCallback,
	"proxyForm":                   proxyFormCallback,
	"launchDmBot":                 launchDmBotCallback,
	"stopDmBot":                   stopDmBotCallback,
	"hotReloadBot":                hotReloadCallback,
	"getLogs":                     getLogsCallback,
	"getConfig":                   getConfigCallback,
}

// HandleMessages is handling function for incoming messages
func HandleMessages(w *astilectron.Window) {
	w.OnMessage(func(m *astilectron.EventMessage) interface{} {
		// Unmarshal
		var i datatypes.MessageIn
		var err error
		if err = m.Unmarshal(&i); err != nil {
			logrus.Errorf("Unmarshaling message %+v failed: %v", *m, err)
			return datatypes.MessageOut{Status: "Error during message reception"}
		}

		// Process message
		config = conf.ImportConfig()
		if callback, ok := CallbackMap[i.Msg]; ok {
			return i.Callback(callback)
		}
		logrus.Errorf("Unexpected message received: \"%s\"", i.Msg)
		return datatypes.MessageOut{Status: datatypes.ERROR, Msg: "Unknown error: Invalid message received"}
	})
	comm.Window = w
}

/* Callback functiosn to handle electron messages */

func resetGlobalSettingsCallback(m *datatypes.MessageIn) datatypes.MessageOut {
	config = conf.ResetBotConfig()
	conf.ExportConfig(config)
	return datatypes.MessageOut{Status: datatypes.SUCCESS, Msg: "Global configuration was successfully reset!"}
}

func clearDataCallback(m *datatypes.MessageIn) datatypes.MessageOut {
	if err := conf.ClearData(); err != nil {
		return datatypes.MessageOut{Status: datatypes.ERROR, Msg: fmt.Sprintf("IGopher data clearing failed! Error: %v", err)}
	}
	return datatypes.MessageOut{Status: datatypes.SUCCESS, Msg: "IGopher data successfully cleared!"}
}

func credentialsFormCallback(m *datatypes.MessageIn) datatypes.MessageOut {
	var err error
	var credentialsConfig confdata.AccountYaml
	// Unmarshal payload
	if err = json.Unmarshal([]byte(m.Payload), &credentialsConfig); err != nil {
		logrus.Errorf("Failed to unmarshal message payload: %v", err)
		return datatypes.MessageOut{Status: datatypes.ERROR, Msg: "Failed to unmarshal message payload."}
	}

	err = validate.Struct(credentialsConfig)
	if err != nil {
		logrus.Warning("Validation issue on credentials form, abort.")
		return datatypes.MessageOut{Status: datatypes.ERROR, Msg: "Validation issue on credentials form, please check given informations."}
	}

	config.Account = credentialsConfig
	conf.ExportConfig(config)
	return datatypes.MessageOut{Status: datatypes.SUCCESS, Msg: "Credentials settings successfully updated!"}
}

func quotasFormCallback(m *datatypes.MessageIn) datatypes.MessageOut {
	var err error
	var quotasConfig confdata.QuotasYaml
	// Unmarshal payload
	if err = json.Unmarshal([]byte(m.Payload), &quotasConfig); err != nil {
		logrus.Errorf("Failed to unmarshal message payload: %v", err)
		return datatypes.MessageOut{Status: datatypes.ERROR, Msg: "Failed to unmarshal message payload."}
	}

	err = validate.Struct(quotasConfig)
	if err != nil {
		logrus.Warning("Validation issue on quotas form, abort.")
		return datatypes.MessageOut{Status: datatypes.ERROR, Msg: "Validation issue on quotas form, please check given informations."}
	}

	config.Quotas = quotasConfig
	conf.ExportConfig(config)
	return datatypes.MessageOut{Status: datatypes.SUCCESS, Msg: "Quotas settings successfully updated!"}
}

func schedulerCallback(m *datatypes.MessageIn) datatypes.MessageOut {
	var err error
	var schedulerConfig confdata.ScheduleYaml
	// Unmarshal payload
	if err = json.Unmarshal([]byte(m.Payload), &schedulerConfig); err != nil {
		logrus.Errorf("Failed to unmarshal message payload: %v", err)
		return datatypes.MessageOut{Status: datatypes.ERROR, Msg: "Failed to unmarshal message payload."}
	}

	err = validate.Struct(schedulerConfig)
	if err != nil {
		logrus.Warning("Validation issue on scheduler form, abort.")
		return datatypes.MessageOut{Status: datatypes.ERROR, Msg: "Validation issue on scheduler form, please check given informations."}
	}

	config.Schedule = schedulerConfig
	conf.ExportConfig(config)
	return datatypes.MessageOut{Status: datatypes.SUCCESS, Msg: "Scheduler settings successfully updated!"}
}

func blacklistFormCallback(m *datatypes.MessageIn) datatypes.MessageOut {
	var err error
	var blacklistConfig confdata.BlacklistYaml
	// Unmarshal payload
	if err = json.Unmarshal([]byte(m.Payload), &blacklistConfig); err != nil {
		logrus.Errorf("Failed to unmarshal message payload: %v", err)
		return datatypes.MessageOut{Status: datatypes.ERROR, Msg: "Failed to unmarshal message payload."}
	}

	err = validate.Struct(blacklistConfig)
	if err != nil {
		logrus.Warning("Validation issue on blacklist form, abort.")
		return datatypes.MessageOut{Status: datatypes.ERROR, Msg: "Validation issue on blacklist form, please check given informations."}
	}

	config.Blacklist = blacklistConfig
	conf.ExportConfig(config)
	return datatypes.MessageOut{Status: datatypes.SUCCESS, Msg: "Blacklist settings successfully updated!"}
}

func dmBotFormCallback(m *datatypes.MessageIn) datatypes.MessageOut {
	var err error
	var dmConfig struct {
		DmTemplates       confdata.SplitStringSlice `json:"dmTemplates" validate:"required"`
		GreetingTemplate  string                    `json:"greetingTemplate"`
		GreetingActivated bool                      `json:"greetingActivation,string"`
		Activated         bool                      `json:"dmActivation,string"`
	}
	// Unmarshal payload
	if err = json.Unmarshal([]byte(m.Payload), &dmConfig); err != nil {
		logrus.Errorf("Failed to unmarshal message payload: %v", err)
		return datatypes.MessageOut{Status: datatypes.ERROR, Msg: "Failed to unmarshal message payload."}
	}

	err = validate.Struct(dmConfig)
	if err != nil {
		logrus.Warning("Validation issue on dm tool form, abort.")
		return datatypes.MessageOut{Status: datatypes.ERROR, Msg: "Validation issue on dm tool form, please check given informations."}
	}

	config.AutoDm.DmTemplates = dmConfig.DmTemplates
	config.AutoDm.Greeting.Template = dmConfig.GreetingTemplate
	config.AutoDm.Greeting.Activated = dmConfig.GreetingActivated
	config.AutoDm.Activated = dmConfig.Activated
	conf.ExportConfig(config)
	return datatypes.MessageOut{Status: datatypes.SUCCESS, Msg: "Dm bot settings successfully updated!"}
}

func dmScrapperFormCallback(m *datatypes.MessageIn) datatypes.MessageOut {
	var err error
	var scrapperConfig confdata.ScrapperYaml
	// Unmarshal payload
	if err = json.Unmarshal([]byte(m.Payload), &scrapperConfig); err != nil {
		logrus.Errorf("Failed to unmarshal message payload: %v", err)
		return datatypes.MessageOut{Status: datatypes.ERROR, Msg: "Failed to unmarshal message payload."}
	}

	err = validate.Struct(scrapperConfig)
	if err != nil {
		logrus.Warning("Validation issue on scrapper form, abort.")
		return datatypes.MessageOut{Status: datatypes.ERROR, Msg: "Validation issue on scrapper form, please check given informations."}
	}

	config.SrcUsers = scrapperConfig
	conf.ExportConfig(config)
	return datatypes.MessageOut{Status: datatypes.SUCCESS, Msg: "Scrapper settings successfully updated!"}
}

func proxyFormCallback(m *datatypes.MessageIn) datatypes.MessageOut {
	var err error
	var proxyConfig proxy.Proxy
	// Unmarshal payload
	if err = json.Unmarshal([]byte(m.Payload), &proxyConfig); err != nil {
		logrus.Errorf("Failed to unmarshal message payload: %v", err)
		return datatypes.MessageOut{Status: datatypes.ERROR, Msg: "Failed to unmarshal message payload."}
	}

	err = validate.Struct(proxyConfig)
	if err != nil {
		logrus.Warning("Validation issue on proxy form, abort.")
		return datatypes.MessageOut{Status: datatypes.ERROR, Msg: "Validation issue on proxy form, please check given informations."}
	}

	config.Selenium.Proxy = proxyConfig
	conf.ExportConfig(config)
	return datatypes.MessageOut{Status: datatypes.SUCCESS, Msg: "Proxy settings successfully updated!"}
}

func launchDmBotCallback(m *datatypes.MessageIn) datatypes.MessageOut {
	var err error
	if err = conf.CheckConfigValidity(); err == nil {
		ctx, cancel = context.WithCancel(context.Background())
		go bot.LaunchBot(ctx)
		return datatypes.MessageOut{Status: datatypes.SUCCESS, Msg: "Dm bot successfully launched!"}
	}
	return datatypes.MessageOut{Status: datatypes.ERROR, Msg: err.Error()}
}

func stopDmBotCallback(m *datatypes.MessageIn) datatypes.MessageOut {
	if bot.ExitedCh != nil {
		cancel()
		res := <-bot.ExitedCh
		if res {
			return datatypes.MessageOut{Status: datatypes.SUCCESS, Msg: "Dm bot successfully stopped!"}
		}
		return datatypes.MessageOut{Status: datatypes.ERROR, Msg: "Error during bot stopping! Please restart IGopher"}
	}
	return datatypes.MessageOut{Status: datatypes.ERROR, Msg: "Bot is in the initialization phase, please wait before trying to stop it."}
}

func hotReloadCallback(m *datatypes.MessageIn) datatypes.MessageOut {
	if automation.BotStruct.Running {
		if bot.HotReloadCh != nil {
			bot.HotReloadCh <- true
			res := <-bot.HotReloadCh
			if res {
				return datatypes.MessageOut{Status: datatypes.SUCCESS, Msg: "Bot hot reload successfully!"}
			}
			return datatypes.MessageOut{Status: datatypes.ERROR, Msg: "Error during bot hot reload! Please restart the bot"}
		}
		return datatypes.MessageOut{Status: datatypes.ERROR, Msg: "Bot is in the initialization phase, please wait before trying to hot reload it."}
	}
	return datatypes.MessageOut{Status: datatypes.ERROR, Msg: "Bot isn't running yet."}
}

func getLogsCallback(m *datatypes.MessageIn) datatypes.MessageOut {
	logs, err := logger.ParseLogsToString()
	if err != nil {
		logrus.Errorf("Can't parse logs: %v", err)
		return datatypes.MessageOut{Status: datatypes.ERROR, Msg: fmt.Sprintf("Can't parse logs: %v", err)}
	}
	logrus.Debug("Logs fetched successfully!")
	return datatypes.MessageOut{Status: datatypes.SUCCESS, Msg: logs}
}

func getConfigCallback(m *datatypes.MessageIn) datatypes.MessageOut {
	config, err := json.Marshal(config)
	if err != nil {
		logrus.Errorf("Can't parse config structure to Json: %v", err)
		return datatypes.MessageOut{Status: datatypes.ERROR, Msg: fmt.Sprintf("Can't parse config structure to Json: %v", err)}
	}
	logrus.Debug("Configuration structure successfully parsed!")
	return datatypes.MessageOut{Status: datatypes.SUCCESS, Msg: string(config)}
}
