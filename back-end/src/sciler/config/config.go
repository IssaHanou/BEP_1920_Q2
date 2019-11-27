package config

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"hash"
	"sort"
	"strings"
)

// ReadJSON transforms json file into config object.
func ReadJSON(input []byte) ReadConfig {
	var config ReadConfig
	err := json.Unmarshal(input, &config)
	if err != nil {
		panic(err.Error())
	}
	newConfig := generateDataStructures(config)
	if !checkConstraintValues(newConfig.Devices, newConfig.Conditions) {
		panic("Could not cast all constraint values to proper type, check types to match with device input.")
	}
	return config
}

// Creates additional structures for conditions and actions
func generateDataStructures(readConfig ReadConfig) WorkingConfig {
	var config WorkingConfig
	for _, d := range readConfig.Devices {
		config.Devices[d.ID] = d
	}
	for _, p := range readConfig.Puzzles {
		for _, r := range p.Rules {
			config.Rules[r.ID] = r
			for _, c := range r.Conditions {
				config.ActionMap[c.TypeID] = r.Actions
				config.ConstraintMap[c.TypeID] = c.Constraints
			}
		}
	}
	return config
}

// Make comparable
// Condition ID?
func checkConstraintValues(devices []Device, condition Condition) ([]ProcessedConstraint, error) {
	var constraints []ProcessedConstraint
	if condition.Type == "device" {
		device := devices[condition]
	}
	for _, con := range c.Constraints {
	}
	return
}

//Takes in duration string in format hh:mm:ss
func formatDuration(duration string) string {
	vars := strings.Split(duration, ":")
	var result string
	if vars[0] != "00" {
		result += vars[0] + " hours and"
	}
	if vars[1] != "00" {
		result += vars[1] + " minutes"
	}
	if vars[2] != "00" {
		result += " and " + vars[2] + " seconds"
	}
	return result
}
