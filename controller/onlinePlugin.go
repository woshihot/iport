package controller

import (
	"isesol.com/iport/agent"
	"isesol.com/iport/controller/plugin"
	"isesol.com/iport/controller/topic"
	"isesol.com/iport/message"
	"isesol.com/iport/options"
	"isesol.com/iport/service"
	"strings"
)

func init() {
	p := NewOnlinePlugin()
	agent.PM.RegisterPlugin(plugin.OnlinePlugin, &p)

}

// 处理上线报文
type OnlinePlugin struct {
	*agent.Super
}

/**
topic :
	设备直连：			Message
	盒子连接：			LocalBoxConnectInitialize
	设备通过盒子连接：	LocalMachineConnectionBegin
*/
func (p OnlinePlugin) ExecLocalMessage(m message.Message) bool {
	//更新路由
	service.UpdateRouting()

	//更新心跳
	service.UpdateHeartBeat()

	//转发云端 或 直接回复
	if options.GetOptions().SendCloud {
		t := topic.LocalMachineConnectionBegin
		if strings.Contains(m.Topic, topic.LocalBoxConnectionInitialize) {
			//除盒子外一律转为topic:LocalMachineConnectionBegin
			t = topic.LocalBoxConnectionInitialize
		}
		agent.CloudPublish(agent.TopicPublishCreate(t).SetPayload(m.ToPayload()))
	} else {
		t := topic.MessageConfirmation
		if !strings.HasPrefix(m.Topic, topic.Message) {
			//除设备外一律回复topic:MessageConfirmationFromAgent
			t = topic.MessageConfirmationFromAgent
		}
		agent.LocalPublish(agent.TopicPublishCreate(t).SetPayload(m.ToPayload()).Channel(m.MachineNo))
	}

	//转发mes

	return true
}


func (p OnlinePlugin) ExecCloudMessage(m message.Message) bool {
	return true
}

//处理本地和云端的0，0报文
func NewOnlinePlugin() agent.Plugin {
	p := new(OnlinePlugin)
	p.TypeOrder(agent.LOCAL, "0", "0")
	p.TypeOrder(agent.CLOUD, "0", "0")
	return p
}
