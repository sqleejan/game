package emsdk

import (
	"encoding/json"
	"strings"
)

// 添加好友
func (c *Client) AddContact(owner, friend string) error {
	url := "users/" + owner + "/contacts/users/" + friend
	_, err := c.sendRequest(url, strings.NewReader(""), "POST")

	return err
}

// 删除好友
func (c *Client) DeleteContact(owner, friend string) error {
	url := "users/" + owner + "/contacts/users/" + friend
	_, err := c.sendRequest(url, strings.NewReader(""), "DELETE")

	return err
}

// 往一个 IM 用户的黑名单中加人
func (c *Client) AddUserToBlackList(owner string, friends []string) error {
	url := "users/" + owner + "/blocks/users/"
	request := struct {
		UserNames []string `json:"usernames"`
	}{
		UserNames: friends,
	}

	data, err := json.Marshal(request)
	if err != nil {
		return err
	}

	_, err = c.sendRequest(url, strings.NewReader(string(data)), "POST")

	return err
}

// 从一个 IM 用户的黑名单中减人
func (c *Client) RemoveUserFromBlackList(owner, blocked string) error {
	url := "users/" + owner + "/blocks/users/" + blocked
	_, err := c.sendRequest(url, strings.NewReader(""), "DELETE")

	return err
}
