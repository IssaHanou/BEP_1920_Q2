package handler

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	logger "github.com/sirupsen/logrus"
	"reflect"
	"sciler/config"
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
	Publish(topic string, message string, retrials int)
	Start()
}

// Handler is a type that mqqt handlers have
type Handler struct {
	Config       config.WorkingConfig
	ConfigFile   string
	Communicator Communicator
}

// NewHandler is the actual MessageHandler
func (handler *Handler) NewHandler(client mqtt.Client, message mqtt.Message) {
	var raw Message
	if err := json.Unmarshal(message.Payload(), &raw); err != nil {
		logger.Errorf("invalid JSON received: %v", err)
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
			logger.Error("message received from ", raw.DeviceID,
				", but no message type could be found for: ", raw.Type)
		}
	}
}

// onStatusMsg is the function to process status messages.
func (handler *Handler) onStatusMsg(raw Message) {
	handler.updateStatus(raw)
	handler.sendStatus(raw.DeviceID)
	handler.HandleEvent(raw.DeviceID)
	handler.sendEventStatus()
}

// onConnectionMsg is the function to process connection messages.
func (handler *Handler) onConnectionMsg(raw Message) {
	contents := raw.Contents.(map[string]interface{})
	device, ok := handler.Config.Devices[raw.DeviceID]
	if !ok {
		logger.Error("connection message received from device " + raw.DeviceID + ", which is not in the config")
	} else {
		logger.Info("connection message received from: ", raw.DeviceID)
		value, ok2 := contents["connection"]
		if !ok2 || reflect.TypeOf(value).Kind() != reflect.Bool {
			logger.Error("received improperly structured connection message from device " + raw.DeviceID)
		} else {
			device.Connection = value.(bool)
			handler.Config.Devices[raw.DeviceID] = device
			logger.Info("setting connection status of ", raw.DeviceID, " to ", value)
			handler.sendStatus(raw.DeviceID)
			if raw.DeviceID == "front-end" && !value.(bool) { // when a front-end disconnect, check if another front-end is connected (maybe multiple front-ends are running
				handler.SendSetup()
			}
		}
	}
}

// onConfirmationMsg is the function to process confirmation messages.
func (handler *Handler) onConfirmationMsg(raw Message) {
	contents := raw.Contents.(map[string]interface{})
	value, ok := contents["completed"]
	if !ok || reflect.TypeOf(value).Kind() != reflect.Bool {
		logger.Errorf("received improperly structured confirmation message from device " + raw.DeviceID)
		return
	}
	original, ok := contents["instructed"]
	if !ok {
		logger.Errorf("received improperly structured confirmation message from device " + raw.DeviceID)
		return
	}
	msg := original.(map[string]interface{})
	instructionContents, err := getMapSlice(msg["contents"])
	if err != nil {
		logger.Errorf(err.Error())
		return
	}

	var instructionString string
	for _, instruction := range instructionContents {
		instructionString += fmt.Sprintf("%s", instruction["instruction"])
		// If original message to which device responded with confirmation was sent by front-end,
		// pass confirmation through
		if instruction["instructed_by"] == "front-end" {
			jsonMessage, _ := json.Marshal(raw)
			handler.Communicator.Publish("front-end", string(jsonMessage), 3)
			logger.Infof("sending confirmation to front-end for instruction %v", instruction["instruction"])
		}
	}

	if !value.(bool) {
		logger.Warn("device " + raw.DeviceID + " did not complete instructions: " +
			instructionString + " at " + raw.TimeSent)
	} else {
		logger.Info("device " + raw.DeviceID + " completed instructions: " +
			instructionString + " at " + raw.TimeSent)
	}

	con, ok := handler.Config.Devices[raw.DeviceID]
	if !ok {
		logger.Errorf("device %s was not found in config", raw.DeviceID)
	} else {
		con.Connection = true
		handler.Config.Devices[raw.DeviceID] = con
	}
}

// onInstructionMsg is the function to process instruction messages.
func (handler *Handler) onInstructionMsg(raw Message) {
	logger.Info("instruction message received from: ", raw.DeviceID)
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
	case "send status":
		handler.onSendStatus()
	case "reset all":
		handler.onResetAll(instructor)
	case "test all":
		handler.onTestAll(instructor)
	case "test device":
		handler.onTestDevice(instruction["device"].(string), instructor)
	case "finish rule":
		handler.onFinishRule(instruction["rule"].(string))
	case "hint":
		handler.onHint(instruction["value"].(string), instructor)
	case "check config":
		handler.processConfig(instruction["config"], "check", "")
	case "use config":
		handler.processConfig(instruction["config"], "use", instruction["file"].(string))
	}
}

// onSendStatus is the function to process the instruction `send status`
func (handler *Handler) onSendStatus() {
	for _, device := range handler.Config.Devices {
		handler.sendStatus(device.ID)
	}
	for _, timer := range handler.Config.Timers {
		handler.sendStatus(timer.ID)
	}
	handler.sendEventStatus()
}

// onResetAll is the function to process the instruction `reset all`
func (handler *Handler) onResetAll(deviceID string) {
	handler.SendInstruction("client-computers", []map[string]string{{
		"instruction":   "reset",
		"instructed_by": deviceID,
	}})
	handler.SendInstruction("front-end", []map[string]string{{
		"instruction":   "reset",
		"instructed_by": deviceID,
	}})
	handler.Config = config.ReadFile(handler.ConfigFile)
	handler.SendSetup()
}

// onTestAll is the function to process the instruction `test all`
func (handler *Handler) onTestAll(instructor string) {
	handler.SendInstruction("client-computers", []map[string]string{{
		"instruction":   "test",
		"instructed_by": instructor,
	}})
}

// onTestDevice is the function to process the instruction `test device`
func (handler *Handler) onTestDevice(deviceID string, instructor string) {
	handler.SendInstruction(deviceID, []map[string]string{{
		"instruction":   "test",
		"instructed_by": instructor,
	}})
}

// onResetAll is the function to process the instruction `finish rule`
func (handler *Handler) onFinishRule(ruleID string) {
	rule, ok := handler.Config.RuleMap[ruleID]
	if !ok {
		logger.Errorf("could not find rule with id %s in map", ruleID)
	}
	rule.Execute(handler)
	handler.sendEventStatus()
}

// onHint is the function to process the instruction `hint`
func (handler *Handler) onHint(hint string, instructor string) {
	handler.SendInstruction("hint", []map[string]string{{
		"instruction":   "hint",
		"value":         hint,
		"instructed_by": instructor,
	}})
}
