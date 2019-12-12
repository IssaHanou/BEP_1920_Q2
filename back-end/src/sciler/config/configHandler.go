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
	config.Devices = make(map[string]*Device)
	for _, readDevice := range readConfig.Devices {
		config.Devices[readDevice.ID] = &(Device{readDevice.ID, readDevice.Description, readDevice.Input,
			readDevice.Output, make(map[string]interface{}), false})
	}
	config.StatusMap = generateStatusMap(&config)
	config.RuleMap = generateRuleMap(&config)
	return config, checkConfig(config)
}

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

func generateRuleMap(config *WorkingConfig) map[string]*Rule {
	ruleMap := make(map[string]*Rule)
	rules := getAllRules(config)

	for _, rule := range rules {
		ruleMap[rule.ID] = rule
	}
	return ruleMap
}

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

// todo make this more efficient
func appendWhenUnique(rules []*Rule, rule *Rule) []*Rule {
	for _, existingRule := range rules {
		if reflect.DeepEqual(*existingRule, *rule) {
			return rules
		}
	}
	return append(rules, rule)
}

// checkConfig is a method that will return an error if the constraints value type is not equal to the device input type specified, the actions type is not equal to the device output type, or some other not allowed json configuration
func checkConfig(config WorkingConfig) error {
	for _, puzzle := range config.Puzzles {
		for _, rule := range puzzle.Event.Rules {
			if err := rule.Conditions.checkConstraints(config); err != nil {
				return err
			}
			if err := checkActions(rule.Actions, config); err != nil {
				return err
			}
		}
	}

	for _, generalEvent := range config.GeneralEvents {
		for _, rule := range generalEvent.Rules {
			if err := rule.Conditions.checkConstraints(config); err != nil {
				return err
			}
			if err := checkActions(rule.Actions, config); err != nil {
				return err
			}
		}
	}
	// todo check uniqueness of all device_id, timer_id and rule_id
	return nil
}

// checkAction is a method that will return an error is the actions value types and instructions are not equal to the device output specifications
func checkActions(actions []Action, config WorkingConfig) error {
	for _, action := range actions {
		switch action.Type {
		case "device":
			{
				if err := checkActionDevice(action, config); err != nil {
					return err
				}
			}
		case "timer":
		default:
			return fmt.Errorf("only device and timer are accepted as type for an action, however type was specified as: %s", action.Type)
		}
	}
	return nil
}

func checkActionDevice(action Action, config WorkingConfig) error {
	if device, ok := config.Devices[action.TypeID]; ok { // checks if device can be found in the map, if so, it is stored in variable device
		for _, actionMessage := range action.Message {
			if outputObject, ok := device.Output[actionMessage.ComponentID]; ok {
				valueType := reflect.TypeOf(actionMessage.Value).Kind()
				if instructionType, ok := outputObject.Instructions[actionMessage.Instruction]; ok {
					switch instructionType {
					case "string":
						{
							if valueType != reflect.String {
								return fmt.Errorf("instruction type string expected but %s found as type of value %v", valueType.String(), actionMessage.Value)
							}
						}
					case "boolean":
						{
							if valueType != reflect.Bool {
								return fmt.Errorf("instruction type boolean expected but %s found as type of value %v", valueType.String(), actionMessage.Value)
							}
						}
					case "numeric":
						{
							if valueType != reflect.Int && valueType != reflect.Float64 {
								return fmt.Errorf("instruction type numeric expected but %s found as type of value %v", valueType.String(), actionMessage.Value)
							}
						}
					case "array":
						{
							if valueType != reflect.Slice {
								return fmt.Errorf("instruction type array/slice expected but %s found as type of value %v", valueType.String(), actionMessage.Value)
							}
						}
					default:
						// todo custom types
						return fmt.Errorf("custom types like: %s, are not yet implemented", instructionType)
					}
				} else {
					return fmt.Errorf("instruction %s not found in map", actionMessage.Instruction)
				}
			} else {
				return fmt.Errorf("component with id %s not found in map", actionMessage.ComponentID)
			}
		}
	} else {
		return fmt.Errorf("device with id %s not found in map", action.TypeID)
	}
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
			Executed:    0,
			Conditions:  nil,
			Actions:     readRule.Actions,
		}
		rule.Conditions = generateLogicalCondition(readRule.Conditions)
		rules = append(rules, rule)
	}
	return rules
}

func generateLogicalCondition(conditions interface{}) LogicalCondition {
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
			panic(fmt.Sprintf("JSON config in wrong format, component_id should be of type string, %v is of type %s", logic["component_id"], reflect.TypeOf(logic["component_id"]).Kind().String()))
		}

		return constraint
	}
	panic(fmt.Sprintf("JSON config in wrong constraint format, conditions: %v, could not be processed", constraints))
}

//TODO multiple errors?
//TODO check device present
//TODO check component present
//TODO check output types
//TODO catch non-existing keys in json
