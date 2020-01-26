package systemtesting

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
	// start broker
	// start latency client

	receiveTimes := make(map[int]int64, 100)

	dir, dirErr := os.Getwd()
	if dirErr != nil {
		logger.Fatal(dirErr)
	}

	// setup back-end
	filename := filepath.Join(dir, "system_test_config.json")
	configurations := config.ReadFile(filename)
	logger.SetLevel(logger.ErrorLevel)

	messageHandler := handler.Handler{Config: configurations, ConfigFile: filename}
	messageHandler.Communicator = communication.NewCommunicator(configurations, func(client mqtt.Client, message mqtt.Message) {
		var msg handler.Message
		if err := json.Unmarshal(message.Payload(), &msg); err != nil {
			logger.Errorf("invalid JSON received: %v", err)
		}

		if msg.Type == "status" {
			contents := msg.Contents.(map[string]interface{})
			receiveTimes[int(contents["ping"].(float64))] = makeTimestamp()
		}

	}, func() {
		messageHandler.SendSetup()
	})
	// start back-end
	messageHandler.Communicator.Start()

	const numberOfMessage = 100

	requestTimes := make(map[int]int64, 100)
	for i := 0; i < numberOfMessage; i++ {

		// Send instruction to Client computer
		requestTimes[i] = makeTimestamp()
		messageHandler.SendComponentInstruction("latency", []config.ComponentInstruction{{
			ComponentID: "ping",
			Instruction: "ping",
			Value:       i,
		}}, "")

	}

	time.Sleep(1 * time.Second)

	var meanLatency float64

	for i := 0; i < numberOfMessage; i++ {
		assert.False(t, receiveTimes[i] == 0, "Did not receive an status message in time")
		latency := receiveTimes[i] - requestTimes[i]
		fmt.Println(receiveTimes[i])
		fmt.Println(requestTimes[i])
		fmt.Println(latency)
		assert.Truef(t, latency <= 100, "the latency of all messages should be <= 100 ms but latency was %v ms", latency)
		meanLatency += float64(latency) / numberOfMessage
	}

	fmt.Println(fmt.Sprintf("%d messages were sent with a mean latency of %v ms", numberOfMessage, meanLatency))
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
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
