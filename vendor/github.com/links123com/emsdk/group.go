package emsdk

import (
	"encoding/json"
	"strings"
)

// 增加群组
func (c *Client) AddGroup(groupName, desc, owner string, public, approval bool, maxUsers int, members []string) (string, error) {
	url := "chatgroups"
	var res string
	var restruct struct {
		Data struct {
			Groupid string
		}
	}
	group := struct {
		GroupName string   `json:"groupname"`
		Desc      string   `json:"desc"`
		Owner     string   `json:"owner"`
		Public    bool     `json:"public"`
		Approval  bool     `json:"approval"`
		MaxUsers  int      `json:"maxusers"`
		Members   []string `json:"members"`
	}{
		GroupName: groupName,
		Desc:      desc,
		Owner:     owner,
		Public:    public,
		Approval:  approval,
		MaxUsers:  maxUsers,
		Members:   members,
	}

	data, err := json.Marshal(group)
	if err != nil {
		return "", err
	}

	res, err = c.sendRequest(url, strings.NewReader(string(data)), "POST")
	if err != nil {
		return "", err
	}

	err = json.Unmarshal([]byte(res), &restruct)
	if err != nil {
		return "", err
	}

	return (restruct.Data.Groupid), err
}

// 修改群组
func (c *Client) ModGroup(groupname, description, groupid string, maxusers int) error {
	url := "chatgroups/" + groupid
	group := struct {
		Groupname   string `json:"groupname"`
		Description string `json:"description"`
		Maxusers    int    `json:"maxusers"`
	}{
		Groupname:   groupname,
		Description: description,
		Maxusers:    maxusers,
	}

	data, err := json.Marshal(group)
	if err != nil {
		return err
	}

	_, err = c.sendRequest(url, strings.NewReader(string(data)), "PUT")

	return err
}

// 删除群组
func (c *Client) DelGroup(groupid string) error {
	url := "chatgroups/" + groupid
	_, err := c.sendRequest(url, strings.NewReader(""), "DELETE")

	return err
}

// 获取APP中所有群组
func (c *Client) FetchAllGroupFromApp() (*GroupList, error) {
	url := "chatgroups/"
	result, err := c.sendRequest(url, strings.NewReader(""), "GET")
	gl := &GroupList{}
	err = json.Unmarshal([]byte(result), gl)
	if err != nil {
		return nil, err
	}
	return gl, nil
}

// 获取群组详情  groupid := "group1,group2..."
func (c *Client) FetchGroupInfo(groupid string) (string, error) {
	url := "chatgroups/" + groupid
	result, err := c.sendRequest(url, strings.NewReader(""), "GET")

	return result, err
}

// 获取群组所有成员
func (c *Client) FetchUserFromGroup(groupid string) (string, error) {
	url := "chatgroups/" + groupid + "/users"
	result, err := c.sendRequest(url, strings.NewReader(""), "GET")

	return result, err
}

// 增加群组成员(单个)
func (c *Client) AddUserToGroup(groupid, username string) error {
	url := "chatgroups/" + groupid + "/users/" + username
	_, err := c.sendRequest(url, strings.NewReader(""), "POST")

	return err
}

// 增加群组成员(批量)
func (c *Client) AddBatchUserToGroup(groupid string, usernames []string) error {
	url := "chatgroups/" + groupid + "/users"
	users := struct {
		Usernames []string `json:"usernames"`
	}{
		Usernames: usernames,
	}

	data, err := json.Marshal(users)
	if err != nil {
		return err
	}

	_, err = c.sendRequest(url, strings.NewReader(string(data)), "POST")

	return err
}

// 移除群组成员(单个)
func (c *Client) DelUserFromGroup(groupid, username string) error {
	url := "chatgroups/" + groupid + "/users/" + username
	_, err := c.sendRequest(url, strings.NewReader(""), "DELETE")

	return err
}

// 移除群组成员(批量) members := "member1,member2..."
func (c *Client) DelBatchUserFromGroup(groupid, members string) error {
	url := "chatgroups/" + groupid + "/users/" + members
	_, err := c.sendRequest(url, strings.NewReader(""), "DELETE")

	return err
}

// 获取一个用户参与的所有群组
func (c *Client) FetchGroupFromUserJoined(username string) (string, error) {
	url := "users/" + username + "/joined_chatgroups"
	result, err := c.sendRequest(url, strings.NewReader(""), "GET")

	return result, err
}

// 修改群组 Owner 为同一 APP 下的其他用户
func (c *Client) ChangeOwner(groupid, username string) error {
	url := "chatgroups/" + groupid
	owner := struct {
		Newowner string `json:"newowner"`
	}{
		Newowner: username,
	}

	data, err := json.Marshal(owner)
	if err != nil {
		return err
	}

	_, err = c.sendRequest(url, strings.NewReader(string(data)), "PUT")

	return err
}

// 查询一个群组黑名单用户名列表
func (c *Client) FetchUserFromBlackList(groupid string) (string, error) {
	url := "chatgroups/" + groupid + "/blocks/users"
	result, err := c.sendRequest(url, strings.NewReader(""), "GET")

	return result, err
}

// 添加一个用户进入一个群组的黑名单
func (c *Client) AddUserToBlack(groupid, username string) error {
	url := "chatgroups/" + groupid + "/blocks/users/" + username
	_, err := c.sendRequest(url, strings.NewReader(""), "POST")

	return err
}

// 批量添加多个用户进入一个群组的黑名单
func (c *Client) AddBatchUserToBlackList(groupid string, usernames []string) error {
	url := "chatgroups/" + groupid + "/blocks/users"
	users := struct {
		Usernames []string `json:"usernames"`
	}{
		Usernames: usernames,
	}

	data, err := json.Marshal(users)
	if err != nil {
		return err
	}

	_, err = c.sendRequest(url, strings.NewReader(string(data)), "POST")

	return err
}

// 从群组黑名单中移除一个用户
func (c *Client) DelUserFromBlackList(groupid, username string) error {
	url := "chatgroups/" + groupid + "/blocks/users/" + username
	_, err := c.sendRequest(url, strings.NewReader(""), "DELETE")

	return err
}

// 从群组黑名单中移除用户(批量) members := "member1,member2..."
func (c *Client) DelBatchUserFromBlackList(groupid, members string) error {
	url := "chatgroups/" + groupid + "/blocks/users/" + members
	_, err := c.sendRequest(url, strings.NewReader(""), "DELETE")

	return err
}

type Group struct {
	Groupid      string `json:"groupid"`
	Groupname    string `json:"groupname"`
	Owner        string `json:"owner"`
	Affiliations int    `json:"affiliations"`
	LastModified int64  `json:"lastModified"`
	Type         string `json:"type"`
	Created      int64  `json:"created"`
}

type GroupList struct {
	Action          string        `json:"action"`
	Application     string        `json:"application"`
	URI             string        `json:"uri"`
	Entities        []interface{} `json:"entities"`
	Data            []Group       `json:"data"`
	Timestamp       int64         `json:"timestamp"`
	Duration        int           `json:"duration"`
	Organization    string        `json:"organization"`
	ApplicationName string        `json:"applicationName"`
	Count           int           `json:"count"`
}
