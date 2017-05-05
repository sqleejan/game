package emsdk

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const retryInterval = 3 * time.Second

var ErrEM EMError

type EMError struct {
	Code        int    `json:"code"`
	Message     string `json:"error"`
	Description string `json:"error_description"`
}

func (e EMError) Error() string {
	return fmt.Sprintf("%v", e.Description)
}

func (c *Client) sendRequest(name string, body io.Reader, method string) (string, error) {
	methodName := c.baseURL + "/" + name
	client := &http.Client{}
	req, err := http.NewRequest(method, methodName, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	if c.adminToken.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.adminToken.AccessToken)
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusServiceUnavailable {
			time.Sleep(retryInterval)
			res, err := client.Do(req)
			if err != nil {
				return "", err
			}

			result, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return "", err
			}
			res.Body.Close()
			return string(result), nil
		}
		ErrEM.Code = res.StatusCode
		json.Unmarshal(result, &ErrEM)
		return "", ErrEM
	}

	return string(result), nil
}
