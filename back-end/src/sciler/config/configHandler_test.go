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
		"component id: color not found in device input",
		func() { ReadFile(filename) },
		"Component used in constraint should be present in device logics")
}

func TestIncorrectTypeCondition(t *testing.T) {
	filename := "../../../resources/testing/testIncorrectTypeCondition.json"
	assert.PanicsWithValue(t,
		"invalid type of condition: mytype",
		func() { ReadFile(filename) },
		"ReadCondition type should be 'device' or 'timer'")
}

func Test_checkConfig(t *testing.T) {

}

func Test_generateDataStructures(t *testing.T) {

}

func Test_generateGeneralEvent(t *testing.T) {

}

func Test_generateGeneralEvents(t *testing.T) {

}

func Test_generateLogicalCondition(t *testing.T) {

}

func Test_generateLogicalConstraint(t *testing.T) {

}

func Test_generatePuzzles(t *testing.T) {

}

func Test_generateRules(t *testing.T) {

}
