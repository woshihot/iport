package agent

import (
	"isesol.com/iport/message"
	"isesol.com/iport/options"
	"strconv"
)

type Plugin interface {
	// 返回是否完全消费此报文，如是则不再继续执行子层级
	ExecLocalMessage(m message.Message) bool
	ExecCloudMessage(m message.Message) bool

	IsAccord(source MessageSource, m message.Message) bool
}

// 每个plugin都有唯一名称
type Super struct {
	Option options.Options
	agent  *Agent
	allows map[MessageSource][]message.TypeOrder
}

func (b *Super) ExecCloudMessage(m message.Message) bool {
	basicInit(b)
	return b.IsAccord(CLOUD, m)
}
func (b *Super) ExecLocalMessage(m message.Message) bool {
	basicInit(b)
	return b.IsAccord(LOCAL, m)
}

func (b *Super) IsAccord(source MessageSource, m message.Message) bool {
	basicInit(b)
	to := message.TypeOrder{message.Rule(strconv.Itoa(m.Type)), message.Rule(strconv.Itoa(m.Order))}
	return to.IsContains(b.allows[source])
}

func (b *Super) TypeOrder(source MessageSource, t, o string) *Super {
	basicInit(b)
	b.allows[source] = append(b.allows[source], message.TypeOrder{message.Rule(t), message.Rule(o)})
	return b
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

func basicInit(b *Super) {
	if nil == b {
		b = &Super{options.NewOption(), nil, make(map[MessageSource][]message.TypeOrder)}
	} else {
		if nil == b.allows {
			b.allows = make(map[MessageSource][]message.TypeOrder)
		}
	}
}
