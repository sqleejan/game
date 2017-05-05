package emsdk

import (
	"encoding/json"
	"strings"
)

func (c *Client) SendMessage(from, targetType string, target []string, msg, ext map[string]string) error {
	url := "messages"
	request := struct {
		TargetType string            `json:"target_type"`
		Target     []string          `json:"target"`
		From       string            `json:"from"`
		Message    map[string]string `json:"msg"`
		Ext        map[string]string `json:"ext"`
	}{
		TargetType: targetType,
		Target:     target,
		From:       from,
		Message:    msg,
		Ext:        ext,
	}

	data, err := json.Marshal(request)
	if err != nil {
		return err
	}

	_, err = c.sendRequest(url, strings.NewReader(string(data)), "POST")

	return err
}
