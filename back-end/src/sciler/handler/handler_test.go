package handler

import (
	"github.com/stretchr/testify/assert"
	"sciler/communication"
	"sciler/config"
	"testing"
)

func getTestHandler() *Handler {
	config := config.WorkingConfig{
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
				Output: map[string]interface{}{
					"testComponent4": false,
					"testComponent5": true,
					"testComponent6": false,
				},
				Status:     map[string]interface{}{},
				Connection: false,
			},
		},
		Rules:         nil,
		ActionMap:     nil,
		ConstraintMap: nil,
	}
	communicator := communication.NewCommunicator(config.General.Host,
		config.General.Port, []string{"back-end", "test"})
	return GetHandler(config, *communicator)
}

func Test_GetHandler(t *testing.T) {
	config := config.WorkingConfig{
		General: config.General{
			Name:     "Test",
			Duration: "1",
			Host:     "localhost",
			Port:     1883,
		},
		Puzzles:       nil,
		GeneralEvents: nil,
		Devices:       nil,
		Rules:         nil,
		ActionMap:     nil,
		ConstraintMap: nil,
	}
	communicator := communication.NewCommunicator(config.General.Host,
		config.General.Port, []string{"back-end", "test"})

	tests := []struct {
		name string
		want *Handler
	}{
		{
			name: "test",
			want: &Handler{
				Config:       config,
				Communicator: *communicator,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetHandler(config, *communicator)
			assert.Equal(t, got.Config, tt.want.Config)
			assert.Equal(t, got.Communicator, tt.want.Communicator)
		})
	}

}

func TestMsgMapperStatus(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "",
		Type:     "status",
		Contents: map[string]interface{}{
			"testComponent0": false,
			"testComponent1": true,
			"testComponent2": false},
	}
	handler.msgMapper(msg)

	assert.Equal(t, true, handler.Config.Devices["TestDevice"].Status["testComponent1"],
		"Device should set connection to true on connection message")

}
func TestOnConnectionMsg(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "",
		Type:     "connection",
		Contents: map[string]interface{}{
			"connection": true,
		},
	}
	handler.onConnectionMsg(msg)

	assert.Equal(t, true, handler.Config.Devices["TestDevice"].Connection,
		"Device should set connection to true on connection message")
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
		TimeSent: "",
		Type:     "connection",
		Contents: map[string]interface{}{
			"connection": true,
		},
	}
	handler.msgMapper(msg)

	assert.Equal(t, true, handler.Config.Devices["TestDevice"].Connection,
		"Device should set connection to true on connection message")
}

func TestOnStatusMsg(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "",
		Type:     "status",
		Contents: map[string]interface{}{
			"testComponent0": false,
			"testComponent1": true,
			"testComponent2": false},
	}
	handler.onStatusMsg(msg)

	assert.Equal(t, true, handler.Config.Devices["TestDevice"].Status["testComponent1"],
		"Device should set connection to true on connection message")
}

func TestOnConfirmationMsg(t *testing.T) {
	assert.Equal(t, true, true,
		"TODO: TestOnConfirmationMsg")
}

func TestOnInstructionMsg(t *testing.T) {
	assert.Equal(t, true, true,
		"TODO: TestOnInstructionMsg")
}

func TestMsgMapper(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "",
		Type:     "test",
		Contents: map[string]interface{}{
			"testComponent0": false,
			"testComponent1": true,
			"testComponent2": false},
	}

	before := handler.Config
	handler.onStatusMsg(msg)
	assert.Equal(t, before, handler.Config,
		"Nothing should have bee changed after an incorrect message type")
}
