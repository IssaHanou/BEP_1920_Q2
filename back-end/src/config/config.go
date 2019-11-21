package config

import (
	"encoding/json"
	"fmt"
	"strings"
)

type General struct {
	Name string `json:"name"`
	Duration string `json:"duration"`
}

// Takes in json with general info of escape room.
func GetFromJson(input []byte) string {
	var data General
	err := json.Unmarshal(input, &data)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	return "Escape room " + data.Name + " should be solved within " +
		formatDuration(data.Duration)
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