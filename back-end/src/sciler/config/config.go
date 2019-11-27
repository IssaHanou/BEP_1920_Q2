package config

import (
	"encoding/json"
	"strings"
)

// Transforms json file into config object.
func ReadJSON(input []byte) Config {
	var data Config
	err := json.Unmarshal(input, &data)
	if err != nil {
		panic(err.Error())
	}
	return data
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
