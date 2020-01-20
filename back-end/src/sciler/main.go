package main

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/rifflock/lfshook"
	logger "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
	"sciler/communication"
	"sciler/config"
	"sciler/handler"
	"time"
)

var topics = []string{"back-end"}

func main() {
	// set maximum number of cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	dir, dirErr := os.Getwd()
	if dirErr != nil {
		logger.Fatal(dirErr)
	}
	// Write to both cmd and file
	writeFile := filepath.Join(dir, "back-end", "output", "log-"+fmt.Sprint(time.Now().Format("02-01-2006--15-04-26"))+".txt")

	// Setup pathmap (maps loglevel to a file, currently all levels are logged in the same file)
	pathMap := lfshook.PathMap{}
	for _, level := range logger.AllLevels {
		pathMap[level] = writeFile
	}
	// create a hook for file for logger
	hook := lfshook.NewHook(pathMap, &logger.TextFormatter{
		FullTimestamp:   true,
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// setting up (colorable) console output
	logger.SetOutput(colorable.NewColorableStdout())
	logger.SetFormatter(&logger.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	// setting up (non colored) file output
	logger.AddHook(hook)

	logger.Info("writing logs to both console and " + writeFile)

	filename := filepath.Join(dir, "back-end", "resources", "production", "room_config.json")
	configurations := config.ReadFile(filename)
	logger.Info("configurations read from: " + filename)
	host := configurations.General.Host
	port := configurations.General.Port

	messageHandler := handler.Handler{Config: configurations, ConfigFile: filename}
	messageHandler.Communicator = communication.NewCommunicator(host, port, topics, messageHandler.NewHandler, func() {
		messageHandler.SendSetup()
	})

	messageHandler.Communicator.Start()

	// prevent exit
	select {}
}
