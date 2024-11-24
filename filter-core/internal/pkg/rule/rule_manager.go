package rule

import (
	"context"
	pb "filter-core/api/v1"
	"filter-core/internal/model/danmu"
	"filter-core/internal/pkg/rule/matcher"
	"filter-core/util/errwarp"
	"filter-core/util/log"
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

func (mng *RuleManager) GetRuleDmChanByRuleId(ctx context.Context, ruleId string) (*danmu.DanmuChannel, error) {
	mng.mu.Lock()
	defer mng.mu.Unlock()
	if rule, ok := mng.ruleMap[ruleId]; ok {
		return rule.dmChan, nil
	}
	return nil, errwarp.Warp(fmt.Sprintf("can't find rule with id: %v", ruleId), nil)
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
func (mng *RuleManager) AddRule(ctx context.Context, name string, dmType int64, paramList []*pb.MatcherParam, actionChs []*danmu.DanmuChannel) error {
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
	rule := newRule(ctx, ruleId, name, danmu.DanmuType(dmType), matcherParamList, actionChs)
	mng.ruleMap[ruleId] = rule
	log.Infoc(ctx, "%v", ruleId)
	return nil
}
