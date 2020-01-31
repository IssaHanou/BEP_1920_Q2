package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileError(t *testing.T) {
	filename := "missing.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"Could not find json file")
}

func TestDurationError(t *testing.T) {
	filename := "../../../resources/testing/testDurationError.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"Incorrect json (duration in int) should panic")
}

func TestDurationNil(t *testing.T) {
	filename := "../../../resources/testing/testDurationMissing.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"Incorrect json (duration in int) should panic")
}

func TestDurationErrorWrongFormat(t *testing.T) {
	filename := "../../../resources/testing/testDurationErrorWrongFormat.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"Incorrect json (duration without unit specification) should panic")
}

func TestDurationTimerError(t *testing.T) {
	filename := "../../../resources/testing/testDurationTimerError.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"Incorrect json (timer duration in int) should panic")
}

func TestDurationTimerErrorWrongFormat(t *testing.T) {
	filename := "../../../resources/testing/testDurationTimerErrorWrongFormat.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"Incorrect json (timer duration with incorrect unit specification) should panic")
}

func TestDeviceInputWrongTypeError(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testDeviceInputWrongTypeError.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"Incorrect json (no input type in string format) should panic")
}

func TestReadFile(t *testing.T) {
	filename := "../../../resources/testing/test_config.json"
	assert.NotPanics(t, func() { ReadFile(filename) })
}

func TestDeviceConstraintNotPresent(t *testing.T) {
	filename := "../../../resources/testing/testDeviceConstraintNotPresent.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"ReadDevice used in constraint should be present in device logics")
}

func TestComponentConstraintNotPresent(t *testing.T) {
	filename := "../../../resources/testing/testComponentConstraintNotPresent.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"Component used in constraint should be present in device logics")
}

func TestIncorrectTypeCondition(t *testing.T) {
	filename := "../../../resources/testing/testIncorrectTypeCondition.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"ReadCondition type should be 'device' or 'timer'")
}

func TestIncorrectConstraintOperation(t *testing.T) {
	filename := "../../../resources/testing/testIncorrectConstraintOperator.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"operator should be 'AND' or 'OR'")
}

func TestIncorrectConditionOperation(t *testing.T) {
	filename := "../../../resources/testing/testIncorrectConditionOperator.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"operator should be 'AND' or 'OR'")
}

func TestWrongConditionStructure(t *testing.T) {
	filename := "../../../resources/testing/testWrongConditionStructure.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"condition should follow the condition or operator format")
}

func TestWrongConstraintStructure(t *testing.T) {
	filename := "../../../resources/testing/testWrongConstraintStructure.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"condition should follow the condition or operator format")
}

func TestWrongComponentIDType(t *testing.T) {
	filename := "../../../resources/testing/testWrongComponentIDType.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"constraint should have a component_id of type string (if the condition is of type device)")
}

func Test_CheckActionWrongType(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckActionWrongType.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"Only device and timer are not supported as action type")

}

func Test_CheckActionWrongDevice(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckActionWrongDevice.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"Cannot perform an action on an unknown device")
}

func Test_CheckActionWrongTimer(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckActionWrongTimer.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"Cannot perform an action on an unknown timer")
}

func Test_CheckActionWrongComponent(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckActionWrongComponent.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"Cannot perform an action on an unknown component")
}

func Test_CheckActionWrongInstruction(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckActionWrongInstruction.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"Cannot perform an action with an unknown instruction")
}

func Test_CheckActionWrongInstructionLabel(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckLabelWrongInstruction.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"Cannot perform an action with an unknown instruction")
}

func Test_CheckActionWrongLabel(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckLabelWrong.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"Cannot perform an action with an unknown instruction")
}

func Test_CheckActionWrongTimerInstruction(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckActionWrongTimerInstruction.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"Cannot perform an action with an unknown instruction")
}

func Test_CheckActionCustom(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckActionCustom.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"custom types are not implemented yet")
}

func Test_CheckActionString(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckActionString.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"custom types are not implemented yet")
}

func Test_CheckActionBoolean(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckActionBoolean.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"instruction type does not match given value in an action")
}

func Test_CheckActionNumeric(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckActionNumeric.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"instruction type does not match given value in an action")
}

func Test_CheckActionArray(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckActionArray.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"instruction type does not match given value in an action")
}

func Test_CheckActionNoValueDevice(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckActionNoValueDevice.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"action does not have a value for front-end set state instruction")
}

func Test_CheckActionNoValueTimer(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckActionNoValueTimer.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"action does not have a value for timer add instruction")
}

func Test_CheckConstraintNoValueDevice(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckConstraintNoValueDevice.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"constraint does not have a value for front-end condition")
}

func Test_CheckConstraintNoValueRule(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckConstraintNoValueRule.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"constraint does not have a value for rule condition")
}

func Test_CheckConstraintNoValueTimer(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckConstraintNoValueTimer.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"constraint does not have a value for timer condition")
}

func Test_GenerateEmptyConditions(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/test_emptyConditions.json"
	assert.NotPanics(t,
		func() { ReadFile(filename) },
		"empty conditions should not panic")
}

func Test_GenerateEmptyConstraints(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/test_emptyConstraints.json"
	assert.NotPanics(t,
		func() { ReadFile(filename) },
		"empty constraints should not panic")
}

func Test_GenerateEmptyConstraintList(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/test_emptyConstraintList.json"
	assert.NotPanics(t,
		func() { ReadFile(filename) },
		"empty constraint list should not panic")
}

func Test_GenerateConditionList(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/test_emptyConditionList.json"
	assert.NotPanics(t,
		func() { ReadFile(filename) },
		"empty condition list should not panic")
}

func Test_GenerateNoConditions(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/test_noConditions.json"
	assert.NotPanics(t,
		func() { ReadFile(filename) },
		"no conditions key in config should not panic")
}

func Test_GenerateNoConstraints(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/test_noConstraints.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"no constraints key in config should panic, because it is an invalid condition")
}

func Test_DoubleRuleId(t *testing.T) {
	filename := "../../../resources/testing/test_doubleRuleId.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"two rules with the same id is not allowed")
}

func Test_DoubleDeviceId(t *testing.T) {
	filename := "../../../resources/testing/test_doubleDeviceId.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"two devices with the same id is not allowed")
}

func Test_DoubleTimerId(t *testing.T) {
	filename := "../../../resources/testing/test_doubleTimerId.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"two timers with the same id is not allowed")
}

func Test_DoubleIDs(t *testing.T) {
	filename := "../../../resources/testing/test_doubleIDs.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"two timers with the same id is not allowed")
}
