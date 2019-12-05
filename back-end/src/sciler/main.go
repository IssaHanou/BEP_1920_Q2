package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sciler/communication"
	"sciler/config"
	"sciler/handler"
	"time"
)

var topics = []string{"back-end", "hint", "connect", "status", "connection", "confirmation", "instruction"}

func main() {
	dir, dirErr := os.Getwd()
	if dirErr != nil {
		logrus.Fatal(dirErr)
	}
	// Write to both cmd and file
	writeFile := dir + "\\back-end\\output\\log-" + fmt.Sprint(time.Now().Format("02-01-2006--15-04-26")) + ".txt"
	file, fileErr := os.Create(writeFile)
	if fileErr != nil {
		logrus.Fatal(fileErr)
	}
	logrus.Info("writing logs to both console and " + writeFile)
	multi := io.MultiWriter(os.Stdout, file)
	logrus.SetOutput(multi)

	filename := dir + "\\back-end\\resources\\room_config.json"
	configurations := config.ReadFile(filename)
	logrus.Info("configurations read from: " + filename)
	host := configurations.General.Host
	port := configurations.General.Port

	communicator := communication.NewCommunicator(host, port, topics)
	messageHandler := handler.GetHandler(configurations, *communicator)
	go communicator.Start(messageHandler.NewHandler)

	// loop for now preventing app to exit
	for {
		time.Sleep(time.Microsecond * time.Duration(250))
	}
}
