package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

// ReadJSON transforms json file into config object.
func ReadJSON(input []byte) WorkingConfig {
	var config ReadConfig
	jsonErr := json.Unmarshal(input, &config)
	if jsonErr != nil {
		panic(jsonErr.Error())
	}
	newConfig, configErr := generateDataStructures(config)
	if configErr != nil {
		panic(configErr.Error())
	}
	return newConfig
}

// Creates additional structures: forms device and rule maps;
// and maps for actions and constraints (retrieved from puzzles and general events), with condition pointer as key
func generateDataStructures(readConfig ReadConfig) (WorkingConfig, error) {
	var config WorkingConfig
	// Copy information from read config to working config.
	config.General = readConfig.General
	config.Puzzles = readConfig.Puzzles
	config.GeneralEvents = readConfig.GeneralEvents
	config.Devices = make(map[string]Device)
	for _, d := range readConfig.Devices {
		config.Devices[d.ID] = d
	}

	// Create additional data structures.
	config.Rules = make(map[string]Rule)
	config.ActionMap = make(map[string][]Action)
	config.ConstraintMap = make(map[string]map[string][]interface{})
	for i, p := range config.Puzzles {
		for j, r := range p.Rules {
			config.Rules[r.ID] = r
			for k, c := range r.Conditions {
				// Set both original pointer's and current pointer's conditions
				config.Puzzles[i].Rules[j].Conditions[k].RuleID = r.ID
				c.RuleID = r.ID
				config.ActionMap[c.GetID()] = r.Actions
				constraints, err := makeConstraints(config.Devices, c)
				if err != nil {
					return config, err
				}
				config.ConstraintMap[c.GetID()] = constraints
			}
		}
	}
	for i, e := range config.GeneralEvents {
		for j, r := range e.Rules {
			config.Rules[r.ID] = r
			for k, c := range r.Conditions {
				// Set both original pointer's and current pointer's conditions
				config.GeneralEvents[i].Rules[j].Conditions[k].RuleID = r.ID
				c.RuleID = r.ID
				config.ActionMap[c.GetID()] = r.Actions
				constraints, err := makeConstraints(config.Devices, c)
				if err != nil {
					return config, err
				}
				config.ConstraintMap[c.GetID()] = constraints
			}
		}
	}
	return config, nil
}

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
				return constraints, errors.New(fmt.Sprint(cons["value"]) + "cannot be cast to string, for time constraint")
			}
			var rgxPat = regexp.MustCompile(`^[0-9]{2}:[0-9]{2}:[0-9]{2}$`)
			if !rgxPat.MatchString(cons["value"].(string)) {
				return constraints, errors.New(fmt.Sprint(cons["value"]) + " did not match pattern 'hh:mm:ss:'")
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
	componentType := device.Input[constraint["component_id"].(string)]
	switch componentType {
	case "string":
		if reflect.TypeOf(constraint["value"]).Kind() != reflect.String {
			return Constraint{}, "", errors.New("Value was not of type string: " + fmt.Sprint(constraint["value"]))
		}
		return ConstraintString{generalConstraint, constraint["value"].(string)}, "string", nil
	case "numeric":
		if reflect.TypeOf(constraint["value"]).Kind() != reflect.Float64 {
			return Constraint{}, "", errors.New("Value was not of type integer: " + fmt.Sprint(constraint["value"]))
		}
		return ConstraintNumeric{generalConstraint, constraint["value"].(float64)}, "numeric", nil
	case "boolean":
		if reflect.TypeOf(constraint["value"]).Kind() != reflect.Bool {
			return Constraint{}, "", errors.New("Value was not of type boolean: " + fmt.Sprint(constraint["value"]))
		}
		return ConstraintBool{generalConstraint, constraint["value"].(bool)}, "boolean", nil
	case "num-array":
		if reflect.TypeOf(constraint["value"]).Kind() != reflect.Slice {
			return Constraint{}, "", errors.New("Value was not of type array: " + fmt.Sprint(constraint["value"]))
		}
		array := constraint["value"].([]interface{})
		var newArray []float64
		for i := range array {
			if reflect.TypeOf(array[i]).Kind() != reflect.Float64 {
				return Constraint{}, "", errors.New("Value was not of type int-array: " + fmt.Sprint(constraint["value"]))
			}
			newArray = append(newArray, array[i].(float64))
		}
		return ConstraintNumericArray{generalConstraint, newArray}, "num-array", nil
	case "string-array":
		if reflect.TypeOf(constraint["value"]).Kind() != reflect.Slice {
			return Constraint{}, "", errors.New("Value was not of type array: " + fmt.Sprint(constraint["value"]))
		}
		array := constraint["value"].([]interface{})
		var newArray []string
		for i := range array {
			if reflect.TypeOf(array[i]).Kind() != reflect.String {
				return Constraint{}, "", errors.New("Value was not of type string-array: " + fmt.Sprint(constraint["value"]))
			}
			newArray = append(newArray, array[i].(string))
		}
		return ConstraintStringArray{generalConstraint, newArray}, "string-array", nil
	case "bool-array":
		if reflect.TypeOf(constraint["value"]).Kind() != reflect.Slice {
			return Constraint{}, "", errors.New("Value was not of type array: " + fmt.Sprint(constraint["value"]))
		}
		array := constraint["value"].([]interface{})
		var newArray []bool
		for i := range array {
			if reflect.TypeOf(array[i]).Kind() != reflect.Bool {
				return Constraint{}, "", errors.New("Value was not of type bool-array: " + fmt.Sprint(constraint["value"]))
			}
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
func createConstraint(constraint ConstraintInfo, includeCompo bool) (Constraint, error) {
	var err error
	if reflect.TypeOf(constraint["comp"]).Kind() != reflect.String {
		err = errors.New("comparison " + fmt.Sprint(constraint["comp"]) + " cannot be cast to string")
	}
	if includeCompo {
		if reflect.TypeOf(constraint["component_id"]).Kind() != reflect.String {
			err = errors.New("component_id " + fmt.Sprint(constraint["component_id"]) + " cannot be cast to string")
		}
		return Constraint{constraint["comp"].(string), constraint["component_id"].(string)}, err
	}
	return Constraint{constraint["comp"].(string), ""}, err
}

// GetConstraintTimer struct object from map
func GetConstraintTimer(config WorkingConfig, conditionID string) []ConstraintTimer {
	var constraintArray []ConstraintTimer
	if config.ConstraintMap[conditionID] == nil {
		panic("Condition ID: " + fmt.Sprint(conditionID) + " not in constraint map")
	}
	var resultArray []interface{} = config.ConstraintMap[conditionID]["timer"]
	for i := range resultArray {
		constraintArray = append(constraintArray, resultArray[i].(ConstraintTimer))
	}
	return constraintArray
}

// GetConstraintBool struct object from map
func GetConstraintBool(config WorkingConfig, conditionID string) []ConstraintBool {
	var constraintArray []ConstraintBool
	if config.ConstraintMap[conditionID] == nil {
		panic("Condition ID: " + fmt.Sprint(conditionID) + " not in constraint map")
	}
	var resultArray []interface{} = config.ConstraintMap[conditionID]["boolean"]
	for i := range resultArray {
		constraintArray = append(constraintArray, resultArray[i].(ConstraintBool))
	}
	return constraintArray
}

// GetConstraintString struct object from map
func GetConstraintString(config WorkingConfig, conditionID string) []ConstraintString {
	var constraintArray []ConstraintString
	if config.ConstraintMap[conditionID] == nil {
		panic("Condition ID: " + fmt.Sprint(conditionID) + " not in constraint map")
	}
	var resultArray []interface{} = config.ConstraintMap[conditionID]["string"]
	for i := range resultArray {
		constraintArray = append(constraintArray, resultArray[i].(ConstraintString))
	}
	return constraintArray
}

// GetConstraintNumeric struct object from map
func GetConstraintNumeric(config WorkingConfig, conditionID string) []ConstraintNumeric {
	var constraintArray []ConstraintNumeric
	if config.ConstraintMap[conditionID] == nil {
		panic("Condition ID: " + fmt.Sprint(conditionID) + " not in constraint map")
	}
	var resultArray []interface{} = config.ConstraintMap[conditionID]["numeric"]
	for i := range resultArray {
		constraintArray = append(constraintArray, resultArray[i].(ConstraintNumeric))
	}
	return constraintArray
}

// GetConstraintStringArray struct object from map
func GetConstraintStringArray(config WorkingConfig, conditionID string) []ConstraintStringArray {
	var constraintArray []ConstraintStringArray
	if config.ConstraintMap[conditionID] == nil {
		panic("Condition ID: " + fmt.Sprint(conditionID) + " not in constraint map")
	}
	var resultArray []interface{} = config.ConstraintMap[conditionID]["string-array"]
	for i := range resultArray {
		constraintArray = append(constraintArray, resultArray[i].(ConstraintStringArray))
	}
	return constraintArray
}

// GetConstraintBoolArray struct object from map
func GetConstraintBoolArray(config WorkingConfig, conditionID string) []ConstraintBoolArray {
	var constraintArray []ConstraintBoolArray
	if config.ConstraintMap[conditionID] == nil {
		panic("Condition ID: " + fmt.Sprint(conditionID) + " not in constraint map")
	}
	var resultArray []interface{} = config.ConstraintMap[conditionID]["bool-array"]
	for i := range resultArray {
		constraintArray = append(constraintArray, resultArray[i].(ConstraintBoolArray))
	}
	return constraintArray
}

// GetConstraintNumArray struct object from map
func GetConstraintNumArray(config WorkingConfig, conditionID string) []ConstraintNumericArray {
	var constraintArray []ConstraintNumericArray
	if config.ConstraintMap[conditionID] == nil {
		panic("Condition ID: " + fmt.Sprint(conditionID) + " not in constraint map")
	}
	var resultArray []interface{} = config.ConstraintMap[conditionID]["num-array"]
	for i := range resultArray {
		constraintArray = append(constraintArray, resultArray[i].(ConstraintNumericArray))
	}
	return constraintArray
}

// GetConstraintCustomType struct object from map
func GetConstraintCustomType(config WorkingConfig, conditionID string) []ConstraintCustomType {
	var constraintArray []ConstraintCustomType
	if config.ConstraintMap[conditionID] == nil {
		panic("Condition ID: " + fmt.Sprint(conditionID) + " not in constraint map")
	}
	var resultArray []interface{} = config.ConstraintMap[conditionID]["custom"]
	for i := range resultArray {
		constraintArray = append(constraintArray, resultArray[i].(ConstraintCustomType))
	}
	return constraintArray
}

//TODO handling multiple OR/AND - new issue
//TODO front-end - new issue
//TODO multiple errors?
//TODO check device present
//TODO check component present
//TODO check output types
//TODO catch non-existing keys in json
