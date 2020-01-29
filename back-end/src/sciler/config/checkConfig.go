package config

import (
	"fmt"
	"reflect"
	"strings"
)

// checkUniqueIDs checks whether all timers, devices and rules have unique ids compared to each other.
func checkUniqueIDs(config *WorkingConfig) []string {
	idList := make(map[string]string, 0) // the value keeps track of type (rule/timer/device) to put in error message
	errorList := make([]string, 0)
	for timerID := range config.Timers { // these timer id's are already checked for uniqueness
		idList[timerID] = "timer"
	}
	idList, deviceErrors := checkDeviceUniqueIDs(idList, config)
	_, ruleErrors := checkRuleUniqueIDs(idList, config)
	return append(errorList, append(deviceErrors, ruleErrors...)...)
}

// checkDeviceUniqueIDs checks whether the devices do not have any ids in common with the timers in the config
func checkDeviceUniqueIDs(idList map[string]string, config *WorkingConfig) (map[string]string, []string) {
	errorList := make([]string, 0)
	for deviceID := range config.Devices {
		if value, ok := idList[deviceID]; ok {
			errorList = append(errorList,
				fmt.Sprintf("level III - implementation error: checking device with id %s, but a %s with id %s already exists",
					deviceID, value, deviceID))
		} else {
			idList[deviceID] = "device"
		}
	}
	return idList, errorList
}

// checkRuleUniqueIDs checks whether the rules do not have any ids in common with the timers or devices in the config
func checkRuleUniqueIDs(idList map[string]string, config *WorkingConfig) (map[string]string, []string) {
	errorList := make([]string, 0)
	for ruleID := range config.RuleMap {
		if value, ok := idList[ruleID]; ok {
			errorList = append(errorList,
				fmt.Sprintf("level III - implementation error: checking rule with id %s, but a %s with id %s already exists",
					ruleID, value, ruleID))
		} else {
			idList[ruleID] = "rule"
		}
	}
	return idList, errorList
}

// checkConfig is a method that will return an error
// if the constraints value type is not equal to the device input type specified,
// the actions type is not equal to the device output type,
// or some other not allowed json configuration
func checkConfig(config WorkingConfig) []string {
	errList := make([]string, 0)
	for _, puzzle := range config.Puzzles {
		for _, rule := range puzzle.Event.Rules {
			if err := rule.Conditions.checkConditions(config, rule.ID); err != nil {
				errList = append(errList, err...)
			}
			errList = append(errList, checkActions(rule.Actions, config, rule.ID)...)
		}
	}

	for _, generalEvent := range config.GeneralEvents {
		for _, rule := range generalEvent.Rules {
			if err := rule.Conditions.checkConditions(config, rule.ID); err != nil {
				errList = append(errList, err...)
			}
			errList = append(errList, checkActions(rule.Actions, config, rule.ID)...)
		}
	}

	for _, rule := range config.ButtonEvents {
		if err := rule.Conditions.checkConditions(config, rule.ID); err != nil {
			errList = append(errList, err...)
		}
		errList = append(errList, checkActions(rule.Actions, config, rule.ID)...)
	}
	return errList
}

// checkAction is a method that will return an error
// if the action's value types and instructions are not equal to the unit's output specifications
func checkActions(actions []Action, config WorkingConfig, ruleID string) []string {
	errorList := make([]string, 0)
	for _, action := range actions {
		switch action.Type {
		case "device":
			errorList = append(errorList, action.checkActionDevice(config, ruleID)...)
		case "timer":
			errorList = append(errorList, action.checkActionTimer(config, ruleID)...)
		case "label":
			errorList = append(errorList, action.checkActionLabel(config, ruleID)...)
		default:
			errorList = append(errorList,
				fmt.Sprintf("level III - implementation error: on rule %s: only device, timer and label are accepted as type for an action, however type was specified as: %s", ruleID, action.Type))
		}
	}
	return errorList
}

// checkActionTimer is a method that checks the config in use for mistakes in the actions of a timer
// if the config does not abide by the manual, a non-empty list of mistakes is returned
func (action Action) checkActionTimer(config WorkingConfig, ruleID string) []string {
	errorList := make([]string, 0)
	if _, ok := config.Timers[action.TypeID]; ok { // checks if timer can be found in the map, if so, it is stored in variable device
		for _, actionMessage := range action.Message {
			switch actionMessage.Instruction {
			case "add", "subtract":
				if err := checkTimeAlterInstruction(actionMessage, []string{ruleID, action.TypeID}); err != "" {
					errorList = append(errorList, err)
				}
				break
			case "start", "pause", "stop", "done":
				break
			default:
				errorList = append(errorList, fmt.Sprintf("level III - implementation error: on rule %s, actions for timer with id %s: instruction '%s' is not defined for a timer", ruleID, action.TypeID, actionMessage.Instruction))
			}
		}
	} else {
		errorList = append(errorList, fmt.Sprintf("level III - implementation error: on rule %s's actions: timer with id %s not found in timer map", ruleID, action.TypeID))
	}
	return errorList
}

// checkTimeAlterInstruction checks if the instruction value is not nil and if its type is a string
// errorParameters contains [ruleID, action.TypeID] to put in the error messages
func checkTimeAlterInstruction(actionMessage ComponentInstruction, errorParameters []string) string {
	if actionMessage.Value == nil {
		return fmt.Sprintf("level III - implementation error: on rule %s, action for timer with id %s: value of action message is nil", errorParameters[0], errorParameters[1])
	}
	valueType := reflect.TypeOf(actionMessage.Value).Kind()
	if valueType != reflect.String {
		return fmt.Sprintf("level III - implementation error: on rule %s, actions for timer with id %s: input type string expected but %s found as type of value %v", errorParameters[0], errorParameters[1], valueType.String(), actionMessage.Value)
	}
	return ""
}

// checkActionDevice is a method that checks the current config for mistakes in the actions of a device
// if the config does not abide by the manual, a non-empty list of mistakes is returned
func (action Action) checkActionDevice(config WorkingConfig, ruleID string) []string {
	errorList := make([]string, 0)
	if device, ok := config.Devices[action.TypeID]; ok { // checks if device can be found in the map, if so, it is stored in variable device
		for _, actionMessage := range action.Message {
			if outputObject, ok := device.Output[actionMessage.ComponentID]; ok {
				if err := config.checkOutputObject(outputObject, actionMessage, []string{ruleID, action.TypeID}); err != "" {
					errorList = append(errorList, err)
				}
			} else {
				errorList = append(errorList, fmt.Sprintf("level III - implementation error: on rule %s, actions for device with id %s: component with id %s not found in map", ruleID, action.TypeID, actionMessage.ComponentID))
			}
		}
	} else {
		errorList = append(errorList, fmt.Sprintf("level III - implementation error: on rule %s, actions for device with id %s: device with id %s not found in device map", ruleID, action.TypeID, action.TypeID))
	}
	return errorList
}

// checkOutputObject checks if the value of action message is not nil and if the instruction type is correct.
// errorParameters contains [ruleID, action.TypeID] to put in the error message
func (config *WorkingConfig) checkOutputObject(outputObject OutputObject, actionMessage ComponentInstruction, errorParameters []string) string {
	if instructionType, ok := outputObject.Instructions[actionMessage.Instruction]; ok {
		if actionMessage.Value == nil {
			return fmt.Sprintf("level III - implementation error: on rule %s, action for device with id %s: value of action message is nil",
				errorParameters[0], errorParameters[1])
		} else if err := config.checkActionInstructionType(reflect.TypeOf(actionMessage.Value).Kind(), instructionType, actionMessage.Value, errorParameters[0]); err != nil {
			return err.Error()
		}
		return ""
	}
	return fmt.Sprintf("level III - implementation error: on rule %s, action for device with id %s: instruction '%s' not found in map",
		errorParameters[0], errorParameters[1], actionMessage.Instruction)
}

// checkActionInstructionType checks if the type op the value of an instruction is the same as the type the instruction requires according to the config
func (config *WorkingConfig) checkActionInstructionType(valueType reflect.Kind, instructionType string, value interface{}, ruleID string) error {
	switch instructionType {
	case "string":
		if valueType != reflect.String {
			return fmt.Errorf("level III - implementation error: on rule %s's actions: instruction type string expected but %s found as type of the value: %v", ruleID, valueType.String(), value)
		}
	case "boolean":
		if valueType != reflect.Bool {
			return fmt.Errorf("level III - implementation error: on rule %s's actions: instruction type boolean expected but %s found as type of the value: %v", ruleID, valueType.String(), value)
		}
	case "numeric":
		if valueType != reflect.Int && valueType != reflect.Float64 {
			return fmt.Errorf("level III - implementation error: on rule %s's actions: instruction type numeric expected but %s found as type of the value: %v", ruleID, valueType.String(), value)
		}
	case "array":
		if valueType != reflect.Slice {
			return fmt.Errorf("level III - implementation error: on rule %s's actions: instruction type array/slice expected but %s found as type of the value: %v", ruleID, valueType.String(), value)
		}
	case "status":
		if valueType != reflect.String {
			return fmt.Errorf("level III - implementation error: on rule %s's actions: instruction type status expected but %s found as type of the value: %v", ruleID, valueType.String(), value)
		} else {
			split := strings.Split(value.(string), ".")
			if len(split) == 2 {
				deviceID := split[0]
				componentID := split[1]
				device := config.Devices[deviceID]
				if device != nil && (device.Input[componentID] != "" || device.Output[componentID].Type != "") {
					break
				}
			}
			return fmt.Errorf("level III - implementation error: on rule %s's actions: instruction type status expected but %v which is not in the form of `deviceID.component`", ruleID, value)
		}
	default:
		return fmt.Errorf("level III - implementation error: on rule %s's actions: custom types of value like: %s, are not yet implemented", ruleID, instructionType)
	}
	return nil
}

// checkActionLabel checks if there is a label with this ID,
// and checks if all components under this label have the correct instructions with a call to checkActionDevice
// if the config does not abide by the manual, a non-empty list of mistakes is returned
func (action Action) checkActionLabel(config WorkingConfig, ruleID string) []string {
	errorList := make([]string, 0)
	if _, ok := config.LabelMap[action.TypeID]; ok { // checks if label can be found in the map, if so, it is stored in variable device
		for _, instruction := range action.Message {
			for _, comp := range config.LabelMap[action.TypeID] {
				instruction.ComponentID = comp.ID
				errorList = append(errorList,
					Action{TypeID: comp.Device.ID, Message: []ComponentInstruction{instruction}}.checkActionDevice(config, ruleID)...)
			}
		}
	} else {
		errorList = append(errorList,
			fmt.Sprintf("level III - implementation error: on rule %s's actions: label with id %s not found in label map",
				ruleID, action.TypeID))
	}
	return errorList
}

// checkValidComparison checks if the comparison is a valid one (one that can be used in a condition)
func checkValidComparison(comparison string) bool {
	comparisonTypesAllowed := []string{"eq", "lt", "gt", "lte", "gte", "contains", "not"}
	for _, comp := range comparisonTypesAllowed {
		if comp == comparison {
			return true
		}
	}
	return false
}

// checkConditions is a method that checks types and comparator operators by
// running through all the conditions in de OrCondition and check all those conditions
func (or OrCondition) checkConditions(config WorkingConfig, ruleID string) []string {
	errorList := make([]string, 0)
	for _, logic := range or.logics {
		err := logic.checkConditions(config, ruleID)
		if err != nil {
			errorList = append(errorList, err...)
		}
	}
	return errorList
}

// checkConditions is a method that checks types and comparator operators by
// running through all the conditions in de AndCondition and check all those conditions
func (and AndCondition) checkConditions(config WorkingConfig, ruleID string) []string {
	errorList := make([]string, 0)
	for _, logic := range and.logics {
		err := logic.checkConditions(config, ruleID)
		if err != nil {
			errorList = append(errorList, err...)
		}
	}
	return errorList
}

// checkConditions is a method that checks the constraints in a condition
// this is different from the checkConditions on an OrCondition or an AndCondition since those contain a list of conditions
func (condition Condition) checkConditions(config WorkingConfig, ruleID string) []string {
	return condition.Constraints.checkConstraints(condition, config, ruleID)
}

// checkConstraints is a method that checks types and comparator operators by
// running through all the conditions in de OrConstraint and check all those constraints
func (or OrConstraint) checkConstraints(condition Condition, config WorkingConfig, ruleID string) []string {
	errorList := make([]string, 0)
	for _, logic := range or.logics {
		err := logic.checkConstraints(condition, config, ruleID)
		if err != nil {
			errorList = append(errorList, err...)
		}
	}
	return errorList
}

// checkConstraints is a method that checks types and comparator operators by
// running through all the conditions in de AndConstraint and check all those constraints
func (and AndConstraint) checkConstraints(condition Condition, config WorkingConfig, ruleID string) []string {
	errorList := make([]string, 0)
	for _, logic := range and.logics {
		err := logic.checkConstraints(condition, config, ruleID)
		if err != nil {
			errorList = append(errorList, err...)
		}
	}
	return errorList
}

// checkConstraints is a method that checks types and comparator operators by
// checking the conditionType and call the correct check for that type
func (constraint Constraint) checkConstraints(condition Condition, config WorkingConfig, ruleID string) []string {
	switch condition.Type {
	case "device":
		return checkConstraintsDevice(condition, config, ruleID, constraint)
	case "timer":
		return checkConstraintsTimer(condition, config, ruleID, constraint)
	case "rule":
		return checkConstraintsRule(condition, config, ruleID, constraint)
	default:
		return []string{fmt.Sprintf("level III - implementation error: on rule %s: invalid type of condition: %v", ruleID, condition.Type)}
	}
}

// checkConstraintsDevice is a method that check all types and comparators of a constraint on a device
func checkConstraintsDevice(condition Condition, config WorkingConfig, ruleID string, constraint Constraint) []string {
	if device, ok := config.Devices[condition.TypeID]; ok { // checks if device can be found in the map, if so, it is stored in variable device
		if inputType, ok := device.Input[constraint.ComponentID]; ok {
			return constraint.checkConstraintsDeviceType(inputType, []string{ruleID, condition.TypeID})
		} else if outputObject, ok := device.Output[constraint.ComponentID]; ok {
			return constraint.checkConstraintsDeviceType(outputObject.Type, []string{ruleID, condition.TypeID})
		} else {
			return []string{fmt.Sprintf("level III - implementation error: on rule %s, constraint on rule with id %s: component with id %s not found in device input or output map",
				ruleID, condition.TypeID, constraint.ComponentID)}
		}
	} else {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s, constraint on rule with id %s: device with id %s not found in device map",
			ruleID, condition.TypeID, condition.TypeID)}
	}
}

// checkConstraintsTimer is a method that check all types and comparators of a constraint on a timer
func checkConstraintsTimer(condition Condition, config WorkingConfig, ruleID string, constraint Constraint) []string {
	if _, ok := config.Timers[condition.TypeID]; ok {
		return checkConstraintsBooleanType([]string{ruleID, condition.TypeID, "timer"}, constraint)
	}
	return []string{fmt.Sprintf("level III - implementation error: on rule %s, constraint on rule with id %s: timer with id %s not found in timer map",
		ruleID, condition.TypeID, condition.TypeID)}
}

// checkConstraintsRule is a method that check all types and comparators of a constraint on a rule
func checkConstraintsRule(condition Condition, config WorkingConfig, ruleID string, constraint Constraint) []string {
	if _, ok := config.RuleMap[condition.TypeID]; ok { // checks if rule can be found in the map, if so, it is stored in variable device
		valueType := reflect.TypeOf(constraint.Value).Kind()
		comparison := constraint.Comparison
		if valueType != reflect.Int && valueType != reflect.Float64 {
			return []string{fmt.Sprintf("level III - implementation error: on rule %s, constraint on rule with id %s: value type numeric expected but %s found as type of te value: %v", ruleID, condition.TypeID, valueType.String(), constraint.Value)}
		} else if !checkValidComparison(comparison) {
			return []string{fmt.Sprintf("level III - implementation error: on rule %s, constraint on rule with id %s: comparison '%s' is not valid", ruleID, condition.TypeID, comparison)}
		} else if comparison == "contains" {
			return []string{fmt.Sprintf("level III - implementation error: on rule %s, constraint on rule with id %s: comparison '%s' not allowed on a rule constraint", ruleID, condition.TypeID, comparison)}
		}
	} else {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s, constraint on rule with id %s: rule with id %s not found in rule map", ruleID, condition.TypeID, condition.TypeID)}
	}
	return make([]string, 0) // all cases for errors are already handled
}

// checkConstraintsDeviceType is a method that checks type defined by the config and check if the constraint has that value type
// errorParameters contains [ruleID, deviceID] to put in error message
func (constraint Constraint) checkConstraintsDeviceType(typeToCheck string, errorParameters []string) []string {
	switch typeToCheck {
	case "string":
		return checkConstraintsDeviceStringType(errorParameters, constraint)
	case "boolean":
		return checkConstraintsBooleanType(append(errorParameters, "device"), constraint)
	case "numeric":
		return checkConstraintsDeviceNumericType(errorParameters, constraint)
	case "array":
		return checkConstraintsDeviceArrayType(errorParameters, constraint)
	default:
		return []string{fmt.Sprintf("level III - implementation error: on rule %s, constraint on device with id %s: custom types of value like: %s, are not yet implemented",
			errorParameters[0], errorParameters[1], typeToCheck)}
	}
}

// checkConstraintsDeviceStringType is a method that returns all error (if any)
// in a constraint of a device with string type constraint
// errorParameters contains [ruleID, deviceID] to put in error message
func checkConstraintsDeviceStringType(errorParameters []string, constraint Constraint) []string {
	valueType := reflect.TypeOf(constraint.Value).Kind()
	comparison := constraint.Comparison
	if valueType != reflect.String {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s, constraint on device with id %s: type string expected but %s found as type of the value: %v",
			errorParameters[0], errorParameters[1], valueType.String(), constraint.Value)}
	}
	if !checkValidComparison(comparison) {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s, constraint on device with id %s: comparison '%s' is not valid",
			errorParameters[0], errorParameters[1], comparison)}
	}
	if comparison != "eq" && comparison != "not" {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s, constraint on device with id %s: comparison %s not allowed on a string constraint",
			errorParameters[0], errorParameters[1], comparison)}
	}
	// all cases for errors are already handled
	return make([]string, 0)
}

// checkConstraintsBooleanType is a method that returns all error (if any)
// in a constraint of a device with boolean type constraint
// errorParameters contains [ruleID, deviceID, 'timer'/'device'] to put in error message
func checkConstraintsBooleanType(errorParameters []string, constraint Constraint) []string {
	valueType := reflect.TypeOf(constraint.Value).Kind()
	comparison := constraint.Comparison
	if valueType != reflect.Bool {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s, constraint on %s with id %s: type boolean expected but %s found as type of the value: %v",
			errorParameters[0], errorParameters[2], errorParameters[1], valueType.String(), constraint.Value)}
	}
	if !checkValidComparison(comparison) {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s, constraint on %s with id %s: comparison '%s' is not valid",
			errorParameters[0], errorParameters[2], errorParameters[1], comparison)}
	}
	if comparison != "eq" {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s, constraint on %s with id %s: comparison %s not allowed on a boolean constraint",
			errorParameters[0], errorParameters[2], errorParameters[1], comparison)}
	}
	// all cases for errors are already handled
	return make([]string, 0)
}

// checkConstraintsDeviceNumericType is a method that returns all error (if any)
// in a constraint of a device with numeric type constraint
// errorParameters contains [ruleID, deviceID] to put in error message
func checkConstraintsDeviceNumericType(errorParameters []string, constraint Constraint) []string {
	valueType := reflect.TypeOf(constraint.Value).Kind()
	comparison := constraint.Comparison
	if valueType != reflect.Int && valueType != reflect.Float64 {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s, constraint on device with id %s: type numeric expected but %s found as type of the value: %v",
			errorParameters[0], errorParameters[1], valueType.String(), constraint.Value)}
	}
	if !checkValidComparison(comparison) {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s, constraint on device with id %s: comparison '%s' is not valid",
			errorParameters[0], errorParameters[1], comparison)}
	}
	if comparison == "contains" {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s, constraint on device with id %s: comparison %s not allowed on a numeric constraint",
			errorParameters[0], errorParameters[1], comparison)}
	}
	// all cases for errors are already handled
	return make([]string, 0)
}

// checkConstraintsDeviceArrayType is a method that returns all error (if any)
// in a constraint of a device with array type constraint
// errorParameters contains [ruleID, deviceID] to put in error message
func checkConstraintsDeviceArrayType(errorParameters []string, constraint Constraint) []string {
	valueType := reflect.TypeOf(constraint.Value).Kind()
	comparison := constraint.Comparison
	if valueType != reflect.Slice {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s, constraint on device with id %s: type array/slice expected but %s found as type of the value: %v",
			errorParameters[0], errorParameters[1], valueType.String(), constraint.Value)}
	}
	if !checkValidComparison(comparison) {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s, constraint on device with id %s: comparison '%s' is not valid",
			errorParameters[0], errorParameters[1], comparison)}
	}
	if comparison != "contains" && comparison != "eq" && comparison != "not" {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s, constraint on device with id %s: comparison %s not allowed on an array constraint",
			errorParameters[0], errorParameters[1], comparison)}
	}
	// all cases for errors are already handled
	return make([]string, 0)
}
