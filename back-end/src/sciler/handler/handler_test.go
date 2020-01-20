package handler

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sciler/communication"
	"sciler/config"
	"testing"
	"time"
)

type CommunicatorMock struct {
	mock.Mock
}

func (communicatorMock *CommunicatorMock) Start() {
	// do nothing
}

func (communicatorMock *CommunicatorMock) Publish(topic string, message string, retrials int) {
	communicatorMock.Called(topic, message, retrials)
}

func getTestHandler() *Handler {
	workingConfig := config.WorkingConfig{
		General: config.General{
			Name:     "Test",
			Duration: "1",
			Host:     "localhost",
			Port:     1883,
		},
		Puzzles: nil,
		Timers: map[string]*config.Timer{
			"TestTimer": &(config.Timer{
				ID:        "TestTimer",
				Duration:  5 * time.Second,
				StartedAt: time.Now(),
				T:         nil,
				State:     "stateIdle",
				Ending:    nil,
				Finished:  false,
			}),
		},
		GeneralEvents: nil,
		Devices: map[string]*config.Device{
			"TestDevice": &(config.Device{
				ID:          "TestDevice",
				Description: "test uitleg",
				Input: map[string]string{
					"testComponent0": "boolean",
					"testComponent1": "numeric",
					"testComponent2": "string",
					"testComponent3": "array",
					"testComponent8": "custom",
				},
				Output: map[string]config.OutputObject{
					"testComponent4": {Type: "boolean", Instructions: map[string]string{}},
					"testComponent5": {Type: "numeric", Instructions: map[string]string{}},
					"testComponent6": {Type: "string", Instructions: map[string]string{}},
					"testComponent7": {Type: "array", Instructions: map[string]string{}},
					"testComponent9": {Type: "custom", Instructions: map[string]string{}},
				},
				Status:     map[string]interface{}{},
				Connection: false,
			}),
		},
	}
	messageHandler := Handler{Config: workingConfig, ConfigFile: "fake file name"}
	communicator := communication.NewCommunicator(workingConfig.General.Host,
		workingConfig.General.Port, []string{"back-end", "test"}, messageHandler.NewHandler, func() {})
	messageHandler.Communicator = communicator
	return &messageHandler
}

////////////////////////////// Helper method tests //////////////////////////////
func TestHandler_SetTimer_Start(t *testing.T) {
	handler := getTestHandler()
	content := config.ComponentInstruction{Instruction: "start"}
	assert.Equal(t, "stateIdle", handler.Config.Timers["TestTimer"].State)
	handler.SetTimer("TestTimer", content)
	assert.Equal(t, "stateActive", handler.Config.Timers["TestTimer"].State)
}

func TestHandler_SetTimer_Stop(t *testing.T) {
	handler := getTestHandler()
	content := config.ComponentInstruction{Instruction: "stop"}
	handler.Config.Timers["TestTimer"].Start(handler)
	assert.Equal(t, "stateActive", handler.Config.Timers["TestTimer"].State)
	handler.SetTimer("TestTimer", content)
	assert.Equal(t, "stateExpired", handler.Config.Timers["TestTimer"].State)
}

func TestHandler_SetTimer_Pause(t *testing.T) {
	handler := getTestHandler()
	content := config.ComponentInstruction{Instruction: "pause"}
	handler.Config.Timers["TestTimer"].Start(handler)
	assert.Equal(t, "stateActive", handler.Config.Timers["TestTimer"].State)
	handler.SetTimer("TestTimer", content)
	assert.Equal(t, "stateIdle", handler.Config.Timers["TestTimer"].State)
}

func TestHandler_SetTimer_Add(t *testing.T) {
	handler := getTestHandler()
	content := config.ComponentInstruction{Instruction: "add", Value: "5s"}
	assert.Equal(t, 5, int(handler.Config.Timers["TestTimer"].Duration.Seconds()))
	handler.SetTimer("TestTimer", content)
	assert.Equal(t, 10, int(handler.Config.Timers["TestTimer"].Duration.Seconds()))

}

func TestHandler_SetTimer_Add_Parse_error(t *testing.T) {
	handler := getTestHandler()
	content := config.ComponentInstruction{Instruction: "add", Value: "5"}
	assert.Equal(t, 5, int(handler.Config.Timers["TestTimer"].Duration.Seconds()))
	handler.SetTimer("TestTimer", content)
	assert.Equal(t, 5, int(handler.Config.Timers["TestTimer"].Duration.Seconds()))

}

func TestHandler_SetTimer_Subtract(t *testing.T) {
	handler := getTestHandler()
	content := config.ComponentInstruction{Instruction: "subtract", Value: "3s"}
	assert.Equal(t, 5, int(handler.Config.Timers["TestTimer"].Duration.Seconds()))
	handler.SetTimer("TestTimer", content)
	assert.Equal(t, 2, int(handler.Config.Timers["TestTimer"].Duration.Seconds()))

}

func TestHandler_SetTimer_Subtract_Parse_Error(t *testing.T) {
	handler := getTestHandler()
	content := config.ComponentInstruction{Instruction: "subtract", Value: "3"}
	assert.Equal(t, 5, int(handler.Config.Timers["TestTimer"].Duration.Seconds()))
	handler.SetTimer("TestTimer", content)
	assert.Equal(t, 5, int(handler.Config.Timers["TestTimer"].Duration.Seconds()))

}

func TestHandler_SetTimer_Done(t *testing.T) {
	handler := getTestHandler()
	content := config.ComponentInstruction{Instruction: "done"}
	handler.Config.Timers["TestTimer"].Start(handler)
	assert.Equal(t, "stateActive", handler.Config.Timers["TestTimer"].State)
	handler.SetTimer("TestTimer", content)
	assert.Equal(t, "stateExpired", handler.Config.Timers["TestTimer"].State)
}

func TestGetStatus(t *testing.T) {
	msg, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{"instruction": "status update"},
		},
	})
	communicatorMock := new(CommunicatorMock)
	handler := Handler{
		Config:       config.ReadFile("../../../resources/testing/test_instruction.json"),
		Communicator: communicatorMock,
	}
	communicatorMock.On("Publish", "display", string(msg), 3)
	handler.GetStatus("display")
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
}

func TestSendInstruction(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	handler := Handler{
		Config:       config.ReadFile("../../../resources/testing/test_instruction.json"),
		Communicator: communicatorMock,
	}
	inst := []config.ComponentInstruction{
		{"display", "hint", "my hint"},
	}
	msg, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: inst,
	})
	communicatorMock.On("Publish", "display", string(msg), 3)
	handler.SendComponentInstruction("display", inst, "")
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
}

func TestSendLabelInstruction(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	handler := Handler{
		Config:       config.ReadFile("../../../resources/testing/test_instruction_label.json"),
		Communicator: communicatorMock,
	}
	inst := []config.ComponentInstruction{
		{"", "hint", "my hint"},
	}
	instMsg := []config.ComponentInstruction{
		{"display1", "hint", "my hint"},
	}
	msg, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: instMsg,
	})
	communicatorMock.On("Publish", "display2", string(msg), 3)
	handler.SendLabelInstruction("display-label1", inst, "")
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
}

func TestSendLabelInstruction_2(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	handler := Handler{
		Config:       config.ReadFile("../../../resources/testing/test_instruction_label.json"),
		Communicator: communicatorMock,
	}
	inst := []config.ComponentInstruction{
		{"", "hint", "my hint"},
	}
	instMsg := []config.ComponentInstruction{
		{"display2", "hint", "my hint"},
	}
	msg, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: instMsg,
	})
	communicatorMock.On("Publish", "display1", string(msg), 3)
	communicatorMock.On("Publish", "display2", string(msg), 3)
	handler.SendLabelInstruction("display-label2", inst, "")
	communicatorMock.AssertNumberOfCalls(t, "Publish", 2)
}

func TestSendInstructionDelay(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	handler := Handler{
		Config:       config.ReadFile("../../../resources/testing/test_instruction.json"),
		Communicator: communicatorMock,
	}
	inst := []config.ComponentInstruction{
		{"display", "hint", "my hint"},
	}
	inst2 := []config.ComponentInstruction{
		{"display", "hint", "my hint 2"},
	}
	msg, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: inst,
	})
	msg2, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: inst2,
	})
	communicatorMock.On("Publish", "display", string(msg), 3)
	communicatorMock.On("Publish", "display", string(msg2), 3)
	go handler.SendComponentInstruction("display", inst, "")
	go handler.SendComponentInstruction("display", inst2, "1s")
	time.Sleep(100 * time.Millisecond)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
	time.Sleep(1 * time.Second)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 2)
}

////////////////////////////// Connection tests //////////////////////////////
func TestOnConnectionMsgFalse(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "connection",
		Contents: map[string]interface{}{
			"connection": false,
		},
	}
	handler.onConnectionMsg(msg)
	assert.Equal(t, false, handler.Config.Devices["TestDevice"].Connection,
		"Device should set connection to false on connection message")
}

func TestOnConnectionMsgOtherDevice(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "WrongDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "connection",
		Contents: map[string]interface{}{
			"connection": true},
	}
	handler.onConnectionMsg(msg)

	_, ok := handler.Config.Devices["WrongDevice"]
	assert.Equal(t, false, ok,
		"Device should not exist in devices because it was not in config")
}

func TestOnConnectionMsgInvalid(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "connection",
		Contents: map[string]interface{}{
			"connection": "true"},
	}
	handler.onConnectionMsg(msg)
	assert.Equal(t, false, handler.Config.Devices["TestDevice"].Connection,
		"Device should not set connection to true on incorrect connection message")
}

func TestOnConnectionMsgInvalid2(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "connection",
		Contents: map[string]interface{}{
			"connected": "true"},
	}
	handler.onConnectionMsg(msg)
	assert.Equal(t, false, handler.Config.Devices["TestDevice"].Connection,
		"Device should not set connection to true on incorrect connection message")
}

func TestMsgMapperConnection(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "connection",
		Contents: map[string]interface{}{
			"connection": false},
	}
	handler.msgMapper(msg)
	assert.Equal(t, false, handler.Config.Devices["TestDevice"].Connection,
		"Device should set connection to true on connection message")
}

func TestOnConnectionMsg(t *testing.T) {
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
		Type:     "connection",
		Contents: map[string]interface{}{
			"connection": true},
	}
	statusMsg := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "status",
		Contents: map[string]interface{}{
			"id":         "controlBoard",
			"connection": true,
			"status":     map[string]interface{}{},
		},
	}
	jsonStatusMsg, _ := json.Marshal(&statusMsg)
	communicatorMock.On("Publish", "front-end", string(jsonStatusMsg), 3)
	assert.False(t, handler.Config.Devices["controlBoard"].Connection,
		"Device connection should be false on default")

	handler.onConnectionMsg(msg)
	assert.True(t, handler.Config.Devices["controlBoard"].Connection,
		"Device should set connection to true on connection message")
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
}

////////////////////////////// Event handling tests //////////////////////////////
func TestHandleSingleEvent(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	handler := Handler{
		Config:       config.ReadFile("../../../resources/testing/test_singleEvent.json"),
		Communicator: communicatorMock,
	}
	msg := Message{
		DeviceID: "controlBoard",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "status",
		Contents: map[string]interface{}{
			"redSwitch":    false,
			"orangeSwitch": false,
			"greenSwitch":  false,
			"slider1":      0,
			"slider2":      0,
			"slider3":      0,
			"mainSwitch":   true,
			"greenLight1":  "off",
			"greenLight2":  "off",
			"greenLight3":  "off",
			"redLight1":    "off",
			"redLight2":    "off",
			"redLight3":    "off",
		},
	}

	messageInstruction, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{
				"component_id": "greenLight1",
				"instruction":  "turnOnOff",
				"value":        true,
			},
		},
	})

	messageStatus, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "status",
		Contents: map[string]interface{}{
			"id": "controlBoard",
			"status": map[string]interface{}{
				"greenLight1":  "off",
				"greenLight2":  "off",
				"greenLight3":  "off",
				"greenSwitch":  false,
				"mainSwitch":   true,
				"orangeSwitch": false,
				"redLight1":    "off",
				"redLight2":    "off",
				"redLight3":    "off",
				"redSwitch":    false,
				"slider1":      0,
				"slider2":      0,
				"slider3":      0,
			},
			"connection": false,
		},
	})

	messageEventStatus, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "event status",
		Contents: []map[string]interface{}{{
			"id":     "mainSwitch flipped",
			"status": true},
		},
	})

	communicatorMock.On("Publish", "front-end", string(messageEventStatus), 3)
	communicatorMock.On("Publish", "front-end", string(messageStatus), 3)
	communicatorMock.On("Publish", "controlBoard", string(messageInstruction), 3)
	handler.msgMapper(msg)
	time.Sleep(10 * time.Millisecond) // Give the goroutine(s) time to finish before asserting number of calls
	communicatorMock.AssertNumberOfCalls(t, "Publish", 3)
	// if this test becomes flaky (only when this test takes longer then 1 second),
	// (message expected includes time...), replace the messages with 'mock.Anything'
}

func TestHandleDoubleEvent(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	handler := Handler{
		Config:       config.ReadFile("../../../resources/testing/test_doubleEvent.json"),
		Communicator: communicatorMock,
	}
	msg := Message{
		DeviceID: "controlBoard",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "status",
		Contents: map[string]interface{}{
			"redSwitch":    false,
			"orangeSwitch": false,
			"greenSwitch":  false,
			"slider1":      0,
			"slider2":      0,
			"slider3":      0,
			"mainSwitch":   true,
			"greenLight1":  "off",
			"greenLight2":  "off",
			"greenLight3":  "off",
			"redLight1":    "off",
			"redLight2":    "off",
			"redLight3":    "off",
		},
	}

	messageInstruction, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{
				"component_id": "greenLight1",
				"instruction":  "turnOnOff",
				"value":        true,
			},
		},
	})

	communicatorMock.On("Publish", "front-end", mock.Anything, 3)
	communicatorMock.On("Publish", "controlBoard", string(messageInstruction), 3)
	handler.msgMapper(msg)
	time.Sleep(10 * time.Millisecond) // Give the goroutine(s) time to finish before asserting number of calls
	communicatorMock.AssertNumberOfCalls(t, "Publish", 3)
	// if this test becomes flaky (only when this test takes longer then 1 second),
	// (message expected includes time...), replace the messages with 'mock.Anything'

	// TODO: to restore test so that the message are checked again:
	//  - duplicate the publish front-end line,
	//  - replace 'mock.Anything' with the correct messages
}

////////////////////////////// Error/irregular behavior tests //////////////////////////////
func TestMsgMapperIllegalType(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "test",
		Contents: map[string]interface{}{
			"testComponent0": false,
			"testComponent1": true,
			"testComponent2": false},
	}

	before := handler.Config
	handler.msgMapper(msg)
	assert.Equal(t, before, handler.Config,
		"Nothing should have been changed after an incorrect message type")
}

func TestLimitRule(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	workingConfig := config.ReadFile("../../../resources/testing/test_singleEvent.json")
	workingConfig.RuleMap["mainSwitch flipped"].Executed = 1
	handler := Handler{
		Config:       workingConfig,
		Communicator: communicatorMock,
	}
	msg := Message{
		DeviceID: "controlBoard",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "status",
		Contents: map[string]interface{}{
			"redSwitch":    false,
			"orangeSwitch": false,
			"greenSwitch":  false,
			"slider1":      0,
			"slider2":      0,
			"slider3":      0,
			"mainSwitch":   true,
			"greenLight1":  "off",
			"greenLight2":  "off",
			"greenLight3":  "off",
			"redLight1":    "off",
			"redLight2":    "off",
			"redLight3":    "off",
		},
	}
	messageStatus, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "status",
		Contents: map[string]interface{}{
			"id": "controlBoard",
			"status": map[string]interface{}{
				"greenLight1":  "off",
				"greenLight2":  "off",
				"greenLight3":  "off",
				"greenSwitch":  false,
				"mainSwitch":   true,
				"orangeSwitch": false,
				"redLight1":    "off",
				"redLight2":    "off",
				"redLight3":    "off",
				"redSwitch":    false,
				"slider1":      0,
				"slider2":      0,
				"slider3":      0,
			},
			"connection": false,
		},
	})
	messageEventStatus, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "event status",
		Contents: []map[string]interface{}{{
			"id":     "mainSwitch flipped",
			"status": true},
		},
	})

	communicatorMock.On("Publish", "front-end", string(messageStatus), 3)
	communicatorMock.On("Publish", "front-end", string(messageEventStatus), 3)
	communicatorMock.On("Publish", "controlBoard", mock.Anything, 3)
	handler.msgMapper(msg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 2)
	// Only publish to front-end for status should be done, no action should be performed
}

func TestGetMapSliceInvalidConfirmation(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "back-end",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "confirmation",
		Contents: []map[string]interface{}{
			{"completed": true,
				"instructed": map[string]interface{}{
					"device_id": "front-end",
					"contents":  "test",
				}},
		},
	}
	assert.Panics(t, func() { handler.onConfirmationMsg(msg) }, "Should throw error with invalid contents")
}

func TestGetMapSliceInvalidInstruction(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "back-end",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "instruction",
		Contents: map[string]interface{}{
			"instruction": "test",
		},
	}
	assert.NotPanics(t, func() { handler.onInstructionMsg(msg) },
		"Should return empty with invalid contents")
}

func TestInstructionFromWrongDevice(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	workingConfig := config.ReadFile("../../../resources/testing/test_config.json")
	handler := Handler{
		Config:       workingConfig,
		ConfigFile:   "../../../resources/testing/test_config.json",
		Communicator: communicatorMock,
	}
	instructionMsg := Message{
		DeviceID: "not front-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{
				"instruction": "hint",
				"value":       "some useful hint",
			},
		},
	}
	jsonHintMessage, _ := json.Marshal(&instructionMsg)
	communicatorMock.On("Publish", "hint", string(jsonHintMessage), 3)
	handler.msgMapper(instructionMsg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 0)
}

func TestSendStatusUnknownDevice(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	workingConfig := config.ReadFile("../../../resources/testing/test_config.json")
	handler := Handler{
		Config:       workingConfig,
		ConfigFile:   "../../../resources/testing/test_config.json",
		Communicator: communicatorMock,
	}
	msg, _ := json.Marshal(Message{
		DeviceID: "TestDevice2",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "status",
		Contents: map[string]interface{}{
			"testComponent0": true,
			"testComponent1": 40,
			"testComponent2": "blue"},
	})

	communicatorMock.On("Publish", "hint", string(msg), 3)
	handler.sendStatus("Unknown device or timer")
	communicatorMock.AssertNumberOfCalls(t, "Publish", 0)
}
