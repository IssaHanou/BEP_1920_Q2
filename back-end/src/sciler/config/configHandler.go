package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"time"
)

// ReadFile reads filename and call readJSON on contents.
func ReadFile(filename string) WorkingConfig {
	dat, err := ioutil.ReadFile(filename)
	errorList := make([]string, 0)
	if err != nil {
		errorList = append(errorList, errors.New("Could not read file "+filename).Error())
		panic(errorList[0])
	}
	config, jsonErrors := ReadJSON(dat)
	errorList = append(errorList, jsonErrors...)
	if len(errorList) > 0 {
		panic(errorList[0])
	}
	return config
}

// ReadJSON transforms json file into config object.
// The string array contains all errors that are found during the transformation.
func ReadJSON(input []byte) (WorkingConfig, []string) {
	var config ReadConfig
	jsonErr := json.Unmarshal(input, &config)
	if jsonErr != nil {
		return WorkingConfig{}, []string{jsonErr.Error()}
	}
	newConfig, configErr := generateDataStructures(config)
	return newConfig, configErr
}

// generateDataStructures creates the working config, with maps to easily lookup objects.
func generateDataStructures(readConfig ReadConfig) (WorkingConfig, []string) {
	var config WorkingConfig
	errorList := make([]string, 0)
	// Copy information from read config to working config.
	config.General = readConfig.General
	config.Cameras = readConfig.Cameras
	newPuzzles, puzzleErrors := generatePuzzles(readConfig.Puzzles, &config)
	config.Puzzles = newPuzzles
	newEvents, eventErrors := generateGeneralEvents(readConfig.GeneralEvents, &config)
	config.GeneralEvents = newEvents
	newButtonEvents, buttonEventErrors := generateButtonEvents(readConfig.ButtonEvents, &config)
	config.ButtonEvents = newButtonEvents

	generateDevices(readConfig.Devices, &config)
	timerErrors := generateTimers(readConfig.Timers, &config)
	errorList = append(errorList, append(buttonEventErrors, append(puzzleErrors, append(eventErrors, timerErrors...)...)...)...)

	if len(errorList) == 0 {
		// if there are errors in config format,
		// wait with creating maps (which use condition type ids)
		// and with checking constraint
		config.StatusMap = generateStatusMap(&config)
		config.RuleMap = generateRuleMap(&config)
		config.LabelMap = generateLabelMap(&config)
		errorList = append(errorList, checkConfig(config)...)
	}
	return config, errorList
}

// generateDevices creates the config devices map which points device id to a device in the WorkingConfig.
// Creates front end device manually as its information is not in `devices` in the configuration file.
// The components are defined as the custom buttons, with boolean status of clicked or not.
func generateDevices(devices []ReadDevice, config *WorkingConfig) {
	config.Devices = make(map[string]*Device)
	for _, readDevice := range devices {
		config.Devices[readDevice.ID] = &(Device{
			readDevice.ID,
			readDevice.Description,
			readDevice.Input,
			readDevice.Output,
			make(map[string]interface{}),
			false,
		})
	}

	input := make(map[string]string)
	status := make(map[string]interface{})
	for _, btn := range config.ButtonEvents {
		input[btn.ID] = "boolean"
		status[btn.ID] = false
	}
	config.Devices["front-end"] = &(Device{
		ID:          "front-end",
		Description: "The operator webapp for managing a escape room",
		Input:       input,
		Output: map[string]OutputObject{
			"gameState": {
				Type: "string",
				Instructions: map[string]string{
					"setState": "string",
				},
			},
		},
		Status:     status,
		Connection: false,
	})
}

// generateTimers creates map with id pointing to timer object for all timer objects and general timer.
// check that all durations are of proper format.
// return error list
func generateTimers(timers []ReadTimer, config *WorkingConfig) []string {
	errorList := make([]string, 0)
	config.Timers = make(map[string]*Timer)
	for _, readTimer := range timers {
		duration, err := time.ParseDuration(readTimer.Duration)
		if err != nil {
			errorList = append(errorList, err.Error())
		} else {
			config.Timers[readTimer.ID] = newTimer(readTimer.ID, duration)
		}
	}
	duration, err := time.ParseDuration(config.General.Duration)
	if err != nil {
		errorList = append(errorList, err.Error())
	} else {
		config.Timers["general"] = newTimer("general", duration)
	}
	return errorList
}

// getAllRules creates rule list of the rule pointers belonging to all events, except button events,
// because those should not be added to status map, only to rule map
func getAllRules(config *WorkingConfig) []*Rule {
	var rules []*Rule
	for _, event := range config.GeneralEvents {
		rules = append(rules, event.GetRules()...)
	}

	for _, event := range config.Puzzles {
		rules = append(rules, event.GetRules()...)
	}

	return rules
}

// generateRuleMap creates rule map with rule id pointing to rule object pointers for all rules of all events
func generateRuleMap(config *WorkingConfig) map[string]*Rule {
	ruleMap := make(map[string]*Rule)
	rules := getAllRules(config)

	for _, event := range config.ButtonEvents {
		rules = append(rules, event)
	}

	for _, rule := range rules {
		ruleMap[rule.ID] = rule
	}
	return ruleMap
}

// generateLabelMap makes a map from a label to a component by checking all components if they have labels
func generateLabelMap(config *WorkingConfig) map[string][]*Component {
	labelMap := make(map[string][]*Component)
	devices := config.Devices
	for _, device := range devices {
		for compID, comp := range device.Output {
			for _, label := range comp.Label {
				component := &Component{Device: device, ID: compID}
				labelMap[label] = appendWhenUniqueComp(labelMap[label], component)
			}
		}
	}
	return labelMap
}

// generateStatusMap creates map from id of rule/timer/device to all the rules
// which have condition based on the rule/timer/device object
func generateStatusMap(config *WorkingConfig) map[string][]*Rule {
	statusMap := make(map[string][]*Rule)
	rules := getAllRules(config)

	for _, rule := range rules {
		for _, id := range rule.Conditions.GetConditionIDs() {
			statusMap[id] = appendWhenUnique(statusMap[id], rule)
		}
	}

	return statusMap
}

// appendWhenUnique is a helper method to only append the rule to the list if it is unique
// todo make this more efficient
func appendWhenUnique(rules []*Rule, rule *Rule) []*Rule {
	for _, existingRule := range rules {
		if reflect.DeepEqual(*existingRule, *rule) {
			return rules
		}
	}
	return append(rules, rule)
}

// appendWhenUnique is a helper method to only append the component to the list if it is unique
func appendWhenUniqueComp(comps []*Component, comp *Component) []*Component {
	for _, existingComp := range comps {
		if reflect.DeepEqual(*existingComp, *comp) {
			return comps
		}
	}
	return append(comps, comp)
}

// checkConfig is a method that will return an error
// if the constraints value type is not equal to the device input type specified,
// the actions type is not equal to the device output type,
// or some other not allowed json configuration
func checkConfig(config WorkingConfig) []string {
	errList := make([]string, 0)
	for _, puzzle := range config.Puzzles {
		for _, rule := range puzzle.Event.Rules {
			if err := rule.Conditions.checkConstraints(config, rule.ID); err != nil {
				errList = append(errList, err...)
			}
			errList = append(errList, checkActions(rule.Actions, config)...)
		}
	}

	for _, generalEvent := range config.GeneralEvents {
		for _, rule := range generalEvent.Rules {
			if err := rule.Conditions.checkConstraints(config, rule.ID); err != nil {
				errList = append(errList, err...)
			}
			errList = append(errList, checkActions(rule.Actions, config)...)
		}
	}

	for _, rule := range config.ButtonEvents {
		if err := rule.Conditions.checkConstraints(config, rule.ID); err != nil {
			errList = append(errList, err...)
		}
		errList = append(errList, checkActions(rule.Actions, config)...)
	}

	// todo check uniqueness of all device_id, timer_id and rule_id
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
				fmt.Sprintf("only device, timer and label are accepted as type for an action, however type was specified as: %s", action.Type))
		}
	}
	return errorList
}

// checkActionTimer is a method that will return an error
// if the action's value types and instructions are not equal to the possibilities of timer's output
func checkActionTimer(action Action, config WorkingConfig) []string {
	errorList := make([]string, 0)
	if _, ok := config.Timers[action.TypeID]; ok { // checks if timer can be found in the map, if so, it is stored in variable device
		for _, actionMessage := range action.Message {
			if actionMessage.Instruction == "add" || actionMessage.Instruction == "subtract" {
				valueType := reflect.TypeOf(actionMessage.Value).Kind()
				if valueType != reflect.String {
					errorList = append(errorList,
						fmt.Sprintf("input type string expected but %s found as type of value %v",
							valueType.String(), actionMessage.Value))
				}
			}
		}
	} else {
		errorList = append(errorList, fmt.Sprintf("timer with id %s not found in map", action.TypeID))
	}
	return errorList
}

// checkAction is a method that will return an error
// if the action's value types and instructions are not equal to the device's output specifications
func checkActionDevice(action Action, config WorkingConfig) []string {
	errorList := make([]string, 0)
	if device, ok := config.Devices[action.TypeID]; ok { // checks if device can be found in the map, if so, it is stored in variable device
		for _, actionMessage := range action.Message {
			if outputObject, ok := device.Output[actionMessage.ComponentID]; ok {
				valueType := reflect.TypeOf(actionMessage.Value).Kind()
				if instructionType, ok := outputObject.Instructions[actionMessage.Instruction]; ok {
					switch instructionType {
					case "string":
						{
							if valueType != reflect.String {
								errorList = append(errorList,
									fmt.Sprintf("instruction type string expected but %s found as type of value %v",
										valueType.String(), actionMessage.Value))
							}
						}
					case "boolean":
						{
							if valueType != reflect.Bool {
								errorList = append(errorList,
									fmt.Sprintf("instruction type boolean expected but %s found as type of value %v",
										valueType.String(), actionMessage.Value))
							}
						}
					case "numeric":
						{
							if valueType != reflect.Int && valueType != reflect.Float64 {
								errorList = append(errorList,
									fmt.Sprintf("instruction type numeric expected but %s found as type of value %v",
										valueType.String(), actionMessage.Value))
							}
						}
					case "array":
						{
							if valueType != reflect.Slice {
								errorList = append(errorList,
									fmt.Sprintf("instruction type array/slice expected but %s found as type of value %v",
										valueType.String(), actionMessage.Value))
							}
						}
					default:
						// todo custom types
						errorList = append(errorList,
							fmt.Sprintf("custom types like: %s, are not yet implemented", instructionType))
					}
				} else {
					errorList = append(errorList,
						fmt.Sprintf("instruction %s not found in map", actionMessage.Instruction))
				}
			} else {
				errorList = append(errorList,
					fmt.Sprintf("component with id %s not found in map", actionMessage.ComponentID))
			}
		}
	} else {
		errorList = append(errorList,
			fmt.Sprintf("device with id %s not found in map", action.TypeID))
	}
	return errorList
}

// checkActionLabel checks if there is a label with this ID,
// and checks if all components under this label have the correct instructions with a call to checkActionDevice
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
			fmt.Sprintf("label with id %s not found in map", action.TypeID))
	}
	return errorList
}

// generatePuzzles creates a list puzzle objects with properly checked inner values (up to constraints)
func generatePuzzles(readPuzzles []ReadPuzzle, config *WorkingConfig) ([]*Puzzle, []string) {
	var result []*Puzzle
	errorList := make([]string, 0)
	for _, readPuzzle := range readPuzzles {
		event, newErrors := generateGeneralEvent(readPuzzle, config)
		puzzle := Puzzle{
			Event: event,
			Hints: readPuzzle.Hints,
		}
		result = append(result, &puzzle)
		errorList = append(errorList, newErrors...)
	}
	return result, errorList
}

// generateGeneralEvents creates a list general events objects with properly checked inner values (up to constraints)
func generateGeneralEvents(readGeneralEvents []ReadGeneralEvent, config *WorkingConfig) ([]*GeneralEvent, []string) {
	var result []*GeneralEvent
	errorList := make([]string, 0)
	for _, readGeneralEvent := range readGeneralEvents {
		newResult, newErrors := generateGeneralEvent(readGeneralEvent, config)
		result = append(result, newResult)
		errorList = append(errorList, newErrors...)
	}
	return result, errorList
}

// generateGeneralEvent creates a the rules for a general events object
func generateGeneralEvent(event ReadEvent, config *WorkingConfig) (*GeneralEvent, []string) {
	rules, errorList := generateRules(event.GetRules(), config)
	return &GeneralEvent{
		Name:  event.GetName(),
		Rules: rules,
	}, errorList
}

// generateButtonEvents creates a list button event objects with properly checked inner values
func generateButtonEvents(buttonEvents []ReadRule, config *WorkingConfig) (map[string]*Rule, []string) {
	newEvents := make(map[string]*Rule)
	rules, errorList := generateRules(buttonEvents, config)
	for _, rule := range rules {
		newEvents[rule.ID] = rule
	}
	return newEvents, errorList
}

// generateRules creates a list rule objects with properly checked inner values
func generateRules(readRules []ReadRule, config *WorkingConfig) ([]*Rule, []string) {
	var rules []*Rule
	errorList := make([]string, 0)
	for _, readRule := range readRules {
		rule := Rule{
			ID:          readRule.ID,
			Description: readRule.Description,
			Limit:       readRule.Limit,
			Executed:    0,
			Conditions:  nil,
			Actions:     readRule.Actions,
		}
		conditions, newErrors := generateLogicalCondition(readRule.Conditions)
		rule.Conditions = conditions
		errorList = append(errorList, newErrors...)
		rules = append(rules, &rule)
	}
	return rules, errorList
}

// generateLogicalConditions creates the LogicalCondition tree of all conditions, which are type checked
func generateLogicalCondition(conditions interface{}) (LogicalCondition, []string) {
	logic := conditions.(map[string]interface{})
	errorList := make([]string, 0)
	if logic["operator"] != nil { // operator
		if logic["operator"] == "AND" {
			and := AndCondition{}
			for _, condition := range logic["list"].([]interface{}) {
				newCondition, newErrors := generateLogicalCondition(condition)
				and.logics = append(and.logics, newCondition)
				errorList = append(errorList, newErrors...)
			}
			return and, errorList
		} else if logic["operator"] == "OR" {
			or := OrCondition{}
			for _, condition := range logic["list"].([]interface{}) {
				newCondition, newErrors := generateLogicalCondition(condition)
				or.logics = append(or.logics, newCondition)
				errorList = append(errorList, newErrors...)
			}
			return or, errorList
		} else {
			return nil, append(errorList,
				fmt.Sprintf("JSON config in wrong format, operator: %v, could not be processed", logic["operator"]))
		}
	} else if logic["type"] != nil && reflect.TypeOf(logic["type"]).Kind() == reflect.String && logic["type_id"] != nil && reflect.TypeOf(logic["type_id"]).Kind() == reflect.String {
		constraints, newErrors := generateLogicalConstraint(logic["constraints"])
		condition := Condition{
			Type:        logic["type"].(string),
			TypeID:      logic["type_id"].(string),
			Constraints: constraints,
		}
		return condition, append(errorList, newErrors...)
	}
	return nil, append(errorList,
		fmt.Sprintf("JSON config in wrong condition format, conditions: %v, could not be processed", conditions))
}

// generateLogicalConstraint creates the LogicalContraint tree of all constraints, which are type checked
func generateLogicalConstraint(constraints interface{}) (LogicalConstraint, []string) {
	logic := constraints.(map[string]interface{})
	errorList := make([]string, 0)
	if logic["operator"] != nil { // operator
		if logic["operator"] == "AND" {
			and := AndConstraint{}
			for _, constraint := range logic["list"].([]interface{}) {
				newConstraint, newErrors := generateLogicalConstraint(constraint)
				and.logics = append(and.logics, newConstraint)
				errorList = append(errorList, newErrors...)
			}
			return and, errorList
		} else if logic["operator"] == "OR" {
			or := OrConstraint{}
			for _, constraint := range logic["list"].([]interface{}) {
				newConstraint, newErrors := generateLogicalConstraint(constraint)
				or.logics = append(or.logics, newConstraint)
				errorList = append(errorList, newErrors...)
			}
			return or, errorList
		} else {
			return nil, append(errorList,
				fmt.Sprintf("JSON config in wrong format, operator: %v, could not be processed", logic["operator"]))
		}
	} else if logic["comparison"] != nil && reflect.TypeOf(logic["comparison"]).Kind() == reflect.String {
		var constraint Constraint
		if logic["component_id"] != nil && reflect.TypeOf(logic["component_id"]).Kind() == reflect.String {
			constraint = Constraint{
				Comparison:  logic["comparison"].(string),
				ComponentID: logic["component_id"].(string),
				Value:       logic["value"],
			}
		} else if logic["component_id"] == nil {
			constraint = Constraint{
				Comparison:  logic["comparison"].(string),
				ComponentID: "",
				Value:       logic["value"],
			}
		} else {
			errorList = append(errorList,
				fmt.Sprintf("JSON config in wrong format, component_id should be of type string, %v is of type %s",
					logic["component_id"], reflect.TypeOf(logic["component_id"]).Kind().String()))
		}

		return constraint, errorList
	}
	return nil, append(errorList,
		fmt.Sprintf("JSON config in wrong constraint format, conditions: %v, could not be processed", constraints))
}
