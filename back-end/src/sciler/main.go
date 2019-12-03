package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"sciler/communication"
	"sciler/config"
	"sciler/handler"
	"time"
)

var topics = []string{"back-end", "status", "hint", "connect"}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		logrus.Fatal(err)
	}
	filename := dir + "/back-end/resources/room_config.json"
	configurations := config.ReadFile(filename)
	logrus.Info("configurations read from: " + filename)
	host := configurations.General.Host
	port := configurations.General.Port

	communicator := communication.NewCommunicator(host, port, topics)
	handler1 := handler.GetHandler(configurations, *communicator)
	go communicator.Start(handler1.NewHandler)

	// loop for now preventing app to exit
	for {
		time.Sleep(time.Microsecond * time.Duration(250))
	}
}
