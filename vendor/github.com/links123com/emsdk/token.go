package emsdk

import (
	"bytes"
	"encoding/json"
)

type adminTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Application string `json:"application"`
}

func (c *Client) getAccessToken() (adminTokenResponse, error) {
	var adminTokenResponse adminTokenResponse
	data := struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}{
		"client_credentials",
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
