package messages

type Message struct {
	Name   string
	Status string
	Data   interface{} // Data is a message
}

func NewMessage() Message {
	m := Message{}
	return m
}

func (m *Message) GetMessage(msg string) *Message {
	m.Data = msg
	return m
}
