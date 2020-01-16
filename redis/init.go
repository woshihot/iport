package redis

var (
	routingKey   = "ROUTING"
	heartBeatKey = "heartBeat"

	lastMessageKey = "lastMessage_"
)

var (
	Routing     *Hash
	LastMessage *Hash
	HeartBeat   *Hash
)

func Init(address string) {

	redisClient := NewClient(address)
	if nil != redisClient {

		Routing = redisClient.NewHash(routingKey)

		LastMessage = redisClient.NewHash(lastMessageKey)

		HeartBeat = redisClient.NewHash(heartBeatKey)
	}

}
