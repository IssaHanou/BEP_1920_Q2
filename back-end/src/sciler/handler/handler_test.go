package handler

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
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

func (communicatorMock *CommunicatorMock) Start(handler mqtt.MessageHandler, onStart func()) {
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
				Finish:    false,
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
	communicator := communication.NewCommunicator(workingConfig.General.Host,
		workingConfig.General.Port, []string{"back-end", "test"})
	return &Handler{workingConfig, "fake file name", communicator}
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
	handler.Config.Timers["TestTimer"].Start(nil)
	assert.Equal(t, "stateActive", handler.Config.Timers["TestTimer"].State)
	handler.SetTimer("TestTimer", content)
	assert.Equal(t, "stateExpired", handler.Config.Timers["TestTimer"].State)
}

func TestHandler_SetTimer_Pause(t *testing.T) {
	handler := getTestHandler()
	content := config.ComponentInstruction{Instruction: "pause"}
	handler.Config.Timers["TestTimer"].Start(nil)
	assert.Equal(t, "stateActive", handler.Config.Timers["TestTimer"].State)
	handler.SetTimer("TestTimer", content)
	assert.Equal(t, "stateIdle", handler.Config.Timers["TestTimer"].State)
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
	handler.SendInstruction("display", inst)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
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

////////////////////////////// Status tests //////////////////////////////
func TestOnStatusMsg(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "status",
		Contents: map[string]interface{}{
			"testComponent0": true,
			"testComponent1": 40,
			"testComponent2": "blue",
			"testComponent3": []int{1, 2, 3},
			"testComponent4": true,
			"testComponent5": 40,
			"testComponent6": "blue",
			"testComponent7": []int{1, 2, 3},
			"testComponent8": "custom",
			"testComponent9": "custom",
		},
	}
	handler.onStatusMsg(msg)

	tests := []struct {
		name      string
		component string
		status    interface{}
	}{
		{
			name:      "component0 test",
			component: "testComponent0",
			status:    true,
		},
		{
			name:      "component1 test",
			component: "testComponent1",
			status:    40,
		},
		{
			name:      "component2 test",
			component: "testComponent2",
			status:    "blue",
		},
		{
			name:      "component3 test",
			component: "testComponent3",
			status:    []int{1, 2, 3},
		},
		{
			name:      "component4 test",
			component: "testComponent4",
			status:    true,
		},
		{
			name:      "component5 test",
			component: "testComponent5",
			status:    40,
		},
		{
			name:      "component6 test",
			component: "testComponent6",
			status:    "blue",
		},
		{
			name:      "component7 test",
			component: "testComponent7",
			status:    []int{1, 2, 3},
		},
		// TODO implement custom
		{
			name:      "component8 test",
			component: "testComponent8",
			status:    nil,
		},
		{
			name:      "component9 test",
			component: "testComponent9",
			status:    nil,
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.status, handler.Config.Devices["TestDevice"].Status[tt.component],
				"Device should set status right")
		})
	}
}

func TestOnStatusMsgOtherDevice(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "WrongDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "status",
		Contents: map[string]interface{}{
			"testComponent0": false,
			"testComponent1": true,
			"testComponent2": false},
	}
	handler.onStatusMsg(msg)

	_, ok := handler.Config.Devices["WrongDevice"]
	assert.Equal(t, false, ok,
		"Device should not exist in devices because it was not in config")
}

func TestOnStatusMsgWrongComponent(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "status",
		Contents: map[string]interface{}{
			"wrongComponent": true},
	}
	handler.onStatusMsg(msg)

	_, ok := handler.Config.Devices["TestDevice"].Status["wrongComponent"]
	assert.Equal(t, false, ok,
		"Component should not exist in device because it was not in config")
}

func TestOnStatusMsgWrongType(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "status",
		Contents: map[string]interface{}{
			"testComponent1": true,
			"testComponent2": 40,
			"testComponent3": "blue",
			"testComponent0": []int{1, 2, 3},
			"testComponent5": true,
			"testComponent6": 40,
			"testComponent7": "blue",
			"testComponent4": []int{1, 2, 3}},
	}
	handler.onStatusMsg(msg)

	tests := []struct {
		name      string
		component string
	}{
		{
			name:      "component0 test",
			component: "testComponent0",
		},
		{
			name:      "component1 test",
			component: "testComponent1",
		},
		{
			name:      "component2 test",
			component: "testComponent2",
		},
		{
			name:      "component3 test",
			component: "testComponent3",
		},
		{
			name:      "component4 test",
			component: "testComponent4",
		},
		{
			name:      "component5 test",
			component: "testComponent5",
		},
		{
			name:      "component6 test",
			component: "testComponent6",
		},
		{
			name:      "component7 test",
			component: "testComponent7",
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := handler.Config.Devices["TestDevice"].Status[tt.component]
			assert.Equal(t, false, ok,
				"component should not been updated in device because it was not the right type")
		})
	}
}

func TestMsgMapperStatus(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "status",
		Contents: map[string]interface{}{
			"testComponent0": true,
			"testComponent1": 40,
			"testComponent2": "blue"},
	}
	handler.msgMapper(msg)

	assert.Equal(t, true, handler.Config.Devices["TestDevice"].Status["testComponent0"],
		"Device should set status to true on component 0")
	assert.Equal(t, 40, handler.Config.Devices["TestDevice"].Status["testComponent1"],
		"Device should set status to 40 on component 1")
	assert.Equal(t, "blue", handler.Config.Devices["TestDevice"].Status["testComponent2"],
		"Device should set status to blue on component 2")
}

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
				"device_id": "front-end",
				"time_sent": "05-12-2019 09:42:10",
				"type":      "instruction",
				"contents": []map[string]interface{}{
					{"instruction": "test"},
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
				"contents":  []map[string]interface{}{{"instruction": "test"}},
				"type":      "instruction",
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
				"contents":  []map[string]interface{}{{"instruction": "test"}},
				"type":      "instruction",
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
				"contents":  []map[string]interface{}{{"instruction": "test"}},
				"type":      "instruction",
			},
		},
	}
	before := handler.Config
	handler.onConfirmationMsg(msg)
	assert.Equal(t, before, handler.Config,
		"Device should alter not config upon invalid confirmation message with no completed value")
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
				"contents":  []map[string]interface{}{{"instruction": "test"}},
				"type":      "instruction",
			},
		},
	}
	before := handler.Config
	handler.onConfirmationMsg(msg)
	assert.Equal(t, before, handler.Config,
		"Device should alter not config upon invalid confirmation message with no instructed object")
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
				"contents":  []map[string]interface{}{{"instruction": "test"}},
				"type":      "instruction",
			},
		},
	}
	before := handler.Config
	handler.onConfirmationMsg(msg)
	assert.Equal(t, before, handler.Config,
		"Device should alter not config upon invalid confirmation message with no boolean completed")
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
		"Device should alter not config upon invalid confirmation message with no instructions list")
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
		"Device should alter not config upon invalid confirmation message with device id not in config")
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
				"contents":  []map[string]interface{}{{"instruction": "test"}},
				"type":      "instruction",
			},
		},
	}
	before := handler.Config
	handler.msgMapper(msg)
	assert.Equal(t, before, handler.Config,
		"Device should alter not config upon confirmation message")
}

////////////////////////////// Instruction tests //////////////////////////////
func TestOnInstructionMsgName(t *testing.T) {
	msg := Message{
		DeviceID: "front-end",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{"instruction": "name"},
		},
	}
	communicatorMock := new(CommunicatorMock)
	handler := Handler{
		Config:       config.ReadFile("../../../resources/testing/test_instruction.json"),
		Communicator: communicatorMock,
	}
	returnMessage, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "name",
		Contents: map[string]interface{}{
			"name": "Escape X",
		},
	})
	communicatorMock.On("Publish", "front-end", string(returnMessage), 3)
	handler.onInstructionMsg(msg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
}

func TestOnInstructionMsgTestAll(t *testing.T) {
	msg := Message{
		DeviceID: "front-end",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{"instruction": "test all"},
		},
	}
	communicatorMock := new(CommunicatorMock)
	handler := Handler{
		Config:       config.ReadFile("../../../resources/testing/test_instruction.json"),
		Communicator: communicatorMock,
	}
	returnMessage, _ := json.Marshal(Message{
		DeviceID: "front-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{"instruction": "test"},
		},
	})
	communicatorMock.On("Publish", "client-computers", string(returnMessage), 3)
	handler.onInstructionMsg(msg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
}

func TestOnInstructionMsgSendStatus(t *testing.T) {
	msg := Message{
		DeviceID: "front-end",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{"instruction": "send status"},
		},
	}
	communicatorMock := new(CommunicatorMock)
	handler := Handler{
		Config:       config.ReadFile("../../../resources/testing/test_instruction_status.json"),
		Communicator: communicatorMock,
	}
	statusMessageDisplay, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "status",
		Contents: map[string]interface{}{
			"id": "display", "connection": false, "status": map[string]interface{}{},
		},
	})
	statusMessageFrontEnd, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "status",
		Contents: map[string]interface{}{
			"id": "front-end", "connection": false, "status": map[string]interface{}{"start": 0, "stop": 0},
		},
	})
	timerMessage, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "time",
		Contents: map[string]interface{}{
			"state":    "stateIdle",
			"duration": 10000,
			"id":       "timer1",
		},
	})
	timerGeneralMessage, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "time",
		Contents: map[string]interface{}{
			"state":    "stateIdle",
			"duration": 60000,
			"id":       "general",
		},
	})
	eventMessage, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "event status",
		Contents: []map[string]interface{}{
			{"id": "rule", "status": true, "description": "My rule"},
		},
	})
	communicatorMock.On("Publish", "front-end", string(eventMessage), 3)
	communicatorMock.On("Publish", "front-end", string(timerMessage), 3)
	communicatorMock.On("Publish", "front-end", string(timerGeneralMessage), 3)
	communicatorMock.On("Publish", "front-end", string(statusMessageDisplay), 3)
	communicatorMock.On("Publish", "front-end", string(statusMessageFrontEnd), 3)
	handler.onInstructionMsg(msg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 5)
}

func TestOnInstructionMsgResetAll(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	workingConfig := config.ReadFile("../../../resources/testing/test_config.json")
	handler := Handler{
		Config:       workingConfig,
		ConfigFile:   "../../../resources/testing/test_config.json",
		Communicator: communicatorMock,
	}
	instructionMsg := Message{
		DeviceID: "front-end",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{"instruction": "reset all"},
		},
	}
	responseMsg := Message{
		DeviceID: "front-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{"instruction": "reset"},
		},
	}
	statusMsg := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "time",
		Contents: map[string]interface{}{
			"id":       "general",
			"duration": 1800000,
			"state":    "stateIdle",
		},
	}

	jsonInstructionMsg, _ := json.Marshal(&responseMsg)
	jsonStatusMsg, _ := json.Marshal(&statusMsg)

	communicatorMock.On("Publish", "client-computers", string(jsonInstructionMsg), 3)
	communicatorMock.On("Publish", "front-end", string(jsonInstructionMsg), 3)
	communicatorMock.On("Publish", "front-end", string(jsonStatusMsg), 3)
	handler.msgMapper(instructionMsg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 3)
}

func TestOnInstructionMsgHint(t *testing.T) {
	msg := Message{
		DeviceID: "front-end",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{"instruction": "hint", "hint": "This is my hint"},
		},
	}
	communicatorMock := new(CommunicatorMock)
	handler := Handler{
		Config:       config.ReadFile("../../../resources/testing/test_instruction.json"),
		Communicator: communicatorMock,
	}
	returnMessage, _ := json.Marshal(Message{
		DeviceID: "front-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{"instruction": "hint", "hint": "This is my hint"},
		},
	})
	communicatorMock.On("Publish", "hint", string(returnMessage), 3)
	handler.onInstructionMsg(msg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
}

func TestOnInstructionMsgFinishRule(t *testing.T) {
	msg := Message{
		DeviceID: "front-end",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{"instruction": "finish rule", "rule": "rule"},
		},
	}
	communicatorMock := new(CommunicatorMock)
	handler := Handler{
		Config:       config.ReadFile("../../../resources/testing/test_instruction.json"),
		Communicator: communicatorMock,
	}
	returnMessage, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "event status",
		Contents: []map[string]interface{}{
			{"id": "rule", "description": "My rule", "status": true},
		},
	})
	communicatorMock.On("Publish", "front-end", string(returnMessage), 3)
	handler.onInstructionMsg(msg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
}

func TestOnInstructionMsgTestDevice(t *testing.T) {
	msg := Message{
		DeviceID: "front-end",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{"instruction": "test device", "device": "display"},
		},
	}
	communicatorMock := new(CommunicatorMock)
	handler := Handler{
		Config:       config.ReadFile("../../../resources/testing/test_instruction.json"),
		Communicator: communicatorMock,
	}
	returnMessage, _ := json.Marshal(Message{
		DeviceID: "front-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{"instruction": "test"},
		},
	})
	communicatorMock.On("Publish", "display", string(returnMessage), 3)
	handler.onInstructionMsg(msg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
}

func TestOnInstructionMsgMapper(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{"instruction": "name"},
		},
	}

	before := handler.Config
	handler.msgMapper(msg)
	assert.Equal(t, before, handler.Config,
		"Nothing should have been changed after an instruction message type")
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
		Contents: []map[string]interface{}{
			{"description": "Als de mainSwitch true is, moet greenLight1 aangaan",
				"id":     "mainSwitch flipped",
				"status": true},
		},
	})

	communicatorMock.On("Publish", "front-end", string(messageEventStatus), 3)
	communicatorMock.On("Publish", "front-end", string(messageStatus), 3)
	communicatorMock.On("Publish", "controlBoard", string(messageInstruction), 3)
	handler.msgMapper(msg)
	communicatorMock.AssertExpectations(t) // if this test becomes flaky (only when this test takes longer then 1 second),
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
		Type:     "event-status",
		Contents: []map[string]interface{}{
			{"description": "Als de mainSwitch true is, gebeurt er niks",
				"id":     "mainSwitch flipped",
				"status": true},
			{"description": "Als rule 'mainSwitch flipped' is gedaan, dan moet greenLight1 aangaan",
				"id":     "weldoen",
				"status": true},
		},
	})

	communicatorMock.On("Publish", "front-end", string(messageEventStatus), 3)
	communicatorMock.On("Publish", "front-end", string(messageStatus), 3)
	communicatorMock.On("Publish", "controlBoard", string(messageInstruction), 3)
	handler.msgMapper(msg)
	communicatorMock.AssertExpectations(t) // if this test becomes flaky (only when this test takes longer then 1 second),
	// (message expected includes time...), replace the messages with 'mock.Anything'
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
		Type:     "event-status",
		Contents: []map[string]interface{}{
			{"description": "Als de mainSwitch true is, moet greenLight1 aangaan",
				"id":     "mainSwitch flipped",
				"status": true},
		},
	})

	communicatorMock.On("Publish", "front-end", string(messageEventStatus), 3)
	communicatorMock.On("Publish", "front-end", string(messageStatus), 3)
	communicatorMock.On("Publish", "controlBoard", mock.Anything, 3)
	handler.msgMapper(msg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 2) // Only publish to front-end for status should be done, no action should be performed
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
	handler.SendStatus("Unknown device or timer")
	communicatorMock.AssertNumberOfCalls(t, "Publish", 0)
}
