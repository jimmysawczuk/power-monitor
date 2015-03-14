package main

import (
	"monitor"
	"web"

	"time"
)

func main() {
	m := monitor.New(5 * time.Second)
	go m.Start()

	w := web.New(&m)
	w.Run(":3000")
}
