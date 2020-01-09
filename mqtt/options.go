package mqtt

import (
	"github.com/eclipse/paho.mqtt.golang"
	"time"
)

type Options struct {
	*mqtt.ClientOptions
}

func NewOptions() *Options {
	o := &Options{mqtt.NewClientOptions()}
	return o.
		AutoReconnect(false).
		CleanSession(true).
		MaxReconnectInterval(3 * time.Second).
		ConnectTimeout(30 * time.Second).
		PingTimeout(10 * time.Second).
		KeepAlive(60 * time.Second).
		WriteTimeout(10 * time.Second)
}

func (o *Options) Broker(b string) *Options {
	o.AddBroker(b)
	return o
}

func (o *Options) ClientID(id string) *Options {
	o.SetClientID(id)
	return o
}

func (o *Options) Username(uName string) *Options {
	o.SetUsername(uName)
	return o
}

func (o *Options) Password(pwd string) *Options {
	o.SetPassword(pwd)
	return o
}

func (o *Options) ConnectHandler(handler *mqtt.OnConnectHandler) *Options {
	o.SetOnConnectHandler(*handler)
	return o
}

func (o *Options) ConnectLostHandler(handler *mqtt.ConnectionLostHandler) *Options {
	o.SetConnectionLostHandler(*handler)
	return o
}

func (o *Options) AutoReconnect(b bool) *Options {
	o.SetAutoReconnect(b)
	return o
}

func (o *Options) CleanSession(b bool) *Options {
	o.SetCleanSession(b)
	return o
}

func (o *Options) MaxReconnectInterval(t time.Duration) *Options {
	o.SetMaxReconnectInterval(t)
	return o
}

func (o *Options) ConnectTimeout(t time.Duration) *Options {
	o.SetConnectTimeout(t)
	return o
}

func (o *Options) PingTimeout(t time.Duration) *Options {
	o.SetPingTimeout(t)
	return o
}

func (o *Options) KeepAlive(t time.Duration) *Options {
	o.SetKeepAlive(t)
	return o
}

func (o *Options) WriteTimeout(t time.Duration) *Options {
	o.SetWriteTimeout(t)
	return o
}

func (o *Options) BinaryWill(topic string, qos byte, payload []byte) *Options {
	o.SetBinaryWill(topic, payload, qos, false)
	return o
}

//
//func (o *Options) Init(broker, id, uName, pwd string, onConnectHandler *mqtt.OnConnectHandler, onConnectLostHandler *mqtt.ConnectionLostHandler) *Options {
//	return o.Broker(broker).
//		ClientID(id).
//		Username(uName).
//		Password(pwd).
//		ConnectHandler(onConnectHandler).
//		ConnectLostHandler(onConnectLostHandler).
//		AutoReconnect(false).
//		CleanSession(true).
//		MaxReconnectInterval(3 * device.Second).
//		ConnectTimeout(30 * device.Second).
//		PingTimeout(10 * device.Second).
//		KeepAlive(60 * device.Second).
//		WriteTimeout(10 * device.Second)
//}
