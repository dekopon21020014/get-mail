package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// struct of message
type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// struct of message payoad
type Payload struct {
	To       string    `json:"to"`
	Messages []Message `json:"messages"`
}

func FetchToLine(message string) {
	url     := "https://api.line.me/v2/bot/message/push" 

	// generating payload of message
	payload := Payload{
		To: os.Getenv("LINE_ID"),
		Messages: []Message{
			{
				Type: "text",
				Text: message,
			},
		},
	}

	// encoding as json
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("encodeing error: %v\n", err)
		return
	}

	// generating a htttp request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Printf("generating requst error: %v\n", err)
		return
	}

	// set up headers
	req.Header.Set("Authorization", "Bearer " + os.Getenv("TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	// creating http client and sneding a request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("http request error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// response handling
	if resp.StatusCode == http.StatusOK {
		fmt.Println("the message was sent successfully")
	} else {
		fmt.Printf("An error was occured while sending message: %s\n", resp.Status)
	}
}
