package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	BASE_URL = "https://slack.com/api/"
	VERBOSE  = false
)

func (slackApi SlackApi) TestAuth() (err error) {
	resp, err := slackApi.get("auth.test", nil)
	if err != nil {
		return
	}

	var slackResp slackAuthResponse
	if err = json.Unmarshal(resp, &slackResp); err != nil {
		return
	}

	if !slackResp.Okay {
		return fmt.Errorf("Bad response from slack: %s", slackResp.Error)
	}
	return nil
}

func (slackApi SlackApi) GetUser(userId string) (SlackUser, error) {
	body, err := slackApi.get("users.info", map[string]string{"user": userId})

	if err != nil {
		return SlackUser{}, err
	}

	var user slackUser
	if err := json.Unmarshal(body, &user); err != nil {
		return SlackUser{}, err
	}
	return user.User, nil
}

func (slackApi SlackApi) PostMessage(channel, message string) {
	slackApi.get("chat.postMessage", map[string]string{
		"channel": channel,
		"text":    message,
	})
}

func NewSlackApi(apiToken string) (api SlackApi) {
	api.token = apiToken
	//	api.TestAuth()
	return
}

func (slackApi SlackApi) startRTM() (start slackStart, err error) {
	body, err := slackApi.get("rtm.start", nil)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, &start); err != nil {
		return
	}
	return
}

func (slackApi SlackApi) verifyApiToken() error {
	if slackApi.token == "" {
		return fmt.Errorf("Api token is empty, make sure you create a slack api using the constructor")
	}
	return nil
}

func (slackApi SlackApi) get(endpoint string, params map[string]string) ([]byte, error) {
	slackApi.verifyApiToken()
	url := BASE_URL + endpoint + "?token=" + slackApi.token

	for key, value := range params {
		paramString := "&" + key + "=" + value
		url += paramString
	}

	if VERBOSE {
		fmt.Println("Making GET request for " + url)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (slackApi SlackApi) whoami() (SlackUser, error) {
	resp, err := slackApi.get("auth.test", nil)
	if err != nil {
		return SlackUser{}, err
	}

	var jsonMap map[string]interface{}
	if err = json.Unmarshal(resp, &jsonMap); err != nil {
		return SlackUser{}, err
	}

	return slackApi.GetUser(jsonMap["user_id"].(string))
}
