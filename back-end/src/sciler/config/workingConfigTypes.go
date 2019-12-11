package config

import (
	"fmt"
	"reflect"
)

// WorkingConfig has additional fields to ReadConfig, with lists of conditions, constraints and actions.
type WorkingConfig struct {
	General       General
	Puzzles       []Puzzle
	GeneralEvents []GeneralEvent
	Devices       map[string]Device
	StatusMap     map[string][]Rule
}

// Rule is a struct that describes how action flow is handled in the escape room.
type Rule struct {
	ID          string
	Description string
	Limit       int
	Executed    int
	Conditions  LogicalCondition
	Actions     []Action
}

// Puzzle is a struct that describes contents of a puzzle.
type Puzzle struct {
	Event GeneralEvent
	Hints []string
}

// GetName returns the name of a Puzzle
func (p Puzzle) GetName() string {
	return p.Event.Name
}

// GetRules returns the rules of a Puzzle
func (p Puzzle) GetRules() []Rule {
	return p.Event.Rules
}

// GeneralEvent defines a general event, like start.
type GeneralEvent struct {
	Name  string
	Rules []Rule
}

// GetName returns the name of a GeneralEvent
func (g GeneralEvent) GetName() string {
	return g.Name
}

// GetRules returns the rules of a GeneralEvent
func (g GeneralEvent) GetRules() []Rule {
	return g.Rules
}

// Event is an interface that both Puzzle and GeneralEvent implement
type Event interface {
	GetName() string
	GetRules() []Rule
}

func compare(param1 interface{}, param2 interface{}, comparision string) bool {
	switch comparision {
	case "eq":
		return reflect.DeepEqual(param1, param2)
	case "lt":
		return param1.(float64) < param2.(float64)
	case "gt":
		return param1.(float64) > param2.(float64)
	case "lte":
		return param1.(float64) <= param2.(float64)
	case "gte":
		return param1.(float64) >= param2.(float64)
	case "contains":
		return contains(param1, param2)
	default:
		panic(fmt.Sprintf("cannot compare on: %s", comparision))
	}
}

func contains(list interface{}, element interface{}) bool {
	slice := reflect.ValueOf(list)
	for i := 0; i < slice.Len(); i++ {
		if element == slice.Index(i).Interface() {
			return true
		}
	}
	return false
}

// Device is a struct for all the devices
type Device struct {
	ID          string
	Description string
	Input       map[string]string
	Output      map[string]OutputObject
	Status      map[string]interface{}
	Connection  bool
}

// Condition is a struct that determines when rules are fired.
type Condition struct {
	Type        string
	TypeID      string
	Constraints LogicalConstraint
}

// GetConditionIDs returns a slice of all condition type IDs in the LogicalCondition
func (condition Condition) GetConditionIDs() []string {
	return []string{condition.TypeID}
}

// checkConstraints is a method that checks types and comparator operators
func (condition Condition) checkConstraints(config WorkingConfig) error {
	// todo check if type id exists
	return condition.Constraints.checkConstraints(condition, config)
}

// Resolve is a method that checks if a condition is met
func (condition Condition) Resolve(config WorkingConfig) bool {
	return condition.Constraints.Resolve(condition, config)
}

// Constraint is main constraint type with comparison and componentID (latter only for device condition).
type Constraint struct {
	Comparison  string
	ComponentID string
	Value       interface{}
}

// checkConstraints is a method that checks types and comparator operators
func (constraint Constraint) checkConstraints(condition Condition, config WorkingConfig) error {
	switch condition.Type {
	case "device":
		{
			if device, ok := config.Devices[condition.TypeID]; ok { // checks if device can be found in the map, if so, it is stored in variable device

				valueType := reflect.TypeOf(constraint.Value).Kind()
				comparision := constraint.Comparison
				if inputType, ok := device.Input[constraint.ComponentID]; ok {
					switch inputType {
					case "string":
						{
							if valueType != reflect.String {
								return fmt.Errorf("input type string expected but %s found as type of value %v", valueType.String(), constraint.Value)
							}
							if comparision != "eq" {
								return fmt.Errorf("comparision %s not allowed on a string", comparision)
							}
						}
					case "boolean":
						{
							if valueType != reflect.Bool {
								return fmt.Errorf("input type boolean expected but %s found as type of value %v", valueType.String(), constraint.Value)
							}
							if comparision != "eq" {
								return fmt.Errorf("comparision %s not allowed on a boolean", comparision)
							}
						}
					case "numeric":
						{
							if valueType != reflect.Int && valueType != reflect.Float64 {
								return fmt.Errorf("input type numeric expected but %s found as type of value %v", valueType.String(), constraint.Value)
							}
							if comparision == "contains" {
								return fmt.Errorf("comparision %s not allowed on a numeric", comparision)
							}
						}
					case "array":
						{
							if valueType != reflect.Slice {
								return fmt.Errorf("input type array/slice expected but %s found as type of value %v", valueType.String(), constraint.Value)
							}
							if comparision != "contains" && comparision != "eq" {
								return fmt.Errorf("comparision %s not allowed on an array", comparision)
							}
						}
					default:
						// todo custom types
						return fmt.Errorf("custom types like: %s, are not yet implemented", inputType)
					}
				} else {
					return fmt.Errorf("component with id %s not found in map", constraint.ComponentID)
				}
			} else {
				return fmt.Errorf("device with id %s not found in map", condition.TypeID)
			}
		}
	case "timer":
		return nil // todo timer
	case "rule":
		return nil // todo rule
	default:
		return fmt.Errorf("invalid type of condition: %v", condition.Type)
	}
	return nil
}

// Resolve is a method that checks if a constraint is met
func (constraint Constraint) Resolve(condition Condition, config WorkingConfig) bool {
	switch condition.Type {
	case "device":
		{
			device := config.Devices[condition.TypeID]
			status := device.Status[constraint.ComponentID]
			return compare(status, constraint.Value, constraint.Comparison)
		}
	case "timer": //todo timer
		panic(fmt.Sprintf("cannot resolve constraint %v because condition.type is an timer type, which is not implemented yet", constraint))
	default:
		panic(fmt.Sprintf("cannot resolve constraint %v because condition.type is an unknown type, this should already be checked when reading in the JSON", constraint))
	}
}

// LogicalCondition is an interface for operators and conditions
type LogicalCondition interface {
	Resolve(config WorkingConfig) bool
	checkConstraints(config WorkingConfig) error
	GetConditionIDs() []string
}

// AndCondition is an operator which implements the LogicalCondition interface
type AndCondition struct {
	logics []LogicalCondition
}

// GetConditionIDs returns a slice of all condition type IDs in the LogicalCondition
func (and AndCondition) GetConditionIDs() []string {
	var IDs []string
	for _, logic := range and.logics {
		IDs = append(IDs, logic.GetConditionIDs()...)
	}
	return IDs
}

// checkConstraints is a method that checks types and comparator operators
func (and AndCondition) checkConstraints(config WorkingConfig) error {
	for _, logic := range and.logics {
		err := logic.checkConstraints(config)
		if err != nil {
			return err
		}
	}
	return nil
}

// Resolve is a method that checks if a condition is met
func (and AndCondition) Resolve(config WorkingConfig) bool {
	result := true
	for _, logic := range and.logics {
		result = result && logic.Resolve(config)
	}
	return result
}

// OrCondition is an operator which implements the LogicalCondition interface
type OrCondition struct {
	logics []LogicalCondition
}

// GetConditionIDs returns a slice of all condition type IDs in the LogicalCondition
func (or OrCondition) GetConditionIDs() []string {
	var IDs []string
	for _, logic := range or.logics {
		IDs = append(IDs, logic.GetConditionIDs()...)
	}
	return IDs
}

// checkConstraints is a method that checks types and comparator operators
func (or OrCondition) checkConstraints(config WorkingConfig) error {
	for _, logic := range or.logics {
		err := logic.checkConstraints(config)
		if err != nil {
			return err
		}
	}
	return nil
}

// Resolve is a method that checks if a condition is met
func (or OrCondition) Resolve(config WorkingConfig) bool {
	result := false
	for _, logic := range or.logics {
		result = result || logic.Resolve(config)
	}
	return result
}

// LogicalConstraint is an interface for operators and constraints
type LogicalConstraint interface {
	Resolve(condition Condition, config WorkingConfig) bool
	checkConstraints(condition Condition, config WorkingConfig) error
}

// AndConstraint is an operator which implement the LogicalConstraint interface
type AndConstraint struct {
	logics []LogicalConstraint
}

// checkConstraints is a method that checks types and comparator operators
func (and AndConstraint) checkConstraints(condition Condition, config WorkingConfig) error {
	for _, logic := range and.logics {
		err := logic.checkConstraints(condition, config)
		if err != nil {
			return err
		}
	}
	return nil
}

// Resolve is a method that checks if a constraint is met
func (and AndConstraint) Resolve(condition Condition, config WorkingConfig) bool { // todo: make lazy
	result := true
	for _, logic := range and.logics {
		result = result && logic.Resolve(condition, config)
	}
	return result
}

// OrConstraint is an operator which implement the LogicalConstraint interface
type OrConstraint struct {
	logics []LogicalConstraint
}

// checkConstraints is a method that checks types and comparator operators
func (or OrConstraint) checkConstraints(condition Condition, config WorkingConfig) error {
	for _, logic := range or.logics {
		err := logic.checkConstraints(condition, config)
		if err != nil {
			return err
		}
	}
	return nil
}

// Resolve is a method that checks if a constraint is met
func (or OrConstraint) Resolve(condition Condition, config WorkingConfig) bool { // todo: make lazy
	result := false
	for _, logic := range or.logics {
		result = result || logic.Resolve(condition, config)
	}
	return result
}
