package handler

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sciler/config"
	"time"
)

// Message is a type that follows the structure all messages have, described in resources/manuals/message_manual.md
type Message struct {
	DeviceID string      `json:"device_id"`
	TimeSent string      `json:"time_sent"`
	Type     string      `json:"type"`
	Contents interface{} `json:"contents"`
}

// Communicator interface is an interface for communication.Communicator, it had the methods needed in the handler
type Communicator interface {
	Publish(topic string, message string, retrials int)
	Start()
}

// Handler is the mqtt.MessageHandler used in the whole S.C.I.L.E.R. system
type Handler struct {
	Config       config.WorkingConfig
	ConfigFile   string
	Communicator Communicator
}

// NewHandler is the actual MessageHandler called when a message is received
// The function processes the JSON payload and logs and error if this fails
// The Message Mapper is called if the JSON is correct to process the contents
func (handler *Handler) NewHandler(client mqtt.Client, message mqtt.Message) {
	defer func() {
		if r := recover(); r != nil {
			logger.Panicf("Recovered panic: %v", r)
		}
	}()
	var raw Message
	if err := json.Unmarshal(message.Payload(), &raw); err != nil {
		logger.Errorf("invalid JSON received: %v", err)
	} else {
		logger.Debugf("message received: %v", raw)
		handler.msgMapper(raw)
	}
}

// msgMapper sends the message through to the right function, filtering on Message.Type
// If the type is not instruction, status, confirmation, or connection, an error is logged
func (handler *Handler) msgMapper(raw Message) {
	switch raw.Type {
	case "instruction":
		{
			handler.onInstructionMsg(raw)
		}
	case "status":
		{
			handler.onStatusMsg(raw)
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
			logger.Errorf("message received from %s, but no message type could be found for: %v", raw.DeviceID, raw.Type)
		}
	}
}

// onStatusMsg is the function to process status messages.
// front-end sends status message with button = true, then immediately button = false,
// and already updates its own state, faster than the back-end messages are received again,
// so we don't send again, as otherwise it has already reset to false, then receives messages for true false again.
func (handler *Handler) onStatusMsg(raw Message) {
	handler.updateStatus(raw)
	if raw.DeviceID != "front-end" {
		handler.sendStatus(raw.DeviceID)
	}
	handler.sendFrontEndStatus(raw)
	handler.HandleEvent(raw.DeviceID)
	handler.sendEventStatus()
}

// onConnectionMsg is the function to process connection messages from devices
// If the device is in the config, and the message is properly structured, the connection status of the device is updated
// After updating the connection status, the new status is send to the front-end
// If the message is from the front-end, the SendSetup function is called
func (handler *Handler) onConnectionMsg(raw Message) {
	contents := raw.Contents.(map[string]interface{})
	device, ok := handler.Config.Devices[raw.DeviceID]
	if !ok {
		logger.Warnf("connection message received from device %s which is not in the config", raw.DeviceID)
	} else {
		logger.Infof("connection message received from: %s", raw.DeviceID)
		value, ok2 := contents["connection"]
		if !ok2 || reflect.TypeOf(value).Kind() != reflect.Bool {
			logger.Errorf("received improperly structured connection message from device %s", raw.DeviceID)
		} else {
			device.Connection = value.(bool)
			handler.Config.Devices[raw.DeviceID] = device
			logger.Infof("setting connection status of %s to %v", raw.DeviceID, value)
			handler.sendStatus(raw.DeviceID)
			if raw.DeviceID == "front-end" && !value.(bool) { // when a front-end disconnect, check if another front-end is connected (maybe multiple front-ends are running
				handler.SendSetup()
			}
		}
	}
}

// onConfirmationMsg is the function to process confirmation messages.
// If the message is properly structured, the success status of the instruction is logged
func (handler *Handler) onConfirmationMsg(raw Message) {
	contents := raw.Contents.(map[string]interface{})
	value, ok := contents["completed"]
	if !ok || reflect.TypeOf(value).Kind() != reflect.Bool {
		logger.Errorf("received improperly structured confirmation message from device %s (no completed key or completed did not carry a boolean value)", raw.DeviceID)
		return
	}
	original, ok := contents["instructed"]
	if !ok {
		logger.Errorf("received improperly structured confirmation message from device %s (no instructed key)", raw.DeviceID)
		return
	}

	if reflect.TypeOf(original) != reflect.TypeOf(map[string]interface{}{}) {
		logger.Errorf("received improperly structured confirmation message from device %s (instructed key did not carry a map value)", raw.DeviceID)
		return
	}
	msg := original.(map[string]interface{})
	instructionContents, err := getMapSlice(msg["contents"])
	if err != nil {
		logger.Errorf(err.Error())
		return
	}
	handler.forwardConfirmation(instructionContents, raw, value.(bool))

	// If a message is received from a device,
	// it can be concluded that the device has positive connection status,
	// and thus it's connection status is set to true
	handler.connected(raw.DeviceID)
}

// forwardConfirmation sends a received confirmation message on to the front-end
func (handler *Handler) forwardConfirmation(instructionContents []map[string]interface{}, raw Message, value bool) {
	var instructionString string
	for _, instruction := range instructionContents {
		// Each instruction is added to a string for proper logging
		instructionString += fmt.Sprintf("%s", instruction["instruction"])
		// If original message to which device responded with confirmation was sent by front-end,
		// pass confirmation through
		if instruction["instructed_by"] == "front-end" {
			jsonMessage, _ := json.Marshal(raw)
			handler.Communicator.Publish("front-end", string(jsonMessage), 3)
			logger.Infof("sending confirmation to front-end for instruction %v", instruction["instruction"])
		}
	}
	if !value {
		logger.Warnf("device %s did not complete instructions: %v at %v", raw.DeviceID,
			instructionString, raw.TimeSent)
	} else {
		logger.Infof("device %s completed instructions: %v at %v", raw.DeviceID,
			instructionString, raw.TimeSent)
	}
}

// onInstructionMsg is the function to process instruction messages
// If the message is properly structured, the instruction in the message is followed
// Currently, only instructions messages from the front-end are supported
// The actions to take are decided by Message.Contents.instruction
func (handler *Handler) onInstructionMsg(raw Message) {
	logger.Infof("instruction message received from: %s", raw.DeviceID)
	instructions, err := getMapSlice(raw.Contents)
	if err != nil {
		logger.Error(err)
		return
	}

	for _, instruction := range instructions {
		if raw.DeviceID == "front-end" {
			handler.handleInstruction(instruction, "front-end")
		} else {
			logger.Warnf("%s, tried to instruct the back-end, only the front-end is allowed to instruct the back-end", raw.DeviceID)
		}
	}
}

// handleInstruction is the function to process an instruction given the ID of the instructor
func (handler *Handler) handleInstruction(instruction map[string]interface{}, instructor string) {
	switch instruction["instruction"] {
	case "send setup":
		handler.SendSetup()
	case "reset all":
		handler.onResetAll(instructor)
	case "test all":
		handler.onTestAll(instructor)
	case "test device":
		handler.onTestDevice(instruction["device"].(string), instructor)
	case "finish rule":
		handler.onFinishRule(instruction["rule"].(string))
	case "hint":
		handler.onHint(instruction, instructor)
	case "check config":
		handler.onCheckConfig(instruction["config"], instruction["name"].(string))
	case "use config":
		handler.onUseConfig(instruction["config"], instruction["file"].(string))
	default:
		logger.Warnf("%s is an unknown instruction", instruction["instruction"])
	}
}

// onResetAll is the function to process the instruction `reset all`
// reset all is instructed when the reset button is clicked in the front-end
func (handler *Handler) onResetAll(deviceID string) {
	handler.sendInstruction("client-computers", []map[string]string{{
		"instruction":   "reset",
		"instructed_by": deviceID,
	}})
	handler.sendInstruction("front-end", []map[string]string{{
		"instruction":   "reset",
		"instructed_by": deviceID,
	}})
	for _, timer := range handler.Config.Timers {
		_ = timer.Stop()
	}
	handler.Config = config.ReadFile(handler.ConfigFile)
	handler.SendSetup()
}

// onTestAll is the function to process the instruction `test all`
// test all is instructed when the test button is clicked in the front-end
func (handler *Handler) onTestAll(instructor string) {
	handler.sendInstruction("client-computers", []map[string]string{{
		"instruction":   "test",
		"instructed_by": instructor,
	}})
}

// onTestDevice is the function to process the instruction `test device`
// test device is instructed when the test button of a device is clicked in the front-end
func (handler *Handler) onTestDevice(deviceID string, instructor string) {
	handler.sendInstruction(deviceID, []map[string]string{{
		"instruction":   "test",
		"instructed_by": instructor,
	}})
}

// onResetAll is the function to process the instruction `finish rule`
// finish rule is instructed when the "puzzel eindigen" button of a rule is clicked in the front-end
func (handler *Handler) onFinishRule(ruleID string) {
	rule, ok := handler.Config.RuleMap[ruleID]
	if !ok {
		logger.Errorf("could not find rule with id %s in map", ruleID)
	} else {
		rule.Execute(handler)
	}
	handler.sendEventStatus()
}

// onHint is the function to process the instruction `hint`
// hint in instructed when the front-end submits a hint
func (handler *Handler) onHint(jsonData map[string]interface{}, instructor string) {
	handler.sendInstruction(jsonData["topic"].(string), []map[string]string{{
		"instruction":   "hint",
		"value":         jsonData["value"].(string),
		"instructed_by": instructor,
	}})
}

// onCheckConfig is the function to process the instruction `check config`
// checks the config and sends a message containing all errors it could find
func (handler *Handler) onCheckConfig(configToRead interface{}, fileName string) {
	message := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "config",
		Contents: map[string]interface{}{
			"name":   fileName,
			"errors": handler.checkConfig(configToRead),
		},
	}
	jsonMessage, _ := json.Marshal(&message)
	handler.Communicator.Publish("front-end", string(jsonMessage), 3)
}

// onUseConfig is the function to process the instruction `use config`
// checks the config and if it found no errors, it saves and uses this config
func (handler *Handler) onUseConfig(configToRead interface{}, fileName string) {
	// check if the config can be used, if so use and send message
	if len(handler.checkConfig(configToRead)) == 0 {
		handler.useConfig(configToRead, fileName)
		message := Message{
			DeviceID: "back-end",
			TimeSent: time.Now().Format("02-01-2006 15:04:05"),
			Type:     "new config",
			Contents: map[string]string{"name": fileName},
		}
		jsonMessage, _ := json.Marshal(&message)
		handler.Communicator.Publish("front-end", string(jsonMessage), 3)
		handler.SendSetup()
	}
}

// useConfig uses new config and save this config to file
// warning this config should be checked first before using the config
func (handler *Handler) useConfig(configToRead interface{}, fileName string) {
	jsonBytes, _ := json.Marshal(configToRead)
	newConfig, _ := config.ReadJSON(jsonBytes)
	dir, _ := os.Getwd()
	fullFileName := filepath.Join(dir, "back-end", "resources", "production", fileName)
	err := ioutil.WriteFile(fullFileName, jsonBytes, 0644)
	if err != nil {
		logger.Error(err)
	}
	handler.Config = newConfig
	handler.ConfigFile = fullFileName
}

// checkConfig checks the config, if it finds any errors in processing the config,
// it will return a list of all found errors, if this slice is empty,
// the config contains no errors and can be used safely
func (handler *Handler) checkConfig(configToRead interface{}) []string {
	errors := make([]string, 0)
	jsonBytes, err := json.Marshal(configToRead)
	if err != nil {
		errors = append(errors, fmt.Sprintf("level I - JSON error: could not unmarshal json, %v", err))
	} else {
		newConfig, errorList := config.ReadJSON(jsonBytes)

		if newConfig.General.Host != handler.Config.General.Host {
			errorList = append(errorList, "level IV - system error: host: different from current host for front and back-end")
		}
		if newConfig.General.Port != handler.Config.General.Port {
			errorList = append(errorList, "level IV - system error: port: different from current port for front and back-end")
		}
		errors = append(errors, errorList...)
	}
	return errors
}
