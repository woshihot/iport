package server

import (
	"isesol.com/iport/agent"
	"isesol.com/iport/controller/plugin"
	"isesol.com/iport/controller/topic"
)

func init() {

	agent.PM.
		Local().
		Topic(topic.Message).
		//收到设备发出的上线事件
		Plugin(plugin.OnlinePlugin).
		HandlerType(agent.Msg)

	agent.PM.
		Local().
		Topic(topic.LocalMachineConnectionBegin).
		//收到设备通过盒子发出的上线事件
		Plugin(plugin.OnlinePlugin).
		HandlerType(agent.Msg)

	agent.PM.
		Local().
		Topic(topic.LocalBoxConnectionInitialize).
		//收到盒子的上线事件
		Plugin(plugin.OnlinePlugin).
		HandlerType(agent.Msg)
}
