package rules

import (
	"filter-core/internal/pkg/danmu"
	"filter-core/internal/pkg/rules/action"
	"filter-core/internal/pkg/rules/matcher"
)

type Rule struct {
	id            string
	name          string
	roomId        int64
	dmType        danmu.DanmuType
	dmMatcher     matcher.DanmuMatcher
	triggerAction action.RuleAction
}
