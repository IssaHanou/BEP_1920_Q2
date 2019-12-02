package config

import (
	"encoding/json"
	"errors"
	"fmt"
)

// CheckJSON checks whether all necessary keys exist.
// It gives errors if necessary keys are missing and warnings if optional keys are missing.
func CheckJSON(input []byte) []error {
	var errorList []error
	jsonMap := make(map[string]map[string]interface{})
	jsonError := json.Unmarshal(input, &jsonMap)
	if jsonError != nil {
		switch jsonError.(type) {
		case *json.SyntaxError:
			errorList = append(errorList, errors.New("error: incorrect input syntax at byte "+
				fmt.Sprint(jsonError.(*json.SyntaxError).Offset)+": "+jsonError.Error()))
		case *json.InvalidUnmarshalError:
			errorList = append(errorList, errors.New("error: invalid type, must be a non-nil pointer"))
		case *json.UnmarshalTypeError:
			errorList = append(errorList, errors.New("error: could not construct GO type of input at byte "+
				fmt.Sprint(jsonError.(*json.UnmarshalTypeError).Offset)))
		default:
			errorList = append(errorList, errors.New("error: error decoding json: "+jsonError.Error()))
		}
		//errorList = append(errorList, errors.New("error: first fix syntax errors, only then can keys and types be checked"))
	} else {
		//mainKeys := []string{"general", "devices", "puzzles", "general_events"}
		//errorList = append(errorList, checkMain(jsonMap, mainKeys)...)
		//generalKeys := []string{"name", "duration", "host", "port"}
		//errorList = append(errorList, checkGeneral(jsonMap["general"], generalKeys)...)
		//deviceKeys := []string{"id", "description", "input", "output"}
		//errorList = append(errorList, checkDevice(jsonMap["devices"], deviceKeys)...)
		//eventKeys := []string{"name", "rules"}
		//errorList = append(errorList, checkEvent(jsonMap["puzzles"], eventKeys, true)...)
		//errorList = append(errorList, checkEvent(jsonMap["general_events"], eventKeys, false)...)

		//ruleKeys := []string{"id", "description", "limit", "conditions", "actions"}
		//conditionKeys := []string{"type", "type_id", "constraints"}
		//constraintKeys := []string{"comparison", "component_id", "value"}
		//actions := []string{"type", "type_id", "message"}
	}
	return errorList
}

func checkMain(jsonMap map[string]map[string]interface{}, keys []string) []error {
	var errorList []error
	for _, key := range keys {
		_, ok := jsonMap[key]
		if !ok {
			if key == "general" {
				errorList = append(errorList, errors.New("error: no key 'general' in configuration"))
			} else {
				errorList = append(errorList, errors.New("warning: no key '"+key+"' in configuration"))
			}
		}
	}
	return errorList
}

func checkGeneral(jsonMap map[string]interface{}, keys []string) []error {
	var errorList []error
	for _, key := range keys {
		_, ok := jsonMap[key]
		if !ok {
			errorList = append(errorList, errors.New("error: no key '"+key+"' in general configuration"))
		}
	}
	return errorList
}

func checkDevice(jsonMap map[string]interface{}, keys []string) []error {
	var errorList []error
	for _, key := range keys {
		_, ok := jsonMap[key]
		if !ok {
			if key == "description" {
				errorList = append(errorList, errors.New("warning: no key 'description' in device configuration"))
			} else {
				errorList = append(errorList, errors.New("error: no key '"+key+"' in device configuration"))
			}
		}
	}
	return errorList
}

func checkEvent(jsonMap map[string]interface{}, keys []string, hints bool) []error {
	var errorList []error
	if hints {
		keys = append(keys, "hints")
	}
	for _, key := range keys {
		_, ok := jsonMap[key]
		if !ok {
			if hints && key == "hints" {
				errorList = append(errorList, errors.New("warning: no key 'hints' in puzzle configuration"))
			} else {
				var config string
				if hints {
					config = "puzzle"
				} else {
					config = "general events"
				}
				errorList = append(errorList, errors.New("error: no key '"+key+"' in "+config+" configuration"))
			}
		}
	}
	return errorList
}
