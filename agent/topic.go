package agent

import (
	"github.com/woshihot/go-lib/utils/str"
	"isesol.com/iport/message"
	"isesol.com/iport/mqtt"
)

type Topic struct {
	name    string
	group   string
	channel string
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

func (t *MQTopic) messageHandler(source MessageSource) mqtt.MessageHandler {
	handler := mqtt.MessageHandler{Topic: t.SubscribeValue(), Qos: t.qos}
	switch t.handlerType {
	case Msg:
		handler.Handler = DefaultMsgHandler(t.Name(), source)
	case Payload:
		handler.Handler = PayloadHandler(t.Name(), source)
	default:
		handler.Handler = DefaultMsgHandler(t.Name(), source)
	}

	return handler
}

func PayloadHandler(topic string, source MessageSource) func(t string, payload []byte) {
	return func(t string, payload []byte) {
		plugins := PM.getPluginsByTopic(topic, source)
		msg := &message.Message{
			Topic:   t,
			Content: string(payload),
		}
		go func() {
			for _, plugin := range plugins {
				execute := executePlugin(plugin, source, msg)
				if execute {
					break
				} else {
					continue
				}
			}
		}()
	}
}

func DefaultMsgHandler(topic string, source MessageSource) func(t string, payload []byte) {
	return func(t string, payload []byte) {
		plugins := PM.getPluginsByTopic(topic, source)
		msg, err := message.NewMessage(&payload)
		if nil == err && nil != msg {
			msg.Topic = t
			go func() {
				for _, plugin := range plugins {
					execute := executePlugin(plugin, source, msg)
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

func executePlugin(plugin *Plugin, source MessageSource, msg *message.Message) bool {
	var execute bool
	if nil != msg {
		if (*plugin).IsAccord(source, *msg) {
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
