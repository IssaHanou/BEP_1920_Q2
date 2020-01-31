package handler

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"sciler/config"
	"testing"
)

////////////////////////////// Confirmation tests //////////////////////////////
func TestOnConfirmationMsgFrontEnd(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	workingConfig := config.ReadFile("../../../resources/testing/test_config.json")
	handler := Handler{
		Config:       workingConfig,
		ConfigFile:   "../../../resources/testing/test_config.json",
		Communicator: communicatorMock,
	}
	msg := Message{
		DeviceID: "controlBoard",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "confirmation",
		Contents: map[string]interface{}{
			"completed": true,
			"instructed": map[string]interface{}{
				"device_id": "back-end",
				"time_sent": "05-12-2019 09:42:10",
				"type":      "instruction",
				"contents": []map[string]interface{}{{
					"instruction":   "test",
					"instructed_by": "front-end"},
				},
			},
		},
	}
	jsonMsg, _ := json.Marshal(&msg)
	communicatorMock.On("Publish", "front-end", string(jsonMsg), 3)
	handler.onConfirmationMsg(msg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
}

func TestOnConfirmationMsgTrue(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "confirmation",
		Contents: map[string]interface{}{
			"completed": true,
			"instructed": map[string]interface{}{
				"device_id": "front-end",
				"time_sent": "05-12-2019 09:42:10",
				"contents": []map[string]interface{}{{
					"instruction": "test"}},
				"type": "instruction",
			},
		},
	}
	handler.onConfirmationMsg(msg)
	assert.Equal(t, true, handler.Config.Devices["TestDevice"].Connection,
		"Device's connection should have been set to true")
}

func TestOnConfirmationMsgFalse(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "confirmation",
		Contents: map[string]interface{}{
			"completed": false,
			"instructed": map[string]interface{}{
				"device_id": "front-end",
				"time_sent": "05-12-2019 09:42:10",
				"contents": []map[string]interface{}{{
					"instruction":   "test",
					"instructed_by": "front-end"}},
				"type": "instruction",
			},
		},
	}
	handler.onConfirmationMsg(msg)
	assert.Equal(t, true, handler.Config.Devices["TestDevice"].Connection,
		"Device's connection should have been set to true")
}

func TestOnConfirmationMsgIncorrect1(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "confirmation",
		Contents: map[string]interface{}{
			"completion": true,
			"instructed": map[string]interface{}{
				"device_id": "back-end",
				"time_sent": "05-12-2019 09:42:10",
				"contents": []map[string]interface{}{{
					"instruction":   "test",
					"instructed_by": "front-end"}},
				"type": "instruction",
			},
		},
	}
	before := handler.Config
	handler.onConfirmationMsg(msg)
	assert.Equal(t, before, handler.Config,
		"Device should not alter config upon invalid confirmation message with no completed value")
}

func TestOnConfirmationMsgIncorrect2(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "confirmation",
		Contents: map[string]interface{}{
			"completed": true,
			"instruction": map[string]interface{}{
				"device_id": "back-end",
				"time_sent": "05-12-2019 09:42:10",
				"type":      "instruction",
				"contents": map[string]interface{}{
					"instructions":  "test",
					"instructed_by": "front-end"},
			},
		},
	}
	before := handler.Config
	handler.onConfirmationMsg(msg)
	assert.Equal(t, before, handler.Config,
		"Device should not alter config upon invalid confirmation message with no instructed object")
}

func TestOnConfirmationMsgIncorrect3(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "confirmation",
		Contents: map[string]interface{}{
			"completed": "true",
			"instructed": map[string]interface{}{
				"device_id": "back-end",
				"time_sent": "05-12-2019 09:42:10",
				"contents": []map[string]interface{}{{
					"instruction":   "test",
					"instructed_by": "front-end"}},
				"type": "instruction",
			},
		},
	}
	before := handler.Config
	handler.onConfirmationMsg(msg)
	assert.Equal(t, before, handler.Config,
		"Device should not alter config upon invalid confirmation message with no boolean completed")
}

func TestOnConfirmationMsgIncorrect4(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "confirmation",
		Contents: map[string]interface{}{
			"completed": true,
			"instructed": map[string]interface{}{
				"device_id": "back-end",
				"time_sent": "05-12-2019 09:42:10",
				"type":      "instruction",
				"contents": map[string]interface{}{
					"instruction": "test",
				},
			},
		},
	}
	before := handler.Config
	handler.onConfirmationMsg(msg)
	assert.Equal(t, before, handler.Config,
		"Device should not alter config upon invalid confirmation message with no instructions list")
}

func TestOnConfirmationMsgIncorrect5(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "display",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "confirmation",
		Contents: map[string]interface{}{
			"completed": true,
			"instructed": map[string]interface{}{
				"device_id": "back-end",
				"time_sent": "05-12-2019 09:42:10",
				"type":      "instruction",
				"contents": []map[string]interface{}{
					{"instruction": "test"},
				},
			},
		},
	}
	before := handler.Config
	handler.onConfirmationMsg(msg)
	assert.Equal(t, before, handler.Config,
		"Device should not alter config upon invalid confirmation message with device id not in config")
}

func TestMsgMapperConfirmation(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "confirmation",
		Contents: map[string]interface{}{
			"completed": true,
			"instructed": map[string]interface{}{
				"device_id": "back-end",
				"time_sent": "05-12-2019 09:42:10",
				"contents": []map[string]interface{}{{
					"instruction":   "test",
					"instructed_by": "front-end"}},
				"type": "instruction",
			},
		},
	}
	before := handler.Config
	handler.msgMapper(msg)
	assert.Equal(t, before, handler.Config,
		"Device should not alter config upon confirmation message")
}
