package emsdk

import (
	"bytes"
	"encoding/json"
	"sync"
)

var (
	client *Client
	once   sync.Once
)

type Client struct {
	clientID     string
	clientSecret string
	baseURL      string
	adminToken   adminTokenResponse
}

func New(orgName, appName, clientID, clientSecret string) (*Client, error) {
	var err error
	once.Do(func() {
		client = &Client{
			baseURL:      "https://a1.easemob.com/" + orgName + "/" + appName,
			clientID:     clientID,
			clientSecret: clientSecret,
		}

		client.adminToken, err = client.getAccessToken()
	})

	return client, err
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (c *Client) GetUserToken(clientID, clientSecret string) (string, error) {
	client = &Client{
		baseURL:      c.baseURL,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
	res, err := client.getClientToken()
	if err != nil {
		return "", err
	}
	return res.AccessToken, nil

}

func (c *Client) getClientToken() (tokenResponse, error) {
	var adminTokenResponse tokenResponse
	data := struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"username"`
		ClientSecret string `json:"password"`
	}{
		"password",
		c.clientID,
		c.clientSecret,
	}

	b, err := json.Marshal(data)
	if err != nil {
		return adminTokenResponse, err
	}

	body := bytes.NewBuffer([]byte(b))
	result, err := c.sendRequest("token", body, "POST")
	if err != nil {
		return adminTokenResponse, err
	}

	json.Unmarshal([]byte(result), &adminTokenResponse)

	return adminTokenResponse, nil
}
