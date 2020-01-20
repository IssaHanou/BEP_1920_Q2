package config

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"reflect"
	"time"
)

// WorkingConfig has additional fields to ReadConfig, with lists of conditions, constraints and actions.
type WorkingConfig struct {
	General       General
	Cameras       []Camera
	Puzzles       []*Puzzle
	GeneralEvents []*GeneralEvent
	Devices       map[string]*Device
	Timers        map[string]*Timer
	StatusMap     map[string][]*Rule
	RuleMap       map[string]*Rule
	LabelMap      map[string][]*Component
}

// Timer is a timer in the escape game
type Timer struct {
	ID        string
	Duration  time.Duration
	StartedAt time.Time
	T         *time.Timer
	State     string
	Ending    func()
	Finished  bool
}

// newTimer create a new timer
func newTimer(id string, d time.Duration) *Timer {
	t := new(Timer)
	t.ID = id
	t.Duration = d
	t.Finished = false
	t.State = "stateIdle"
	return t
}

// GetTimeLeft gets time left of timer
func (t *Timer) GetTimeLeft() (time.Duration, string) {
	left := 0 * time.Second
	if t.State == "stateActive" {
		dif := time.Now().Sub(t.StartedAt)
		left = t.Duration - dif
	} else if t.State == "stateIdle" {
		left = t.Duration
	}
	return left, t.State
}

// Start starts Timer that executes handler.HandleEvent(t.ID) after duration d.
// Can not start if timer is not state Idle
func (t *Timer) Start(handler InstructionSender) error {
	if t.State != "stateIdle" {
		return fmt.Errorf("timer %v does not have an Idle state and can not be started", t.ID)
	}
	t.StartedAt = time.Now()
	t.State = "stateActive"
	t.Ending = func() {
		t.State = "stateExpired"
		t.Finished = true
		logger.Infof("timer %v finished", t.ID)
		handler.HandleEvent(t.ID)
	}
	t.T = time.AfterFunc(t.Duration, t.Ending)
	logger.Infof("timer %v started for %v", t.ID, t.Duration)
	return nil
}

// Pause sets the timer to Idle and sets new duration to the time left
// can not pause if timer is not state Active
func (t *Timer) Pause() error {
	if t.State != "stateActive" {
		return fmt.Errorf("timer %v does not have a Active state and can not be paused", t.ID)
	}
	t.T.Stop()
	t.Duration, _ = t.GetTimeLeft()
	t.State = "stateIdle"
	logger.Infof("timer paused with %v left", t.Duration)
	return nil
}

// AddSubTime add or subtract time to a timer
// can only subtract if that time is equally left
// can only add or subtract time if the state is not Expired
func (t *Timer) AddSubTime(handler InstructionSender, time time.Duration, add bool) error {
	if t.State == "stateIdle" {
		if add {
			t.Duration = t.Duration + time
			logger.Infof("timer %v added %v to duration and now has a duration of %v", t.ID, time, t.Duration)
		} else {
			if t.Duration > time {
				t.Duration = t.Duration - time
				logger.Infof("timer %v subtracted %v to duration and now has a duration of %v", t.ID, time, t.Duration)
			} else {
				return fmt.Errorf("timer %v could not subtract %v since there is only %v left", t.ID, time, t.Duration)
			}
		}
		return nil
	} else if t.State == "stateActive" {
		t.Pause()
		err := t.AddSubTime(handler, time, add)
		t.Start(handler)
		return err
	} else if t.State == "stateExpired" {
		return fmt.Errorf("timer %v could not be edited since it is already Expired", t.ID)
	}
	return fmt.Errorf("timer %v could not be edited because there is something wrong with its state: %v", t.ID, t.State)
}

// Stop make a timer stop
// can not stop timer that has state Expired
func (t *Timer) Stop() error {
	if t.State == "stateExpired" {
		return fmt.Errorf("timer %v is already Expired and can not be stopped again", t.ID)
	}
	if t.State == "stateIdle" {
		t.Duration = 0 * time.Second
	} else {
		t.T.Stop()
	}
	t.State = "stateExpired"

	logger.Infof("timer %v stopped and set to Expired without handling it's actions", t.ID)
	return nil
}

// Done finishes the timer as if it ran out of time
// can not finish a timer that is already Expired
func (t *Timer) Done() error {
	if t.State == "stateExpired" {
		return fmt.Errorf("timer %v is already Expired and can not be finished again", t.ID)
	}
	t.Ending()
	t.T.Stop()
	logger.Infof("timer %v stopped and set to Expired, actions are being handled", t.ID)
	return nil
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

// InstructionSender is an interface needed for preventing cyclic imports
type InstructionSender interface {
	SendComponentInstruction(string, []ComponentInstruction, string)
	SetTimer(string, ComponentInstruction)
	HandleEvent(string)
	SendLabelInstruction(string, []ComponentInstruction, string)
}

// Finished is a method that checks is the a rule have been finished, meaning if it reached its maximum number of executions
func (r *Rule) Finished() bool {
	return r.Executed == r.Limit
}

// Execute performs all actions of a rule
func (r *Rule) Execute(handler InstructionSender) {
	for _, action := range r.Actions {
		go action.Execute(handler)
	}
	r.Executed++
	logger.Infof("Executed rule %s", r.ID)
	handler.HandleEvent(r.ID)
}

// Puzzle is a struct that describes contents of a puzzle.
type Puzzle struct {
	Event *GeneralEvent
	Hints []string
}

// GetName returns the name of a Puzzle
func (p *Puzzle) GetName() string {
	return p.Event.Name
}

// GetRules returns the rules of a Puzzle
func (p *Puzzle) GetRules() []*Rule {
	var rules []*Rule
	for _, rule := range p.Event.Rules {
		rules = append(rules, rule)
	}
	return rules
}

// GeneralEvent defines a general event, like start.
type GeneralEvent struct {
	Name  string
	Rules []*Rule
}

// GetName returns the name of a GeneralEvent
func (g *GeneralEvent) GetName() string {
	return g.Name
}

// GetRules returns the rules of a GeneralEvent
func (g *GeneralEvent) GetRules() []*Rule {
	var rules []*Rule
	for _, rule := range g.Rules {
		rules = append(rules, rule)
	}
	return rules
}

// Event is an interface that both Puzzle and GeneralEvent implement
type Event interface {
	GetName() string
	GetRules() []*Rule
}

func compare(param1 interface{}, param2 interface{}, comparision string) bool {
	if param1 == nil {
		return false
	}
	switch comparision {
	case "eq":
		if reflect.TypeOf(param1).Kind() == reflect.Int || reflect.TypeOf(param2).Kind() == reflect.Int {
			return numericToFloat64(param1) == numericToFloat64(param2)
		}
		return reflect.DeepEqual(param1, param2)
	case "lt":
		return numericToFloat64(param1) < numericToFloat64(param2)
	case "gt":
		return numericToFloat64(param1) > numericToFloat64(param2)
	case "lte":
		return numericToFloat64(param1) <= numericToFloat64(param2)
	case "gte":
		return numericToFloat64(param1) >= numericToFloat64(param2)
	case "contains":
		return contains(param1, param2)
	case "not":
		if reflect.TypeOf(param1).Kind() == reflect.Int || reflect.TypeOf(param2).Kind() == reflect.Int {
			return numericToFloat64(param1) != numericToFloat64(param2)
		}
		return !reflect.DeepEqual(param1, param2)
	default:
		// This case is already handled to give error in checkConstraint
		return false
	}
}

// numericToFloat64 checks if numeric value is in or float64 and returns float64
func numericToFloat64(input interface{}) float64 {
	switch input.(type) {
	case float64:
		return input.(float64)
	case int:
		return float64(input.(int))
	default:
		// This case is already handled to give error in checkConstraint
		return 0
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

// Component is a struct that links a component id with a device id
type Component struct {
	ID     string
	Device *Device
}

// GetConditionIDs returns a slice of all condition type IDs in the LogicalCondition
func (condition Condition) GetConditionIDs() []string {
	return []string{condition.TypeID}
}

// checkConstraints is a method that checks types and comparator operators
func (condition Condition) checkConstraints(config WorkingConfig, ruleID string) []string {
	// todo check if type id exists
	return condition.Constraints.checkConstraints(condition, config, ruleID)
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
func (constraint Constraint) checkConstraints(condition Condition, config WorkingConfig, ruleID string) []string {
	switch condition.Type {
	case "device":
		{
			if device, ok := config.Devices[condition.TypeID]; ok { // checks if device can be found in the map, if so, it is stored in variable device

				valueType := reflect.TypeOf(constraint.Value).Kind()
				comparison := constraint.Comparison
				if inputType, ok := device.Input[constraint.ComponentID]; ok {
					switch inputType {
					case "string":
						{
							if valueType != reflect.String {
								return []string{fmt.Sprintf("on rule %s: input type string expected but %s found as type of value %v", ruleID, valueType.String(), constraint.Value)}
							}
							if !CheckValidComparison(comparison) {
								return []string{fmt.Sprintf("on rule %s: comparison %s is not valid", ruleID, comparison)}
							}
							if comparison != "eq" && comparison != "not" {
								return []string{fmt.Sprintf("on rule %s: comparison %s not allowed on a string", ruleID, comparison)}
							}
						}
					case "boolean":
						{
							if valueType != reflect.Bool {
								return []string{fmt.Sprintf("on rule %s: input type boolean expected but %s found as type of value %v", ruleID, valueType.String(), constraint.Value)}
							}
							if !CheckValidComparison(comparison) {
								return []string{fmt.Sprintf("on rule %s: comparison %s is not valid", ruleID, comparison)}
							}
							if comparison != "eq" {
								return []string{fmt.Sprintf("on rule %s: comparison %s not allowed on a boolean", ruleID, comparison)}
							}
						}
					case "numeric":
						{
							if valueType != reflect.Int && valueType != reflect.Float64 {
								return []string{fmt.Sprintf("on rule %s: input type numeric expected but %s found as type of value %v", ruleID, valueType.String(), constraint.Value)}
							}
							if !CheckValidComparison(comparison) {
								return []string{fmt.Sprintf("on rule %s: comparison %s is not valid", ruleID, comparison)}
							}
							if comparison == "contains" {
								return []string{fmt.Sprintf("on rule %s: comparison %s not allowed on a numeric", ruleID, comparison)}
							}
						}
					case "array":
						{
							if valueType != reflect.Slice {
								return []string{fmt.Sprintf("on rule %s: input type array/slice expected but %s found as type of value %v", ruleID, valueType.String(), constraint.Value)}
							}
							if !CheckValidComparison(comparison) {
								return []string{fmt.Sprintf("on rule %s: comparison %s is not valid", ruleID, comparison)}
							}
							if comparison != "contains" && comparison != "eq" && comparison != "not" {
								return []string{fmt.Sprintf("on rule %s: comparison %s not allowed on an array", ruleID, comparison)}
							}
						}
					default:
						// todo custom types
						return []string{fmt.Sprintf("on rule %s: custom types like: %s, are not yet implemented", ruleID, inputType)}
					}
				} else {
					return []string{fmt.Sprintf("on rule %s: component with id %s not found in map", ruleID, constraint.ComponentID)}
				}
			} else {
				return []string{fmt.Sprintf("on rule %s: device with id %s not found in map", ruleID, condition.TypeID)}
			}
		}
	case "timer":
		if _, ok := config.Timers[condition.TypeID]; ok {
			valueType := reflect.TypeOf(constraint.Value).Kind()
			comparison := constraint.Comparison
			if valueType != reflect.Bool {
				return []string{fmt.Sprintf("on rule %s: input type boolean expected but %s found as type of value %v", ruleID, valueType.String(), constraint.Value)}
			}
			if !CheckValidComparison(comparison) {
				return []string{fmt.Sprintf("on rule %s: comparison %s is not valid", ruleID, comparison)}
			}
			if comparison != "eq" {
				return []string{fmt.Sprintf("on rule %s: comparison %s not allowed on a boolean", ruleID, comparison)}
			}

		} else {
			return []string{fmt.Sprintf("on rule %s: timer with id %s not found in map", ruleID, condition.TypeID)}
		}
	case "rule":
		if _, ok := config.RuleMap[condition.TypeID]; ok { // checks if rule can be found in the map, if so, it is stored in variable device
			valueType := reflect.TypeOf(constraint.Value).Kind()
			comparison := constraint.Comparison
			if valueType != reflect.Int && valueType != reflect.Float64 {
				return []string{fmt.Sprintf("on rule %s: value type numeric expected but %s found as type of value %v", ruleID, valueType.String(), constraint.Value)}
			}
			if !CheckValidComparison(comparison) {
				return []string{fmt.Sprintf("on rule %s: comparison %s is not valid", ruleID, comparison)}
			}
			if comparison == "contains" {
				return []string{fmt.Sprintf("on rule %s: comparison %s not allowed on rule", ruleID, comparison)}
			}
		} else {
			return []string{fmt.Sprintf("on rule %s: rule with id %s not found in map", ruleID, condition.TypeID)}
		}
	default:
		return []string{fmt.Sprintf("on rule %s: invalid type of condition: %v", ruleID, condition.Type)}
	}
	// all cases for errors are already handled
	return make([]string, 0)
}

// CheckValidComparison checks if the comparison is a valid one
func CheckValidComparison(comparison string) bool {
	comparisonTypesAllowed := []string{"eq", "lt", "gt", "lte", "gte", "contains", "not"}
	for _, comp := range comparisonTypesAllowed {
		if comp == comparison {
			return true
		}
	}
	return false
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
	case "rule":
		{
			rule := config.RuleMap[condition.TypeID]
			return compare(rule.Executed, constraint.Value, constraint.Comparison)
		}
	case "timer":
		{
			timer := config.Timers[condition.TypeID]
			return compare(timer.Finished, constraint.Value, constraint.Comparison)

		}
	default:
		// This case is already handled to give error in checkConstraint
		return false
	}
}

// LogicalCondition is an interface for operators and conditions
type LogicalCondition interface {
	Resolve(config WorkingConfig) bool
	checkConstraints(config WorkingConfig, ruleID string) []string
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
func (and AndCondition) checkConstraints(config WorkingConfig, ruleID string) []string {
	errorList := make([]string, 0)
	for _, logic := range and.logics {
		err := logic.checkConstraints(config, ruleID)
		if err != nil {
			errorList = append(errorList, err...)
		}
	}
	return errorList
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
func (or OrCondition) checkConstraints(config WorkingConfig, ruleID string) []string {
	errorList := make([]string, 0)
	for _, logic := range or.logics {
		err := logic.checkConstraints(config, ruleID)
		if err != nil {
			errorList = append(errorList, err...)
		}
	}
	return errorList
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
	checkConstraints(condition Condition, config WorkingConfig, ruleID string) []string
}

// AndConstraint is an operator which implement the LogicalConstraint interface
type AndConstraint struct {
	logics []LogicalConstraint
}

// checkConstraints is a method that checks types and comparator operators
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

// Resolve is a method that checks if a constraint is met
func (or OrConstraint) Resolve(condition Condition, config WorkingConfig) bool { // todo: make lazy
	result := false
	for _, logic := range or.logics {
		result = result || logic.Resolve(condition, config)
	}
	return result
}
