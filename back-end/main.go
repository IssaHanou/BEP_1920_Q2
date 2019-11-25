package main

import (
	"fmt"
	"github.com/IssaHanou/BEP_1920_Q2/back-end/config"
	"./communication"
	"time"
)

const host string = "localhost"
const port string = "1883"
var topics = []string{"back-end", "status", "hint"}

func main() {
	fmt.Println("Starting server")
	data := []byte(`{
            "name": "My Awesome Escape",
            "duration": "00:30:00"
    }`)
	fmt.Println(config.GetFromJSON(data))

	communicator := communication.NewCommunicator(host, port, topics)
	communicator.Start()
	for {
		time.Sleep(time.Microsecond * time.Duration(250))
	}
}

