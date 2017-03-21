package main

import (
	"github.com/jimmysawczuk/power-monitor/monitor"
	"github.com/jimmysawczuk/power-monitor/web"

	"log"
	"net/http"
	"time"
)

const (
	releaseModeRelease = "release"
	releaseModeDebug   = "debug"
)

var releaseMode = releaseModeDebug

func main() {
	m := monitor.New(5 * time.Second)
	go m.Start()

	http.Handle("/", web.GetRouter(&m))

	listen := ":3000"

	log.Printf("Starting web server in %s mode on %s", releaseMode, listen)
	http.ListenAndServe(listen, nil)
}
