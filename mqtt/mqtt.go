package mqtt

import (
	"errors"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/woshihot/go-lib/utils/log"
)

var (
	MQTTErrTag = "[mqtt-error]"
)

type MqClient struct {
	mqtt.Client
	*Options
}

//是否连接
func (mq *MqClient) Connected() bool {

	var result bool
	if nil != mq.Client {
		result = mq.IsConnected()
	} else {
		result = false
	}
	return result
}

//连接
func (mq *MqClient) Connect() error {

	if token := mq.Client.Connect(); token.Wait() && token.Error() != nil {
		log.EF(MQTTErrTag, "mqttConnect fail clientId=[%s],error=%s\n", mq.ClientOptions.ClientID, token.Error())
		return token.Error()
	}
	return nil
}

//推送
func (mq *MqClient) Publish(topic string, qos byte, payload interface{}) error {
	if mq.Connected() {
		token := mq.Client.Publish(topic, qos, false, payload)
		if nil != token.Error() {
			return token.Error()
		}
	} else {
		return errors.New("mqtt not connect")
	}
	return nil
}

func NewClient(o *Options) *MqClient {
	return &MqClient{mqtt.NewClient(o.ClientOptions), o}
}
