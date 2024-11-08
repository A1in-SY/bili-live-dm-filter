package rule

import (
	pb "filter-core/api/v1"
	"filter-core/internal/pkg/action"
	"filter-core/internal/pkg/rule/matcher"
	"filter-core/util/errwarp"
	"fmt"
	"strconv"
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
func (mng *RuleManager) AddRule(name string, dmType int64, paramList []*pb.MatcherParam, actionList []action.RuleAction) error {
	mng.mu.Lock()
	defer mng.mu.Unlock()
	ruleId := fmt.Sprintf("rule_%v", time.Now().UnixMilli())
	matcherParamList := make([]*matcher.MatcherParam, 0)
	for _, param := range paramList {
		var value interface{}
		var err error
		switch param.BaseType {
		case int64(matcher.BaseMatcherTypeString):
			value = param.Value
		case int64(matcher.BaseMatcherTypeInt64):
			value, err = strconv.ParseInt(param.Value, 10, 64)
			if err != nil {
				return errwarp.Warp("parse base int matcher type fail", err)
			}
		default:
			return errwarp.Warp(fmt.Sprintf("unknown base matcher type: %v", param.BaseType), nil)
		}
		matcherParamList = append(matcherParamList, &matcher.MatcherParam{
			Param:     param.Param,
			BaseType:  matcher.BaseMatcherType(param.BaseType),
			MatchMode: matcher.BaseMatchMode(param.MatchMode),
			Value:     value,
		})
	}
	rule := NewRule(ruleId, name, dmType, matcherParamList, actionList)
	mng.ruleMap[ruleId] = rule
	return nil
}
