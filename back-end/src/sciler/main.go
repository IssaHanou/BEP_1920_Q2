package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"sciler/communication"
	"sciler/config"
	"sciler/handler"
	"time"
)

var topics = []string{"back-end"}

func main() {
	dir, dirErr := os.Getwd()
	if dirErr != nil {
		logrus.Fatal(dirErr)
	}
	// Write to both cmd and file
	writeFile := filepath.Join(dir, "back-end", "output", "log-"+fmt.Sprint(time.Now().Format("02-01-2006--15-04-26"))+".txt")
	file, fileErr := os.Create(writeFile)
	if fileErr != nil {
		logrus.Fatal(fileErr)
	}
	logrus.Info("writing logs to both console and " + writeFile)
	multi := io.MultiWriter(os.Stdout, file)
	logrus.SetOutput(multi)

	filename := filepath.Join(dir, "back-end", "resources", "room_config.json")
	configurations := config.ReadFile(filename)
	logrus.Info("configurations read from: " + filename)
	host := configurations.General.Host
	port := configurations.General.Port

	communicator := communication.NewCommunicator(host, port, topics)

	messageHandler := handler.Handler{Config: configurations, ConfigFile: filename, Communicator: communicator}
	go communicator.Start(messageHandler.NewHandler)

	// todo move code below to better location
	//Set up front-end
	time.Sleep(5 * time.Second)
	messageHandler.SendStatus("general")
	for _, value := range messageHandler.Config.Devices {
		messageHandler.SendStatus(value.ID)
		messageHandler.GetStatus(value.ID)
	}

	// prevent exit
	select {}
}
