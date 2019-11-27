package config

/**
Config specifies all configuration elements of an escape room.
*/
type Config struct {
	General       General        `json:"general"`
	Devices       []Device       `json:"devices"`
	Puzzles       []Puzzle       `json:"puzzles"`
	GeneralEvents []GeneralEvent `json:"general_events"`
}

/**
General is a struct that describes the configurations of an escape room.
*/
type General struct {
	Name     string `json:"name"`
	Duration string `json:"duration"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

/**
Device is a struct that describes the configurations of a device in the room.
*/
type Device struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Input       IOObject `json:"input"`
	Output      IOObject `json:"output"`
}

/**
Puzzle is a struct that describes contents of a puzzle.
*/
type Puzzle struct {
	Name  string   `json:"name"`
	Rules []Rule   `json:"rules"`
	Hints []string `json:"hints"`
}

/**
GeneralEvent defines a general event, like start.
*/
type GeneralEvent struct {
	Name  string `json:"name"`
	Rules []Rule `json:"rules"`
}

/**
Rule is a struct that describes how action flow is handled in the escape room.
*/
type Rule struct {
	ID          string      `json:"id"`
	Description string      `json:"description"`
	Limit       int         `json:"limit"`
	Conditions  []Condition `json:"conditions"`
	Actions     []Action    `json:"actions"`
}

/**
Condition is a struct that determines when rules are fired.
*/
type Condition struct {
	Type        string       `json:"type"`
	ID          string       `json:"id"`
	constraints []Constraint `json:"constraints"`
}

/**
Action is a struct that determines what happens when a rule is fired.
*/
type Action struct {
	Type    string        `json:"type"`
	ID      string        `json:"id"`
	Message ActionMessage `json:"message"`
}

/**
Constraint specifies a conditions.
*/
type Constraint struct {
	Comparison  string `json:"comparison"`
	Value       string `json:"value"`
	ComponentId string `json:"component_id"`
}

/**
Messsage can be sent across clients of the brokers.
*/
type ActionMessage struct {
	Output IOObject `json:"output"`
}

type IOObject map[string]interface{}
