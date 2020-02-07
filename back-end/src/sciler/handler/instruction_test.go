package handler

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"sciler/config"
	"testing"
	"time"
)

////////////////////////////// Instruction tests //////////////////////////////
func TestInstructionSetup(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	workingConfig := config.ReadFile("../../../resources/testing/test_setup.json")
	handler := Handler{
		Config:       workingConfig,
		ConfigFile:   "../../../resources/testing/test_setup.json",
		Communicator: communicatorMock,
	}
	instructionMsg := Message{
		DeviceID: "front-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{
				"instruction": "send setup",
			},
		},
	}
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
	statusInstructionMsg, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{
				"instruction": "status update",
			},
		},
	})
	statusMessage, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "status",
		Contents: map[string]interface{}{
			"id":         "telephone",
			"connection": false,
			"status":     map[string]interface{}{},
		},
	})
	statusMessageFrontEnd, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "status",
		Contents: map[string]interface{}{
			"id":         "front-end",
			"connection": false,
			"status": map[string]interface{}{
				"start":     false,
				"gameState": "gereed",
				"hintLog":   []string{},
			},
		},
	})
	messageEventStatus, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "event status",
		Contents: []map[string]interface{}{
			{
				"id":     "correctSequence",
				"status": false,
			},
		},
	})

	communicatorMock.On("Publish", "front-end", mock.AnythingOfType("string"), 3).Once()
	communicatorMock.On("Publish", "front-end", string(timerGeneralMessage), 3).Once()
	communicatorMock.On("Publish", "front-end", string(statusMessage), 3).Once()
	communicatorMock.On("Publish", "telephone", string(statusInstructionMsg), 3).Once()
	communicatorMock.On("Publish", "front-end", string(statusMessageFrontEnd), 3).Once()
	communicatorMock.On("Publish", "front-end", string(statusMessageFrontEnd), 3).Once()
	communicatorMock.On("Publish", "front-end", string(statusInstructionMsg), 3).Once()
	communicatorMock.On("Publish", "front-end", string(messageEventStatus), 3).Once()
	communicatorMock.On("Publish", "time", string(timerGeneralMessage), 3).Once()
	handler.msgMapper(instructionMsg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 8)
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
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{{
			"instruction":   "reset",
			"instructed_by": "front-end"},
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

	communicatorMock.On("Publish", "client-computers", string(jsonInstructionMsg), 3).Once()
	communicatorMock.On("Publish", "front-end", string(jsonInstructionMsg), 3).Once()
	communicatorMock.On("Publish", "front-end", string(jsonStatusMsg), 3).Once()
	communicatorMock.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("string"), 3) // all calls from sendStatus update (tested in another test)
	handler.msgMapper(instructionMsg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 14)
}

func TestOnInstructionMsgTestAll(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	workingConfig := config.ReadFile("../../../resources/testing/test_config.json")
	handler := Handler{
		Config:       workingConfig,
		Communicator: communicatorMock,
	}
	instructionMsg := Message{
		DeviceID: "front-end",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{"instruction": "test all"},
		},
	}
	responseMsg, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{{
			"instruction":   "test",
			"instructed_by": "front-end"},
		},
	})
	communicatorMock.On("Publish", "client-computers", string(responseMsg), 3)
	handler.onInstructionMsg(instructionMsg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
}

func TestOnInstructionMsgTestDevice(t *testing.T) {
	msg := Message{
		DeviceID: "front-end",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "instruction",
		Contents: []map[string]interface{}{{
			"instruction": "test device",
			"device":      "display"},
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
		Type:     "instruction",
		Contents: []map[string]interface{}{{
			"instruction":   "test",
			"instructed_by": "front-end"},
		},
	})
	communicatorMock.On("Publish", "display", string(returnMessage), 3)
	handler.onInstructionMsg(msg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
}

func TestOnInstructionMsgFinishRule(t *testing.T) {
	msg := Message{
		DeviceID: "front-end",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "instruction",
		Contents: []map[string]interface{}{{
			"instruction": "finish rule",
			"rule":        "rule"},
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
		Contents: []map[string]interface{}{{
			"id":     "rule",
			"status": true},
		},
	})
	instHintMessage, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{{
			"component_id": "display",
			"instruction":  "hint",
			"value":        "it tested!"},
		},
	})
	instTimerMessage, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "time",
		Contents: map[string]interface{}{
			"duration": 1800000,
			"id":       "general",
			"state":    "stateIdle",
		},
	})
	communicatorMock.On("Publish", "front-end", string(instTimerMessage), 3)
	communicatorMock.On("Publish", "time", string(instTimerMessage), 3)
	communicatorMock.On("Publish", "front-end", string(returnMessage), 3)
	communicatorMock.On("Publish", "display", string(instHintMessage), 3)
	handler.onInstructionMsg(msg)
	time.Sleep(10 * time.Millisecond) // Give the goroutine(s) time to finish before asserting number of calls
	communicatorMock.AssertNumberOfCalls(t, "Publish", 4)
}

func TestOnInstructionMsgFinishRuleLabel(t *testing.T) {
	msg := Message{
		DeviceID: "front-end",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "instruction",
		Contents: []map[string]interface{}{{
			"instruction": "finish rule",
			"rule":        "rule"},
		},
	}
	communicatorMock := new(CommunicatorMock)
	handler := Handler{
		Config:       config.ReadFile("../../../resources/testing/test_instruction_label.json"),
		Communicator: communicatorMock,
	}
	instMessage, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{{
			"value":        "test",
			"component_id": "display1",
			"instruction":  "hint"},
		},
	},
	)
	communicatorMock.On("Publish", "display2", string(instMessage), 3).Once()
	communicatorMock.On("Publish", "front-end", mock.Anything, 3).Once()
	handler.onInstructionMsg(msg)
	time.Sleep(10 * time.Millisecond) // Give the goroutine(s) time to finish before asserting number of calls
	communicatorMock.AssertNumberOfCalls(t, "Publish", 2)
}

func TestOnInstructionMsgHint(t *testing.T) {
	msg := Message{
		DeviceID: "front-end",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{
				"instruction": "hint",
				"value":       "This is my hint",
				"topic":       "hint",
			},
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
		Type:     "instruction",
		Contents: []map[string]interface{}{{
			"instruction":   "hint",
			"value":         "This is my hint",
			"instructed_by": "front-end"},
		},
	})
	communicatorMock.On("Publish", "hint", string(returnMessage), 3)
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

func TestInstructionCheckConfigNoErrors(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	fileName := "../../../resources/testing/test_config.json"
	workingConfig := config.ReadFile(fileName)
	handler := Handler{
		Config:       workingConfig,
		ConfigFile:   fileName,
		Communicator: communicatorMock,
	}
	jsonFile, _ := ioutil.ReadFile(fileName)
	configToTest := make(map[string]interface{})
	if err := json.Unmarshal(jsonFile, &configToTest); err != nil {
		assert.FailNow(t, "cannot create instruction message")
	}
	instructionMsg := Message{
		DeviceID: "front-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{
				"instruction": "check config",
				"config":      configToTest,
				"name":        "testFile",
			},
		},
	}
	returnMsg := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "config",
		Contents: map[string]interface{}{
			"errors": []string{},
			"name":   "testFile",
		},
	}
	jsonMessage, _ := json.Marshal(&returnMsg)
	communicatorMock.On("Publish", "front-end", string(jsonMessage), 3)
	handler.msgMapper(instructionMsg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
}

func TestInstructionCheckConfigWithErrors(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	workingConfig := config.ReadFile("../../../resources/testing/test_config.json")
	handler := Handler{
		Config:       workingConfig,
		ConfigFile:   "../../../resources/testing/test_config.json",
		Communicator: communicatorMock,
	}
	jsonFile, _ := ioutil.ReadFile("../../../resources/testing/test_config_errors.json")
	configToTest := make(map[string]interface{})
	if err := json.Unmarshal(jsonFile, &configToTest); err != nil {
		assert.FailNow(t, "cannot create instruction message")
	}
	instructionMsg := Message{
		DeviceID: "front-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{
			{
				"instruction": "check config",
				"config":      configToTest,
				"name":        "testFile",
			},
		},
	}
	returnMsg := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "config",
		Contents: map[string]interface{}{
			"errors": []string{
				"level II - format error: time: missing unit in duration 30",
				"level II - format error: time: unknown unit x in duration 10x",
				"level IV - system error: host: different from current host for front and back-end",
				"level IV - system error: port: different from current port for front and back-end",
			},
			"name": "testFile",
		},
	}
	jsonMessage, _ := json.Marshal(&returnMsg)
	communicatorMock.On("Publish", "front-end", string(jsonMessage), 3)
	handler.msgMapper(instructionMsg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
}

func TestInstructionUseConfig(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	fileName := "../../../resources/testing/test_config.json"
	workingConfig := config.ReadFile(fileName)
	handler := Handler{
		Config:       workingConfig,
		ConfigFile:   fileName,
		Communicator: communicatorMock,
	}
	jsonFile, _ := ioutil.ReadFile(fileName)
	configToTest := make(map[string]interface{})
	if err := json.Unmarshal(jsonFile, &configToTest); err != nil {
		assert.FailNow(t, "cannot create instruction message")
	}
	instructionMsg := Message{
		DeviceID: "front-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: []map[string]interface{}{{
			"instruction": "use config",
			"config":      configToTest,
			"file":        "new_file.json"},
		},
	}
	returnMsg := Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "new config",
		Contents: map[string]interface{}{"name": "new_file.json"},
	}
	jsonMessage, _ := json.Marshal(&returnMsg)
	communicatorMock.On("Publish", "front-end", string(jsonMessage), 3).Once()
	communicatorMock.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("string"), 3) // sendSetup tested in another test
	handler.msgMapper(instructionMsg)
	communicatorMock.AssertNumberOfCalls(t, "Publish", 13)
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

func TestSendFrontEndInstruction(t *testing.T) {
	communicatorMock := new(CommunicatorMock)
	handler := Handler{
		Config:       config.ReadFile("../../../resources/testing/test_instruction_frontend.json"),
		Communicator: communicatorMock,
	}
	instMsg := []config.ComponentInstruction{
		{"gameState", "set state", "newState"},
	}
	msg, _ := json.Marshal(Message{
		DeviceID: "back-end",
		TimeSent: time.Now().Format("02-01-2006 15:04:05"),
		Type:     "instruction",
		Contents: instMsg,
	})
	communicatorMock.On("Publish", "front-end", string(msg), 3)
	handler.SendComponentInstruction("front-end", instMsg, "")
	communicatorMock.AssertNumberOfCalls(t, "Publish", 1)
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
