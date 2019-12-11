package handler

import (
	"encoding/json"
	"fmt"
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
		logrus.Errorf("invalid JSON received: %v", err)
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
			handler.SendStatus(raw.DeviceID)
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
	device, ok := handler.Config.Devices[raw.DeviceID]
	if !ok {
		logrus.Warn("connection message received from device " + raw.DeviceID + ", which is not in the config")
	} else {
		logrus.Info("connection message received from: ", raw.DeviceID)
		value, ok2 := raw.Contents["connection"]
		if !ok2 || reflect.TypeOf(value).Kind() != reflect.Bool {
			logrus.Error("received improperly structured connection message from device " + raw.DeviceID)
		}
		device.Connection = value.(bool)
		handler.Config.Devices[raw.DeviceID] = device
		handler.SendStatus(raw.DeviceID)
	}
}

//onStatusMsg is the function to process status messages.
func (handler *Handler) onStatusMsg(raw Message) {
	_, ok := handler.Config.Devices[raw.DeviceID]
	if !ok {
		logrus.Warn("status message received from device " + raw.DeviceID + ", which is not in the config")
	} else {
		logrus.Info("status message received from: " + raw.DeviceID + ", status: " + fmt.Sprint(raw.Contents))
		for k, v := range raw.Contents {
			handler.Config.Devices[raw.DeviceID].Status[k] = v
		}
	}
}

//onConfirmationMsg is the function to process confirmation messages.
func (handler *Handler) onConfirmationMsg(raw Message) {
	value, ok := raw.Contents["completed"]
	if !ok || reflect.TypeOf(value).Kind() != reflect.Bool {
		logrus.Error("received improperly structured confirmation message from device " + raw.DeviceID)
	} else {
		original, ok := raw.Contents["instructed"]
		if !ok {
			logrus.Error("received improperly structured confirmation message from device " + raw.DeviceID)
		} else {
			msg := original.(map[string]interface{})
			if !value.(bool) {
				logrus.Error("device " + raw.DeviceID + " did not complete instruction " +
					fmt.Sprint(msg["contents"].(map[string]interface{})["instruction"]) +
					" at " + raw.TimeSent)
			} else {
				logrus.Info("device " + raw.DeviceID + " completed instruction " +
					fmt.Sprint(msg["contents"].(map[string]interface{})["instruction"]) +
					" at " + raw.TimeSent)
			}
			// If original message to which device responded with confirmation was sent by front-end,
			// pass confirmation through
			if msg["device_id"] == "front-end" {
				jsonMessage, _ := json.Marshal(raw)
				handler.Communicator.Publish("front-end", string(jsonMessage), 3)
			}
		}
	}
	con := handler.Config.Devices[raw.DeviceID]
	con.Connection = true
	handler.Config.Devices[raw.DeviceID] = con
}

// SendStatus sends all status and connection data of a device to the front-end.
// Information retrieved from config.
func (handler *Handler) SendStatus(deviceID string) {
	message := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "status",
		Contents: map[string]interface{}{
			"id":         handler.Config.Devices[deviceID].ID,
			"status":     handler.Config.Devices[deviceID].Status,
			"connection": handler.Config.Devices[deviceID].Connection,
		},
	}

	jsonMessage, err := json.Marshal(&message)
	if err != nil {
		logrus.Errorf("error occurred while constructing message to publish: %v", err)
	} else {
		logrus.Info("sending status data to front-end: " + fmt.Sprint(message.Contents))
		handler.Communicator.Publish("front-end", string(jsonMessage), 3)
	}

}

// openDoorBeun is the test function for developers to test the door and switch combo.
func (handler *Handler) openDoorBeun(raw Message) {
	logrus.Info("checking if door needs to open based on received status message")
	if raw.DeviceID == "controlBoard" {
		var instruction string
		if raw.Contents["mainSwitch"] == float64(0) {
			instruction = "turn off"
		} else if raw.Contents["mainSwitch"] == float64(1) {
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
				logrus.Errorf("error occurred while constructing message to publish: %v", err)
			} else {
				handler.Communicator.Publish("door", string(jsonMessage), 3)
			}
		}
	}
}

//onInstructionMsg is the function to process instruction messages.
func (handler *Handler) onInstructionMsg(raw Message) {
	logrus.Info("instruction message received from: ", raw.DeviceID, " with instruction: ", raw.Contents["instruction"])
	if raw.Contents["instruction"] == "test all" && raw.DeviceID == "front-end" { // TODO maybe switch again
		message := Message{
			DeviceID: raw.DeviceID,
			TimeSent: time.Now().Format("02-01-2006 15:04:05"),
			Type:     "instruction",
			Contents: map[string]interface{}{
				"instruction": "test",
			},
		}
		jsonMessage, err := json.Marshal(&message)
		if err != nil {
			logrus.Errorf("error occurred while constructing message to publish: %v", err)
		} else {
			handler.Communicator.Publish("client-computers", string(jsonMessage), 3)
		}
	}
	if raw.Contents["instruction"] == "send status" && raw.DeviceID == "front-end" {
		for _, value := range handler.Config.Devices {
			handler.SendStatus(value.ID)
		}
	}
	if raw.Contents["instruction"] == "hint" && raw.DeviceID == "front-end" {
		message := Message{
			DeviceID: raw.DeviceID,
			TimeSent: time.Now().Format("02-01-2006 15:04:05"),
			Type:     "instruction",
			Contents: raw.Contents,
		}
		jsonMessage, err := json.Marshal(&message)
		if err != nil {
			logrus.Errorf("error occurred while constructing message to publish: %v", err)
		} else {
			handler.Communicator.Publish("hint", string(jsonMessage), 3)
		}
	}
}
