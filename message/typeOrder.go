package message

import (
	"strconv"
	"strings"
)

type Rule string

const (
	more = "+"
	less = "-"
	and  = "-"
	or   = ","
)

func (r Rule) isMore(code int) bool {
	rs := string(r)
	im := strings.HasSuffix(rs, more) && strings.Count(rs, more) == 1
	if im {
		imCode, err := strconv.Atoi(strings.TrimSuffix(rs, more))
		if nil == err {
			return code > imCode
		}
	}
	return false
}

func (r Rule) isLess(code int) bool {
	rs := string(r)
	il := strings.HasSuffix(rs, less) && strings.Count(rs, less) == 1
	if il {
		ilCode, err := strconv.Atoi(strings.TrimSuffix(rs, less))
		if nil == err {
			return code < ilCode
		}
	}
	return false
}

func (r Rule) isInterval(code int) bool {
	rs := string(r)
	ii := strings.Count(rs, and) == 1 && !strings.HasSuffix(rs, and) && !strings.HasPrefix(rs, and)
	if ii {
		index := strings.Index(rs, and)
		lessCode, err := strconv.Atoi(rs[:index])
		if nil != err {
			return false
		}
		moreCode, err := strconv.Atoi(rs[index+1:])
		if nil != err {
			return false
		}
		return code > lessCode && code < moreCode
	}
	return strings.Count(rs, and) == 1 && !strings.HasSuffix(rs, and)
}

func (r Rule) Same(to string) bool {
	if "" == r || string(r) == to {
		return true
	}
	toCode, err := strconv.Atoi(to)
	if err != nil {
		return false
	}
	var result bool
	rules := r.Split()
	for _, rl := range rules {
		if string(rl) == to || rl.isMore(toCode) || rl.isLess(toCode) || rl.isInterval(toCode) {
			result = true
			break
		}
	}
	return result

}

func (r Rule) Split() []Rule {
	rls := strings.Split(string(r), or)
	var rules []Rule
	for _, rl := range rls {
		rules = append(rules, Rule(rl))
	}
	return rules
}

type TypeOrder struct {
	Type  Rule
	Order Rule
}

func (t TypeOrder) Accord(typeOrder TypeOrder) bool {

	return t.Type.Same(string(typeOrder.Type)) && t.Order.Same(string(typeOrder.Order))
}

func (t TypeOrder) IsContains(typeOrders []TypeOrder) bool {

	result := false
	if nil != typeOrders {
		for _, value := range typeOrders {
			if t.Accord(value) {
				result = true
			}
		}
	}

	return result
}
