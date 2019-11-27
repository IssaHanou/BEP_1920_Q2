package config

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func readFile(t *testing.T) ReadConfig {
	filename := "../../../resources/room_config.json"
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error("Could not read file room_config.json")
	}
	config := ReadJSON(dat)
	return config
}

func TestConfig(t *testing.T) {
	config := readFile(t)
	assert.Equal(t, config.General.Duration, "00:30:00", "Duration was not correct")
}

func TestDeviceInput(t *testing.T) {
	config := readFile(t)
	input := IOObject{
		"redSwitch":    "boolean",
		"orangeSwitch": "boolean",
		"greenSwitch":  "boolean",
		"slider1":      "integer",
		"slider2":      "integer",
		"slider3":      "integer",
		"mainSwitch":   "boolean",
	}
	assert.Equal(t, config.Devices[0].Input,
		input,
		"Incorrect input of device")
}

func TestActionOutput(t *testing.T) {
	config := readFile(t)
	input := IOObject{
		"greenLight1": false,
		"greenLight2": true,
		"greenLight3": false,
		"redLight1":   false,
		"redLight2":   false,
		"redLight3":   false,
	}
	assert.Equal(t, config.Puzzles[1].Rules[0].Actions[1].Message.Output,
		input,
		"Incorrect message output in action")
}

//TODO Condition inlezen
//TODO Condition or arrays
func TestConstraint(t *testing.T) {
	config := readFile(t)
	assert.Equal(t, config.Puzzles[1].Rules[0].Conditions[0].Constraints[1].ComponentID,
		"greenSwitch",
		"Incorrect message output in action")
}
