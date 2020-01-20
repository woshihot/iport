package options

import "time"

var iportOptions Options
var defaultOptions Options

func init() {
	defaultOptions = Options{
		OperationInfo: OperationInfo{
			SendCloud:                   true,
			SendThirdPart:               false,
			ThreeCodeMaxReconnectTimes:  3,
			ThreeCodeReconnectFrequency: 5 * time.Second,
		},
		AgentInfo: AgentInfo{
			ApiHost:       "https://api.isesol.com",
			CloudMqttAddr: "mqttAgent.isesol.com:1883",
			ApiServer:     "/agentServer",
		},
		BoxInfo: BoxInfo{
			BoxName:    "",
			BoxMac:     "",
			BoxLicense: "",
		},
	}
}

func NewOption() Options {
	return defaultOptions
}

type Options struct {
	OperationInfo // 配置iport库行为
	AgentInfo     // agent上下层地址信息
	BoxInfo       // 盒子的相关信息
}

type AgentInfo struct {
	ApiHost       string // 云端host
	ApiServer     string // 云端server(三码)
	CloudMqttAddr string // 云端mqtt地址
	LocalMqttAddr string // 本地mqtt地址
}

type BoxInfo struct {
	BoxName    string // 盒子号
	BoxMac     string // 盒子三码时mac地址
	BoxLicense string // 盒子加密狗号
}

type OperationInfo struct {
	SendCloud     bool // 是否转发报文到上层
	SendThirdPart bool // 是否发给三方

	ThreeCodeMaxRetryTimes  int           // 三码失败重试次数
	ThreeCodeRetryFrequency time.Duration // 三码失败重试频率

	MqttRetryFrequency time.Duration // mqtt连接重试频率
}
