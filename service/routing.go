package service

import (
	jsoniter "github.com/json-iterator/go"
	"isesol.com/iport/redis"
	"strconv"
	"strings"
)

func UpdateRouting(field string, routing Routing) {
	if nil != redis.Routing {
		routingJSON, _ := jsoniter.Marshal(routing)
		redis.Routing.HSet(field, routingJSON)
	}
}

type Routing struct {
	ChannelName    string         `json:"channelName"`
	CloudAgentName string         `json:"cloudAgentName"`
	MachineID      string         `json:"machineId"`
	MachineType    string         `json:"machineType"`
	RegisterDate   string         `json:"registerDate"`
	Version        RoutingVersion `json:"version"`
	Token          string         `json:"token"`
}

type RoutingVersion string

func (v RoutingVersion) After(bv RoutingVersion) bool {
	if "" == v {
		return false
	}
	if "" == bv {
		return true
	}
	currentVersion := strings.Split(string(v), ".")
	compareVersion := strings.Split(string(bv), ".")
	lcuv := len(currentVersion)
	lcov := len(compareVersion)
	maxL := lcuv
	if lcov > maxL {
		maxL = lcov
	}
	for i := 0; i < maxL; i++ {
		var (
			cuvi = 0
			covi = 0
		)
		if i < lcuv {
			cuvi, _ = strconv.Atoi(currentVersion[i])
		}
		if i < lcov {
			covi, _ = strconv.Atoi(compareVersion[i])
		}
		if cuvi > covi {
			return true
		}
	}
	return false

}

func (v RoutingVersion) Before(bv RoutingVersion) bool {
	if "" == v {
		return false
	}
	if "" == bv {
		return true
	}
	currentVersion := strings.Split(string(v), ".")
	compareVersion := strings.Split(string(bv), ".")
	lcuv := len(currentVersion)
	lcov := len(compareVersion)
	maxL := lcuv
	if lcov > maxL {
		maxL = lcov
	}
	for i := 0; i < maxL; i++ {
		var (
			cuvi = 0
			covi = 0
		)
		if i < lcuv {
			cuvi, _ = strconv.Atoi(currentVersion[i])
		}
		if i < lcov {
			covi, _ = strconv.Atoi(compareVersion[i])
		}
		if cuvi < covi {
			return true
		}
	}
	return false

}

func (v RoutingVersion) Equal(bv RoutingVersion) bool {
	return v == bv
}
