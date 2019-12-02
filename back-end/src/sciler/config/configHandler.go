package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
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
	config.Devices = make(map[string]Device)
	for _, d := range readConfig.Devices {
		config.Devices[d.ID] = d
	}

	// Create additional data structures.
	config.Puzzles = []Puzzle{}
	config.GeneralEvents = []Event{}
	config.Rules = make(map[string]Rule)
	config.ActionMap = make(map[string][]Action)
	config.ConstraintMap = make(map[string]map[string][]interface{})

	var eventList []*Event
	for _, p := range readConfig.Puzzles {
		newEvent := Event{p.Name, p.Rules}
		newPuzzle := Puzzle{newEvent, p.Hints}
		config.Puzzles = append(config.Puzzles, newPuzzle)
		eventList = append(eventList, &newEvent)
	}
	for _, e := range readConfig.GeneralEvents {
		newEvent := Event{e.Name, e.Rules}
		config.GeneralEvents = append(config.GeneralEvents, newEvent)
		eventList = append(eventList, &newEvent)
	}
	for _, e := range eventList {
		for j, r := range (*e).Rules {
			config.Rules[r.ID] = r
			actions, err := checkActions(config.Devices, r.Actions)
			if err != nil {
				return config, err
			}
			for k, c := range r.Conditions {
				// Set both original pointer's and current pointer's conditions
				(*e).Rules[j].Conditions[k].RuleID = r.ID
				//config.Puzzles[i].Rules[j].Conditions[k].RuleID = r.ID
				c.RuleID = r.ID
				config.ActionMap[c.GetID()] = actions
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
			actions, err := checkActions(config.Devices, r.Actions)
			if err != nil {
				return config, err
			}
			for k, c := range r.Conditions {
				// Set both original pointer's and current pointer's conditions
				config.GeneralEvents[i].Rules[j].Conditions[k].RuleID = r.ID
				c.RuleID = r.ID
				config.ActionMap[c.GetID()] = actions
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

func checkActions(devices map[string]Device, actions []Action) ([]Action, error) {
	for _, a := range actions {
		output := a.Message.Output
		if a.Type == "timer" {
			instruction, ok := output["instruction"]
			if !ok {
				return actions, errors.New("timer should have an instruction defined")
			}
			if instruction != "stop" && instruction != "subtract" {
				return actions, errors.New("timer should have an instruction defined, which is either stop or subtract")
			}
			if instruction == "subtract" {
				value, ok2 := output["value"]
				if !ok2 {
					return actions, errors.New("timer with subtract instruction should have value")
				}
				if reflect.TypeOf(value).Kind() != reflect.String {
					return actions, errors.New("timer with subtract instruction should have value in string format")
				}
				var rgxPat = regexp.MustCompile(`^[0-9]{2}:[0-9]{2}:[0-9]{2}$`)
				if !rgxPat.MatchString(value.(string)) {
					return actions, errors.New(value.(string) + " did not match pattern 'hh:mm:ss'")
				}
			}
		} else if a.Type == "device" {
			device, ok := devices[a.TypeID]
			if !ok {
				return actions, errors.New("device with id " + a.TypeID + " not found in map")
			}
			for key, value := range output {
				expectedType, ok2 := device.Output[key]
				if !ok2 {
					return actions, errors.New("component id: " + key + " not found in device input")
				}
				// value should be of type specified in device output
				err := CheckComponentType(expectedType, value)
				if err != nil {
					return actions, err
				}
			}
		} else {
			return actions, errors.New("invalid type of action: " + a.Type)
		}
	}
	return actions, nil
}

// CheckComponentType checks if value is of componentType and returns error if not.
func CheckComponentType(componentType interface{}, value interface{}) error {
	switch componentType {
	case "string":
		if reflect.TypeOf(value).Kind() != reflect.String {
			return errors.New("Value was not of type string: " + fmt.Sprint(value))
		}
	case "numeric":
		if reflect.TypeOf(value).Kind() != reflect.Float64 {
			return errors.New("Value was not of type numeric: " + fmt.Sprint(value))
		}
	case "boolean":
		if reflect.TypeOf(value).Kind() != reflect.Bool {
			return errors.New("Value was not of type boolean: " + fmt.Sprint(value))
		}
	case "num-array":
		if reflect.TypeOf(value).Kind() != reflect.Slice {
			return errors.New("Value was not of type array: " + fmt.Sprint(value))
		}
		array := value.([]interface{})
		for i := range array {
			if reflect.TypeOf(array[i]).Kind() != reflect.Float64 {
				return errors.New("Value was not of type num-array: " + fmt.Sprint(value))
			}
		}
	case "string-array":
		if reflect.TypeOf(value).Kind() != reflect.Slice {
			return errors.New("Value was not of type array: " + fmt.Sprint(value))
		}
		array := value.([]interface{})
		for i := range array {
			if reflect.TypeOf(array[i]).Kind() != reflect.String {
				return errors.New("Value was not of type string-array: " + fmt.Sprint(value))
			}
		}
	case "bool-array":
		if reflect.TypeOf(value).Kind() != reflect.Slice {
			return errors.New("Value was not of type array: " + fmt.Sprint(value))
		}
		array := value.([]interface{})
		for i := range array {
			if reflect.TypeOf(array[i]).Kind() != reflect.Bool {
				return errors.New("Value was not of type bool-array: " + fmt.Sprint(value))
			}
		}
	default:
		// Value is of custom type, must be checked at client computer level
		return nil
	}
	return nil
}

//TODO handling multiple OR/AND - new issue
//TODO front-end - new issue
//TODO multiple errors?
//TODO check device present
//TODO check component present
//TODO check output types
//TODO catch non-existing keys in json
