package matcher

import (
	"go.uber.org/zap"
	"strings"
)

type BaseMatcherType int64
type BaseMatchMode int64

const (
	BaseMatcherTypeString BaseMatcherType = 1
	BaseMatcherTypeInt64  BaseMatcherType = 2
)

type stringMatcher struct {
	value string
	mode  BaseMatchMode
}

const (
	stringMatchModeEqual   BaseMatchMode = 1
	stringMatchModeContain BaseMatchMode = 2
	stringMatchModeRegex   BaseMatchMode = 3
)

func (m *stringMatcher) isBaseMatch(a interface{}) bool {
	if m == nil {
		return true
	}
	target, ok := a.(string)
	if !ok {
		return false
	}
	switch m.mode {
	case stringMatchModeEqual:
		return target == m.value
	case stringMatchModeContain:
		return strings.Contains(target, m.value)
	case stringMatchModeRegex:
		return true
	default:
		zap.S().Errorf("unsupported string match mode: %v", m.mode)
		return false
	}
}

type int64Matcher struct {
	value int64
	mode  BaseMatchMode
}

const (
	int64MatchModeEqual          BaseMatchMode = 1
	int64MatchModeNotEqual       BaseMatchMode = 2
	int64MatchModeGreater        BaseMatchMode = 3
	int64MatchModeLess           BaseMatchMode = 4
	int64MatchModeGreaterOrEqual BaseMatchMode = 5
	int64MatchModeLessOrEqual    BaseMatchMode = 6
)

func (m *int64Matcher) isBaseMatch(a interface{}) bool {
	if m == nil {
		return true
	}
	target, ok := a.(int64)
	if !ok {
		return false
	}
	switch m.mode {
	case int64MatchModeEqual:
		return target == m.value
	case int64MatchModeNotEqual:
		return target != m.value
	case int64MatchModeGreater:
		return target > m.value
	case int64MatchModeLess:
		return target < m.value
	case int64MatchModeGreaterOrEqual:
		return target >= m.value
	case int64MatchModeLessOrEqual:
		return target <= m.value
	default:
		zap.S().Errorf("unsupported int64 match mode: %v", m.mode)
		return false
	}
}
