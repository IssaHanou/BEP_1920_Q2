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
					"testComponent3": "array",
				},
				Output: map[string]config.OutputObject{
					"testComponent4": {Type: "boolean", Instructions: map[string]string{}},
					"testComponent5": {Type: "numeric", Instructions: map[string]string{}},
					"testComponent6": {Type: "string", Instructions: map[string]string{}},
					"testComponent7": {Type: "array", Instructions: map[string]string{}},
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
