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

const host string = "192.168.178.82"
const port string = "1883"

var topics = []string{"back-end", "status", "hint", "connect"}

type Message struct {
	DeviceId string      `json:"device_id"`
	TimeSent string      `json:"time_sent"`
	Type     string      `json:"type"`
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
			{
				if raw.Contents["instruction"] == "test all" { // todo maybe switch again
					logrus.Info("instruction -> test all")
					message := Message{
						DeviceId: "back-end",
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
						communicator.Publish("test", string(jsonMessage))
					}
				}
			}
		case "status":
			{
				logrus.Info("status")
			}
		case "confirmation":
			{
				logrus.Info("confirmation")
			}
		case "connection":
			{
				logrus.Info("connection")
			}
		}
	})

	//loop for now preventing app to exit
	for {
		time.Sleep(time.Microsecond * time.Duration(250))
	}
}
