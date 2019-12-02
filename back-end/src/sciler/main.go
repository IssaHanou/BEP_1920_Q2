package main

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"os"
	"reflect"
	"sciler/communication"
	"sciler/config"
	"time"
)

var topics = []string{"back-end", "status", "hint", "connect"}

// Message is a type that follows the structure all messages have, described in resources/message_manual.md
type Message struct {
	DeviceID string                 `json:"device_id"`
	TimeSent string                 `json:"time_sent"`
	Type     string                 `json:"type"`
	Contents map[string]interface{} `json:"contents"`
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		logrus.Fatal(err)
	}
	filename := dir + "/back-end/resources/room_config.json"
	configurations := config.ReadFile(filename)
	logrus.Info("configurations read from: " + filename)
	host := configurations.General.Host
	port := configurations.General.Port

	communicator := communication.NewCommunicator(host, port, topics)
	go communicator.Start(func(client mqtt.Client, message mqtt.Message) {
		// TODO: Make advanced message handler which acts according to the events / configuration
		var raw Message
		if err := json.Unmarshal(message.Payload(), &raw); err != nil {
			logrus.Errorf("Invalid JSON received: %v", err)
		}
		switch raw.Type {
		case "instruction":
			logrus.Info("instruction  message received")
			{
				if raw.Contents["instruction"] == "test all" && raw.DeviceID == "front-end" { // TODO maybe switch again
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
				value, ok := raw.Contents["completed"]
				if !ok || reflect.TypeOf(value).Kind() != reflect.Bool {
					logrus.Error("Received improperly structured confirmation message from device " + raw.DeviceID)
				} else if !value.(bool) {
					value2, ok2 := raw.Contents["instructed"]
					if !ok2 {
						logrus.Error("Device " + raw.DeviceID + " did not complete instruction at " +
							raw.TimeSent)
					} else {
						logrus.Error("Device " + raw.DeviceID + " did not complete instruction with type " +
							value2.(string) + " at " + raw.TimeSent)
					}
				} else {
					value2, ok2 := raw.Contents["instructed"]
					if !ok2 {
						logrus.Info("Device " + raw.DeviceID + " completed instruction at " + raw.TimeSent)
					} else {
						logrus.Info("Device " + raw.DeviceID + " completed instruction with type " +
							value2.(string) + " at " + raw.TimeSent)
					}
				}
			}
		case "connection":
			{
				logrus.Info("connection message received")
			}
		default:
			{
				logrus.Info("received unsupported message of type " + raw.Type)
			}
		}
	})

	// loop for now preventing app to exit
	for {
		time.Sleep(time.Microsecond * time.Duration(250))
	}
}
