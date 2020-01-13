package message

import (
	"github.com/json-iterator/go"
	"github.com/woshihot/go-lib/utils/log"
	"github.com/woshihot/go-lib/utils/mapstructure"
)

var (
	MessageErrTag = "[Message-error]"
)

type Message struct {
	ID        string `json:"id"`
	Type      int    `json:"type"`
	Order     int    `json:"order"`
	Content   string `json:"content"`
	MachineNo string `json:"machineNo"`
	Dest      string `json:"dest"`
	Source    string `json:"source"`
	Encode    bool   `json:"encode"`

	//来源TOPIC
	Topic string `json:"-"`
}

func (m *Message) ToJSON() string {
	return string(m.ToPayload()[:])
}

func (m *Message) ToPayload() []byte {
	payload, _ := jsoniter.Marshal(m)
	return payload
}

func (m *Message) SetType(t int) *Message {
	m.Type = t
	return m
}

func (m *Message) SetID(id string) *Message {
	m.ID = id
	return m
}

func (m *Message) SetOrder(o int) *Message {
	m.Order = o
	return m
}

func (m *Message) SetContent(c string) *Message {
	m.Content = c
	return m
}

func (m *Message) SetMachineNo(mNo string) *Message {
	m.MachineNo = mNo
	return m
}
func (m *Message) GetContent() (Content, error) {
	var content Content
	err := jsoniter.Unmarshal([]byte(m.Content), &content)
	if nil != err {
		log.EF(MessageErrTag, "message get content error = %s\n", err.Error())
		return content, err
	}
	return content, nil
}

func (m *Message) GetData() (map[string]interface{}, error) {
	content, err := m.GetContent()
	if nil != err {
		return nil, err
	}
	return content.GetData()
}

func (m *Message) FormatData(dPointer interface{}) error {
	content, err := m.GetContent()
	if nil != err {
		return err
	}
	return content.FormatData(dPointer)
}

type Content struct {
	CmdId string `json:"cmdId"`
	Data  string `json:"data"`
}

func (c Content) ToString() string {
	payload, _ := jsoniter.Marshal(c)
	return string(payload[:])
}
func (c *Content) GetData() (map[string]interface{}, error) {
	var data map[string]interface{}
	err := jsoniter.Unmarshal([]byte(c.Data), &data)
	if nil != err {
		log.EF(MessageErrTag, "content get data error = %s\n", err.Error())
		return nil, err
	}
	return data, nil
}

func (c *Content) FormatData(dPointer interface{}) error {
	mp, err := c.GetData()
	if err != nil {
		return err
	} else {
		if nil != mp {
			return mapstructure.MapToStructure(mp, dPointer)
		} else {
			return nil
		}
	}
}

func NewMessage(payload *[]byte) (*Message, error) {
	defer func() {
		if err := recover(); err != nil {
			log.EF(MessageErrTag, "payload deserialize error %s\n", string((*payload)[:]))
		}
	}()
	var message Message
	err := jsoniter.Unmarshal(*payload, &message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}
