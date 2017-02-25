package main

import (
	"github.com/jimmysawczuk/power-monitor/monitor"
	"github.com/jimmysawczuk/power-monitor/web"

	"log"
	"net/http"
	"time"
)

func main() {
	m := monitor.New(5 * time.Second)
	go m.Start()

	http.Handle("/", web.GetRouter(&m))
	log.Println("Starting webserver")
	http.ListenAndServe(":3000", nil)
}
