package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileError(t *testing.T) {
	filename := "missing.json"
	assert.PanicsWithValue(t,
		"Could not read file missing.json",
		func() { ReadFile(filename) },
		"Could not find json file")
}

func TestDurationError(t *testing.T) {
	filename := "../../../resources/testing/testDurationError.json"
	assert.PanicsWithValue(t,
		"json: cannot unmarshal number into Go struct field General.general.duration of type string",
		func() { ReadFile(filename) },
		"Incorrect json (duration in int) should panic")
}

func TestDeviceInputWrongTypeError(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testDeviceInputWrongTypeError.json"
	assert.PanicsWithValue(t,
		"json: cannot unmarshal number into Go struct field ReadDevice.devices.input of type string",
		func() { ReadFile(filename) },
		"Incorrect json (no input type in string format) should panic")
}

func TestReadFile(t *testing.T) {
	filename := "../../../resources/testing/test_config.json"
	assert.NotPanics(t, func() { ReadFile(filename) })
}

func TestDeviceConstraintNotPresent(t *testing.T) {
	filename := "../../../resources/testing/testDeviceConstraintNotPresent.json"
	assert.PanicsWithValue(t,
		"device with id non existing not found in map",
		func() { ReadFile(filename) },
		"ReadDevice used in constraint should be present in device logics")
}

func TestComponentConstraintNotPresent(t *testing.T) {
	filename := "../../../resources/testing/testComponentConstraintNotPresent.json"
	assert.PanicsWithValue(t,
		"device with id non existing not found in map",
		func() { ReadFile(filename) },
		"Component used in constraint should be present in device logics")
}

func TestIncorrectTypeCondition(t *testing.T) {
	filename := "../../../resources/testing/testIncorrectTypeCondition.json"
	assert.PanicsWithValue(t,
		"invalid type of condition: not device or timer",
		func() { ReadFile(filename) },
		"ReadCondition type should be 'device' or 'timer'")
}

func TestIncorrectConstraintOperation(t *testing.T) {
	filename := "../../../resources/testing/testIncorrectConstraintOperator.json"
	assert.PanicsWithValue(t,
		"JSON config in wrong format, operator: non existing operator, could not be processed",
		func() { ReadFile(filename) },
		"operator should be 'AND' or 'OR'")
}

func TestIncorrectConditionOperation(t *testing.T) {
	filename := "../../../resources/testing/testIncorrectConditionOperator.json"
	assert.PanicsWithValue(t,
		"JSON config in wrong format, operator: non existing operator, could not be processed",
		func() { ReadFile(filename) },
		"operator should be 'AND' or 'OR'")
}

func TestWrongConditionStructure(t *testing.T) {
	filename := "../../../resources/testing/testWrongConditionStructure.json"
	assert.PanicsWithValue(t,
		"JSON config in wrong condition format, conditions: map[non existing:non existing], could not be processed",
		func() { ReadFile(filename) },
		"condition should follow the condition or operator format")
}

func TestWrongConstraintStructure(t *testing.T) {
	filename := "../../../resources/testing/testWrongConstraintStructure.json"
	assert.PanicsWithValue(t,
		"JSON config in wrong constraint format, conditions: map[non existing:non existing], could not be processed",
		func() { ReadFile(filename) },
		"condition should follow the condition or operator format")
}

func TestWrongComponentIDType(t *testing.T) {
	filename := "../../../resources/testing/testWrongComponentIDType.json"
	assert.PanicsWithValue(t,
		"JSON config in wrong format, component_id should be of type string, 6 is of type float64",
		func() { ReadFile(filename) },
		"constraint should have a component_id of type string (if the condition is of type device)")
}
