package handler

import (
	"github.com/stretchr/testify/assert"
	"sciler/communication"
	"sciler/config"
	"testing"
)

func getTestHandler() *Handler {
	workingConfig := config.WorkingConfig{
		General: config.General{
			Name:     "Test",
			Duration: "1",
			Host:     "localhost",
			Port:     1883,
		},
		Puzzles:       nil,
		GeneralEvents: nil,
		Devices: map[string]config.Device{
			"TestDevice": {
				ID:          "TestDevice",
				Description: "test uitleg",
				Input: map[string]string{
					"testComponent0": "boolean",
					"testComponent1": "numeric",
					"testComponent2": "string",
				},
				Output: map[string]config.OutputObject{
					"testComponent4": {Type: "string", Instructions: map[string]string{}},
					"testComponent5": {Type: "string", Instructions: map[string]string{}},
					"testComponent6": {Type: "string", Instructions: map[string]string{}},
				},
				Status:     map[string]interface{}{},
				Connection: false,
			},
		},
	}
	communicator := communication.NewCommunicator(workingConfig.General.Host,
		workingConfig.General.Port, []string{"back-end", "test"})
	return GetHandler(workingConfig, *communicator)
}

func Test_GetHandler(t *testing.T) {
	workingConfig := config.WorkingConfig{
		General: config.General{
			Name:     "Test",
			Duration: "1",
			Host:     "localhost",
			Port:     1883,
		},
		Puzzles:       nil,
		GeneralEvents: nil,
		Devices:       nil,
	}
	communicator := communication.NewCommunicator(workingConfig.General.Host,
		workingConfig.General.Port, []string{"back-end", "test"})

	tests := []struct {
		name string
		want *Handler
	}{
		{
			name: "test",
			want: &Handler{
				Config:       workingConfig,
				Communicator: *communicator,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetHandler(workingConfig, *communicator)
			assert.Equal(t, got.Config, tt.want.Config)
			assert.Equal(t, got.Communicator, tt.want.Communicator)
		})
	}
}

func TestOnConnectionMsg(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "connection",
		Contents: map[string]interface{}{
			"connection": true},
	}
	handler.onConnectionMsg(msg)
	assert.Equal(t, true, handler.Config.Devices["TestDevice"].Connection,
		"Device should set connection to true on connection message")
}

func TestMsgMapperConnectionOtherDevice(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice2",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "connection",
		Contents: map[string]interface{}{
			"connection": true},
	}
	handler.msgMapper(msg)

	_, ok := handler.Config.Devices["TestDevice2"]
	assert.Equal(t, false, ok,
		"Device should not exist in devices because it was not in config")
}

func TestMsgMapperConnectionOtherDeviceInvalid(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice2",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "connection",
		Contents: map[string]interface{}{
			"connection": "true"},
	}
	handler.onConnectionMsg(msg)
	assert.Equal(t, false, handler.Config.Devices["TestDevice"].Connection,
		"Device should not set connection to true on incorrect connection message")
}

func TestMsgMapperConnectionOtherDeviceInvalid2(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice2",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "connection",
		Contents: map[string]interface{}{
			"connected": "true"},
	}
	handler.onConnectionMsg(msg)
	assert.Equal(t, false, handler.Config.Devices["TestDevice"].Connection,
		"Device should not set connection to true on incorrect connection message")
}

func TestOnConnectionMsgFalse(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "",
		Type:     "connection",
		Contents: map[string]interface{}{
			"connection": false,
		},
	}
	handler.onConnectionMsg(msg)

	assert.Equal(t, false, handler.Config.Devices["TestDevice"].Connection,
		"Device should set connection to false on connection message")
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

func TestOnStatusMsg(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "status",
		Contents: map[string]interface{}{
			"testComponent0": false,
			"testComponent1": true,
			"testComponent2": false},
	}
	handler.onStatusMsg(msg)

	assert.Equal(t, true, handler.Config.Devices["TestDevice"].Status["testComponent1"],
		"Device should set status to true on component 1")
}

func TestOnStatusMsgOtherDevice(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice2",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "status",
		Contents: map[string]interface{}{
			"testComponent0": false,
			"testComponent1": true,
			"testComponent2": false},
	}
	handler.onStatusMsg(msg)

	_, ok := handler.Config.Devices["TestDevice2"]
	assert.Equal(t, false, ok,
		"Device should not exist in devices because it was not in config")
}

func TestMsgMapperStatus(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "status",
		Contents: map[string]interface{}{
			"testComponent0": false,
			"testComponent1": true,
			"testComponent2": false},
	}
	handler.msgMapper(msg)

	assert.Equal(t, true, handler.Config.Devices["TestDevice"].Status["testComponent1"],
		"Device should set status for component 1 to true on status message")

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
				"device_id": "back-end",
				"time_sent": "05-12-2019 09:42:10",
				"contents":  map[string]interface{}{"instructions": "test"},
				"type":      "instructions",
			},
		},
	}
	assert.NotPanics(t, func() { handler.onConfirmationMsg(msg) },
		"Device should return valid confirmation message")
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
				"device_id": "back-end",
				"time_sent": "05-12-2019 09:42:10",
				"contents":  map[string]interface{}{"instructions": "test"},
				"type":      "instructions",
			},
		},
	}
	assert.NotPanics(t, func() { handler.onConfirmationMsg(msg) },
		"Device should return valid confirmation message")
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
				"contents":  map[string]interface{}{"instructions": "test"},
				"type":      "instructions",
			},
		},
	}
	assert.NotPanics(t, func() { handler.onConfirmationMsg(msg) },
		"Device should not panic on incorrect json with no completed key")
}

func TestOnConfirmationMsgIncorrect2(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "05-12-2019 09:42:10",
		Type:     "confirmation",
		Contents: map[string]interface{}{
			"completed": true,
			"instructions": map[string]interface{}{
				"device_id": "back-end",
				"time_sent": "05-12-2019 09:42:10",
				"contents":  map[string]interface{}{"instructions": "test"},
				"type":      "instructions",
			},
		},
	}
	assert.NotPanics(t, func() { handler.onConfirmationMsg(msg) },
		"Device should not panic on incorrect json with no instructed key")
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
				"contents":  map[string]interface{}{"instructions": "test"},
				"type":      "instructions",
			},
		},
	}
	assert.NotPanics(t, func() { handler.onConfirmationMsg(msg) },
		"Device should not panic on json with no boolean completed value")
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
				"contents":  map[string]interface{}{"instructions": "test"},
				"type":      "instructions",
			},
		},
	}
	before := handler.Config
	handler.msgMapper(msg)
	assert.Equal(t, before, handler.Config,
		"Device should not config with confirmation message")

}

func TestOnInstructionMsg(t *testing.T) {
	assert.Equal(t, true, true,
		"TODO: TestOnInstructionMsg")
}

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
		"Nothing should have bee changed after an incorrect message type")
}
