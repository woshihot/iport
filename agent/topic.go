package agent

import (
	"github.com/woshihot/go-lib/utils/str"
	"isesol.com/iport/message"
	"isesol.com/iport/mqtt"
	"strings"
)

type Topic struct {
	name    string
	group   string
	channel string
}

func ParseTopic(s string) *Topic {

	ts := strings.Split(s, "/")
	if len(ts) == 3 {
		return &Topic{ts[0], ts[1], ts[2]}
	} else if len(ts) == 2 {
		return &Topic{ts[0], "", ts[1]}
	} else {
		return &Topic{s, "", ""}
	}

}

func (t *Topic) SubscribeValue() string {

	result := t.name
	if "" != t.group {
		result += "/" + t.group
	}
	if "" != t.channel {
		result += "/" + t.channel
	}

	return result
}

func (t *Topic) Group(g string) *Topic {
	t.group = g
	return t
}

func (t *Topic) Channel(c string) *Topic {
	t.channel = c
	return t
}

func (t *Topic) Absolutely() *Topic {
	t.group = ""
	t.channel = ""
	return t
}

func (t *Topic) Relative() *Topic {
	t.group = ""
	t.channel = "#"
	return t
}

func (t *Topic) Name() string {
	return t.name
}

func (t *Topic) GetChannel() string {
	return t.channel
}

func (t *Topic) GetGroup() string {
	return t.group
}

type MQTopic struct {
	Topic
	plugins     []string
	qos         byte
	handlerType HandlerType
}

func (t *MQTopic) Group(g string) *MQTopic {
	t.Topic.Group(g)
	return t
}

func (t *MQTopic) Channel(c string) *MQTopic {
	t.Topic.Channel(c)
	return t
}

func (t *MQTopic) Absolutely() *MQTopic {
	t.Topic.Absolutely()

	return t
}

func (t *MQTopic) Relative() *MQTopic {
	t.Topic.Relative()
	return t
}

func (t *MQTopic) Qos(qos byte) *MQTopic {
	t.qos = qos
	return t
}

func (t *MQTopic) Plugin(plugin ...string) *MQTopic {
	t.plugins = str.AppendSet(t.plugins, plugin...)
	return t
}

func (t *MQTopic) HandlerType(handlerType HandlerType) *MQTopic {
	t.handlerType = handlerType
	return t
}

func (t *MQTopic) messageHandler(source MessageSource, plugins []*Plugin) mqtt.MessageHandler {
	handler := mqtt.MessageHandler{Topic: t.SubscribeValue(), Qos: t.qos}
	switch t.handlerType {
	case Msg:
		handler.Handler = DefaultMsgHandler(t.Name(), source, plugins)
	case Payload:
		handler.Handler = PayloadHandler(t.Name(), source, plugins)
	default:
		handler.Handler = DefaultMsgHandler(t.Name(), source, plugins)
	}

	return handler
}

func PayloadHandler(topic string, source MessageSource, plugins []*Plugin) func(t string, payload []byte) {
	return func(t string, payload []byte) {
		msg := &message.Message{
			Topic:   t,
			Content: string(payload),
		}
		go func() {
			for _, plugin := range plugins {
				execute := executePlugin(plugin, source, topic, msg)
				if execute {
					break
				} else {
					continue
				}
			}
		}()
	}
}

func DefaultMsgHandler(topic string, source MessageSource, plugins []*Plugin) func(t string, payload []byte) {
	return func(t string, payload []byte) {
		msg, err := message.NewMessage(&payload)
		if nil == err && nil != msg {
			msg.Topic = t
			go func() {
				for _, plugin := range plugins {
					execute := executePlugin(plugin, source, topic, msg)
					if execute {
						break
					} else {
						continue
					}
				}
			}()
		}
	}
}

func executePlugin(plugin *Plugin, source MessageSource, topic string, msg *message.Message) bool {
	var execute bool
	if nil != msg {
		if (*plugin).IsAccord(source, topic, *msg) {
			switch source {
			case CLOUD:
				execute = (*plugin).ExecCloudMessage(*msg)
			case LOCAL:
				execute = (*plugin).ExecLocalMessage(*msg)
			}
		}
	}
	return execute
}
