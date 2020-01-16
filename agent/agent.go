package agent

import (
	"github.com/woshihot/go-lib/utils/log"
	"isesol.com/iport/options"
	"strings"
)

type MessageSource int

const (
	CLOUD MessageSource = iota
	LOCAL
)

func LocalPublish(publish *TopicPublish, logFlag ...string) {

	go func() {
		if LocalMqtt.Connected() {
			localPublishLock.Lock()
			defer localPublishLock.Unlock()
			if "" == publish.Topic.group {
				publish.Group("x")
			}
			topic, qos, payload := publish.Run()
			log.DF(flagToTag(logFlag...), "send-cloud topic = %s , msg = %s\n", topic, string(payload))
			LocalMqtt.Publish(topic, qos, payload)
		}
	}()

}

func CloudPublish(publish *TopicPublish, logFlag ...string) {
	go func() {
		if CloudMqtt.Connected() {
			cloudPublishLock.Lock()
			defer cloudPublishLock.Unlock()
			if "" == publish.Topic.group && "" == publish.Topic.channel {
				//使用初始化时的group和channel向上层push
				publish.Group(CloudGroup).Channel(options.GetOptions().BoxName)
			}
			topic, qos, payload := publish.Run()
			log.DF(flagToTag(logFlag...), "send-cloud topic = %s , msg = %s\n", topic, string(payload))
			CloudMqtt.Publish(topic, qos, payload)
		}
	}()
}

func flagToTag(logFlag ...string) string {
	tag := strings.Join(logFlag, "-")
	if "" != tag {
		tag = "[ " + tag + " ]"
	}
	return tag
}

type TopicPublish struct {
	Topic   *Topic
	Qos     byte
	Payload []byte
}

func TopicPublishCreate(topic string) *TopicPublish {
	return &TopicPublish{Topic: &Topic{name: topic}, Qos: byte(0), Payload: nil}
}

func (t *TopicPublish) Group(group string) *TopicPublish {
	t.Topic.Group(group)
	return t
}

func (t *TopicPublish) Channel(channel string) *TopicPublish {
	t.Topic.Channel(channel)
	return t
}

func (t *TopicPublish) SetQos(b byte) *TopicPublish {
	t.Qos = b
	return t
}

func (t *TopicPublish) SetPayload(payload []byte) *TopicPublish {
	t.Payload = payload
	return t
}

func (t *TopicPublish) Run() (string, byte, []byte) {
	return t.Topic.SubscribeValue(), t.Qos, t.Payload
}
