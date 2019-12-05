package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
)

// ReadFile reads filename and call readJSON on contents.
func ReadFile(filename string) WorkingConfig {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(errors.New("Could not read file " + filename).Error())
	}
	config := ReadJSON(dat)
	return config
}

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
	config.Puzzles = generatePuzzles(readConfig.Puzzles)
	config.GeneralEvents = generateGeneralEvents(readConfig.GeneralEvents)
	config.Devices = make(map[string]Device)
	for _, d := range readConfig.Devices {
		config.Devices[d.ID] = Device{d.ID, d.Description, d.Input,
			d.Output, make(map[string]interface{}), false}
	}

	return config, checkConfig(config)
}

// checkConfig is a method that will return an error if the constraints value type is not equal to the device input type specified, the actions type is not equal to the device output type, or some other not allowed json configuration
func checkConfig(config WorkingConfig) error {
	for _, puzzle := range config.Puzzles {
		for _, rule := range puzzle.Event.Rules {
			err := rule.Conditions.CheckConstraints(config)
			if err != nil {
				return err
			}
		}
	}

	for _, generalEvent := range config.GeneralEvents {
		for _, rule := range generalEvent.Rules {
			err := rule.Conditions.CheckConstraints(config)
			if err != nil {
				return err
			}
		}
	}

	// todo add check for actions
	return nil
}

func generatePuzzles(readPuzzles []ReadPuzzle) []Puzzle {
	var result []Puzzle
	for _, readPuzzle := range readPuzzles {
		puzzle := Puzzle{
			Event: generateGeneralEvent(readPuzzle),
			Hints: readPuzzle.Hints,
		}
		result = append(result, puzzle)
	}
	return result
}

func generateGeneralEvents(readGeneralEvents []ReadGeneralEvent) []GeneralEvent {
	var result []GeneralEvent
	for _, readGeneralEvent := range readGeneralEvents {
		result = append(result, generateGeneralEvent(readGeneralEvent))
	}
	return result
}

func generateGeneralEvent(event ReadEvent) GeneralEvent {
	return GeneralEvent{
		Name:  event.GetName(),
		Rules: generateRules(event.GetRules()),
	}
}

func generateRules(readRules []ReadRule) []Rule {
	var rules []Rule
	for _, readRule := range readRules {
		rule := Rule{
			ID:          readRule.ID,
			Description: readRule.Description,
			Limit:       readRule.Limit,
			Conditions:  nil,
			Actions:     readRule.Actions,
		}
		rule.Conditions = generateLogicalCondition(readRule.Conditions)
		rules = append(rules, rule)
	}
	return rules
}

func generateLogicalCondition(conditions interface{}) LogicalCondition { // todo check types
	logic := conditions.(map[string]interface{})
	if logic["operator"] != nil { // operator
		if logic["operator"] == "AND" {
			and := AndCondition{}
			for _, condition := range logic["list"].([]interface{}) {
				and.logics = append(and.logics, generateLogicalCondition(condition))
			}
			return and
		} else if logic["operator"] == "OR" {
			or := OrCondition{}
			for _, condition := range logic["list"].([]interface{}) {
				or.logics = append(or.logics, generateLogicalCondition(condition))
			}
			return or
		} else {
			panic(fmt.Sprintf("JSON config in wrong format, operator: %v, could not be processed", logic["operator"]))
		}
	} else if logic["type"] != nil && reflect.TypeOf(logic["type"]).Kind() == reflect.String && logic["type_id"] != nil && reflect.TypeOf(logic["type_id"]).Kind() == reflect.String {
		condition := Condition{
			Type:        logic["type"].(string),
			TypeID:      logic["type_id"].(string),
			Constraints: generateLogicalConstraint(logic["constraints"]),
		}

		return condition
	}
	panic(fmt.Sprintf("JSON config in wrong condition format, conditions: %v, could not be processed", conditions))
}

func generateLogicalConstraint(constraints interface{}) LogicalConstraint {
	logic := constraints.(map[string]interface{})
	if logic["operator"] != nil { // operator
		if logic["operator"] == "AND" {
			and := AndConstraint{}
			for _, constraint := range logic["list"].([]interface{}) {
				and.logics = append(and.logics, generateLogicalConstraint(constraint))
			}
			return and
		} else if logic["operator"] == "OR" {
			or := OrConstraint{}
			for _, constraint := range logic["list"].([]interface{}) {
				or.logics = append(or.logics, generateLogicalConstraint(constraint))
			}
			return or
		} else {
			panic(fmt.Sprintf("JSON config in wrong format, operator: %v, could not be processed", logic["operator"]))
		}
	} else if logic["comp"] != nil && reflect.TypeOf(logic["comp"]).Kind() == reflect.String {
		var constraint Constraint
		if logic["component_id"] != nil && reflect.TypeOf(logic["component_id"]).Kind() == reflect.String {
			constraint = Constraint{
				Comparison:  logic["comp"].(string),
				ComponentID: logic["component_id"].(string),
				Value:       logic["value"],
			}
		} else if logic["component_id"] == nil {
			constraint = Constraint{
				Comparison:  logic["comp"].(string),
				ComponentID: "",
				Value:       logic["value"],
			}
		} else {
			panic(fmt.Sprintf("JSON config in wrong format, component_id should be of type string, %v is of type %s", logic["component_id"], reflect.TypeOf(logic["component_id"]).Kind().String()))
		}

		return constraint
	}
	panic(fmt.Sprintf("JSON config in wrong constraint format, conditions: %v, could not be processed", constraints))
}

//func checkActions(devices map[string]Device, actions []Action) ([]Action, error) {
//	for _, a := range actions {
//		output := a.Message.Output
//		if a.Type == "timer" {
//			instruction, ok := output["instruction"]
//			if !ok {
//				return actions, errors.New("timer should have an instruction defined")
//			}
//			if instruction != "stop" && instruction != "subtract" {
//				return actions, errors.New("timer should have an instruction defined, which is either stop or subtract")
//			}
//			if instruction == "subtract" {
//				value, ok2 := output["value"]
//				if !ok2 {
//					return actions, errors.New("timer with subtract instruction should have value")
//				}
//				if reflect.TypeOf(value).Kind() != reflect.String {
//					return actions, errors.New("timer with subtract instruction should have value in string format")
//				}
//				var rgxPat = regexp.MustCompile(`^[0-9]{2}:[0-9]{2}:[0-9]{2}$`)
//				if !rgxPat.MatchString(value.(string)) {
//					return actions, errors.New(value.(string) + " did not match pattern 'hh:mm:ss'")
//				}
//			}
//		} else if a.Type == "device" {
//			device, ok := devices[a.TypeID]
//			if !ok {
//				return actions, errors.New("device with id " + a.TypeID + " not found in map")
//			}
//			for key, value := range output {
//				expectedType, ok2 := device.Output[key]
//				if !ok2 {
//					return actions, errors.New("component id: " + key + " not found in device input")
//				}
//				// value should be of type specified in device output
//				err := CheckComponentType(expectedType, value)
//				if err != nil {
//					return actions, err
//				}
//			}
//		} else {
//			return actions, errors.New("invalid type of action: " + a.Type)
//		}
//	}
//	return actions, nil
//}

//// CheckComponentType checks if value is of componentType and returns error if not.
//func CheckComponentType(componentType interface{}, value interface{}) error {
//	switch componentType {
//	case "string":
//		if reflect.TypeOf(value).Kind() != reflect.String {
//			return errors.New("Value was not of type string: " + fmt.Sprint(value))
//		}
//	case "numeric":
//		if reflect.TypeOf(value).Kind() != reflect.Float64 {
//			return errors.New("Value was not of type numeric: " + fmt.Sprint(value))
//		}
//	case "boolean":
//		if reflect.TypeOf(value).Kind() != reflect.Bool {
//			return errors.New("Value was not of type boolean: " + fmt.Sprint(value))
//		}
//	case "num-array":
//		if reflect.TypeOf(value).Kind() != reflect.Slice {
//			return errors.New("Value was not of type array: " + fmt.Sprint(value))
//		}
//		array := value.([]interface{})
//		for i := range array {
//			if reflect.TypeOf(array[i]).Kind() != reflect.Float64 {
//				return errors.New("Value was not of type num-array: " + fmt.Sprint(value))
//			}
//		}
//	case "string-array":
//		if reflect.TypeOf(value).Kind() != reflect.Slice {
//			return errors.New("Value was not of type array: " + fmt.Sprint(value))
//		}
//		array := value.([]interface{})
//		for i := range array {
//			if reflect.TypeOf(array[i]).Kind() != reflect.String {
//				return errors.New("Value was not of type string-array: " + fmt.Sprint(value))
//			}
//		}
//	case "bool-array":
//		if reflect.TypeOf(value).Kind() != reflect.Slice {
//			return errors.New("Value was not of type array: " + fmt.Sprint(value))
//		}
//		array := value.([]interface{})
//		for i := range array {
//			if reflect.TypeOf(array[i]).Kind() != reflect.Bool {
//				return errors.New("Value was not of type bool-array: " + fmt.Sprint(value))
//			}
//		}
//	default:
//		// Value is of custom type, must be checked at client computer level
//		return nil
//	}
//	return nil
//}

//TODO handling multiple OR/AND - new issue
//TODO front-end - new issue
//TODO multiple errors?
//TODO check device present
//TODO check component present
//TODO check output types
//TODO catch non-existing keys in json
