package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/////////////////////////////////// Test errors
func TestDurationError(t *testing.T) {
	filename := "../../../resources/testing/testDurationError.json"
	assert.PanicsWithValue(t,
		"json: cannot unmarshal number into Go struct field General.general.duration of type string",
		func() { ReadFile(filename) },
		"Incorrect json (duration in int) should panic")
}

func TestFileError(t *testing.T) {
	filename := "missing.json"
	assert.PanicsWithValue(t,
		"Could not read file missing.json",
		func() { ReadFile(filename) },
		"Could not find json file")
}

/////////////////////////////////// Wrong reference to device
func TestDeviceInputWrongTypeError(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testDeviceInputWrongTypeError.json"
	assert.PanicsWithValue(t,
		"json: cannot unmarshal number into Go struct field ReadDevice.devices.input of type string",
		func() { ReadFile(filename) },
		"Incorrect json (no input type in string format) should panic")
}

// Currently not relevant as device output is allowed to take in a map of type/instruction
//func TestDeviceOutputWrongTypeError(t *testing.T) {
//	filename := "../../../resources/testing/wrong-types/testDeviceOutputWrongTypeError.json"
//	assert.PanicsWithValue(t,
//		"json: cannot unmarshal number into Go struct field ReadDevice.devices.output of type string",
//		func() { ReadFile(filename) },
//		"Incorrect json (no output type in string format) should panic")
//}

func TestDeviceActionNotPresent(t *testing.T) {
	filename := "../../../resources/testing/testDeviceActionNotPresent.json"
	assert.PanicsWithValue(t,
		"device with id telephone not found in map",
		func() { ReadFile(filename) },
		"ReadDevice used in action should be present in device list")
}

func TestComponentActionNotPresent(t *testing.T) {
	filename := "../../../resources/testing/testComponentActionNotPresent.json"
	assert.PanicsWithValue(t,
		"component id: light not found in device input",
		func() { ReadFile(filename) },
		"Component used in action should be present in device list")
}

func TestDeviceConstraintNotPresent(t *testing.T) {
	filename := "../../../resources/testing/testDeviceConstraintNotPresent.json"
	assert.PanicsWithValue(t,
		"device with id telephone not found in map",
		func() { ReadFile(filename) },
		"ReadDevice used in constraint should be present in device list")
}

func TestComponentConstraintNotPresent(t *testing.T) {
	filename := "../../../resources/testing/testComponentConstraintNotPresent.json"
	assert.PanicsWithValue(t,
		"component id: color not found in device input",
		func() { ReadFile(filename) },
		"Component used in constraint should be present in device list")
}

func TestIncorrectTypeCondition(t *testing.T) {
	filename := "../../../resources/testing/testIncorrectTypeCondition.json"
	assert.PanicsWithValue(t,
		"invalid type of condition: mytype",
		func() { ReadFile(filename) },
		"Condition type should be 'device' or 'timer'")
}

func TestIncorrectTypeAction(t *testing.T) {
	filename := "../../../resources/testing/testIncorrectTypeAction.json"
	assert.PanicsWithValue(t,
		"invalid type of action: my-type",
		func() { ReadFile(filename) },
		"Condition type should be 'device' or 'timer'")
}

func TestCheckActionComponentType(t *testing.T) {
	filename := "../../../resources/testing/testCheckActionComponentType.json"
	assert.PanicsWithValue(t,
		"Value was not of type string: 0",
		func() { ReadFile(filename) },
		"Condition type should be 'device' or 'timer'")
}

/////////////////////////////////// Test timer actions
func TestTimerNoInstruction(t *testing.T) {
	filename := "../../../resources/testing/timer/testTimerNoInstruction.json"
	assert.PanicsWithValue(t,
		"timer should have an instruction defined",
		func() { ReadFile(filename) },
		"Timer should have an instruction defined")
}

func TestTimerIncorrectInstruction(t *testing.T) {
	filename := "../../../resources/testing/timer/testTimerIncorrectInstruction.json"
	assert.PanicsWithValue(t,
		"timer should have an instruction defined, which is either stop or subtract",
		func() { ReadFile(filename) },
		"Timer should have a correct instruction defined: stop or subtract")
}

func TestTimerNoSubtractValue(t *testing.T) {
	filename := "../../../resources/testing/timer/testTimerNoSubtractValue.json"
	assert.PanicsWithValue(t,
		"timer with subtract instruction should have value",
		func() { ReadFile(filename) },
		"Timer should have a value defined for subtract instruction")
}

func TestTimerIncorrectSubtractValue(t *testing.T) {
	filename := "../../../resources/testing/timer/testTimerIncorrectSubtractValue.json"
	assert.PanicsWithValue(t,
		"timer with subtract instruction should have value in string format",
		func() { ReadFile(filename) },
		"Timer should have a string value defined for subtract instruction")
}

func TestTimerInvalidSubtractValue(t *testing.T) {
	filename := "../../../resources/testing/timer/testTimerInvalidSubtractValue.json"
	assert.PanicsWithValue(t,
		"test did not match pattern 'hh:mm:ss'",
		func() { ReadFile(filename) },
		"Timer should have a value defined for subtract instruction in format 'hh:mm:ss'")
}

/////////////////////////////////// Wrong inputs
func TestTimerIncorrectPattern(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testTimerIncorrectPattern.json"
	assert.PanicsWithValue(t,
		"30:00 did not match pattern 'hh:mm:ss'",
		func() { ReadFile(filename) },
		"Time should be entered in correct pattern")
}

func TestTimerIncorrectType(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testTimerIncorrectType.json"
	assert.PanicsWithValue(t,
		"30 cannot be cast to string, for time constraint",
		func() { ReadFile(filename) },
		"Timer should take in string value")
}

func TestStringIncorrectType(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testStringIncorrectType.json"
	assert.PanicsWithValue(t,
		"Value was not of type string: 30",
		func() { ReadFile(filename) },
		"If device input specifies string type, constraint value should be a string")
}

func TestNumericIncorrectType(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testNumericIncorrectType.json"
	assert.PanicsWithValue(t,
		"Value was not of type numeric: 30",
		func() { ReadFile(filename) },
		"If device input specifies numeric type, constraint value should be a numeric")
}

func TestBooleanIncorrectType(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testBooleanIncorrectType.json"
	assert.PanicsWithValue(t,
		"Value was not of type boolean: 30",
		func() { ReadFile(filename) },
		"If device input specifies boolean type, constraint value should be a boolean")
}

func TestNumArrayIncorrectType(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testNumArrayIncorrectType.json"
	assert.PanicsWithValue(t,
		"Value was not of type num-array: [true]",
		func() { ReadFile(filename) },
		"If device input specifies num array type, constraint value should be a num array")
}

func TestNumArrayIncorrectType2(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testNumArrayIncorrectType2.json"
	assert.PanicsWithValue(t,
		"Value was not of type array: test",
		func() { ReadFile(filename) },
		"If device input specifies num array type, constraint value should be a num array")
}

func TestStringArrayIncorrectType(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testStringArrayIncorrectType.json"
	assert.PanicsWithValue(t,
		"Value was not of type string-array: [30]",
		func() { ReadFile(filename) },
		"If device input specifies string array type, constraint value should be a string array")
}

func TestStringArrayIncorrectType2(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testStringArrayIncorrectType2.json"
	assert.PanicsWithValue(t,
		"Value was not of type array: 30",
		func() { ReadFile(filename) },
		"If device input specifies string array type, constraint value should be a string array")
}

func TestBoolArrayIncorrectType(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testBoolArrayIncorrectType.json"
	assert.PanicsWithValue(t,
		"Value was not of type bool-array: [30]",
		func() { ReadFile(filename) },
		"If device input specifies bool array type, constraint value should be a bool array")
}

func TestBoolArrayIncorrectType2(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testBoolArrayIncorrectType2.json"
	assert.PanicsWithValue(t,
		"Value was not of type array: 30",
		func() { ReadFile(filename) },
		"If device input specifies bool array type, constraint value should be a bool array")
}

func TestIncorrectComparisonType(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testIncorrectTypeComparison.json"
	assert.PanicsWithValue(t,
		"comparison '30' cannot be cast to string",
		func() { ReadFile(filename) },
		"Comparison should be entered in string format")
}

func TestInvalidComparisonType(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testInvalidTypeComparison.json"
	assert.PanicsWithValue(t,
		"comparison should be 'eq', 'lt', 'lte', 'gt', 'gte' or 'contains'",
		func() { ReadFile(filename) },
		"Comparison should be a correct type")
}

func TestIncorrectComponentIDType(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testIncorrectTypeComponentID.json"
	assert.PanicsWithValue(t,
		"component_id 30 cannot be cast to string",
		func() { ReadFile(filename) },
		"Component ID should be entered in string format")
}

/////////////////////////////////// Missing condition keys
func TestNoComponentIDForDeviceConstraintBool(t *testing.T) {
	filename := "../../../resources/testing/test_config.json"
	config := ReadFile(filename)
	assert.PanicsWithValue(t,
		"Condition ID: ID not in constraint map",
		func() { GetConstraintBool(config, "ID") },
		"Panic thrown when ID not found for bool")
	assert.PanicsWithValue(t,
		"Condition ID: ID not in constraint map",
		func() { GetConstraintString(config, "ID") },
		"Panic thrown when ID not found for string")
	assert.PanicsWithValue(t,
		"Condition ID: ID not in constraint map",
		func() { GetConstraintNumeric(config, "ID") },
		"Panic thrown when ID not found for numeric")
	assert.PanicsWithValue(t,
		"Condition ID: ID not in constraint map",
		func() { GetConstraintStringArray(config, "ID") },
		"Panic thrown when ID not found for string array")
	assert.PanicsWithValue(t,
		"Condition ID: ID not in constraint map",
		func() { GetConstraintNumArray(config, "ID") },
		"Panic thrown when ID not found for num array")
	assert.PanicsWithValue(t,
		"Condition ID: ID not in constraint map",
		func() { GetConstraintBoolArray(config, "ID") },
		"Panic thrown when ID not found for bool array")
	assert.PanicsWithValue(t,
		"Condition ID: ID not in constraint map",
		func() { GetConstraintCustomType(config, "ID") },
		"Panic thrown when ID not found for custom type")
	assert.PanicsWithValue(t,
		"Condition ID: ID not in constraint map",
		func() { GetConstraintTimer(config, "ID") },
		"Panic thrown when ID not found for timer")
}
