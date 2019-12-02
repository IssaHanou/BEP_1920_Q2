package config

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

// Forms constraint to processed constraint, which includes type. Value is also checked to be of that type.
func makeConstraints(devices map[string]Device, condition Condition) (map[string][]interface{}, error) {
	var constraints = make(map[string][]interface{})
	for _, cons := range condition.Constraints {
		// If the cons is enforced on a device, then the cons value type must match the input type.
		if condition.Type == "device" {
			generalConstraint, constraintError := createConstraint(cons, true)
			if constraintError != nil {
				return constraints, constraintError
			}
			device, ok := devices[condition.TypeID]
			if !ok {
				return constraints, errors.New("device with id " + condition.TypeID + " not found in map")
			}
			newConstraint, constraintType, err := checkDeviceType(cons, generalConstraint, device)
			if err != nil {
				return constraints, err
			}
			constraints[constraintType] = append(constraints[constraintType], newConstraint)
		} else
		// If the cons is enforced on a timer, then the cons value type must be in format "hh:mm:ss".
		if condition.Type == "timer" {
			generalConstraint, constraintError := createConstraint(cons, false)
			if constraintError != nil {
				return constraints, constraintError
			}
			if reflect.TypeOf(cons["value"]).Kind() != reflect.String {
				return constraints, errors.New(fmt.Sprint(cons["value"]) + " cannot be cast to string, for time constraint")
			}
			var rgxPat = regexp.MustCompile(`^[0-9]{2}:[0-9]{2}:[0-9]{2}$`)
			if !rgxPat.MatchString(cons["value"].(string)) {
				return constraints, errors.New(fmt.Sprint(cons["value"]) + " did not match pattern 'hh:mm:ss'")
			}
			constraints["timer"] = append(constraints["timer"],
				ConstraintTimer{generalConstraint, cons["value"].(string)})
		} else {
			return constraints, errors.New("invalid type of condition: " + condition.Type)
		}
	}
	return constraints, nil
}

// Checks constraint value to be of type, retrieved from device input.
func checkDeviceType(constraint ConstraintInfo, generalConstraint Constraint, device Device) (interface{}, string, error) {
	componentType, ok := device.Input[constraint["component_id"].(string)]
	if !ok {
		return Constraint{}, "", errors.New("component id: " + constraint["component_id"].(string) + " not found in device input")
	}
	err := CheckComponentType(componentType, constraint["value"])
	if err != nil {
		return Constraint{}, "", err
	}
	switch componentType {
	case "string":
		return ConstraintString{generalConstraint, constraint["value"].(string)}, "string", nil
	case "numeric":
		return ConstraintNumeric{generalConstraint, constraint["value"].(float64)}, "numeric", nil
	case "boolean":
		return ConstraintBool{generalConstraint, constraint["value"].(bool)}, "boolean", nil
	case "num-array":
		array := constraint["value"].([]interface{})
		var newArray []float64
		for i := range array {
			newArray = append(newArray, array[i].(float64))
		}
		return ConstraintNumericArray{generalConstraint, newArray}, "num-array", nil
	case "string-array":
		array := constraint["value"].([]interface{})
		var newArray []string
		for i := range array {
			newArray = append(newArray, array[i].(string))
		}
		return ConstraintStringArray{generalConstraint, newArray}, "string-array", nil
	case "bool-array":
		array := constraint["value"].([]interface{})
		var newArray []bool
		for i := range array {
			newArray = append(newArray, array[i].(bool))
		}
		return ConstraintBoolArray{generalConstraint, newArray}, "bool-array", nil
	default:
		// Value is of custom type, must be checked at client computer level
		return ConstraintCustomType{generalConstraint, constraint["value"]}, "custom", nil
	}
}

// Check if constraint can be constructed, with string comparison and string component id.
// includeCompo specifies whether to include the component id in the constraint, which is not the case for timer type.
func createConstraint(constraint ConstraintInfo, includeComponentID bool) (Constraint, error) {
	var err error
	if reflect.TypeOf(constraint["comparison"]).Kind() != reflect.String {
		return Constraint{}, errors.New("comparison '" + fmt.Sprint(constraint["comparison"]) + "' cannot be cast to string")
	}
	if !checkComparison(constraint["comparison"].(string)) {
		return Constraint{}, errors.New("comparison should be 'eq', 'lt', 'lte', 'gt', 'gte' or 'contains'")
	}
	if includeComponentID {
		if reflect.TypeOf(constraint["component_id"]).Kind() != reflect.String {
			return Constraint{}, errors.New("component_id " + fmt.Sprint(constraint["component_id"]) + " cannot be cast to string")
		}
		return Constraint{constraint["comparison"].(string), constraint["component_id"].(string)}, err
	}
	return Constraint{constraint["comparison"].(string), ""}, err
}

// Check if comparison is valid
func checkComparison(comparison string) bool {
	possibleComparisons := []string{"eq", "lt", "lte", "gt", "gte", "contains"}
	for _, ex := range possibleComparisons {
		if ex == comparison {
			return true
		}
	}
	return false
}

// GetConstraintTimer struct object from map
func GetConstraintTimer(config WorkingConfig, conditionID string) []ConstraintTimer {
	var constraintArray []ConstraintTimer
	if config.ConstraintMap[conditionID] == nil {
		panic(errors.New("Condition ID: " + conditionID + " not in constraint map").Error())
	}
	resultArray := config.ConstraintMap[conditionID]["timer"]
	for i := range resultArray {
		constraintArray = append(constraintArray, resultArray[i].(ConstraintTimer))
	}
	return constraintArray
}

// GetConstraintBool struct object from map
func GetConstraintBool(config WorkingConfig, conditionID string) []ConstraintBool {
	var constraintArray []ConstraintBool
	if config.ConstraintMap[conditionID] == nil {
		panic(errors.New("Condition ID: " + conditionID + " not in constraint map").Error())
	}
	resultArray := config.ConstraintMap[conditionID]["boolean"]
	for i := range resultArray {
		constraintArray = append(constraintArray, resultArray[i].(ConstraintBool))
	}
	return constraintArray
}

// GetConstraintString struct object from map
func GetConstraintString(config WorkingConfig, conditionID string) []ConstraintString {
	var constraintArray []ConstraintString
	if config.ConstraintMap[conditionID] == nil {
		panic(errors.New("Condition ID: " + conditionID + " not in constraint map").Error())
	}
	resultArray := config.ConstraintMap[conditionID]["string"]
	for i := range resultArray {
		constraintArray = append(constraintArray, resultArray[i].(ConstraintString))
	}
	return constraintArray
}

// GetConstraintNumeric struct object from map
func GetConstraintNumeric(config WorkingConfig, conditionID string) []ConstraintNumeric {
	var constraintArray []ConstraintNumeric
	if config.ConstraintMap[conditionID] == nil {
		panic(errors.New("Condition ID: " + conditionID + " not in constraint map").Error())
	}
	resultArray := config.ConstraintMap[conditionID]["numeric"]
	for i := range resultArray {
		constraintArray = append(constraintArray, resultArray[i].(ConstraintNumeric))
	}
	return constraintArray
}

// GetConstraintStringArray struct object from map
func GetConstraintStringArray(config WorkingConfig, conditionID string) []ConstraintStringArray {
	var constraintArray []ConstraintStringArray
	if config.ConstraintMap[conditionID] == nil {
		panic(errors.New("Condition ID: " + conditionID + " not in constraint map").Error())
	}
	resultArray := config.ConstraintMap[conditionID]["string-array"]
	for i := range resultArray {
		constraintArray = append(constraintArray, resultArray[i].(ConstraintStringArray))
	}
	return constraintArray
}

// GetConstraintBoolArray struct object from map
func GetConstraintBoolArray(config WorkingConfig, conditionID string) []ConstraintBoolArray {
	var constraintArray []ConstraintBoolArray
	if config.ConstraintMap[conditionID] == nil {
		panic(errors.New("Condition ID: " + conditionID + " not in constraint map").Error())
	}
	resultArray := config.ConstraintMap[conditionID]["bool-array"]
	for i := range resultArray {
		constraintArray = append(constraintArray, resultArray[i].(ConstraintBoolArray))
	}
	return constraintArray
}

// GetConstraintNumArray struct object from map
func GetConstraintNumArray(config WorkingConfig, conditionID string) []ConstraintNumericArray {
	var constraintArray []ConstraintNumericArray
	if config.ConstraintMap[conditionID] == nil {
		panic(errors.New("Condition ID: " + conditionID + " not in constraint map").Error())
	}
	resultArray := config.ConstraintMap[conditionID]["num-array"]
	for i := range resultArray {
		constraintArray = append(constraintArray, resultArray[i].(ConstraintNumericArray))
	}
	return constraintArray
}

// GetConstraintCustomType struct object from map
func GetConstraintCustomType(config WorkingConfig, conditionID string) []ConstraintCustomType {
	var constraintArray []ConstraintCustomType
	if config.ConstraintMap[conditionID] == nil {
		panic(errors.New("Condition ID: " + conditionID + " not in constraint map").Error())
	}
	resultArray := config.ConstraintMap[conditionID]["custom"]
	for i := range resultArray {
		constraintArray = append(constraintArray, resultArray[i].(ConstraintCustomType))
	}
	return constraintArray
}
