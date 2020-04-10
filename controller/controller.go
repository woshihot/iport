package controller

import (
	"isesol.com/iport/agent"
	"isesol.com/iport/controller/plugin"
	"isesol.com/iport/controller/topic"
	"isesol.com/iport/message"
	"isesol.com/iport/options"
)

func CreatePlugin(a *agent.Agent, o options.Options) {

	pluginMap := plugin.GetPlugins()
	for name, createFunc := range pluginMap {
		p := createFunc(a, o)
		a.RegisterPlugin(name, &p)
	}

}

func AddLocalControl(a *agent.Agent, group, channel string) {
	a.Local().
		Topic(topic.Message).
		// 收到设备发出的上线事件
		Plugin(plugin.OnlinePlugin).
		HandlerType(agent.Msg)

	a.Local().
		Topic(topic.LocalMachineConnectionBegin).
		// 收到设备通过盒子发出的上线事件
		Plugin(plugin.OnlinePlugin).
		HandlerType(agent.Msg)

	a.Local().
		Topic(topic.LocalBoxConnectionInitialize).
		// 收到盒子的上线事件
		Plugin(plugin.OnlinePlugin).
		HandlerType(agent.Msg)
	a.Local().GroupChannel("", group, channel)
}

func AddCloudControl(a *agent.Agent, group, channel string) {
	a.Cloud().
		Topic(topic.MessageConfirmationFromAgent).
		// 盒子发的 0，0 回复
		Plugin(plugin.OnlinePlugin).
		HandlerType(agent.Msg)

	a.Cloud().
		Topic(topic.Command).
		// 云端发的 0, 0 回复
		Plugin(plugin.OnlinePlugin).
		HandlerType(agent.Msg)
	a.Cloud().GroupChannel("", group, channel)
}

func OfflinePublish(group, channel string) agent.TopicPublish {
	return *agent.TopicPublishCreate(topic.LocalBoxConnectionLost).
		Channel(channel).
		Group(group).
		SetPayload(
			(&message.Message{
				MachineNo: channel,
				Encode:    false,
				ID:        "",
				Content:   "",
				Type:      -1,
				Order:     -1,
			}).ToPayload())
}

func OnlinePublish(token, group, channel string) agent.TopicPublish {
	return *agent.TopicPublishCreate(topic.LocalBoxConnectionInitialize).
		Channel(channel).
		Group(group).
		SetPayload(
			(&message.Message{
				MachineNo: channel,
				Encode:    false,
				ID:        "",
				Content: onlineContent{
					Token:       token,
					Version:     agent.IportVersion,
					MachineType: agent.BoxMachineType,
				}.ToString(),
				Type:  0,
				Order: 0,
			}).ToPayload())
}
