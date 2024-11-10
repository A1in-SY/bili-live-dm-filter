package action

import (
	"bytes"
	"context"
	"encoding/json"
	"filter-core/util/errwarp"
	"io"
	"net/http"
	"strconv"
)

type qqPrivateAction struct {
	cli    *http.Client
	url    string
	userId int64
}

func newQQPrivateAction(extraInfo map[string]string) (*qqPrivateAction, error) {
	url := extraInfo["url"]
	userId, err := strconv.ParseInt(extraInfo["user_id"], 10, 64)
	if err != nil {
		return nil, errwarp.Warp("qq private action parse user_id fail", err)
	}
	return &qqPrivateAction{
		cli:    http.DefaultClient,
		url:    url,
		userId: userId,
	}, nil
}

func (a *qqPrivateAction) doAction(ctx context.Context, content string) error {
	m1 := map[string]interface{}{
		"user_id": a.userId,
		"message": content,
	}
	data, _ := json.Marshal(m1)
	req, _ := http.NewRequest(http.MethodPost, a.url, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.cli.Do(req)
	if err != nil {
		return errwarp.Warp("do http call fail", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errwarp.Warp("read http resp body fail", err)
	}
	m2 := make(map[string]interface{})
	err = json.Unmarshal(body, &m2)
	if err != nil {
		return errwarp.Warp("unmarshal resp body fail", err)
	}
	switch m2["status"].(type) {
	case string:
		if m2["status"] != "ok" {
			return errwarp.Warp("send qq private message req fail", nil)
		}
		return nil
	default:
		return errwarp.Warp("wrong type of qq private message resp status", nil)
	}
}
