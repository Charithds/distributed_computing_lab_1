package main

import (
	"fmt"
	"os"
	"strings"

	"example.com/client"
	"example.com/server"
)

func main() {
	paramLen := len(os.Args)
	if paramLen != 2 {
		fmt.Println("give parameters")
		os.Exit(1)
	}

	application := os.Args[1]
	if strings.EqualFold(application, "server") {
		fmt.Println("Server starting")
		server.Start()
	} else {
		fmt.Println("Client starting")
		client.Start()
	}
}
