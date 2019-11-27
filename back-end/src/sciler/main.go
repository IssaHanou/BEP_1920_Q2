package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"sciler/communication"
	"sciler/config"
	"time"
)

const host string = "localhost"
const port string = "1883"

var topics = []string{"back-end", "status", "hint", "connect"}

func main() {
	fmt.Println("Starting server")
	data := []byte(`{
            "name": "My Awesome Escape",
            "duration": "00:30:00"
    }`)
	fmt.Println(config.GetFromJSON(data))

	communicator := communication.NewCommunicator(host, port, topics)
	communicator.Start(func(client mqtt.Client, message mqtt.Message) {
		// Todo: Make advanced message handler which acts according to the events / configuration
		logrus.Info(string(message.Payload()))
	})

	//loop for now preventing app to exit
	for {
		time.Sleep(time.Microsecond * time.Duration(250))
	}
}
