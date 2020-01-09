package agent

import (
	"isesol.com/iport/mqtt"
)

type HandlerType int

const (
	Msg HandlerType = iota
)

var (
	PM = PluginManager{make(map[string]*Plugin), make(map[string]*MQTopic), make(map[string]*MQTopic)}
)
//插件管理类
type PluginManager struct {
	plugins map[string]*Plugin
	l       TopicInterface
	c       TopicInterface
}

//注册插件
func (p *PluginManager) RegisterPlugin(pluginName string, plugin *Plugin) *PluginManager {
	if nil == p.plugins {
		p.plugins = make(map[string]*Plugin)
	}
	p.plugins[pluginName] = plugin
	return p
}

//获取插件
func (p *PluginManager) getPlugin(key string) (*Plugin, bool) {
	v, ok := p.plugins[key]
	return v, ok
}

//通过topic获取调用的插件
func (p *PluginManager) getPluginsByTopic(topic string,source MessageSource) []*Plugin {
	var result []*Plugin
	t := p.topicInterface(source)

	if nil != t {
		keys := t.getPlugins(topic)
		if len(keys) > 0 {
			for _, k := range keys {
				plugin, ok := p.getPlugin(k)
				if ok && nil != plugin {
					result = append(result, plugin)
				}
			}
		}
	}
	return result
}
//获取local topic-plugin转换接口
func (p *PluginManager) Local() TopicInterface {
	if nil != p.l {
		return p.l
	}
	p.l = make(map[string]*MQTopic)
	return p.l
}

//获取cloud topic-plugin转换接口
func (p *PluginManager) Cloud() TopicInterface {
	if nil != p.c {
		return p.c
	}
	p.c = make(TopicInterface)
	return p.c
}
func (p *PluginManager) topicInterface(source MessageSource) TopicInterface {
	var t TopicInterface
	switch source {
	case CLOUD:
		t = p.Cloud()
	case LOCAL:
		t = p.Local()
	default:
		t = make(TopicInterface)
	}
	return t
}

//将注册信息转换为mqtt的订阅handler
func (p *PluginManager) Handlers(source MessageSource) *mqtt.OnConnectHandler {
	msgHandler := p.topicInterface(source).Topics(source)
	connectHandler := mqtt.NewOnConnectHandler()
	connectHandler.Handler(msgHandler...)

	return connectHandler
}

type TopicInterface map[string]*MQTopic

func (t TopicInterface) Topic(name string) *MQTopic {
	topic, ok := t[name]
	if ok {
		return topic
	} else {
		t[name] = &MQTopic{Topic: Topic{name,"",""}, qos: 0, handlerType: Msg}
		return t[name]
	}
}
func (t TopicInterface) TopicRelative(name string) TopicInterface {
	return t.GroupChannel(name, "", "#")
}
func (t TopicInterface) TopicAbsolutely(name string) TopicInterface {
	return t.GroupChannel(name, "", "")
}
func (t TopicInterface) GroupChannel(name, group, channel string) TopicInterface {
	if "" == name {
		for _, v := range t {
			v.Group(group).Channel(channel)
		}
	} else {
		t.Topic(name).Group(group).Channel(channel)
	}
	return t
}


func (t TopicInterface) Topics(source MessageSource) []mqtt.MessageHandler {
	var result []mqtt.MessageHandler
	for _, v := range t {
		result = append(result, v.messageHandler(source))
	}
	return result
}


func (t TopicInterface) AddTopic(topic Topic) TopicInterface {
	changeTopic(t, topic, nil, nil)
	return t
}

func (t TopicInterface) TopicHandler(topic Topic, handlerType HandlerType) TopicInterface {
	changeTopic(t, topic, nil, &handlerType)
	return t
}

func (t TopicInterface) TopicPlugins(topic Topic, plugins ...string) TopicInterface {
	changeTopic(t, topic, nil, nil, plugins...)
	return t
}

func (t TopicInterface) TopicQos(topic Topic, qos byte) TopicInterface {
	changeTopic(t, topic, &qos, nil)
	return t
}

func changeTopic(topics map[string]*MQTopic, t Topic, qos *byte, handlerType *HandlerType, plugins ...string) {
	topic, ok := topics[t.Name()]
	if !ok {
		topic = &MQTopic{Topic: t, qos: byte(0), handlerType: Msg, plugins: []string{}}
	}
	if nil != handlerType {
		topic.handlerType = *handlerType
	}
	if nil != qos {
		topic.qos = *qos
	}
	if len(plugins) > 0 {
		topic.plugins = plugins
	}
	topics[t.Name()] = topic
}

func (t TopicInterface) getPlugins(name string) []string {
	topic, ok := t[name]
	if ok {
		return topic.plugins
	}
	return []string{}
}
