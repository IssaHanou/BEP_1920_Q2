package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
)

// WorkingConfig has additional fields to ReadConfig, with lists of conditions, constraints and actions.
type WorkingConfig struct {
	General       General
	Puzzles       []Puzzle
	GeneralEvents []Event
	Devices       map[string]Device
}

// Rule is a struct that describes how action flow is handled in the escape room.
type Rule struct {
	ID          string
	Description string
	Limit       int
	Conditions  LogicalCondition
	Actions     []Action
}

// Puzzle is a struct that describes contents of a puzzle.
type Puzzle struct {
	Event Event
	Hints []string
}

// Event defines a general event, like start.
type Event struct {
	Name  string `json:"name"`
	Rules []Rule `json:"rules"`
}

func compare(param1 interface{}, param2 interface{}, comparision string) bool { // todo check types
	switch comparision {
	case "eq":
		return param1 == param2
	case "lt":
		return param1.(float64) < param2.(float64)
	case "gt":
		return param1.(float64) > param2.(float64)
	case "lte":
		return param1.(float64) <= param2.(float64)
	case "gte":
		return param1.(float64) >= param2.(float64)
	case "contains":
		return contains(param1.([]interface{}), param2)
	default:
		logrus.Errorf("Cannot compare on: %s", comparision)
		return false
	}
}

func contains(list []interface{}, element interface{}) bool {
	for _, item := range list {
		if element == item {
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
	Output      OutputObject
	Status      map[string]interface{}
	Connection  bool
}

// Condition is a struct that determines when rules are fired.
type Condition struct {
	Type        string            `json:"type"`
	TypeID      string            `json:"type_id"`
	Constraints LogicalConstraint `json:"constraints"`
}

// CheckConstraints is a method that checks types and comparator operators
func (condition Condition) CheckConstraints(config WorkingConfig) error {
	return condition.Constraints.CheckConstraints(condition, config)
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

// CheckConstraints is a method that checks types and comparator operators
func (constraint Constraint) CheckConstraints(condition Condition, config WorkingConfig) error {
	switch condition.Type {
	case "device":
		{
			device := config.Devices[condition.TypeID]
			valueType := reflect.TypeOf(constraint.Value).Kind()
			comparision := constraint.Comparison
			switch device.Input[constraint.ComponentID] {
			case "string":
				{
					if valueType != reflect.String {
						return fmt.Errorf("input type string expected but %s found as type of %v", valueType.String(), constraint.Value)
					}
					if comparision != "eq" {
						return fmt.Errorf("comparision %s not allowed on a string", comparision)
					}
				}
			case "boolean":
				{
					if valueType != reflect.Bool {
						return fmt.Errorf("input type boolean expected but %s found as type of %v", valueType.String(), constraint.Value)
					}
					if comparision != "eq" {
						return fmt.Errorf("comparision %s not allowed on a bool", comparision)
					}
				}
			case "numeric":
				{
					if valueType != reflect.Int && valueType != reflect.Float64 {
						return fmt.Errorf("input type numeric expected but %s found as type of %v", valueType.String(), constraint.Value)
					}
					if comparision == "contains" {
						return fmt.Errorf("comparision %s not allowed on a numeric", comparision)
					}
				}
			case "array":
				{
					if valueType != reflect.Array {
						return fmt.Errorf("input type array expected but %s found as type of %v", valueType.String(), constraint.Value)
					}
					if comparision != "contains" {
						return fmt.Errorf("comparision %s not allowed on an array", comparision)
					}
				}
			default:
				// todo custom types
			}
		}
	case "timer":
		return nil // todo
	default:
		return fmt.Errorf("cannot resolve constraint %v because condition.type is an unknown type", constraint)
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
	case "timer": //todo
		return false
	default:
		logrus.Errorf("cannot resolve constraint %v because condition.type is an unknown type", constraint)
		return false
	}
}

// LogicalCondition is an interface for operators and conditions
type LogicalCondition interface {
	Resolve(config WorkingConfig) bool
	CheckConstraints(config WorkingConfig) error
}

// AndCondition is an operator which implements the LogicalCondition interface
type AndCondition struct {
	logics []LogicalCondition
}

// CheckConstraints is a method that checks types and comparator operators
func (and AndCondition) CheckConstraints(config WorkingConfig) error {
	for _, logic := range and.logics {
		err := logic.CheckConstraints(config)
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

// CheckConstraints is a method that checks types and comparator operators
func (or OrCondition) CheckConstraints(config WorkingConfig) error {
	for _, logic := range or.logics {
		err := logic.CheckConstraints(config)
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
	CheckConstraints(condition Condition, config WorkingConfig) error
}

// AndConstraint is an operator which implement the LogicalConstraint interface
type AndConstraint struct {
	logics []LogicalConstraint
}

// CheckConstraints is a method that checks types and comparator operators
func (and AndConstraint) CheckConstraints(condition Condition, config WorkingConfig) error {
	for _, logic := range and.logics {
		err := logic.CheckConstraints(condition, config)
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

// CheckConstraints is a method that checks types and comparator operators
func (or OrConstraint) CheckConstraints(condition Condition, config WorkingConfig) error {
	for _, logic := range or.logics {
		err := logic.CheckConstraints(condition, config)
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

//todo resolve on constraint
