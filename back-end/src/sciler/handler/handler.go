package handler

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"os"
	"sciler/communication"
	"sciler/config"
	"time"
)

// Message is a type that follows the structure all messages have, described in resources/message_manual.md
type Message struct {
	DeviceID string                 `json:"device_id"`
	TimeSent string                 `json:"time_sent"`
	Type     string                 `json:"type"`
	Contents map[string]interface{} `json:"contents"`
}

// Handler is a type that mqqt handlers have
type Handler struct {
	Config       config.WorkingConfig
	Communicator communication.Communicator
}

// Create an instance of Handler
func GetHandler(workingConfig config.WorkingConfig, communicator communication.Communicator) *Handler {
	return &Handler{workingConfig, communicator}
}

// NewHandler is the actual MessageHandler
func (handler *Handler) NewHandler(client mqtt.Client, message mqtt.Message) {
	// TODO: Make advanced message handler which acts according to the events / configuration
	var raw Message
	if err := json.Unmarshal(message.Payload(), &raw); err != nil {
		logrus.Errorf("Invalid JSON received: %v", err)
	}
	switch raw.Type {
	case "instruction":
		{
			logrus.Info("instruction message received from: ", raw.DeviceID)
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
					Communicator.Publish("test", string(jsonMessage), 3)
				}
			}
		}
	case "status":
		{
			logrus.Info("status  message received from: ", raw.DeviceID)
			for k, v := range raw.Contents {
				configurations.Devices[raw.DeviceID].Status[k] = v
			}

			//logrus.Print(configurations.Devices[raw.DeviceID])
			//logrus.Info("status message received")
			//if raw.DeviceID == "controlBoard1" {
			//	var instruction string
			//	if raw.Contents["switch1"] == "1" {
			//		instruction = "turn off"
			//	} else if raw.Contents["switch1"] == "0" {
			//		instruction = "turn on"
			//	}
			//	if instruction != "" {
			//		message := Message{
			//			DeviceID: "back-end",
			//			TimeSent: time.Now().Format("02-01-2006 15:04:05"),
			//			Type:     "instruction",
			//			Contents: map[string]interface{}{
			//				"instruction": instruction,
			//			},
			//		}
			//		jsonMessage, err := json.Marshal(&message)
			//		if err != nil {
			//			logrus.Errorf("Error occurred while constructing message to publish: %v", err)
			//		} else {
			//			communicator.Publish("test", string(jsonMessage), 3)
			//		}
			//	}
			// }

		}
	case "confirmation":
		{
			logrus.Info("confirmation message received from: ", raw.DeviceID)
		}
	case "connection":
		{
			on_connection_msg(raw)
		}
	}

}

func on_connection_msg(raw json.RawMessage) {
	logrus.Info("connection message received from:", raw.DeviceID)
	con := configurations.Devices[raw.DeviceID]
	con.Connection = true
	configurations.Devices[raw.DeviceID] = con
}
