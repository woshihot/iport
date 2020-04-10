package iport

import (
	"github.com/woshihot/go-lib/do"
	"isesol.com/iport/agent"
	"isesol.com/iport/controller"
	"isesol.com/iport/options"
)

const (
	IportStart = iota
	IportStop
)

var TagDebug = "[iport-debug]"

type Iport struct {
	o      options.Options
	status IportStatus
	a      *agent.Agent
}

func Create(o options.Options) *Iport {
	return &Iport{
		o:      o,
		status: 0,
		a:      agent.Initial(o.BoxInfo),
	}
}

func (i Iport) PrintOptions() {

}

func (i Iport) Start() {
	i.PrintOptions()
	i.status = IportStart
	controller.CreatePlugin(i.a, i.o)
	controller.AddLocalControl(i.a, "x", "")
	// 连接本地mqtt
	i.a.LocalConnect(i.o.LocalMqttAddr, *i.a.PluginManager.Handlers(agent.LOCAL), i.o.MqttRetryFrequency)
	// 连接云端mqtt

	if i.o.SendCloud {
		i.login()
		controller.AddCloudControl(i.a, i.a.Group, i.o.BoxName)
		i.a.CloudConnect(i.o.CloudMqttAddr, i.o.BoxName, i.a.Handlers(agent.CLOUD), controller.OfflinePublish(i.a.Group, i.a.CloudToken), i.o.MqttRetryFrequency)
		// 发送0，0
		onlinePublish := controller.OnlinePublish(i.a.CloudToken, i.a.Group, i.o.BoxName)
		i.a.CloudPublish(&onlinePublish)
		// id 报文重发
	}
}

func (i Iport) IsStart() bool {
	return IportStart == i.status
}

type IportStatus int

func (i Iport) login() {
	url := i.o.ApiHost + i.o.ApiServer + ThirdCodeVerifyApi
	logFunc := func() {
		resp := cloudLogin(url, i.o.BoxName, i.o.BoxMac, i.o.BoxLicense)
		//TODO 登录判断有问题
		if "" != resp.Token && "" != resp.GroupName {
			i.a.Logged(resp.Token, resp.GroupName)
		}
	}
	checkFunc := func() bool {
		return i.IsStart() && i.a.IsLogged() && i.o.SendCloud
	}

	do.DoSthUntilTimes(logFunc, i.o.ThreeCodeMaxRetryTimes+1, i.o.ThreeCodeRetryFrequency, checkFunc)
}
