package matcher

import (
	"filter-core/internal/model/danmu"
)

type MatcherParam struct {
	Param string
	// 标识value类型，由MatcherInfo给出后透传返回
	BaseType  BaseMatcherType
	MatchMode BaseMatchMode
	// service层为string，在pkg层根据type解析
	Value interface{}
}

type MatcherInfo struct {
	Param    string          `json:"param,omitempty"`
	Type     BaseMatcherType `json:"type,omitempty"`
	ModeList []BaseMatchMode `json:"mode_list,omitempty"`
}

type baseMatcher interface {
	isBaseMatch(a interface{}) bool
}

type DanmuMatcher interface {
	IsDanmuMatch(dm *danmu.Danmu) bool
	GetMatcherInfo() []*MatcherInfo
}

func NewDanmuMatcher(t danmu.DanmuType, paramList []*MatcherParam) DanmuMatcher {
	switch t {
	case danmu.DanmuTypeDANMUMSG:
		return newDanmuMsgMatcher(paramList)
	default:
		return nil
	}
}
