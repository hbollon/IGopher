package datatypes

import "encoding/json"

type MsgState string

const (
	SUCCESS MsgState = "Success"
	ERROR   MsgState = "Error"
	INFO    MsgState = "Info"
)

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

func (m *MessageIn) Callback(callback func(m *MessageIn) MessageOut) MessageOut {
	return callback(m)
}
