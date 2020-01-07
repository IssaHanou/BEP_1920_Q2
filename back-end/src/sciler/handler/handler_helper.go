package handler

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
	"sciler/config"
	"time"
)

// SendSetUp sends the general set-up information to the front-end.
// This includes the name, all hints and event descriptions
// Statuses are also sent
func (handler *Handler) SendSetUp() {
	message := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "setup",
		Contents: map[string]interface{}{
			"name":   handler.Config.General.Name,
			"hints":  handler.GetHints(),
			"events": handler.GetEventDescriptions(),
		},
	}
	jsonMessage, _ := json.Marshal(&message)
	handler.Communicator.Publish("front-end", string(jsonMessage), 3)
	handler.SendStatus("general")
	for _, value := range handler.Config.Devices {
		handler.SendStatus(value.ID)
		handler.GetStatus(value.ID)
	}
	handler.SendEventStatus()
}

// SendInstruction sends a list of instructions to a client
func (handler *Handler) SendInstruction(clientID string, instructions []config.ComponentInstruction) {
	message := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: instructions,
	}
	jsonMessage, _ := json.Marshal(&message)
	logrus.Infof("sending instruction data to %s: %s", clientID, fmt.Sprint(message.Contents))
	handler.Communicator.Publish(clientID, string(jsonMessage), 3)
}

// UpdateStatus is the function to process status messages.
func (handler *Handler) UpdateStatus(raw Message) {
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

// SendStatus sends all status and connection data of a device to the front-end.
// Information retrieved from config.
func (handler *Handler) SendStatus(deviceID string) {
	var message Message
	if device, ok := handler.Config.Devices[deviceID]; ok {
		message = Message{
			DeviceID: "back-end",
			TimeSent: time.Now().Format("02-01-2006 15:04:05"),
			Type:     "status",
			Contents: map[string]interface{}{
				"id":         device.ID,
				"status":     device.Status,
				"connection": device.Connection,
			},
		}
	} else if timer, ok2 := handler.Config.Timers[deviceID]; ok2 {
		duration, _ := timer.GetTimeLeft()
		message = Message{
			DeviceID: "back-end",
			TimeSent: time.Now().Format("02-01-2006 15:04:05"),
			Type:     "time",
			Contents: map[string]interface{}{
				"id":       timer.ID,
				"duration": duration.Milliseconds(),
				"state":    timer.State,
			},
		}
	} else {
		logrus.Errorf("error occurred while sending status of %s, since it is not recognised as a device or timer", deviceID)
		return
	}
	jsonMessage, _ := json.Marshal(&message)
	logrus.Info("sending status data to front-end: " + fmt.Sprint(message.Contents))
	handler.Communicator.Publish("front-end", string(jsonMessage), 3)
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

// SendEventStatus sends the status of events to the front-end
func (handler *Handler) SendEventStatus() {
	status := handler.GetEventStatus()
	message := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "event status",
		Contents: status,
	}
	jsonMessage, _ := json.Marshal(&message)
	logrus.Info("sending event status to front-end")
	handler.Communicator.Publish("front-end", string(jsonMessage), 3)
}

// GetEventStatus returns json list with json objects with keys ["id", "status"]
// status is json object with key ruleName and value true (if executed == limit) or false
func (handler *Handler) GetEventStatus() []map[string]interface{} {
	var list []map[string]interface{}
	for _, rule := range handler.Config.RuleMap {
		var status = make(map[string]interface{})
		status["id"] = rule.ID
		status["status"] = rule.Finished()
		list = append(list, status)
	}
	return list
}

// GetHints returns map of hints with puzzle name as key and list of hints for that puzzle as value
func (handler *Handler) GetHints() map[string][]string {
	hints := make(map[string][]string)
	for _, puzzle := range handler.Config.Puzzles {
		hints[puzzle.Event.Name] = puzzle.Hints
	}
	return hints
}

// GetHints returns map of hints with puzzle name as key and list of hints for that puzzle as value
func (handler *Handler) GetEventDescriptions() map[string]string {
	events := make(map[string]string)
	for _, rule := range handler.Config.RuleMap {
		events[rule.ID] = rule.Description
	}
	return events
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
	jsonMessage, _ := json.Marshal(&message)
	logrus.Info("sending status request to client computer: ", deviceID, fmt.Sprint(message.Contents))
	handler.Communicator.Publish(deviceID, string(jsonMessage), 3)
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
		logrus.Warnf("error occurred while reading timer instruction message: %v", instructions.Instruction)
	}
	handler.SendStatus(timerID)
}

// compareType compares a reflect.Kind and a string type and returns an error if not the same
func compareType(valueType reflect.Kind, inputType string) error {
	switch inputType {
	case "string":
		{
			if valueType != reflect.String {
				return fmt.Errorf("status type string expected but %s found as type", valueType.String())
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

func getMapSlice(input interface{}) ([]map[string]interface{}, error) {
	bytes, _ := json.Marshal(input)
	var output []map[string]interface{} // dirty trick to go from interface{} to []map[string]interface{}
	err := json.Unmarshal(bytes, &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}
