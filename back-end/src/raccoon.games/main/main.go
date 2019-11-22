package main

import (
	"../config"
	"fmt"
	"os"
)

func main() {
	os.Chdir("/src")
	fmt.Println("Starting server")
	data := []byte(`{
            "name": "My Awesome Escape",
            "duration": "00:30:00"
    }`)
	fmt.Println(config.GetFromJson(data))
}
