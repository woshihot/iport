package topic

// topic
const (
	Message                      = "Message"                      //消息
	MessageConfirmation          = "MessageConfirmation"          //给设备的消息回复
	MessageConfirmationFromAgent = "MessageConfirmationFromAgent" //给盒子的消息回复

	LocalMachineConnectionBegin  = "LocalMachineConnectionBegin"  //设备登录转发到云端
	LocalBoxConnectionLost       = "LocalBoxConnectionLost"       //盒子掉线时发送给云端
	LocalBoxConnectionInitialize = "LocalBoxConnectionInitialize" //盒子连接上时发送给云端

	Command = "Command" //指令

)
