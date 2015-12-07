package main

import (
	"fmt"
	"github.com/adufrene/gobot"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Configuration struct {
	Token string `yaml:"apiToken"`
}

func main() {
	apiToken, err := readApiToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading api token: %s\n", err.Error())
		os.Exit(1)
	}
	gobot := gobot.NewGobot(apiToken)
	gobot.RegisterPresenceChangeFunction(changedPresence)
	gobot.RegisterAllMessageFunction(printMessage)
	gobot.RegisterMessageFunction(echoMessage)
	if err = gobot.Listen(); err != nil {
		fmt.Fprintf(os.Stderr, "Error while listening: %s\n", err.Error())
		os.Exit(1)
	}
}

func readApiToken() (string, error) {
	file, err := ioutil.ReadFile("configuration.yaml")
	if err != nil {
		return "", err
	}
	var conf Configuration
	if err = yaml.Unmarshal(file, &conf); err != nil {
		return "", err
	}
	return conf.Token, nil
}

func echoMessage(slackApi gobot.SlackApi, message gobot.Message) {
	slackApi.PostMessage(message.Channel, message.Text)
}

func printMessage(slackApi gobot.SlackApi, message gobot.Message) {
	fmt.Printf("%s: %s\n", message.User, message.Text)
}

func changedPresence(slackApi gobot.SlackApi, presenceChange gobot.PresenceChange) {
	user, err := slackApi.GetUser(presenceChange.User)
	if err == nil {
		fmt.Println(user.Name + " changed status to " + presenceChange.Presence)
	}
}
