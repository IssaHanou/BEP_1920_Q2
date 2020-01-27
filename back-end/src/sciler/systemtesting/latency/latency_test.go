package latency

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	logger "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sciler/communication"
	"sciler/config"
	"sciler/handler"
	"testing"
	"time"
)

// Make sure when running windows, the setup in systemtesting/README.md is followed

func TestLatency(t *testing.T) {
	// edit mosquitto.conf bind address to local ip of the Pi
	// start broker with `mosquitto -c mosquitto.conf`
	// edit latency_config.json host to local ip of the Pi
	// start latency client with python3 latency.py

	// run this test

	logger.Info("starting latency test")
	receiveTimes := make(map[int]int64, 100)

	dir, dirErr := os.Getwd()
	if dirErr != nil {
		logger.Fatal(dirErr)
	}

	// setup back-end
	filename := filepath.Join(dir, "system_test_config.json")
	configurations := config.ReadFile(filename)
	logger.SetLevel(logger.ErrorLevel)

	loop := true

	messageHandler := handler.Handler{Config: configurations, ConfigFile: filename}
	messageHandler.Communicator = communication.NewCommunicator(configurations, func(client mqtt.Client, message mqtt.Message) {
		var msg handler.Message
		if err := json.Unmarshal(message.Payload(), &msg); err != nil {
			logger.Errorf("invalid JSON received: %v", err)
		}

		if msg.Type == "status" {
			contents := msg.Contents.(map[string]interface{})
			i := int(contents["ping"].(float64))
			if receiveTimes[i] == 0 {
				receiveTimes[i] = makeTimestamp()
				loop = false
			}
		}

	}, func() {
		messageHandler.SendSetup()
	})
	// start back-end
	messageHandler.Communicator.Start()

	const numberOfMessage = 100
	var meanLatency float64
	var minLatency int64
	var maxLatency int64

	requestTimes := make(map[int]int64, 100)
	for i := 0; i < numberOfMessage; i++ {
		// Send instruction to Client computer
		requestTimes[i] = makeTimestamp()
		messageHandler.SendComponentInstruction("latency", []config.ComponentInstruction{{
			ComponentID: "ping",
			Instruction: "ping",
			Value:       i,
		}}, "")

		// wait for status update
		for loop {
			time.Sleep(1 * time.Nanosecond)
		}

		assert.False(t, receiveTimes[i] == 0, "Did not receive an status message in time")
		latency := receiveTimes[i] - requestTimes[i]
		assert.Truef(t, latency <= 100, "the latency of all messages should be <= 100 ms but latency was %v ms", latency)
		meanLatency += float64(latency) / numberOfMessage
		if i == 0 {
			minLatency = latency
			maxLatency = latency
		} else {
			if minLatency > latency {
				minLatency = latency
			}
			if maxLatency < latency {
				maxLatency = latency
			}
		}
		loop = true
	}

	fmt.Println(fmt.Sprintf("%d messages were sent with a mean latency of %v ms", numberOfMessage, meanLatency))
	fmt.Println(fmt.Sprintf("with a minimal latency of %v ms", minLatency))
	fmt.Println(fmt.Sprintf("with a maximal latency of %v ms", maxLatency))
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func executeScript(script string) {
	if runtime.GOOS == "windows" {
		out, err := exec.Command("bash", script).Output()
		if err != nil {
			logger.Fatal(err)
		}
		fmt.Println(string(out))
	} else if runtime.GOOS == "linux" {
		out, err := exec.Command("./" + script).Output()
		if err != nil {
			logger.Fatal(err)
		}
		fmt.Println(string(out))
	}
}
