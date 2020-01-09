package agent

type MessageSource int

const (
	CLOUD MessageSource = iota
	LOCAL
)

func LocalPublish(publish *TopicPublish) {

	go func() {
		if LocalMqtt.Connected() {
			localPublishLock.Lock()
			defer localPublishLock.Unlock()
			if "" == publish.Topic.group {
				publish.Group("x")
			}
			LocalMqtt.Publish(publish.Run())
		}
	}()

}


func CloudPublish(publish *TopicPublish) {
	go func() {
		if CloudMqtt.Connected() {
			cloudPublishLock.Lock()
			defer cloudPublishLock.Unlock()
			if "" == publish.Topic.group && "" == publish.Topic.channel {
				//使用初始化时的group和channel向上层push
				publish.Group(CloudGroup).Channel(CloudChannel)
			}

			CloudMqtt.Publish(publish.Run())
		}
	}()
}


type TopicPublish struct {
	Topic   *Topic
	Qos     byte
	Payload []byte
}

func TopicPublishCreate(topic string) *TopicPublish {
	return &TopicPublish{Topic: &Topic{name: topic}, Qos: byte(0), Payload: nil}
}

func (t *TopicPublish) Group(group string) *TopicPublish {
	t.Topic.Group(group)
	return t
}

func (t *TopicPublish) Channel(channel string) *TopicPublish {
	t.Topic.Channel(channel)
	return t
}

func (t *TopicPublish) SetQos(b byte) *TopicPublish {
	t.Qos = b
	return t
}

func (t *TopicPublish) SetPayload(payload []byte) *TopicPublish {
	t.Payload = payload
	return t
}

func (t *TopicPublish) Run() (string, byte, []byte) {
	return t.Topic.SubscribeValue(), t.Qos, t.Payload
}
