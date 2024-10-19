package action

type RuleAction interface {
	DoAction(abstract string) error
}

type RuleActionType int64

const (
	RuleActionTypeQQPrivate = 1
	RuleActionTypeQQGroup   = 2
	RuleActionTypeWebhook   = 3
)

type ActionParam struct {
	Type RuleActionType
}
