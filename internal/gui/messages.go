package gui

import (
	"encoding/json"

	"github.com/asticode/go-astilectron"
	"github.com/go-playground/validator/v10"
	"github.com/hbollon/igopher"
	"github.com/sirupsen/logrus"
)

const (
	SUCCESS = "Success"
	FAIL    = "Fail"
)

var (
	config   igopher.BotConfigYaml
	validate = validator.New()
)

// MessageOut represents a message going out
type MessageOut struct {
	Msg     string      `json:"msg"`
	Payload interface{} `json:"payload,omitempty"`
}

// MessageIn represents a message going in
type MessageIn struct {
	Msg     string          `json:"msg"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

func handleMessages(w *astilectron.Window) {
	config = igopher.ImportConfig()
	w.OnMessage(func(m *astilectron.EventMessage) interface{} {
		// Unmarshal
		var i MessageIn
		var err error
		if err = m.Unmarshal(&i); err != nil {
			logrus.Errorf("Unmarshaling message %+v failed: %v", *m, err)
			return MessageOut{Msg: "Error during message reception"}
		}

		// Process message
		switch i.Msg {
		case "resetGlobalDefaultSettings":
			return i.resetGlobalSettingsCallback()

		case "igCredentialsForm":
			return i.credentialsFormCallback()

		case "quotasForm":
			return i.quotasFormCallback()

		case "schedulerForm":
			return i.schedulerCallback()

		case "blacklistForm":
			return i.blacklistFormCallback()

		case "dmSettingsForm":
			return i.dmBotFormCallback()

		case "dmUserScrappingSettingsForm":
			return i.dmScrapperFormCallback()

		default:
			logrus.Error("Unexpected message received.")
			return MessageOut{Msg: FAIL}
		}
	})
}

func (m *MessageIn) resetGlobalSettingsCallback() MessageOut {
	config = igopher.ResetBotConfig()
	igopher.ExportConfig(config)
	return MessageOut{Msg: SUCCESS}
}

func (m *MessageIn) credentialsFormCallback() MessageOut {
	var err error
	var credentialsConfig igopher.AccountYaml
	// Unmarshal payload
	if err = json.Unmarshal([]byte(m.Payload), &credentialsConfig); err != nil {
		logrus.Errorf("Failed to unmarshal message payload: %v", err)
		return MessageOut{Msg: FAIL}
	}

	err = validate.Struct(credentialsConfig)
	if err != nil {
		logrus.Warning("Validation issue on credentials form, abort.")
		return MessageOut{Msg: FAIL}
	}

	config.Account = credentialsConfig
	igopher.ExportConfig(config)
	return MessageOut{Msg: SUCCESS}
}

func (m *MessageIn) quotasFormCallback() MessageOut {
	var err error
	var quotasConfig igopher.QuotasYaml
	// Unmarshal payload
	if err = json.Unmarshal([]byte(m.Payload), &quotasConfig); err != nil {
		logrus.Errorf("Failed to unmarshal message payload: %v", err)
		return MessageOut{Msg: FAIL}
	}

	err = validate.Struct(quotasConfig)
	if err != nil {
		logrus.Warning("Validation issue on quotas form, abort.")
		return MessageOut{Msg: FAIL}
	}

	config.Quotas = quotasConfig
	igopher.ExportConfig(config)
	return MessageOut{Msg: SUCCESS}
}

func (m *MessageIn) schedulerCallback() MessageOut {
	var err error
	var schedulerConfig igopher.ScheduleYaml
	// Unmarshal payload
	if err = json.Unmarshal([]byte(m.Payload), &schedulerConfig); err != nil {
		logrus.Errorf("Failed to unmarshal message payload: %v", err)
		return MessageOut{Msg: FAIL}
	}

	err = validate.Struct(schedulerConfig)
	if err != nil {
		logrus.Warning("Validation issue on scheduler form, abort.")
		return MessageOut{Msg: FAIL}
	}

	config.Schedule = schedulerConfig
	igopher.ExportConfig(config)
	return MessageOut{Msg: SUCCESS}
}

func (m *MessageIn) blacklistFormCallback() MessageOut {
	var err error
	var blacklistConfig igopher.BlacklistYaml
	// Unmarshal payload
	if err = json.Unmarshal([]byte(m.Payload), &blacklistConfig); err != nil {
		logrus.Errorf("Failed to unmarshal message payload: %v", err)
		return MessageOut{Msg: FAIL}
	}

	err = validate.Struct(blacklistConfig)
	if err != nil {
		logrus.Warning("Validation issue on blacklist form, abort.")
		return MessageOut{Msg: FAIL}
	}

	config.Blacklist = blacklistConfig
	igopher.ExportConfig(config)
	return MessageOut{Msg: SUCCESS}
}

func (m *MessageIn) dmBotFormCallback() MessageOut {
	var err error
	var dmConfig igopher.AutoDmYaml
	// Unmarshal payload
	if err = json.Unmarshal([]byte(m.Payload), &dmConfig); err != nil {
		logrus.Errorf("Failed to unmarshal message payload: %v", err)
		return MessageOut{Msg: FAIL}
	}

	err = validate.Struct(dmConfig)
	if err != nil {
		logrus.Warning("Validation issue on dm tool form, abort.")
		return MessageOut{Msg: FAIL}
	}

	config.AutoDm = dmConfig
	igopher.ExportConfig(config)
	return MessageOut{Msg: SUCCESS}
}

func (m *MessageIn) dmScrapperFormCallback() MessageOut {
	var err error
	var scrapperConfig igopher.ScrapperYaml
	// Unmarshal payload
	if err = json.Unmarshal([]byte(m.Payload), &scrapperConfig); err != nil {
		logrus.Errorf("Failed to unmarshal message payload: %v", err)
		return MessageOut{Msg: FAIL}
	}

	err = validate.Struct(scrapperConfig)
	if err != nil {
		logrus.Warning("Validation issue on blacklist form, abort.")
		return MessageOut{Msg: FAIL}
	}

	config.SrcUsers = scrapperConfig
	igopher.ExportConfig(config)
	return MessageOut{Msg: SUCCESS}
}
