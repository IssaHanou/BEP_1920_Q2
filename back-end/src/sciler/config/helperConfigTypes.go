package config

// WorkingConfig has additional fields to ReadConfig, with lists of conditions, constraints and actions.
type WorkingConfig struct {
	General       General
	Puzzles       []Puzzle
	GeneralEvents []GeneralEvent
	//TODO maps and concurrency
	Devices   map[string]Device
	Rules     map[string]Rule
	ActionMap map[string][]Action
	// first key is condition pointer, second is type of constraint:
	// string, bool, numeric, string-array, bool-array, num-array, custom, timer
	// ex: ConstraintMap[conditionID][boolean] = []ConstraintBool
	ConstraintMap map[string]map[string][]interface{}
}

// Constraint is main constraint type with comparison and componentID (latter only for device condition).
type Constraint struct {
	Comparison  string
	ComponentID string
}

// ConstraintTimer is struct for constraint with timer value: "hh:mm:ss".
type ConstraintTimer struct {
	Constraint
	Value string
}

// ConstraintNumeric is struct for constraint with numeric value.
type ConstraintNumeric struct {
	Constraint
	Value float64
}

// ConstraintString is struct for constraint with string value.
type ConstraintString struct {
	Constraint
	Value string
}

// ConstraintBool is struct for constraint with boolean value.
type ConstraintBool struct {
	Constraint
	Value bool
}

// ConstraintStringArray is struct for constraint with string array value.
type ConstraintStringArray struct {
	Constraint
	Value []string
}

// ConstraintNumericArray is struct for constraint with numeric array value.
type ConstraintNumericArray struct {
	Constraint
	Value []float64
}

// ConstraintBoolArray is struct for constraint with bool array value.
type ConstraintBoolArray struct {
	Constraint
	Value []bool
}

// ConstraintCustomType is struct for constraint with custom type value.
type ConstraintCustomType struct {
	Constraint
	Value interface{}
}
