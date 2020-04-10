package controller

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/woshihot/go-lib/utils/time/timefmt"
	"isesol.com/iport/agent"
	"isesol.com/iport/controller/plugin"
	"isesol.com/iport/controller/topic"
	"isesol.com/iport/message"
	"isesol.com/iport/options"
	"isesol.com/iport/service"
	"strings"
)

func init() {
	plugin.RegisterPlugin(plugin.OnlinePlugin, NewOnlinePlugin)

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
	content := getOnlineContent(m)
	// 更新路由
	service.UpdateRouting(m.MachineNo, service.Routing{
		ChannelName:    agent.ParseTopic(m.Topic).GetChannel(),
		CloudAgentName: "cloudMqtt",
		MachineID:      m.MachineNo,
		MachineType:    content.MachineType,
		RegisterDate:   timefmt.CurrTimeGet(),
		Version:        content.Version,
		Token:          content.Token,
	})

	// 更新心跳
	service.UpdateHeartBeat(m.MachineNo)

	// 转发云端 或 直接回复
	if p.Option.SendCloud {
		t := topic.LocalMachineConnectionBegin
		if strings.Contains(m.Topic, topic.LocalBoxConnectionInitialize) {
			// 除盒子外一律转为topic:LocalMachineConnectionBegin
			t = topic.LocalBoxConnectionInitialize
		}
		p.Agent.CloudPublish(agent.TopicPublishCreate(t).SetPayload(m.ToPayload()), plugin.OnlinePlugin)
	} else {
		t := topic.MessageConfirmation
		if !strings.HasPrefix(m.Topic, topic.Message) {
			// 除设备外一律回复topic:MessageConfirmationFromAgent
			t = topic.MessageConfirmationFromAgent
		}
		p.Agent.LocalPublish(agent.TopicPublishCreate(t).SetPayload(m.ToPayload()).Channel(m.MachineNo), plugin.OnlinePlugin)
	}

	// 转发第三方
	if p.Option.SendThirdPart {
		service.SendThirdPart(m)
	}
	return true
}

func getOnlineContent(m message.Message) onlineContent {
	var content onlineContent
	err := jsoniter.Unmarshal([]byte(m.Content), &content)
	if nil != err {
		return onlineContent{}
	}

	if content.MachineType == "" {
		content.MachineType = agent.DefaultMachineType
	}
	return content
}

func (p OnlinePlugin) ExecCloudMessage(m message.Message) bool {
	if p.Option.BoxName == m.MachineNo {
		// 更新自身路由
		service.UpdateRouting(m.MachineNo, service.Routing{
			ChannelName:    m.MachineNo,
			CloudAgentName: p.Agent.Group,
			MachineID:      m.MachineNo,
			MachineType:    agent.BoxMachineType,
			RegisterDate:   timefmt.CurrTimeGet(),
			Version:        agent.IportVersion,
			Token:          p.Agent.CloudToken,
		})
		// 开始心跳

	}

	// 透传
	p.Agent.LocalPublish(agent.TopicPublishCreate(topic.Command).SetPayload(m.ToPayload()).Channel(m.MachineNo), plugin.OnlinePlugin)

	return true
}

// 处理本地和云端的0，0报文
func NewOnlinePlugin(a *agent.Agent, o options.Options) agent.Plugin {
	p := OnlinePlugin{&agent.Super{Option: o, Agent: a}}
	p.TypeOrder(agent.LOCAL, "0", "0").
		TypeOrder(agent.CLOUD, "0", "0")

	return p
}

type onlineContent struct {
	Version     service.RoutingVersion `json:"version"`
	Token       string                 `json:"token"`
	MachineType string                 `json:"machineType"`
}

func (o onlineContent) ToString() string {
	b, _ := jsoniter.Marshal(o)
	return string(b)
}
