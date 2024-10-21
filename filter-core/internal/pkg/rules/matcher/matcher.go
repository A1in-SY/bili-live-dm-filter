package matcher

import (
	danmu2 "filter-core/internal/model/danmu"
)

type MatcherParams struct {
	List []*MatcherParamItem
}

type MatcherParamItem struct {
	Param string
	Type  BaseMatcherType
	Mode  BaseMatchMode
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
