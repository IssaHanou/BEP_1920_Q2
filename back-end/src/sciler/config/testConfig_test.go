package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/////////////////////////////////// Test correct test_config
func TestGeneralInformation(t *testing.T) {
	filename := "../../../resources/testing/test_config.json"
	config := ReadFile(filename)
	general := General{"Escape X", "00:30:00", "192.0.0.84", 1883}
	assert.Equal(t, general, config.General,
		"General information should be correct")
}

func TestPuzzleSize(t *testing.T) {
	filename := "../../../resources/testing/test_config.json"
	config := ReadFile(filename)
	assert.Equal(t, 3, len(config.Puzzles), "Should have read two puzzles")
}

func TestDeviceInput(t *testing.T) {
	filename := "../../../resources/testing/test_config.json"
	config := ReadFile(filename)
	input := map[string]string{
		"redSwitch":    "boolean",
		"orangeSwitch": "boolean",
		"greenSwitch":  "boolean",
		"slider1":      "numeric",
		"slider2":      "numeric",
		"slider3":      "numeric",
		"mainSwitch":   "boolean",
	}
	assert.Equal(t, input, config.Devices["controlBoard"].Input,
		"Input of device should be retrieved correctly from the devices map")
}

func TestActionOutput(t *testing.T) {
	filename := "../../../resources/testing/test_config.json"
	config := ReadFile(filename)
	output := OutputObject{
		"greenLight1": false,
		"greenLight2": true,
		"greenLight3": false,
		"redLight1":   false,
		"redLight2":   false,
		"redLight3":   false,
	}
	assert.Equal(t, output, config.Puzzles[1].Rules[0].Actions[1].Message.Output,
		"Message from action should be of OutputObject type, retrieved through puzzles in config")
}

func TestRulesMap(t *testing.T) {
	filename := "../../../resources/testing/test_config.json"
	config := ReadFile(filename)
	assert.Equal(t, "De juiste volgorde van cijfers moet gedraaid worden.",
		config.Rules["correctSequence"].Description,
		"Description from rule should retrieved correctly through rules map")
}

func TestGeneralEvents(t *testing.T) {
	filename := "../../../resources/testing/test_config.json"
	config := ReadFile(filename)
	assert.Equal(t, "Start", config.GeneralEvents[0].Name,
		"Name of general event should be retrieved correctly")
}

func TestComponentIDForDeviceConstraintNum(t *testing.T) {
	filename := "../../../resources/testing/test_config.json"
	config := ReadFile(filename)
	list := GetConstraintNumeric(config, config.Puzzles[1].Rules[0].Conditions[0].GetID())
	assert.Equal(t, 1, len(list), "There should be one numeric constraint for this condition")
	assert.Equal(t, "slider2", list[0].ComponentID,
		"ComponentID should be retrieved correctly through puzzle, rule, condition, constraint")
}

func TestComponentIDForDeviceConstraintBool(t *testing.T) {
	filename := "../../../resources/testing/test_config.json"
	config := ReadFile(filename)
	list := GetConstraintBool(config, config.Puzzles[1].Rules[0].Conditions[0].GetID())
	assert.Equal(t, 3, len(list), "There should be three boolean constraints for this condition")
	assert.Equal(t, "greenSwitch", list[1].ComponentID,
		"ComponentID should be retrieved correctly through puzzle, rule, condition, constraint")
}

func TestComponentIDForDeviceConstraintDouble(t *testing.T) {
	filename := "../../../resources/testing/test_config.json"
	config := ReadFile(filename)
	list := GetConstraintNumeric(config, config.Puzzles[2].Rules[0].Conditions[0].GetID())
	assert.Equal(t, 1, len(list), "There should be one numeric constraint for this condition")
	assert.Equal(t, "numeric", list[0].ComponentID,
		"ComponentID should be retrieved correctly through puzzle, rule, condition, constraint")
	assert.Equal(t, 2.5, list[0].Value,
		"Custom type constraint should have correctly retrieved value")
}

func TestComponentIDForDeviceConstraintString(t *testing.T) {
	filename := "../../../resources/testing/test_config.json"
	config := ReadFile(filename)
	list := GetConstraintString(config, config.Puzzles[2].Rules[0].Conditions[0].GetID())
	assert.Equal(t, 1, len(list), "There should be one string constraint for this condition")
	assert.Equal(t, "string", list[0].ComponentID,
		"ComponentID should be retrieved correctly through puzzle, rule, condition, constraint")
	assert.Equal(t, "mystring", list[0].Value,
		"Custom type constraint should have correctly retrieved value")
}

func TestComponentIDForDeviceConstraintNumArray(t *testing.T) {
	filename := "../../../resources/testing/test_config.json"
	config := ReadFile(filename)
	list := GetConstraintNumArray(config, config.Puzzles[2].Rules[0].Conditions[0].GetID())
	assert.Equal(t, 1, len(list), "There should be one num-array constraint for this condition")
	assert.Equal(t, "num-array", list[0].ComponentID,
		"ComponentID should be retrieved correctly through puzzle, rule, condition, constraint")
	assert.Equal(t, []float64{0, -1, 0.5, 25, 9.12}, list[0].Value,
		"Custom type constraint should have correctly retrieved value")
}

func TestComponentIDForDeviceConstraintStringArray(t *testing.T) {
	filename := "../../../resources/testing/test_config.json"
	config := ReadFile(filename)
	list := GetConstraintStringArray(config, config.Puzzles[2].Rules[0].Conditions[0].GetID())
	assert.Equal(t, 1, len(list), "There should be one num-array constraint for this condition")
	assert.Equal(t, "string-array", list[0].ComponentID,
		"ComponentID should be retrieved correctly through puzzle, rule, condition, constraint")
	assert.Equal(t, []string{"mystring1", "mystring2"}, list[0].Value,
		"Custom type constraint should have correctly retrieved value")
}

func TestComponentIDForDeviceConstraintBoolArray(t *testing.T) {
	filename := "../../../resources/testing/test_config.json"
	config := ReadFile(filename)
	list := GetConstraintBoolArray(config, config.Puzzles[2].Rules[0].Conditions[0].GetID())
	assert.Equal(t, 1, len(list), "There should be one bool-array constraint for this condition")
	assert.Equal(t, "bool-array", list[0].ComponentID,
		"ComponentID should be retrieved correctly through puzzle, rule, condition, constraint")
	assert.Equal(t, []bool{true, false, true}, list[0].Value,
		"Custom type constraint should have correctly retrieved value")
}

func TestComponentIDForDeviceConstraintCustom(t *testing.T) {
	filename := "../../../resources/testing/test_config.json"
	config := ReadFile(filename)
	list := GetConstraintCustomType(config, config.Puzzles[2].Rules[0].Conditions[0].GetID())
	assert.Equal(t, 1, len(list), "There should be one custom type constraint for this condition")
	assert.Equal(t, "custom", list[0].ComponentID,
		"ComponentID should be retrieved correctly through puzzle, rule, condition, constraint")
	output := map[string]interface{}{"instruction": "test"}
	assert.Equal(t, output, list[0].Value,
		"Custom type constraint should have correctly retrieved value")
}

func TestNoComponentForTimeConstraint(t *testing.T) {
	filename := "../../../resources/testing/test_config.json"
	config := ReadFile(filename)
	list := GetConstraintTimer(config, config.GeneralEvents[0].Rules[0].Conditions[0].GetID())
	assert.Equal(t, 1, len(list), "There should be one timer constraint for this condition")
	assert.Equal(t, "", list[0].ComponentID,
		"There should not be a component id for a timer constraint, it should be ''")
	assert.Equal(t, "00:01:00", list[0].Value,
		"Timer constraint should have correctly retrieved value")
}

/////////////////////////////////// Edge case behavior
func TestDeviceInputCustomType(t *testing.T) {
	filename := "../../../resources/testing/testDeviceCustomType.json"
	config := ReadFile(filename)
	for key, value := range config.Devices["telephone"].Input {
		assert.Equal(t, "turningWheel", key,
			"Id of component should be key in input map")
		assert.Equal(t, "my-type", value,
			"Custom type of component should be value in input map")
	}
}

func TestDeviceOutputCustomType(t *testing.T) {
	filename := "../../../resources/testing/testDeviceCustomType.json"
	config := ReadFile(filename)
	for key, value := range config.Devices["telephone"].Output {
		assert.Equal(t, "audio", key,
			"Id of component should be key in input map")
		assert.Equal(t, "string", value,
			"Custom type of component should be value in input map")
	}
}
