package main

import (
	"fmt"
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
	communicator.Start()

	//loop for now preventing app to exit
	for {
		time.Sleep(time.Microsecond * time.Duration(250))
	}
}
