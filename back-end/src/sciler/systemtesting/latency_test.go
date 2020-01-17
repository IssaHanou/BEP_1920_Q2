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
	executeScript("clientComputer.sh")
	executeScript("localBroker.sh")
	dir, dirErr := os.Getwd()
	if dirErr != nil {
		logger.Fatal(dirErr)
	}

	// setup back-end
	filename := filepath.Join(dir, "system_test_config.json")
	configurations := config.ReadFile(filename)
	logger.Info("configurations read from: " + filename)
	host := configurations.General.Host
	port := configurations.General.Port

	var responseTime time.Time

	messageHandler := handler.Handler{Config: configurations, ConfigFile: filename}
	messageHandler.Communicator = communication.NewCommunicator(host, port, []string{"back-end"}, func(client mqtt.Client, message mqtt.Message) {
		var msg handler.Message
		if err := json.Unmarshal(message.Payload(), &msg); err != nil {
			logger.Errorf("invalid JSON received: %v", err)
		}

		if msg.Type == "confirmation" {
			responseTime = time.Now()
		}

	}, func() {
		messageHandler.SendSetup()
	})
	// start back-end
	messageHandler.Communicator.Start()

	// Send instruction to Client computer
	requestTime := time.Now()
	messageHandler.SendInstruction("latency", []map[string]string{{
		"instruction":   "ping",
		"instructed_by": "back-end",
	}})

	time.Sleep(time.Duration(1) * time.Second)
	latency := responseTime.Sub(requestTime)
	logger.Info(latency)
	assert.True(t, latency.Round(time.Millisecond).Milliseconds() <= 100)
}

func executeScript(script string) {
	if runtime.GOOS == "windows" {
		out, err := exec.Command("bash", script).Output()
		if err != nil {
			logger.Fatal(err)
		}
		fmt.Println(string(out))
	} else if runtime.GOOS == "linux" {
		out, err := exec.Command(script).Output()
		if err != nil {
			logger.Fatal(err)
		}
		fmt.Println(string(out))
	}
}
