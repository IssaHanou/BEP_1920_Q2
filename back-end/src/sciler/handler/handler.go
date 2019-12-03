package handler

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"reflect"
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

//GetHandler creates an instance of Handler
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
	handler.msgMapper(raw)
}

// msgMapper sends the right message through to the right function
func (handler *Handler) msgMapper(raw Message) {
	switch raw.Type {
	case "instruction":
		{
			handler.onInstructionMsg(raw)
		}
	case "status":
		{
			handler.onStatusMsg(raw)
			handler.openDoorBeun(raw)
		}
	case "confirmation":
		{
			handler.onConfirmationMsg(raw)
		}
	case "connection":
		{
			handler.onConnectionMsg(raw)
		}
	default:
		{
			logrus.Error("message received from ", raw.DeviceID,
				", but no message type could be found for: ", raw.Type)
		}
	}

}

//onConnectionMsg is the function to process connection messages.
func (handler *Handler) onConnectionMsg(raw Message) {
	logrus.Info("connection message received from:", raw.DeviceID)
	con := handler.Config.Devices[raw.DeviceID]
	con.Connection = true
	handler.Config.Devices[raw.DeviceID] = con
}

//onStatusMsg is the function to process status messages.
func (handler *Handler) onStatusMsg(raw Message) {
	logrus.Info("status message received from: ", raw.DeviceID)
	for k, v := range raw.Contents {
		handler.Config.Devices[raw.DeviceID].Status[k] = v
	}
}

//onConfirmationMsg is the function to process confirmation messages.
func (handler *Handler) onConfirmationMsg(raw Message) {
	logrus.Info("confirmation message received from: ", raw.DeviceID)
	value, ok := raw.Contents["completed"]
	if !ok || reflect.TypeOf(value).Kind() != reflect.Bool {
		logrus.Error("Received improperly structured confirmation message from device " + raw.DeviceID)
	} else if !value.(bool) {
		logrus.Error("Device " + raw.DeviceID + " did not complete instruction with type " +
			raw.Contents["instructed"].(string) + " at " + raw.TimeSent)
	} else {
		logrus.Info("Device " + raw.DeviceID + " completed instruction with type " +
			raw.Contents["instructed"].(string) + " at " + raw.TimeSent)
	}
}

//openDoorBeun is the test function for developers to test the door and switch combo
func (handler *Handler) openDoorBeun(raw Message) {
	logrus.Info("status message received, checking if door needs to open.")
	if raw.DeviceID == "controlBoard" {
		var instruction string
		if raw.Contents["mainSwitch"] == "1" {
			instruction = "turn off"
		} else if raw.Contents["mainSwitch"] == "0" {
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
				handler.Communicator.Publish("test", string(jsonMessage), 3)
			}
		}
	}
}

//onInstructionMsg is the function to process instruction messages.
func (handler *Handler) onInstructionMsg(raw Message) {
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
			handler.Communicator.Publish("test", string(jsonMessage), 3)
		}
	}
}
