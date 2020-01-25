package config

import (
	"fmt"
	"reflect"
)

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
			errList = append(errList, checkActions(rule.Actions, config)...)
		}
	}

	for _, generalEvent := range config.GeneralEvents {
		for _, rule := range generalEvent.Rules {
			if err := rule.Conditions.checkConditions(config, rule.ID); err != nil {
				errList = append(errList, err...)
			}
			errList = append(errList, checkActions(rule.Actions, config)...)
		}
	}

	for _, rule := range config.ButtonEvents {
		if err := rule.Conditions.checkConditions(config, rule.ID); err != nil {
			errList = append(errList, err...)
		}
		errList = append(errList, checkActions(rule.Actions, config)...)
	}

	// TODO check uniqueness of all device_id, timer_id and rule_id
	return errList
}

// checkAction is a method that will return an error
// if the action's value types and instructions are not equal to the unit's output specifications
func checkActions(actions []Action, config WorkingConfig) []string {
	errorList := make([]string, 0)
	for _, action := range actions {
		switch action.Type {
		case "device":
			errorList = append(errorList, checkActionDevice(action, config)...)
		case "timer":
			errorList = append(errorList, checkActionTimer(action, config)...)
		case "label":
			errorList = append(errorList, checkActionLabel(action, config)...)
		default:
			errorList = append(errorList,
				fmt.Sprintf("level III - implementation error: only device, timer and label are accepted as type for an action, however type was specified as: %s", action.Type))
		}
	}
	return errorList
}

// checkActionTimer is a method that checks the config in use for mistakes in the actions of a timer
// if the config does not abide by the manual, a non-empty list of mistakes is returned
func checkActionTimer(action Action, config WorkingConfig) []string {
	errorList := make([]string, 0)
	if _, ok := config.Timers[action.TypeID]; ok { // checks if timer can be found in the map, if so, it is stored in variable device
		for _, actionMessage := range action.Message {
			switch actionMessage.Instruction {
			case "add", "subtract":
				valueType := reflect.TypeOf(actionMessage.Value).Kind()
				if valueType != reflect.String {
					errorList = append(errorList,
						fmt.Sprintf("level III - implementation error: input type string expected but %s found as type of value %v",
							valueType.String(), actionMessage.Value))
				}
				break
			case "start", "pause", "stop", "done":
				break
			default:
				errorList = append(errorList, fmt.Sprintf("level III - implementation error: instruction %s is not defined for a timer", actionMessage.Instruction))
			}
		}
	} else {
		errorList = append(errorList, fmt.Sprintf("level III - implementation error: timer with id %s not found in map", action.TypeID))
	}
	return errorList
}

// checkActionDevice is a method that checks the current config for mistakes in the actions of a device
// if the config does not abide by the manual, a non-empty list of mistakes is returned
func checkActionDevice(action Action, config WorkingConfig) []string {
	errorList := make([]string, 0)
	if device, ok := config.Devices[action.TypeID]; ok { // checks if device can be found in the map, if so, it is stored in variable device
		for _, actionMessage := range action.Message {
			if outputObject, ok := device.Output[actionMessage.ComponentID]; ok {
				if instructionType, ok := outputObject.Instructions[actionMessage.Instruction]; ok {
					if err := checkActionInstructionType(reflect.TypeOf(actionMessage.Value).Kind(), instructionType, actionMessage.Value); err != nil {
						errorList = append(errorList, err.Error())
					}
				} else {
					errorList = append(errorList, fmt.Sprintf("level III - implementation error: instruction %s not found in map", actionMessage.Instruction))
				}
			} else {
				errorList = append(errorList, fmt.Sprintf("level III - implementation error: component with id %s not found in map", actionMessage.ComponentID))
			}
		}
	} else {
		errorList = append(errorList, fmt.Sprintf("level III - implementation error: device with id %s not found in map", action.TypeID))
	}
	return errorList
}

// checkActionInstructionType checks if the type op the value of an instruction is the same as the type the instruction requires according to the config
func checkActionInstructionType(valueType reflect.Kind, instructionType string, value interface{}) error {
	switch instructionType {
	case "string":
		if valueType != reflect.String {
			return fmt.Errorf("level III - implementation error: instruction type string expected but %s found as type of value %v",
				valueType.String(), value)
		}
	case "boolean":
		if valueType != reflect.Bool {
			return fmt.Errorf("level III - implementation error: instruction type boolean expected but %s found as type of value %v",
				valueType.String(), value)
		}
	case "numeric":
		if valueType != reflect.Int && valueType != reflect.Float64 {
			return fmt.Errorf("level III - implementation error: instruction type numeric expected but %s found as type of value %v",
				valueType.String(), value)
		}
	case "array":
		if valueType != reflect.Slice {
			return fmt.Errorf("level III - implementation error: instruction type array/slice expected but %s found as type of value %v",
				valueType.String(), value)
		}
	default:
		return fmt.Errorf("level III - implementation error: custom types like: %s, are not yet implemented", instructionType)
	}
	return nil
}

// checkActionLabel checks if there is a label with this ID,
// and checks if all components under this label have the correct instructions with a call to checkActionDevice
// if the config does not abide by the manual, a non-empty list of mistakes is returned
func checkActionLabel(action Action, config WorkingConfig) []string {
	errorList := make([]string, 0)
	if _, ok := config.LabelMap[action.TypeID]; ok { // checks if label can be found in the map, if so, it is stored in variable device
		for _, instruction := range action.Message {
			for _, comp := range config.LabelMap[action.TypeID] {
				instruction.ComponentID = comp.ID
				errorList = append(errorList,
					checkActionDevice(Action{TypeID: comp.Device.ID, Message: []ComponentInstruction{instruction}}, config)...)
			}
		}
	} else {
		errorList = append(errorList,
			fmt.Sprintf("level III - implementation error: label with id %s not found in map", action.TypeID))
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
			return constraint.checkConstraintsDeviceType(inputType, ruleID)
		} else if outputObject, ok := device.Output[constraint.ComponentID]; ok {
			return constraint.checkConstraintsDeviceType(outputObject.Type, ruleID)
		} else {
			return []string{fmt.Sprintf("level III - implementation error: on rule %s: component with id %s not found in map", ruleID, constraint.ComponentID)}
		}
	} else {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s: device with id %s not found in map", ruleID, condition.TypeID)}
	}
}

// checkConstraintsTimer is a method that check all types and comparators of a constraint on a timer
func checkConstraintsTimer(condition Condition, config WorkingConfig, ruleID string, constraint Constraint) []string {
	if _, ok := config.Timers[condition.TypeID]; ok {
		return checkConstraintsBooleanType(ruleID, constraint)
	}
	return []string{fmt.Sprintf("level III - implementation error: on rule %s: timer with id %s not found in map", ruleID, condition.TypeID)}
}

// checkConstraintsRule is a method that check all types and comparators of a constraint on a rule
func checkConstraintsRule(condition Condition, config WorkingConfig, ruleID string, constraint Constraint) []string {
	if _, ok := config.RuleMap[condition.TypeID]; ok { // checks if rule can be found in the map, if so, it is stored in variable device
		valueType := reflect.TypeOf(constraint.Value).Kind()
		comparison := constraint.Comparison
		if valueType != reflect.Int && valueType != reflect.Float64 {
			return []string{fmt.Sprintf("level III - implementation error: on rule %s: value type numeric expected but %s found as type of value %v", ruleID, valueType.String(), constraint.Value)}
		}
		if !checkValidComparison(comparison) {
			return []string{fmt.Sprintf("level III - implementation error: on rule %s: comparison %s is not valid", ruleID, comparison)}
		}
		if comparison == "contains" {
			return []string{fmt.Sprintf("level III - implementation error: on rule %s: comparison %s not allowed on rule", ruleID, comparison)}
		}
	} else {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s: rule with id %s not found in map", ruleID, condition.TypeID)}
	}
	// all cases for errors are already handled
	return make([]string, 0)
}

// checkConstraintsDeviceType is a method that checks type defined by the config and check if the constraint has that value type
func (constraint Constraint) checkConstraintsDeviceType(typeToCheck string, ruleID string) []string {
	switch typeToCheck {
	case "string":
		return checkConstraintsDeviceStringType(ruleID, constraint)
	case "boolean":
		return checkConstraintsBooleanType(ruleID, constraint)
	case "numeric":
		return checkConstraintsDeviceNumericType(ruleID, constraint)
	case "array":
		return checkConstraintsDeviceArrayType(ruleID, constraint)
	default:
		return []string{fmt.Sprintf("level III - implementation error: on rule %s: custom types like: %s, are not yet implemented", ruleID, typeToCheck)}
	}
}

// checkConstraintsDeviceStringType is a method that returns all error (if any)
// in a constraint of a device with string type constraint
func checkConstraintsDeviceStringType(ruleID string, constraint Constraint) []string {
	valueType := reflect.TypeOf(constraint.Value).Kind()
	comparison := constraint.Comparison
	if valueType != reflect.String {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s: type string expected but %s found as type of value %v", ruleID, valueType.String(), constraint.Value)}
	}
	if !checkValidComparison(comparison) {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s: comparison %s is not valid", ruleID, comparison)}
	}
	if comparison != "eq" && comparison != "not" {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s: comparison %s not allowed on a string", ruleID, comparison)}
	}
	// all cases for errors are already handled
	return make([]string, 0)
}

// checkConstraintsBooleanType is a method that returns all error (if any)
// in a constraint of a device with boolean type constraint
func checkConstraintsBooleanType(ruleID string, constraint Constraint) []string {
	valueType := reflect.TypeOf(constraint.Value).Kind()
	comparison := constraint.Comparison
	if valueType != reflect.Bool {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s: type boolean expected but %s found as type of value %v", ruleID, valueType.String(), constraint.Value)}
	}
	if !checkValidComparison(comparison) {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s: comparison %s is not valid", ruleID, comparison)}
	}
	if comparison != "eq" {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s: comparison %s not allowed on a boolean", ruleID, comparison)}
	}
	// all cases for errors are already handled
	return make([]string, 0)
}

// checkConstraintsDeviceNumericType is a method that returns all error (if any)
// in a constraint of a device with numeric type constraint
func checkConstraintsDeviceNumericType(ruleID string, constraint Constraint) []string {
	valueType := reflect.TypeOf(constraint.Value).Kind()
	comparison := constraint.Comparison
	if valueType != reflect.Int && valueType != reflect.Float64 {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s: type numeric expected but %s found as type of value %v", ruleID, valueType.String(), constraint.Value)}
	}
	if !checkValidComparison(comparison) {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s: comparison %s is not valid", ruleID, comparison)}
	}
	if comparison == "contains" {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s: comparison %s not allowed on a numeric", ruleID, comparison)}
	}
	// all cases for errors are already handled
	return make([]string, 0)
}

// checkConstraintsDeviceArrayType is a method that returns all error (if any)
// in a constraint of a device with array type constraint
func checkConstraintsDeviceArrayType(ruleID string, constraint Constraint) []string {
	valueType := reflect.TypeOf(constraint.Value).Kind()
	comparison := constraint.Comparison
	if valueType != reflect.Slice {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s: type array/slice expected but %s found as type of value %v", ruleID, valueType.String(), constraint.Value)}
	}
	if !checkValidComparison(comparison) {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s: comparison %s is not valid", ruleID, comparison)}
	}
	if comparison != "contains" && comparison != "eq" && comparison != "not" {
		return []string{fmt.Sprintf("level III - implementation error: on rule %s: comparison %s not allowed on an array", ruleID, comparison)}
	}
	// all cases for errors are already handled
	return make([]string, 0)
}
