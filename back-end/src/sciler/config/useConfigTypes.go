package config

// WorkingConfig has additional fields to ReadConfig, with lists of conditions, constraints and actions.
type WorkingConfig struct {
	General       General
	Puzzles       []Puzzle
	GeneralEvents []GeneralEvent
	Devices       map[string]Device
	Rules         map[string]Rule
	ConstraintMap map[Condition][]ProcessedConstraint
	ActionMap     map[Condition][]Action
}

// ProcessedConstraint has additional type and the value has been checked to be of Type.
type ProcessedConstraint struct {
	Comparison  string
	Value       string
	ComponentID string
	Type        string
}
