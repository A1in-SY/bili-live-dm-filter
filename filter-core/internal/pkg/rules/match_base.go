package rules

import (
	"go.uber.org/zap"
	"strings"
)

type baseMatch interface {
	isBaseMatch(a interface{}) bool
}

type stringMatch struct {
	value string
	mode  stringMatchMode
}

type stringMatchMode int64

const (
	stringMatchModeEqual stringMatchMode = iota
	stringMatchModeContain
	stringMatchModeRegex
)

func (m *stringMatch) isBaseMatch(a interface{}) bool {
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

type int64Match struct {
	value int64
	mode  int64MatchMode
}

type int64MatchMode int64

const (
	int64MatchModeEqual int64MatchMode = iota
	int64MatchModeNotEqual
	int64MatchModeGreater
	int64MatchModeLess
	int64MatchModeGreaterOrEqual
	int64MatchModeLessOrEqual
)

func (m *int64Match) isBaseMatch(a interface{}) bool {
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
