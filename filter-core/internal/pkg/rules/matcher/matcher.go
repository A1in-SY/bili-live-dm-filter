package matcher

import "filter-core/internal/pkg/danmu"

type RuleParam struct {
	Param string
	Type  BaseMatcherType
	Mode  BaseMatchMode
	Value string
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
	//GetMatcherInfo() *MatcherInfo
	IsDanmuMatch(dm *danmu.Danmu) bool
}
