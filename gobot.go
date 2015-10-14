package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
)

const (
	API_TOKEN = "xoxb-12456280647-vKyUQvjTf3BH2aUhUux2qWF3"
)

type slackStart struct {
	Okay bool   `json:"ok"`
	URL  string `json:"url"`
}

func main() {
	resp, err := http.Get("https://slack.com/api/rtm.start?token=" + API_TOKEN)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	//	fmt.Println(string(body))

	var start slackStart
	err = json.Unmarshal(body, &start)

	if err != nil {
		panic(err)
	}

	requestHeader := http.Header{
		"Origin":               {"http://127.0.0.1:80"},
		"WebSocket-Extensions": {"permessage-deflate; client_max_window_bits, x-webkit-deflate-frame"},
	}

	var defaultDialer *websocket.Dialer
	conn, resp, err := defaultDialer.Dial(start.URL, requestHeader)

	if err != nil {
		panic(err)
	}

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v, %s\n", msgType, msg)
	}
}
