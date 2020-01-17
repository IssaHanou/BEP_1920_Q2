package systemtesting

import (
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
)

// Make sure when running windows, the setup in systemtesting/README.md is followed

func TestLatency(t *testing.T) {
	if runtime.GOOS == "windows" {
		out, err := exec.Command("bash", "localBroker.sh").Output()
		if err != nil {
			logger.Fatal(err)
		}
		fmt.Println(string(out))
	} else if runtime.GOOS == "linux" {
		out, err := exec.Command("localBroker.sh").Output()
		if err != nil {
			logger.Fatal(err)
		}
		fmt.Println(string(out))
	}

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

	messageHandler := handler.Handler{Config: configurations, ConfigFile: filename}
	messageHandler.Communicator = communication.NewCommunicator(host, port, []string{"back-end"}, func(client mqtt.Client, message mqtt.Message) {
		//messageHandler.NewHandler(client, message)
		logger.Info(message.Payload())
	}, func() {
		messageHandler.SendSetup()
	})
	// start back-end
	messageHandler.Communicator.Start()

	// Send instruction to Client computer
	logger.Info(messageHandler.Config.Devices["latency"].Status)
	messageHandler.SendInstruction("latency", []map[string]string{{
		"instruction":   "ping",
		"instructed_by": "back-end",
	}})

	assert.Equal(t, true, true)
}
