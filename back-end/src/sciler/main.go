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

func main() {
	defer func() {
		if r := recover(); r != nil {
			logger.Panicf("Recovered panic: %v", r)
		}
	}()

	// set maximum number of cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	dir, dirErr := os.Getwd()
	if dirErr != nil {
		logger.Fatal(dirErr)
	}

	// configure logger
	setupLogger(dir)

	// start Application
	startApplication(dir)

	// prevent exit
	select {}
}

func startApplication(dir string) {
	filename := filepath.Join(dir, "back-end", "resources", "production", "room_config.json")
	configurations := config.ReadFile(filename)
	logger.Infof("configurations read from: %v", filename)

	messageHandler := handler.Handler{Config: configurations, ConfigFile: filename}
	messageHandler.Communicator = communication.NewCommunicator(configurations, messageHandler.NewHandler, func() {
		messageHandler.SendSetup()
	})

	logger.Infof("attempting to connect to broker at %s on port %v", configurations.General.Host, configurations.General.Port)
	messageHandler.Communicator.Start()
}

// setupLogger configures the logger such that both to file and console log messages are printed in the correct format
func setupLogger(dir string) {
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
	logger.SetLevel(logger.PanicLevel)
	logger.SetFormatter(&logger.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	// setting up (non colored) file output
	logger.AddHook(hook)

	logger.Infof("writing logs to both console and %v", writeFile)
}
