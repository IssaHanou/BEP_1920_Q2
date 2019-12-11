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
	DeviceID string      `json:"device_id"`
	TimeSent string      `json:"time_sent"`
	Type     string      `json:"type"`
	Contents interface{} `json:"contents"`
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
			//handler.openDoorBeun(raw)
			handler.SendStatus(raw.DeviceID)
			handler.handleEvent(raw.DeviceID)
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
	contents := raw.Contents.(map[string]interface{})
	device, ok := handler.Config.Devices[raw.DeviceID]
	if !ok {
		logrus.Warn("connection message received from device " + raw.DeviceID + ", which is not in the config")
	} else {
		logrus.Info("connection message received from: ", raw.DeviceID)
		value, ok2 := contents["connection"]
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
	contents := raw.Contents.(map[string]interface{})
	_, ok := handler.Config.Devices[raw.DeviceID]
	if !ok {
		logrus.Warn("status message received from device " + raw.DeviceID + ", which is not in the config")
	} else {
		logrus.Info("status message received from: " + raw.DeviceID + ", status: " + fmt.Sprint(raw.Contents))
		for k, v := range contents {
			handler.Config.Devices[raw.DeviceID].Status[k] = v
		}
	}
}

//onConfirmationMsg is the function to process confirmation messages.
func (handler *Handler) onConfirmationMsg(raw Message) {
	contents := raw.Contents.(map[string]interface{})
	value, ok := contents["completed"]
	if !ok || reflect.TypeOf(value).Kind() != reflect.Bool {
		logrus.Error("received improperly structured confirmation message from device " + raw.DeviceID)
	} else {
		original, ok := contents["instructed"]
		if !ok {
			logrus.Error("received improperly structured confirmation message from device " + raw.DeviceID)
		} else {
			msg := original.(map[string]interface{})
			contents, err := getMapSlice(msg["contents"])
			if err != nil {
				logrus.Error(err)
				return
			}

			var instructionString string
			for i, instruction := range contents {
				instructionString += fmt.Sprintf("%d: %s ", i, instruction["instruction"])
			}

			if !value.(bool) {
				logrus.Error("device " + raw.DeviceID + " did not complete instructions: " +
					instructionString + "at " + raw.TimeSent)
			} else {
				logrus.Info("device " + raw.DeviceID + " completed instructions: " +
					instructionString + "at " + raw.TimeSent)
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
	contents := raw.Contents.(map[string]interface{})
	logrus.Info("checking if door needs to open based on received status message")
	if raw.DeviceID == "controlBoard" {
		var instruction bool
		if contents["mainSwitch"] == float64(0) {
			instruction = false
		} else if contents["mainSwitch"] == float64(1) {
			instruction = true
		} else {
			return
		}

		message := Message{
			DeviceID: "back-end",
			TimeSent: time.Now().Format("02-01-2006 15:04:05"),
			Type:     "instruction",
			Contents: []map[string]interface{}{
				{
					"instruction": "open",
					"value":       instruction,
				},
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

//onInstructionMsg is the function to process instruction messages.
func (handler *Handler) onInstructionMsg(raw Message) {
	logrus.Info("instruction message received from: ", raw.DeviceID)

	instructions, err := getMapSlice(raw.Contents)
	if err != nil {
		logrus.Error(err)
		return
	}

	for _, instruction := range instructions {
		if instruction["instruction"] == "test all" && raw.DeviceID == "front-end" { // TODO maybe switch again
			message := Message{
				DeviceID: raw.DeviceID,
				TimeSent: time.Now().Format("02-01-2006 15:04:05"),
				Type:     "instruction",
				Contents: []map[string]interface{}{
					{"instruction": "test"},
				},
			}
			jsonMessage, err := json.Marshal(&message)
			if err != nil {
				logrus.Errorf("error occurred while constructing message to publish: %v", err)
			} else {
				handler.Communicator.Publish("client-computers", string(jsonMessage), 3)
			}
		}
		if instruction["instruction"] == "send status" && raw.DeviceID == "front-end" {
			for _, value := range handler.Config.Devices {
				handler.SendStatus(value.ID)
			}
		}
		if instruction["instruction"] == "hint" && raw.DeviceID == "front-end" {
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
}

func (handler *Handler) handleEvent(id string) {

}

func getMapSlice(input interface{}) ([]map[string]interface{}, error) {
	bytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	var output []map[string]interface{} // dirty trick to go from interface{} to []map[string]interface{}
	err = json.Unmarshal(bytes, &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}
