package igopher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/asticode/go-astilectron"
	"github.com/go-playground/validator/v10"
	"github.com/hbollon/igopher/internal/proxy"
	"github.com/sirupsen/logrus"
)

type MsgState string

const (
	SUCCESS MsgState = "Success"
	ERROR   MsgState = "Error"
	INFO    MsgState = "Info"
)

var (
	window                          *astilectron.Window
	config                          BotConfigYaml
	validate                        = validator.New()
	reloadCh, hotReloadCh, exitedCh chan bool
	ctx                             context.Context
	cancel                          context.CancelFunc
)

// CallbackMap is a map of callback functions for each message
var CallbackMap = map[string]func(*MessageIn) MessageOut{
	"resetGlobalDefaultSettings":  (*MessageIn).resetGlobalSettingsCallback,
	"clearAllData":                (*MessageIn).clearDataCallback,
	"igCredentialsForm":           (*MessageIn).credentialsFormCallback,
	"quotasForm":                  (*MessageIn).quotasFormCallback,
	"schedulerForm":               (*MessageIn).schedulerCallback,
	"blacklistForm":               (*MessageIn).blacklistFormCallback,
	"dmSettingsForm":              (*MessageIn).dmBotFormCallback,
	"dmUserScrappingSettingsForm": (*MessageIn).dmScrapperFormCallback,
	"proxyForm":                   (*MessageIn).proxyFormCallback,
	"launchDmBot":                 (*MessageIn).launchDmBotCallback,
	"stopDmBot":                   (*MessageIn).stopDmBotCallback,
	"hotReloadBot":                (*MessageIn).hotReloadCallback,
	"getLogs":                     (*MessageIn).getLogsCallback,
	"getConfig":                   (*MessageIn).getConfigCallback,
}

// MessageOut represents a message for electron (going out)
type MessageOut struct {
	Status  MsgState    `json:"status"`
	Msg     string      `json:"msg"`
	Payload interface{} `json:"payload,omitempty"`
}

// MessageIn represents a message from electron (going in)
type MessageIn struct {
	Msg     string          `json:"msg"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// IsElectronRunning checks if electron is running
func IsElectronRunning() bool {
	return window != nil
}

// SendMessageToElectron will send a message to Electron Gui and execute a callback
// Callback function is optional
func SendMessageToElectron(msg MessageOut, callbacks ...astilectron.CallbackMessage) {
	if IsElectronRunning() {
		window.SendMessage(msg, callbacks...)
	}
}

// HandleMessages is handling function for incoming messages
func HandleMessages(w *astilectron.Window) {
	w.OnMessage(func(m *astilectron.EventMessage) interface{} {
		// Unmarshal
		var i MessageIn
		var err error
		if err = m.Unmarshal(&i); err != nil {
			logrus.Errorf("Unmarshaling message %+v failed: %v", *m, err)
			return MessageOut{Status: "Error during message reception"}
		}

		// Process message
		config = ImportConfig()
		if callback, ok := CallbackMap[i.Msg]; ok {
			return callback(&i)
		}
		logrus.Errorf("Unexpected message received: \"%s\"", i.Msg)
		return MessageOut{Status: ERROR, Msg: "Unknown error: Invalid message received"}
	})
	window = w
}

/* Callback functiosn to handle electron messages */

func (m *MessageIn) resetGlobalSettingsCallback() MessageOut {
	config = ResetBotConfig()
	ExportConfig(config)
	return MessageOut{Status: SUCCESS, Msg: "Global configuration was successfully reset!"}
}

func (m *MessageIn) clearDataCallback() MessageOut {
	if err := ClearData(); err != nil {
		return MessageOut{Status: ERROR, Msg: fmt.Sprintf("IGopher data clearing failed! Error: %v", err)}
	}
	return MessageOut{Status: SUCCESS, Msg: "IGopher data successfully cleared!"}
}

func (m *MessageIn) credentialsFormCallback() MessageOut {
	var err error
	var credentialsConfig AccountYaml
	// Unmarshal payload
	if err = json.Unmarshal([]byte(m.Payload), &credentialsConfig); err != nil {
		logrus.Errorf("Failed to unmarshal message payload: %v", err)
		return MessageOut{Status: ERROR, Msg: "Failed to unmarshal message payload."}
	}

	err = validate.Struct(credentialsConfig)
	if err != nil {
		logrus.Warning("Validation issue on credentials form, abort.")
		return MessageOut{Status: ERROR, Msg: "Validation issue on credentials form, please check given informations."}
	}

	config.Account = credentialsConfig
	ExportConfig(config)
	return MessageOut{Status: SUCCESS, Msg: "Credentials settings successfully updated!"}
}

func (m *MessageIn) quotasFormCallback() MessageOut {
	var err error
	var quotasConfig QuotasYaml
	// Unmarshal payload
	if err = json.Unmarshal([]byte(m.Payload), &quotasConfig); err != nil {
		logrus.Errorf("Failed to unmarshal message payload: %v", err)
		return MessageOut{Status: ERROR, Msg: "Failed to unmarshal message payload."}
	}

	err = validate.Struct(quotasConfig)
	if err != nil {
		logrus.Warning("Validation issue on quotas form, abort.")
		return MessageOut{Status: ERROR, Msg: "Validation issue on quotas form, please check given informations."}
	}

	config.Quotas = quotasConfig
	ExportConfig(config)
	return MessageOut{Status: SUCCESS, Msg: "Quotas settings successfully updated!"}
}

func (m *MessageIn) schedulerCallback() MessageOut {
	var err error
	var schedulerConfig ScheduleYaml
	// Unmarshal payload
	if err = json.Unmarshal([]byte(m.Payload), &schedulerConfig); err != nil {
		logrus.Errorf("Failed to unmarshal message payload: %v", err)
		return MessageOut{Status: ERROR, Msg: "Failed to unmarshal message payload."}
	}

	err = validate.Struct(schedulerConfig)
	if err != nil {
		logrus.Warning("Validation issue on scheduler form, abort.")
		return MessageOut{Status: ERROR, Msg: "Validation issue on scheduler form, please check given informations."}
	}

	config.Schedule = schedulerConfig
	ExportConfig(config)
	return MessageOut{Status: SUCCESS, Msg: "Scheduler settings successfully updated!"}
}

func (m *MessageIn) blacklistFormCallback() MessageOut {
	var err error
	var blacklistConfig BlacklistYaml
	// Unmarshal payload
	if err = json.Unmarshal([]byte(m.Payload), &blacklistConfig); err != nil {
		logrus.Errorf("Failed to unmarshal message payload: %v", err)
		return MessageOut{Status: ERROR, Msg: "Failed to unmarshal message payload."}
	}

	err = validate.Struct(blacklistConfig)
	if err != nil {
		logrus.Warning("Validation issue on blacklist form, abort.")
		return MessageOut{Status: ERROR, Msg: "Validation issue on blacklist form, please check given informations."}
	}

	config.Blacklist = blacklistConfig
	ExportConfig(config)
	return MessageOut{Status: SUCCESS, Msg: "Blacklist settings successfully updated!"}
}

func (m *MessageIn) dmBotFormCallback() MessageOut {
	var err error
	var dmConfig struct {
		DmTemplates       SplitStringSlice `json:"dmTemplates" validate:"required"`
		GreetingTemplate  string           `json:"greetingTemplate"`
		GreetingActivated bool             `json:"greetingActivation,string"`
		Activated         bool             `json:"dmActivation,string"`
	}
	// Unmarshal payload
	if err = json.Unmarshal([]byte(m.Payload), &dmConfig); err != nil {
		logrus.Errorf("Failed to unmarshal message payload: %v", err)
		return MessageOut{Status: ERROR, Msg: "Failed to unmarshal message payload."}
	}

	err = validate.Struct(dmConfig)
	if err != nil {
		logrus.Warning("Validation issue on dm tool form, abort.")
		return MessageOut{Status: ERROR, Msg: "Validation issue on dm tool form, please check given informations."}
	}

	config.AutoDm.DmTemplates = dmConfig.DmTemplates
	config.AutoDm.Greeting.Template = dmConfig.GreetingTemplate
	config.AutoDm.Greeting.Activated = dmConfig.GreetingActivated
	config.AutoDm.Activated = dmConfig.Activated
	ExportConfig(config)
	return MessageOut{Status: SUCCESS, Msg: "Dm bot settings successfully updated!"}
}

func (m *MessageIn) dmScrapperFormCallback() MessageOut {
	var err error
	var scrapperConfig ScrapperYaml
	// Unmarshal payload
	if err = json.Unmarshal([]byte(m.Payload), &scrapperConfig); err != nil {
		logrus.Errorf("Failed to unmarshal message payload: %v", err)
		return MessageOut{Status: ERROR, Msg: "Failed to unmarshal message payload."}
	}

	err = validate.Struct(scrapperConfig)
	if err != nil {
		logrus.Warning("Validation issue on scrapper form, abort.")
		return MessageOut{Status: ERROR, Msg: "Validation issue on scrapper form, please check given informations."}
	}

	config.SrcUsers = scrapperConfig
	ExportConfig(config)
	return MessageOut{Status: SUCCESS, Msg: "Scrapper settings successfully updated!"}
}

func (m *MessageIn) proxyFormCallback() MessageOut {
	var err error
	var proxyConfig proxy.Proxy
	// Unmarshal payload
	if err = json.Unmarshal([]byte(m.Payload), &proxyConfig); err != nil {
		logrus.Errorf("Failed to unmarshal message payload: %v", err)
		return MessageOut{Status: ERROR, Msg: "Failed to unmarshal message payload."}
	}

	err = validate.Struct(proxyConfig)
	if err != nil {
		logrus.Warning("Validation issue on proxy form, abort.")
		return MessageOut{Status: ERROR, Msg: "Validation issue on proxy form, please check given informations."}
	}

	config.Selenium.Proxy = proxyConfig
	ExportConfig(config)
	return MessageOut{Status: SUCCESS, Msg: "Proxy settings successfully updated!"}
}

func (m *MessageIn) launchDmBotCallback() MessageOut {
	var err error
	if err = CheckConfigValidity(); err == nil {
		ctx, cancel = context.WithCancel(context.Background())
		go launchBot(ctx)
		return MessageOut{Status: SUCCESS, Msg: "Dm bot successfully launched!"}
	}
	return MessageOut{Status: ERROR, Msg: err.Error()}
}

func (m *MessageIn) stopDmBotCallback() MessageOut {
	if exitedCh != nil {
		cancel()
		res := <-exitedCh
		if res {
			return MessageOut{Status: SUCCESS, Msg: "Dm bot successfully stopped!"}
		}
		return MessageOut{Status: ERROR, Msg: "Error during bot stopping! Please restart IGopher"}
	}
	return MessageOut{Status: ERROR, Msg: "Bot is in the initialization phase, please wait before trying to stop it."}
}

func (m *MessageIn) hotReloadCallback() MessageOut {
	if BotStruct.running {
		if hotReloadCh != nil {
			hotReloadCh <- true
			res := <-hotReloadCh
			if res {
				return MessageOut{Status: SUCCESS, Msg: "Bot hot reload successfully!"}
			}
			return MessageOut{Status: ERROR, Msg: "Error during bot hot reload! Please restart the bot"}
		}
		return MessageOut{Status: ERROR, Msg: "Bot is in the initialization phase, please wait before trying to hot reload it."}
	}
	return MessageOut{Status: ERROR, Msg: "Bot isn't running yet."}
}

func (m *MessageIn) getLogsCallback() MessageOut {
	logs, err := parseLogsToString()
	if err != nil {
		logrus.Errorf("Can't parse logs: %v", err)
		return MessageOut{Status: ERROR, Msg: fmt.Sprintf("Can't parse logs: %v", err)}
	}
	logrus.Debug("Logs fetched successfully!")
	return MessageOut{Status: SUCCESS, Msg: logs}
}

func (m *MessageIn) getConfigCallback() MessageOut {
	config, err := json.Marshal(config)
	if err != nil {
		logrus.Errorf("Can't parse config structure to Json: %v", err)
		return MessageOut{Status: ERROR, Msg: fmt.Sprintf("Can't parse config structure to Json: %v", err)}
	}
	logrus.Debug("Configuration structure successfully parsed!")
	return MessageOut{Status: SUCCESS, Msg: string(config)}
}
