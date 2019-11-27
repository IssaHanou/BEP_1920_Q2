package config

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestGetFromJson(t *testing.T) {
	filename := "../../../resources/room_config.json"
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error("Could not read file room_config.json")
	}
	result := ReadJSON(dat)
	assert.Equal(t, result.General.Duration, "00:30:00", "Duration was not correct")
}
