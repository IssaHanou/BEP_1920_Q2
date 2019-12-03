package config

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
)

// ReadConfig specifies all configuration elements of an escape room.
type ReadConfig struct {
	General       General        `json:"general"`
	Devices       []ReadDevice   `json:"devices"`
	Puzzles       []Puzzle       `json:"puzzles"`
	GeneralEvents []GeneralEvent `json:"general_events"`
}

// General is a struct that describes the configurations of an escape room.
type General struct {
	Name     string `json:"name"`
	Duration string `json:"duration"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

// ReadDevice is a struct that describes the configurations of a device in the room.
type ReadDevice struct {
	ID          string            `json:"id"`
	Description string            `json:"description"`
	Input       map[string]string `json:"input"`
	Output      OutputObject      `json:"output"`
}

// Puzzle is a struct that describes contents of a puzzle.
type Puzzle struct {
	Name  string   `json:"name"`
	Rules []Rule   `json:"rules"`
	Hints []string `json:"hints"`
}

// GeneralEvent defines a general event, like start.
type GeneralEvent struct {
	Name  string `json:"name"`
	Rules []Rule `json:"rules"`
}

// Rule is a struct that describes how action flow is handled in the escape room.
type Rule struct {
	ID          string      `json:"id"`
	Description string      `json:"description"`
	Limit       int         `json:"limit"`
	Conditions  []Condition `json:"conditions"`
	Actions     []Action    `json:"actions"`
}

// Condition is a struct that determines when rules are fired.
type Condition struct {
	Type        string           `json:"type"`
	TypeID      string           `json:"type_id"`
	Constraints []ConstraintInfo `json:"constraints"`
	RuleID      string
}

// ConstraintInfo is a general map allowing to read input constraints, which are later parsed to real constraint objects.
type ConstraintInfo map[string]interface{}

// GetID returns hash of condition, limited to the first 24 characters
func (condition *Condition) GetID() string {
	hasher := sha512.New()
	toHash := condition.RuleID + condition.TypeID + fmt.Sprint(condition.Constraints)
	hasher.Write([]byte(toHash))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return hash[0:24]
}

// Action is a struct that determines what happens when a rule is fired.
type Action struct {
	Type    string        `json:"type"`
	TypeID  string        `json:"type_id"`
	Message ActionMessage `json:"message"`
}

// ActionMessage can be sent across clients of the brokers.
type ActionMessage struct {
	Output OutputObject `json:"output"`
}

// OutputObject contains a map defining either input or output.
type OutputObject map[string]interface{}
