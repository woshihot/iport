package mqtt

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/woshihot/go-lib/utils/log"
)

var (
	HandlerErrTag = "[MqttHandler-error]"
)

type OnConnectHandler struct {
	msgHandlers []MessageHandler
}

func (oc *OnConnectHandler) ExecTopic(topic string, payload []byte) {

	for _, h := range oc.msgHandlers {
		log.Df("msgTopic = %s , handler topic = %s", topic, h.Topic)
		if h.Topic == topic {

			h.Handler(topic, payload)
		}
	}
}

type MessageHandler struct {
	Topic   string
	Qos     byte
	Handler func(topic string, payload []byte)
}

func NewOnConnectHandler() *OnConnectHandler {
	return &OnConnectHandler{[]MessageHandler{}}
}

func (oc *OnConnectHandler) Create() mqtt.OnConnectHandler {
	return func(client mqtt.Client) {
		for _, handler := range oc.msgHandlers {
			var callback mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
				fmt.Printf("topic = %s,msg=%s\n", message.Topic(), message.Payload())
				go handler.Handler(message.Topic(), message.Payload())
			}
			token := client.Subscribe(handler.Topic, handler.Qos, callback)
			if token.Wait() && token.Error() != nil {
				log.EF(HandlerErrTag, "mqtt subscribe topic =%s error=%s\n", handler.Topic, token.Error())
			}
		}
	}
}

func (oc *OnConnectHandler) Handler(handler ...MessageHandler) *OnConnectHandler {
	oc.msgHandlers = append(oc.msgHandlers, handler...)
	return oc
}

func NewConnectionLostHandler(f func(err error)) mqtt.ConnectionLostHandler {
	return func(client mqtt.Client, e error) {
		f(e)
	}
}

func DefaultConnectionLostHandler(mqttName string) func(err error) {
	return func(err error) {
		log.EF(HandlerErrTag, "%s mqtt connect lost error=%s\n", mqttName, err.Error())
	}
}
