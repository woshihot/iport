package agent

import (
	"github.com/woshihot/go-lib/do"
	"github.com/woshihot/go-lib/utils/log"
	"isesol.com/iport/mqtt"
	"isesol.com/iport/options"
	"strings"
	"sync"
	"time"
)

type MessageSource int

const (
	CLOUD MessageSource = iota
	LOCAL
)

func (m MessageSource) ToString() string {
	if CLOUD == m {
		return "cloud"
	} else {
		return "local"
	}
}

var (
	localPublishLock = new(sync.Mutex)
	cloudPublishLock = new(sync.Mutex)

	TagError = "[agent-error]"
	TagDebug = "[agent-debug]"
)

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

type Agent struct {
	cloudClient *mqtt.MqClient
	localClient *mqtt.MqClient

	status  agentStatus
	boxInfo options.BoxInfo
	PluginManager
	agentConnectInfo
}

func Initial(info options.BoxInfo) *Agent {
	a := &Agent{
		boxInfo: info,
	}
	a.Initial()
	return a
}

type agentConnectInfo struct {
	CloudToken string
	Group      string
}

func (a *Agent) IsLogged() bool {
	return a.status >= agentLogged
}

func (a *Agent) IsConnected() bool {
	return a.status >= agentConnected
}

func (a *Agent) Logged(token, group string) {
	a.CloudToken = token
	a.Group = group
	a.status = agentLogged
	log.DF(TagDebug, "%s login ,token =%s ,group =%s\n", a.boxInfo.BoxName, token, group)
}

func (a *Agent) LocalConnect(broker string, connectHandler mqtt.OnConnectHandler, reConnectFrequency time.Duration) {
	connectFunc := func() {
		cHandler := connectHandler.Create()
		connectionLostHandler := mqtt.NewConnectionLostHandler(mqtt.DefaultConnectionLostHandler("local"))
		o := mqtt.NewOptions().
			Broker(broker).
			ClientID(a.boxInfo.BoxName).
			Password("123456").
			Username("listener").
			ConnectHandler(&cHandler).
			ConnectLostHandler(&connectionLostHandler)
		client := mqtt.NewClient(o)

		err := client.Connect()
		if nil == err {
			a.localClient = client
			log.DF(TagDebug, "local mqtt connect\n")
		} else {
			log.EF(TagError, "local connect error = %s\n", err.Error())
		}
	}
	check := func() bool {
		return nil != a.localClient && a.localClient.Connected()
	}
	do.DoSthUntil(connectFunc, reConnectFrequency, check)
}

func (a *Agent) CloudConnect(broker, boxName string, connectHandler *mqtt.OnConnectHandler, willPayload TopicPublish, reConnectFrequency time.Duration) {
	connectFunc := func() {
		cHandler := connectHandler.Create()
		willPayload.Group(a.Group).Channel(boxName)
		connectionLostHandler := mqtt.NewConnectionLostHandler(mqtt.DefaultConnectionLostHandler("cloud"))
		o := mqtt.NewOptions().
			Broker(broker).
			ClientID(boxName).
			Password(a.CloudToken).
			Username(boxName).
			ConnectHandler(&cHandler).
			ConnectLostHandler(&connectionLostHandler).
			BinaryWill(willPayload.Run())
		client := mqtt.NewClient(o)
		err := client.Connect()
		if nil == err {
			a.cloudClient = client
			a.status = agentConnected
			log.DF(TagDebug, "cloud mqtt connect\n")
		} else {
			log.EF(TagError, "cloud connect error = %s\n", err.Error())
		}
	}

	checkFunc := func() bool {
		return a.IsConnected() && nil != a.cloudClient && a.cloudClient.Connected()
	}
	do.DoSthUntil(connectFunc, reConnectFrequency, checkFunc)

}

func (a *Agent) LocalPublish(publish *TopicPublish, logFlag ...string) {
	go func() {
		if a.localClient.Connected() {
			localPublishLock.Lock()
			defer localPublishLock.Unlock()
			if "" == publish.Topic.group {
				publish.Group("x")
			}
			topic, qos, payload := publish.Run()
			err := a.localClient.Publish(topic, qos, payload)
			if nil != err {
				log.EF(flagToTag(logFlag...), "send-local topic =%s error = %s", topic, err.Error())
			} else {
				log.DF(flagToTag(logFlag...), "send-local topic = %s , msg = %s\n", topic, string(payload))
			}
		} else {
			log.EF(flagToTag(logFlag...), "local mqtt not connect")
		}
	}()
}

func (a *Agent) CloudPublish(publish *TopicPublish, logFlag ...string) {
	go func() {
		if a.cloudClient.Connected() {
			cloudPublishLock.Lock()
			defer cloudPublishLock.Unlock()
			if "" == publish.Topic.group && "" == publish.Topic.channel {
				// 使用初始化时的group和channel向上层push
				publish.Group(a.Group).Channel(a.boxInfo.BoxName)
			}
			topic, qos, payload := publish.Run()
			err := a.cloudClient.Publish(topic, qos, payload)
			if nil != err {
				log.EF(flagToTag(logFlag...), "send-cloud topic =%s error = %s", topic, err.Error())
			} else {
				log.DF(flagToTag(logFlag...), "send-cloud topic = %s , msg = %s\n", topic, string(payload))
			}
		} else {
			log.EF(flagToTag(logFlag...), "cloud mqtt not connect")
		}
	}()
}

type agentStatus int

const (
	agentLogged = iota
	agentConnected
)
