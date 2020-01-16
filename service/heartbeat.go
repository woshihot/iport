package service

import (
	"isesol.com/iport/redis"
	"time"
)

func UpdateHeartBeat(machineNo string) {
	if nil != redis.HeartBeat {
		now := time.Now().Unix()
		redis.HeartBeat.HSet(machineNo, now)
	}

}
