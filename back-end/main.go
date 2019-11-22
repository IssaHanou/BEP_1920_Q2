package main

import (
	"fmt"
	"github.com/IssaHanou/BEP_1920_Q2/back-end/config"
)

func main() {
	fmt.Println("Starting server")
	data := []byte(`{
            "name": "My Awesome Escape",
            "duration": "00:30:00"
    }`)
	fmt.Println(config.GetFromJSON(data))
}
