package handler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	handler.UpdateStatus(msg)

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
	handler.UpdateStatus(msg)

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
	handler.UpdateStatus(msg)

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
	handler.UpdateStatus(msg)

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
