package main

import (
	"monitor"

	"time"
)

func main() {
	m := monitor.New(2 * time.Second)
	go m.Start()

	for {
		// nothing
		time.Sleep(10 * time.Millisecond)
	}
}
