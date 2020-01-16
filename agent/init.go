package agent

import (
	"isesol.com/iport/mqtt"
	"log"
	"sync"
)

var (
	AgentErrTag  = "[Agent-error]"
	LocalMqtt    mqtt.MqClient
	CloudMqtt    mqtt.MqClient
	CloudGroup   string
	CloudToken   string
	CloudChannel string
)

var (
	localPublishLock = new(sync.Mutex)
	cloudPublishLock = new(sync.Mutex)
)

func LocalInit(broker, clientId, username, pwd string) {
	PM.Local().TopicRelative("")
	connectHandler := PM.Handlers(LOCAL).Create()
	connectionLostHandler := mqtt.NewConnectionLostHandler(mqtt.DefaultConnectionLostHandler)
	options := mqtt.NewOptions().
		Broker(broker).
		ClientID(clientId).
		Password(pwd).
		Username(username).
		ConnectHandler(&connectHandler).
		ConnectLostHandler(&connectionLostHandler)
	client := mqtt.NewClient(options)

	err := client.Connect()
	if nil == err {
		LocalMqtt = *client
	} else {
		log.Panicf(AgentErrTag+"%s\n", err.Error())
	}
}

func CloudInit(broker, clientId, username, pwd, groupName, channel string, willPayload TopicPublish) {
	PM.Cloud().GroupChannel("", groupName, channel)
	CloudGroup = groupName
	CloudToken = pwd
	willPayload.Group(groupName).Channel(channel)

	connectHandler := PM.Handlers(CLOUD).Create()
	connectionLostHandler := mqtt.NewConnectionLostHandler(mqtt.DefaultConnectionLostHandler)

	options := mqtt.NewOptions().
		Broker(broker).
		ClientID(clientId).
		Password(pwd).
		Username(username).
		ConnectHandler(&connectHandler).
		ConnectLostHandler(&connectionLostHandler).
		BinaryWill(willPayload.Run())
	client := mqtt.NewClient(options)

	err := client.Connect()

	if nil == err {
		LocalMqtt = *client
	} else {
		log.Panicf(AgentErrTag+"%s\n", err.Error())
	}

}
