package rules

import (
	"filter-core/internal/pkg/action"
	"filter-core/internal/pkg/rules/matcher"
	"fmt"
	"sync"
	"time"
)

type RuleManager struct {
	ruleMap map[string]*Rule
	mu      sync.Mutex
}

func NewRuleManager() *RuleManager {
	return &RuleManager{
		ruleMap: make(map[string]*Rule),
		mu:      sync.Mutex{},
	}
}

func (mng *RuleManager) GetRuleByRuleId(ruleId string) *Rule {
	mng.mu.Lock()
	defer mng.mu.Unlock()
	return mng.ruleMap[ruleId]
}

func (mng *RuleManager) GetRuleList() []*Rule {
	mng.mu.Lock()
	defer mng.mu.Unlock()
	rules := make([]*Rule, 0)
	for _, rule := range mng.ruleMap {
		rules = append(rules, rule)
	}
	return rules
}

// 这里接受pb，解析为params
func (mng *RuleManager) AddRule(name string, dmType int64, actionList []action.RuleAction) error {
	mng.mu.Lock()
	defer mng.mu.Unlock()
	ruleId := fmt.Sprintf("rule_%v", time.Now().UnixMilli())
	rule := NewRule(ruleId, name, dmType, &matcher.MatcherParams{}, actionList)
	mng.ruleMap[ruleId] = rule
	return nil
}
