package handler

import (
	"bytes"
	"encoding/gob"
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
	Start(handler mqtt.MessageHandler, onStart func())
	Publish(topic string, message string, retrials int)
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
			handler.updateStatus(raw)
			handler.sendStatus(raw.DeviceID)
			handler.HandleEvent(raw.DeviceID)
			handler.sendEventStatus()
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
			handler.sendStatus(raw.DeviceID)
		}
	}
}

// onConfirmationMsg is the function to process confirmation messages.
func (handler *Handler) onConfirmationMsg(raw Message) {
	contents := raw.Contents.(map[string]interface{})
	value, ok := contents["completed"]
	if !ok || reflect.TypeOf(value).Kind() != reflect.Bool {
		logrus.Errorf("received improperly structured confirmation message from device " + raw.DeviceID)
		return
	}
	original, ok := contents["instructed"]
	if !ok {
		logrus.Errorf("received improperly structured confirmation message from device " + raw.DeviceID)
		return
	}
	msg := original.(map[string]interface{})
	instructionContents, err := getMapSlice(msg["contents"])
	if err != nil {
		logrus.Errorf(err.Error())
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
			logrus.Infof("sending confirmation to front-end for instruction %v", instruction["instruction"])
		}
	}

	if !value.(bool) {
		logrus.Warn("device " + raw.DeviceID + " did not complete instructions: " +
			instructionString + " at " + raw.TimeSent)
	} else {
		logrus.Info("device " + raw.DeviceID + " completed instructions: " +
			instructionString + " at " + raw.TimeSent)
	}

	con, ok := handler.Config.Devices[raw.DeviceID]
	if !ok {
		logrus.Errorf("device %s was not found in config", raw.DeviceID)
	} else {
		con.Connection = true
		handler.Config.Devices[raw.DeviceID] = con
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
		if raw.DeviceID == "front-end" {
			switch instruction["instruction"] {
			case "send setup":
				{
					handler.SendSetup()
				}
			case "send status":
				{
					for _, device := range handler.Config.Devices {
						handler.sendStatus(device.ID)
					}
					for _, timer := range handler.Config.Timers {
						handler.sendStatus(timer.ID)
					}
					handler.sendEventStatus()
				}
			case "reset all":
				{
					handler.SendInstruction("client-computers", []map[string]string{{
						"instruction":   "reset",
						"instructed_by": raw.DeviceID,
					}})
					handler.SendInstruction("front-end", []map[string]string{{
						"instruction":   "reset",
						"instructed_by": raw.DeviceID,
					}})
					handler.Config = config.ReadFile(handler.ConfigFile)
					handler.sendStatus("general")
				}
			case "test all":
				{
					handler.SendInstruction("client-computers", []map[string]string{{
						"instruction":   "test",
						"instructed_by": raw.DeviceID,
					}})
				}
			case "test device":
				{
					handler.SendInstruction(instruction["device"].(string), []map[string]string{{
						"instruction":   "test",
						"instructed_by": raw.DeviceID,
					}})
				}
			case "finish rule":
				{
					ruleToFinish := instruction["rule"].(string)
					rule, ok := handler.Config.RuleMap[ruleToFinish]
					if !ok {
						logrus.Errorf("could not find rule with id %s in map", ruleToFinish)
					}
					rule.Execute(handler)
					handler.sendEventStatus()
				}
			case "hint":
				{
					message := Message{
						DeviceID: "back-end",
						TimeSent: time.Now().Format("02-01-2006 15:04:05"),
						Type:     "instruction",
						Contents: []map[string]interface{}{{
							"instruction":   "hint",
							"value":         instruction["value"],
							"instructed_by": raw.DeviceID},
						},
					}
					jsonMessage, _ := json.Marshal(&message)
					handler.Communicator.Publish("hint", string(jsonMessage), 3)
				}
			case "check config":
				{
					handler.ProcessConfig(instruction["config"], "check")
				}
			case "use config":
				{
					handler.ProcessConfig(instruction["config"], "use")
				}
			}
		} else {
			logrus.Warnf("%s, tried to instruct the back-end, only the front-end is allowed to instruct the back-end", raw.DeviceID)
		}
	}
}

func getBytes(key map[string]interface{}) []byte {
	var buf bytes.Buffer
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		panic(err.Error())
	}
	return buf.Bytes()
}
