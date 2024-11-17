package main

import (
	"log"
	"time"
)

func main() {
	start := time.Now()
	defer func() { log.Printf("Probe took %s", time.Since(start)) }()

	if err := probe(); err != nil {
		log.Fatal(err)
	}
}
