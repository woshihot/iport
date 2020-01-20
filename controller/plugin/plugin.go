package plugin

import (
	"isesol.com/iport/agent"
	"isesol.com/iport/options"
)

const (
	OnlinePlugin = "OnlinePlugin"
)

var (
	pm = make(map[string]PluginCreateFunc)
)

func RegisterPlugin(name string, createFunc PluginCreateFunc) {
	pm[name] = createFunc
}
func GetPlugins() map[string]PluginCreateFunc {
	return pm
}

type PluginCreateFunc func(agent *agent.Agent, options options.Options) agent.Plugin
