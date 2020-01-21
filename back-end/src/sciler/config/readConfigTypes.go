package config

// ReadConfig specifies all configuration elements of an escape room.
type ReadConfig struct {
	General       General            `json:"general"`
	Cameras       []Camera           `json:"cameras"`
	Devices       []ReadDevice       `json:"devices"`
	Timers        []ReadTimer        `json:"timers"`
	Puzzles       []ReadPuzzle       `json:"puzzles"`
	GeneralEvents []ReadGeneralEvent `json:"general_events"`
}

// General is a struct that describes the configurations of an escape room.
type General struct {
	Name     string `json:"name"`
	Duration string `json:"duration"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

// Camera is a struct describing camera feed with name of camera and link to the feed
type Camera struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

// ReadDevice is a struct that describes the configurations of a device in the room.
type ReadDevice struct {
	ID          string                  `json:"id"`
	Description string                  `json:"description"`
	Input       map[string]string       `json:"input"`
	Output      map[string]OutputObject `json:"output"`
}

// ReadTimer is a struct that describes time and id of timers
type ReadTimer struct {
	ID       string `json:"id"`
	Duration string `json:"duration"`
}

// ReadPuzzle is a struct that describes contents of a puzzle.
type ReadPuzzle struct {
	Name  string     `json:"name"`
	Rules []ReadRule `json:"rules"`
	Hints []string   `json:"hints"`
}

// GetName returns the name of a ReadPuzzle
func (r ReadPuzzle) GetName() string {
	return r.Name
}

// GetRules returns the readRules of a ReadPuzzle
func (r ReadPuzzle) GetRules() []ReadRule {
	return r.Rules
}

// ReadGeneralEvent defines a general event, like start.
type ReadGeneralEvent struct {
	Name  string     `json:"name"`
	Rules []ReadRule `json:"rules"`
}

// GetName returns the name of a ReadGeneralEvent
func (r ReadGeneralEvent) GetName() string {
	return r.Name
}

// GetRules returns the readRules of a ReadGeneralEvent
func (r ReadGeneralEvent) GetRules() []ReadRule {
	return r.Rules
}

// ReadEvent is an interface that both ReadPuzzle and ReadGeneralEvent implement
type ReadEvent interface {
	GetName() string
	GetRules() []ReadRule
}

// ReadOperator defines a object that takes an operator and combines a logics of other operators or conditions
type ReadOperator struct {
	Operator string        `json:"operator"`
	List     []interface{} `json:"logics"`
}

// ReadRule is a struct that describes how action flow is handled in the escape room.
type ReadRule struct {
	ID          string      `json:"id"`
	Description string      `json:"description"`
	Limit       int         `json:"limit"`
	Conditions  interface{} `json:"conditions"`
	Actions     []Action    `json:"actions"`
}

// ReadCondition is a struct that determines when rules are fired.
type ReadCondition struct {
	Type        string           `json:"type"`
	TypeID      string           `json:"type_id"`
	Constraints []ConstraintInfo `json:"constraints"`
	RuleID      string
}

// ConstraintInfo is a general map allowing to read input constraints, which are later parsed to real constraint objects.
type ConstraintInfo map[string]interface{}

// Action is a struct that determines what happens when a rule is fired.
type Action struct {
	Type    string                 `json:"type"`
	TypeID  string                 `json:"type_id"`
	Message []ComponentInstruction `json:"message"`
	Delay   string                 `json:"delay"`
}

// Execute is a method that performs the action
// TODO test this
func (action Action) Execute(handler InstructionSender) {
	switch action.Type { // this cannot be any other Type than device, timer or label, (checked in checkActions function)
	case "device":
		{
			handler.SendComponentInstruction(action.TypeID, action.Message, action.Delay)
		}
	case "timer":
		{
			handler.SetTimer(action.TypeID, action.Message[0])
		}
	case "label":
		{
			handler.SendLabelInstruction(action.TypeID, action.Message, action.Delay)
		}
	}

}

// ComponentInstruction can be sent across clients of the brokers.
type ComponentInstruction struct {
	ComponentID string      `json:"component_id"`
	Instruction string      `json:"instruction"`
	Value       interface{} `json:"value"`
}

// OutputObject contains a map defining either input or output.
type OutputObject struct {
	Type         string            `json:"type"`
	Instructions map[string]string `json:"instructions"`
	Label        []string          `json:"label"`
}
