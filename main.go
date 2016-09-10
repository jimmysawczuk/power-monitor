package main

import (
	"github.com/jimmysawczuk/power-monitor/monitor"
	"github.com/jimmysawczuk/power-monitor/web"

	"time"
)

func main() {
	m := monitor.New(5 * time.Second)
	go m.Start()

	w := web.New(&m)
	w.Run(":3000")
}
