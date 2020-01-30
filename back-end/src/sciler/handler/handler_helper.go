package handler

import (
	"encoding/json"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"reflect"
	"sciler/config"
	"strings"
	"time"
)

// SendSetup sends the general set-up information to the front-end.
// This includes:
// a message with the name, all hints, event descriptions, cameras and config buttons
// a message for each device and the general timer
// and message to send the event statuses
func (handler *Handler) SendSetup() {
	message := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "setup",
		Contents: map[string]interface{}{
			"name":    handler.Config.General.Name,
			"hints":   handler.getHints(),
			"events":  handler.getEventDescriptions(),
			"cameras": handler.getCameras(),
			"buttons": handler.getButtons(),
		},
	}
	jsonMessage, _ := json.Marshal(&message)
	handler.Communicator.Publish("front-end", string(jsonMessage), 3)
	logger.Info("published setup data to front-end")
	handler.sendStatus("general")
	for _, value := range handler.Config.Devices {
		handler.sendStatus(value.ID)
		// After sending device information, a request to the device is send to publish the current status,
		// its return message will be handled in a new pipeline
		// This way, the front-end will receive all device information from the config,
		// and after that the current status from all the connected devices
		handler.GetStatus(value.ID)
	}
	handler.sendEventStatus()
}

// SendComponentInstruction sends a list of instructions to a client, with a delay if given a valid duration.
// This function is called in response to the execution of a rule
// param clientID  is the target for the message (topic to publish to)
// param instructions are the contends with instructions for the device
// param delay is the possible delay it will wait to send the message
// if delay is not properly structured (XhXmXs), no error will be given, the function will continue as if no delay was given
func (handler *Handler) SendComponentInstruction(clientID string, instructions []config.ComponentInstruction, delay string) {
	message := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: instructions,
	}
	// If the instruction is to reset the status of a front-end button, update its status in the config.
	if clientID == "front-end" {
		for _, instruction := range instructions {
			handler.Config.Devices["front-end"].Status[instruction.ComponentID] = instruction.Value
		}
	}
	jsonMessage, _ := json.Marshal(&message)
	delayDur, err := time.ParseDuration(delay)
	if err == nil {
		go func() {
			logger.Infof("waiting %s to send instruction data to %s: %s", delay, clientID, fmt.Sprint(message.Contents))
			time.Sleep(delayDur)
			logger.Infof("sending instruction data to %s after waiting %s: %s", clientID, delay, fmt.Sprint(message.Contents))
			handler.Communicator.Publish(clientID, string(jsonMessage), 3)
		}()
	} else {
		logger.Infof("sending instructions to %s: %s", clientID, fmt.Sprint(message.Contents))
		handler.Communicator.Publish(clientID, string(jsonMessage), 3)
	}
}

// PrepareMessage scans a message and if the instruction is of type status, if so the value of the message is replaced by the status of a device
func (handler *Handler) PrepareMessage(typeID string, messages []config.ComponentInstruction) []config.ComponentInstruction {
	res := make([]config.ComponentInstruction, len(messages))
	device := handler.Config.Devices[typeID]
	for i, message := range messages {
		msg := message
		instructionType := device.Output[message.ComponentID].Instructions[message.Instruction]
		if instructionType == "status" { // when the instruction type is status
			split := strings.Split(message.Value.(string), ".")
			deviceID := split[0]
			componentID := split[1]
			status := handler.Config.Devices[deviceID].Status[componentID]
			msg.Value = status // set status of message to retrieved status
		}
		res[i] = msg
	}
	return res
}

// SendLabelInstruction provides the action with a componentID from de LabelMap and a device to send it to
// This function is called in response to the execution of a rule
// SendComponentInstruction is called for each component in the LabelMap under labelID
// Each instruction needs a componentID to execute on the device,
// this is why the original instructions can not be immediately passed to the SendComponentInstruction function
// param labelID is the target label for the message
// param instructions are the contends with instructions for the device
// param delay is the possible delay it will wait to send the message
func (handler *Handler) SendLabelInstruction(labelID string, instructions []config.ComponentInstruction, delay string) {
	for _, instruction := range instructions {
		for _, comp := range handler.Config.LabelMap[labelID] {
			instruction.ComponentID = comp.ID
			handler.SendComponentInstruction(comp.Device.ID, []config.ComponentInstruction{instruction}, delay)
		}
	}
}

// sendInstruction sends a list of instructions to a client
// This function is called in reaction to a instruction message to be back-end
// param cliendID  is the target for the message (topic to publish to)
// param instructions are the contends with instructions for the device
func (handler *Handler) sendInstruction(clientID string, instructions []map[string]string) {
	message := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: instructions,
	}
	jsonMessage, _ := json.Marshal(&message)
	logger.Infof("sending instructions to %s: %s", clientID, fmt.Sprint(message.Contents))
	handler.Communicator.Publish(clientID, string(jsonMessage), 3)
}

// updateStatus is the function to process status messages.
// If the device is in the config, and the status types are correct, the device status gets updated
// if the device is the front-end, handleFrontEndStatus() is called
func (handler *Handler) updateStatus(raw Message) {
	contents := raw.Contents.(map[string]interface{})
	if device, ok := handler.Config.Devices[raw.DeviceID]; ok {
		logger.Infof("status message received from: %s", raw.DeviceID)
		if device.ID == "front-end" {
			handler.handleFrontEndStatus(contents)
		}

		for k, v := range contents {
			err := handler.checkStatusType(*device, v, k)
			if err != nil {
				logger.Error(err)
			} else {
				handler.Config.Devices[raw.DeviceID].Status[k] = v
			}
		}
	} else {
		logger.Warnf("status message received from device %s which is not in the config", raw.DeviceID)
	}
}

// sendStatus sends all status and connection data of device/timer deviceID to the front-end
// For devices the status of components and the connection status is send
// For timers the duration left and the state are send
// param deviceID can be the ID of a device or a timer
func (handler *Handler) sendStatus(ID string) {
	if device, ok := handler.Config.Devices[ID]; ok {
		handler.sendStatusDevice(device)
	} else if timer, ok2 := handler.Config.Timers[ID]; ok2 {
		handler.sendStatusTimer(timer)
	} else {
		logger.Errorf("error occurred while sending status of %s, since it is not recognised as a device or timer", ID)
		return
	}
}

// sendStatusDevice sends status of a device to the front-end
func (handler *Handler) sendStatusDevice(device *config.Device) {
	message := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "status",
		Contents: map[string]interface{}{
			"id":         device.ID,
			"status":     device.Status,
			"connection": device.Connection,
		},
	}
	jsonMessage, _ := json.Marshal(&message)
	logger.Infof("sending status data to front-end: %v", message.Contents)
	handler.Communicator.Publish("front-end", string(jsonMessage), 3)
}

// sendStatusTimer sends status of a timer to the front-end
func (handler *Handler) sendStatusTimer(timer *config.Timer) {
	duration, _ := timer.GetTimeLeft()
	message := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "time",
		Contents: map[string]interface{}{
			"id":       timer.ID,
			"duration": duration.Milliseconds(),
			"state":    timer.State,
		},
	}
	jsonMessage, _ := json.Marshal(&message)
	logger.Infof("sending status data to front-end: %v", message.Contents)
	handler.Communicator.Publish("front-end", string(jsonMessage), 3)
	handler.Communicator.Publish("time", string(jsonMessage), 3)
}

// HandleEvent is a function that checks and possible executes all rules according to the given device/rule/timer
// The StatusMap has all the rules containing id in a condition
// if the rule is below its execution limit and all conditions resolve to true, the rule actions will be executed
// param id is the ID of a device, rule, or timer which has had a status update
func (handler *Handler) HandleEvent(id string) {
	if rules, ok := handler.Config.StatusMap[id]; ok {
		for _, rule := range rules {
			if (rule.Executed < rule.Limit || rule.Limit == 0) && rule.Conditions.Resolve(handler.Config) {
				rule.Execute(handler)
			}
		}
	}
}

// sendEventStatus sends the status of events to the front-end
// The status contents is a list of eventIDs with a boolean finish status
func (handler *Handler) sendEventStatus() {
	status := handler.getEventStatus()
	message := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "event status",
		Contents: status,
	}
	jsonMessage, _ := json.Marshal(&message)
	logger.Infof("sending event status to front-end")
	handler.Communicator.Publish("front-end", string(jsonMessage), 3)
}

// Handles status updates from front-end specifically as these will only be affected by button events,
// which are handled differently.
// For all components of the front-end status, check if the status is not empty, and whether it is a boolean (gameState is string).
// If the button is pressed (new status = true) and its status was false,
// the rule that belongs to the pressed button is executed
func (handler *Handler) handleFrontEndStatus(contents map[string]interface{}) {
	device, _ := handler.Config.Devices["front-end"]
	for component, status := range device.Status {
		rule, _ := handler.Config.RuleMap[component]
		if newStatus, ok := contents[component]; ok && reflect.TypeOf(newStatus).Kind() == reflect.Bool {
			if newStatus.(bool) && !status.(bool) {
				rule.Execute(handler)
			}
		}
	}
}

// Sends the front-end the new status of disabled buttons
func (handler *Handler) sendFrontEndStatus(message Message) {
	returnButtons := handler.getButtons()
	newMessage := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "front-end status",
		Contents: returnButtons,
	}
	jsonMessage, _ := json.Marshal(&newMessage)
	logger.Infof("sending front-end button status")
	handler.Communicator.Publish("front-end", string(jsonMessage), 3)
}

// getEventStatus returns a json list with json objects with keys ["id", "status"]
// status is json object with key ruleName and value true (if executed == limit) or false
func (handler *Handler) getEventStatus() []map[string]interface{} {
	var list []map[string]interface{}
	for _, rule := range handler.Config.EventRuleMap {
		var status = make(map[string]interface{})
		status["id"] = rule.ID
		status["status"] = rule.Finished()
		list = append(list, status)
	}
	return list
}

// getHints returns a map of hints with puzzle name as key and list of hints for that puzzle as value
func (handler *Handler) getHints() map[string][]string {
	hints := make(map[string][]string)
	for _, puzzle := range handler.Config.Puzzles {
		hints[puzzle.Event.Name] = puzzle.Hints
	}
	return hints
}

// getEventDescriptions returns a map of hints with puzzle name as key and list of hints for that puzzle as value
func (handler *Handler) getEventDescriptions() map[string]string {
	events := make(map[string]string)
	for _, rule := range handler.Config.EventRuleMap {
		events[rule.ID] = rule.Description
	}
	return events
}

// getCameras returns a map with camera name and camera link
func (handler *Handler) getCameras() []map[string]string {
	var cameras []map[string]string
	for _, camera := range handler.Config.Cameras {
		result := make(map[string]string)
		result["name"] = camera.Name
		result["link"] = camera.Link
		cameras = append(cameras, result)
	}
	return cameras
}

// getButtons returns a list with button names
func (handler *Handler) getButtons() []map[string]interface{} {
	var buttons []map[string]interface{}
	for _, btn := range handler.Config.ButtonEvents {
		rule, _ := handler.Config.RuleMap[btn.ID]
		button := make(map[string]interface{})
		button["id"] = btn.ID
		button["disabled"] = !btn.Conditions.Resolve(handler.Config) || rule.Finished()
		buttons = append(buttons, button)
	}
	return buttons
}

// GetStatus asks devices to send status
// An "instruction": "status update" is published to the topic deviceID
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
	logger.Infof("sending status request to client computer: %s", deviceID)
	handler.Communicator.Publish(deviceID, string(jsonMessage), 3)
}

// SetTimer processes timer actions
// param timerID is the ID of the timer that needs to change state
// param instructions has a instructions.Instruction and instructions.Value to do the instruction
func (handler *Handler) SetTimer(timerID string, instructions config.ComponentInstruction) {
	err := error(nil)
	switch instructions.Instruction {
	case "start":
		err = handler.Config.Timers[timerID].Start(handler)
	case "pause":
		err = handler.Config.Timers[timerID].Pause()
	case "add":
		duration, durErr := time.ParseDuration(fmt.Sprintf("%v", instructions.Value))
		if durErr == nil {
			err = handler.Config.Timers[timerID].AddSubTime(handler, duration, true)
		} else {
			err = fmt.Errorf("could not parse %v to duration to add for timer %v", instructions.Value, timerID)
		}
	case "subtract":
		duration, durErr := time.ParseDuration(fmt.Sprintf("%v", instructions.Value))
		if durErr == nil {
			err = handler.Config.Timers[timerID].AddSubTime(handler, duration, false)
		} else {
			err = fmt.Errorf("could not parse %v to duration to subtract for timer %v", instructions.Value, timerID)
		}
	case "stop":
		err = handler.Config.Timers[timerID].Stop()
	case "done":
		err = handler.Config.Timers[timerID].Done(handler)
	default:
		err = fmt.Errorf("error occurred while reading timer instruction message: %v", instructions.Instruction)
	}
	if err != nil {
		logger.Error(err)
	}
	handler.sendStatus(timerID)
}

// compareType compares a reflect.Kind and a string type and returns an error if not the same
func compareType(valueType reflect.Kind, inputType string) error {
	switch inputType {
	case "string":
		if valueType != reflect.String {
			return fmt.Errorf("status type string expected but %s found as type", valueType.String())
		}
	case "boolean":
		if valueType != reflect.Bool {
			return fmt.Errorf("status type boolean expected but %s found as type", valueType.String())
		}
	case "numeric":
		if valueType != reflect.Int && valueType != reflect.Float64 {
			return fmt.Errorf("status type numeric expected but %s found as type", valueType.String())
		}
	case "array":
		if valueType != reflect.Slice {
			return fmt.Errorf("status type array/slice expected but %s found as type", valueType.String())
		}
	default:
		return fmt.Errorf("custom types like: %s, are not yet implemented", inputType)
	}
	return nil
}

// checkStatusType checks if the type of the status change is correct for the component
// param device the config device where component is located
// param status the incoming status needed to be type-checked
// param component the component that should get that status
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

// getMapSlice is a dirty trick to go from `interface{}` to `[]map[string]interface{}`
func getMapSlice(input interface{}) ([]map[string]interface{}, error) {
	bytes, _ := json.Marshal(input)
	var output []map[string]interface{}
	err := json.Unmarshal(bytes, &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// connected is a method that sets a device to connected
func (handler *Handler) connected(deviceID string) {
	device, ok := handler.Config.Devices[deviceID]
	if !ok {
		logger.Warnf("device %s was not found in config", deviceID)
	} else {
		device.Connection = true
		handler.Config.Devices[deviceID] = device
	}
}
