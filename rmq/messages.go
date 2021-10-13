package rmq

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Type HandleMessageType
	Data []byte
}

func NewHandleMessageType(m interface{}) HandleMessageType {
	return HandleMessageType(fmt.Sprintf("%T", m))
}

func NewMessage(m interface{}) *Message {
	body, _ := json.Marshal(m)
	return &Message{
		Type: NewHandleMessageType(m),
		Data: body,
	}
}

func (m *Message) Bytes() []byte {
	b, _ := json.Marshal(m)
	return b
}
