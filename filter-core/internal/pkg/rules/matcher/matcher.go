package matcher

import (
	danmu2 "filter-core/internal/model/danmu"
)

type MatcherParams struct {
	List []*MatcherParamItem
}

type MatcherParamItem struct {
	Param string
	// 标识value类型，由MatcherInfo给出后透传返回
	Type BaseMatcherType
	Mode BaseMatchMode
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
	IsDanmuMatch(dm *danmu2.Danmu) bool
}

func NewDanmuMatcher(t danmu2.DanmuType, params *MatcherParams) DanmuMatcher {
	switch t {
	case danmu2.DanmuTypeDANMUMSG:
		return newDanmuMsgMatcher(params.List)
	default:
		return nil
	}
}
