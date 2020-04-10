package agent

import (
	"fmt"
	"isesol.com/iport/message"
	"isesol.com/iport/options"
	"strconv"
)

type Plugin interface {
	// 返回是否完全消费此报文，如是则不再继续执行子层级
	ExecLocalMessage(m message.Message) bool
	ExecCloudMessage(m message.Message) bool

	IsAccord(source MessageSource, topic string, m message.Message) bool
}

// 每个plugin都有唯一名称
type Super struct {
	Option     options.Options
	Agent      *Agent
	allows     map[MessageSource][]message.TypeOrder
	checkTopic func(topic string) bool
}

func (b *Super) ExecCloudMessage(m message.Message) bool {

	return false
}
func (b *Super) ExecLocalMessage(m message.Message) bool {
	return false
}

func (b *Super) IsAccord(source MessageSource, topic string, m message.Message) bool {
	fmt.Printf("isAccord source = %s, message type = %d,,message order =%d\n", source.ToString(), m.Type, m.Order)
	basicInit(b)
	topicCheck := true
	if b.checkTopic != nil {
		topicCheck = b.checkTopic(topic)
	}
	to := message.TypeOrder{message.Rule(strconv.Itoa(m.Type)), message.Rule(strconv.Itoa(m.Order))}

	return topicCheck && to.IsContains(b.allows[source])
}

func (b *Super) TypeOrder(source MessageSource, t, o string) *Super {
	basicInit(b)
	b.allows[source] = append(b.allows[source], message.TypeOrder{message.Rule(t), message.Rule(o)})
	return b
}

func (b *Super) CheckTopic(checkFunc func(topic string) bool) {
	b.checkTopic = checkFunc
}
func (b *Super) TypeOrders(source MessageSource, to ...message.TypeOrder) *Super {
	basicInit(b)
	b.allows[source] = append(b.allows[source], to...)
	return b
}

func (b *Super) Options(o options.Options) *Super {
	b.Option = o
	return b
}

func (b *Super) SetAgent(a *Agent) *Super {
	b.Agent = a
	return b
}

func basicInit(b *Super) {
	if nil == b.allows {
		b.allows = make(map[MessageSource][]message.TypeOrder)
	}
}
