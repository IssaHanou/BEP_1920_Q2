package main

import (
    "../config"
    "fmt"
)

func main() {
    fmt.Println("Starting server")
    data := []byte(`{
            "name": "My Awesome Escape",
            "duration": "00:30:00"
    }`)
    fmt.Println(config.GetFromJson(data))
}
