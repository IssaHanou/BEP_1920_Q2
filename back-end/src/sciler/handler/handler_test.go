package handler

import (
	"encoding/json"
	"fmt"
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
		Rules:         nil,
		ActionMap:     nil,
		ConstraintMap: nil,
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
			"testComponent4": false,
			"testComponent5": true,
			"testComponent6": false},
	}
	handler.onConnectionMsg(msg)

	assert.Equal(t, true, handler.Config.Devices["TestDevice"].Connection,
		"Device should set connection to true on connection message")
}

func TestMsgMapperConnection(t *testing.T) {
	handler := getTestHandler()
	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "",
		Type:     "connection",
		Contents: map[string]interface{}{
			"testComponent4": false,
			"testComponent5": true,
			"testComponent6": false},
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
	handler := getTestHandler()
	text := []byte(`{
		"completed": true,
		"instructed": {
			"instruction": "test"
		}
	}`)
	var data map[string]interface{}
	if err := json.Unmarshal(text, &data); err != nil {
		assert.Fail(t, "Should be valid json message")
	}
	fmt.Println(data)

	msg := Message{
		DeviceID: "TestDevice",
		TimeSent: "",
		Type:     "confirmation",
		Contents: data,
	}
	assert.NotPanics(t, func() { handler.onConfirmationMsg(msg) },
		"Device should return valid confirmation message")
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
