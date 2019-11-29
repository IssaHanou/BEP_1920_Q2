package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"sciler/communication"
	"sciler/config"
	"time"
)

const host string = "localhost"
const port string = "1883"

var topics = []string{"back-end", "status", "hint", "connect"}

// Message is a type that follows the structure all messages have, described in resources/message_manual.md
type Message struct {
	DeviceID string                 `json:"device_id"`
	TimeSent string                 `json:"time_sent"`
	Type     string                 `json:"type"`
	Contents map[string]interface{} `json:"contents"`
}

func main() {
	fmt.Println("Starting server")
	data := []byte(`{
            "name": "My Awesome Escape",
            "duration": "00:30:00"
    }`)
	fmt.Println(config.GetFromJSON(data))

	communicator := communication.NewCommunicator(host, port, topics)
	go communicator.Start(func(client mqtt.Client, message mqtt.Message) {
		// Todo: Make advanced message handler which acts according to the events / configuration
		var raw Message
		if err := json.Unmarshal(message.Payload(), &raw); err != nil {
			logrus.Errorf("Invalid JSON received: %v", err)
		}
		switch raw.Type {
		case "instruction":
			logrus.Info("instruction  message received")
			{
				if raw.Contents["instruction"] == "test all" && raw.DeviceID == "front-end" { // todo maybe switch again
					message := Message{
						DeviceID: "back-end",
						TimeSent: time.Now().Format("02-01-2006 15:04:05"),
						Type:     "instruction",
						Contents: map[string]interface{}{
							"instruction": "test",
						},
					}
					jsonMessage, err := json.Marshal(&message)
					if err != nil {
						logrus.Errorf("Error occurred while constructing message to publish: %v", err)
					} else {
						communicator.Publish("test", string(jsonMessage), 3)
					}
				}
			}
		case "status":
			{
				logrus.Info("status message received")
				if raw.DeviceID == "controlBoard1" {
					var instruction string
					if raw.Contents["switch1"] == "1" {
						instruction = "turn off"
					} else if raw.Contents["switch1"] == "0" {
						instruction = "turn on"
					}
					if instruction != "" {
						message := Message{
							DeviceID: "back-end",
							TimeSent: time.Now().Format("02-01-2006 15:04:05"),
							Type:     "instruction",
							Contents: map[string]interface{}{
								"instruction": instruction,
							},
						}
						jsonMessage, err := json.Marshal(&message)
						if err != nil {
							logrus.Errorf("Error occurred while constructing message to publish: %v", err)
						} else {
							communicator.Publish("test", string(jsonMessage), 3)
						}
					}
				}

			}
		case "confirmation":
			{
				logrus.Info("confirmation message received")
			}
		case "connection":
			{
				logrus.Info("connection message received")
			}
		}
	})

	//loop for now preventing app to exit
	for {
		time.Sleep(time.Microsecond * time.Duration(250))
	}
}
