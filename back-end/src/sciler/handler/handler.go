package handler

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"reflect"
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

// Communicator interface is an interface for mqtt communication
type Communicator interface {
	Start(handler mqtt.MessageHandler)
	Publish(topic string, message string, retrials int)
}

// Handler is a type that mqqt handlers have
type Handler struct {
	Config       config.WorkingConfig
	Communicator Communicator
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
			handler.SendStatus(raw.DeviceID)
			handler.HandleEvent(raw.DeviceID)
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

// onConnectionMsg is the function to process connection messages.
func (handler *Handler) onConnectionMsg(raw Message) {
	contents := raw.Contents.(map[string]interface{})
	device, ok := handler.Config.Devices[raw.DeviceID]
	if !ok {
		logrus.Error("connection message received from device " + raw.DeviceID + ", which is not in the config")
	} else {
		logrus.Info("connection message received from: ", raw.DeviceID)
		value, ok2 := contents["connection"]
		if !ok2 || reflect.TypeOf(value).Kind() != reflect.Bool {
			logrus.Error("received improperly structured connection message from device " + raw.DeviceID)
		} else {
			device.Connection = value.(bool)
			handler.Config.Devices[raw.DeviceID] = device
			logrus.Info("setting connection status of ", raw.DeviceID, " to ", value)
			handler.SendStatus(raw.DeviceID)
		}
	}
}

// compareType compares a reflect.Kind and a string type and returns an error if not the same
func compareType(valueType reflect.Kind, inputType string) error {
	switch inputType {
	case "string":
		{
			if valueType != reflect.String {
				return fmt.Errorf("status type string expected but %sfound as type", valueType.String())
			}
		}
	case "boolean":
		{
			if valueType != reflect.Bool {
				return fmt.Errorf("status type boolean expected but %s found as type", valueType.String())
			}
		}
	case "numeric":
		{
			if valueType != reflect.Int && valueType != reflect.Float64 {
				return fmt.Errorf("status type numeric expected but %s found as type", valueType.String())
			}
		}
	case "array":
		{
			if valueType != reflect.Slice {
				return fmt.Errorf("status type array/slice expected but %s found as type", valueType.String())
			}
		}
	default:
		// todo custom types
		return fmt.Errorf("custom types like: %s, are not yet implemented", inputType)
	}
	return nil
}

// checkStatusType checks if the type of the status change is correct for the component
func (handler *Handler) checkStatusType(device config.Device, status interface{}, component string) error {
	valueType := reflect.TypeOf(status).Kind()
	if inputType, ok := device.Input[component]; ok {
		if err := compareType(valueType, inputType); err != nil {
			return fmt.Errorf("%v with status %v for component %s", err, status, component)
		}
	} else if output, ok2 := device.Output[component]; ok2 {
		if err := compareType(valueType, output.Type); err != nil {
			return fmt.Errorf("%v with status %v for component %s", err, status, component)
		}
	} else {
		return fmt.Errorf("status message received from component %s, which is not in the config under device %s", component, device.ID)
	}
	return nil
}

//onStatusMsg is the function to process status messages.
func (handler *Handler) onStatusMsg(raw Message) {
	contents := raw.Contents.(map[string]interface{})
	if device, ok := handler.Config.Devices[raw.DeviceID]; ok {
		logrus.Info("status message received from: " + raw.DeviceID + ", status: " + fmt.Sprint(raw.Contents))
		for k, v := range contents {
			err := handler.checkStatusType(*device, v, k)
			if err != nil {
				logrus.Error(err)
			} else {
				handler.Config.Devices[raw.DeviceID].Status[k] = v
			}
		}
	} else {
		logrus.Error("status message received from device ", raw.DeviceID, ", which is not in the config")
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
			for _, instruction := range contents {
				instructionString += fmt.Sprintf("%s", instruction["instruction"])
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
	var message Message
	if _, ok := handler.Config.Devices[deviceID]; ok {
		message = Message{
			DeviceID: "back-end",
			TimeSent: time.Now().Format("02-01-2006 15:04:05"),
			Type:     "status",
			Contents: map[string]interface{}{
				"id":         handler.Config.Devices[deviceID].ID,
				"status":     handler.Config.Devices[deviceID].Status,
				"connection": handler.Config.Devices[deviceID].Connection,
			},
		}
	} else if _, ok2 := handler.Config.Timers[deviceID]; ok2 {
		status, _ := handler.Config.Timers[deviceID].GetTimeLeft()
		logrus.Info(status.Milliseconds())
		message = Message{
			DeviceID: "back-end",
			TimeSent: time.Now().Format("02-01-2006 15:04:05"),
			Type:     "time",
			Contents: map[string]interface{}{
				"id":     handler.Config.Timers[deviceID].ID,
				"status": status.Milliseconds(),
				"state":  handler.Config.Timers[deviceID].State,
			},
		}
	}
	jsonMessage, err := json.Marshal(&message)
	if err != nil {
		logrus.Errorf("error occurred while constructing message to publish: %v", err)
	} else {
		logrus.Info("sending status data to front-end: " + fmt.Sprint(message.Contents))
		handler.Communicator.Publish("front-end", string(jsonMessage), 3)
	}
}

// SendInstruction sends a list of instructions to a client
func (handler *Handler) SendInstruction(clientID string, instructions []config.ComponentInstruction) {
	message := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: instructions,
	}

	jsonMessage, err := json.Marshal(&message)
	if err != nil {
		logrus.Errorf("error occurred while constructing message to publish: %v", err)
	} else {
		logrus.Infof("sending instruction data to %s: %s", clientID, fmt.Sprint(message.Contents))
		handler.Communicator.Publish(clientID, string(jsonMessage), 3)
	}
}

// onInstructionMsg is the function to process instruction messages.
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

// HandleEvent is a function that checks and possible executes all rules according to the given (device/rule/timer) id
func (handler *Handler) HandleEvent(id string) {
	if rules, ok := handler.Config.StatusMap[id]; ok {
		for _, rule := range rules {
			if rule.Executed < rule.Limit && rule.Conditions.Resolve(handler.Config) {
				rule.Execute(handler)
			}
		}
	}
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

// GetStatus asks devices to send status
func (handler *Handler) GetStatus(deviceID string) {
	message := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{"instruction": "status update"},
		},
	}

	jsonMessage, err := json.Marshal(&message)
	if err != nil {
		logrus.Errorf("error occurred while constructing message to publish: %v", err)
	} else {
		logrus.Info("sending status request to client computer: ", deviceID, fmt.Sprint(message.Contents))
		handler.Communicator.Publish(deviceID, string(jsonMessage), 3)
	}
}

// SetTimer starts given timer
func (handler *Handler) SetTimer(timerID string, instructions config.ComponentInstruction) {
	switch instructions.Instruction {
	case "start":
		handler.Config.Timers[timerID].Start(handler)
	case "pause":
		handler.Config.Timers[timerID].Pause()
	case "add": // TODO: implement timer Add
	case "subtract": // TODO: implement timer subtract
	case "stop":
		handler.Config.Timers[timerID].Stop()
	default:
		logrus.Warn("error occurred while reading timer instruction message: %v", instructions.Instruction)
	}
	handler.SendStatus(timerID)

}
