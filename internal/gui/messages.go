package gui

import (
	"encoding/json"

	"github.com/asticode/go-astilectron"
	"github.com/go-playground/validator/v10"
	"github.com/hbollon/igopher"
	"github.com/sirupsen/logrus"
)

type SucessState string

const (
	SUCCESS SucessState = "Success"
	ERROR   SucessState = "Error"
)

var (
	config   igopher.BotConfigYaml
	validate = validator.New()
)

// MessageOut represents a message going out
type MessageOut struct {
	Status  SucessState `json:"status"`
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
			return MessageOut{Status: "Error during message reception"}
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
			return MessageOut{Status: ERROR}
		}
	})
}

func (m *MessageIn) resetGlobalSettingsCallback() MessageOut {
	config = igopher.ResetBotConfig()
	igopher.ExportConfig(config)
	return MessageOut{Status: SUCCESS, Msg: "Global configuration was successfully reseted!"}
}

func (m *MessageIn) credentialsFormCallback() MessageOut {
	var err error
	var credentialsConfig igopher.AccountYaml
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
	igopher.ExportConfig(config)
	return MessageOut{Status: SUCCESS, Msg: "Credentials settings successfully updated!"}
}

func (m *MessageIn) quotasFormCallback() MessageOut {
	var err error
	var quotasConfig igopher.QuotasYaml
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
	igopher.ExportConfig(config)
	return MessageOut{Status: SUCCESS, Msg: "Quotas settings successfully updated!"}
}

func (m *MessageIn) schedulerCallback() MessageOut {
	var err error
	var schedulerConfig igopher.ScheduleYaml
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
	igopher.ExportConfig(config)
	return MessageOut{Status: SUCCESS, Msg: "Scheduler settings successfully updated!"}
}

func (m *MessageIn) blacklistFormCallback() MessageOut {
	var err error
	var blacklistConfig igopher.BlacklistYaml
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
	igopher.ExportConfig(config)
	return MessageOut{Status: SUCCESS, Msg: "Blacklist settings successfully updated!"}
}

func (m *MessageIn) dmBotFormCallback() MessageOut {
	var err error
	var dmConfig igopher.AutoDmYaml
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

	config.AutoDm = dmConfig
	igopher.ExportConfig(config)
	return MessageOut{Status: SUCCESS, Msg: "Dm bot settings successfully updated!"}
}

func (m *MessageIn) dmScrapperFormCallback() MessageOut {
	var err error
	var scrapperConfig igopher.ScrapperYaml
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
	igopher.ExportConfig(config)
	return MessageOut{Status: SUCCESS, Msg: "Scrapper settings successfully updated!"}
}
