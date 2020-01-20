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

type Iport struct {
	o      options.Options
	status IportStatus
	a      *agent.Agent
}

func Create(o options.Options) *Iport {
	return &Iport{
		o:      o,
		status: 0,
		a: &agent.Agent{
			nil, nil, 0, o.BoxInfo, nil, nil,
		},
	}
}

func (i Iport) Start() {
	controller.CreatePlugin(i.a, i.o)

	i.status = IportStart

	// 连接本地mqtt
	i.a.LocalConnect(i.o.LocalMqttAddr, *i.a.PluginManager.Handlers(agent.LOCAL), i.o.MqttRetryFrequency)
	// 连接云端mqtt
	if i.o.SendCloud {
		i.login()
		i.a.CloudConnect(i.o.CloudMqttAddr, i.o.BoxName, *i.a.Handlers(agent.CLOUD), controller.OfflinePublish(i.a.Group, i.a.CloudToken), i.o.MqttRetryFrequency)
		// 发送0，0

		// 开始心跳

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
		if "" != resp.Token && "" != resp.GroupName {
			i.a.Logged(resp.Token, resp.GroupName)
		}
	}
	checkFunc := func() bool {
		return i.IsStart() && i.a.IsLogged() && i.o.SendCloud
	}

	do.DoSthUntilTimes(logFunc, i.o.ThreeCodeMaxRetryTimes+1, i.o.ThreeCodeRetryFrequency, checkFunc)
}
