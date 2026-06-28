package main

import (
	"cal-salary/core/server"
	"log"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatalf("Server terminated with error: %v", err)
	}
}
