package main

import (
	"github.com/jimmysawczuk/power-monitor/monitor"
	"github.com/jimmysawczuk/power-monitor/web"

	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	releaseModeRelease = "release"
	releaseModeDebug   = "debug"
)

var releaseMode = releaseModeDebug
var port = 3000

func main() {
	m := monitor.New(5 * time.Second)
	go m.Start()

	listen := fmt.Sprintf(":%d", port)

	cert, _ := tls.X509KeyPair(
		MustAsset("certificate.pem"),
		MustAsset("key.pem"),
	)

	srv := &http.Server{
		Addr:    listen,
		Handler: web.GetRouter(&m),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	log.Printf("Starting web server in %s mode on %s:", releaseMode, listen)
	srv.ListenAndServeTLS("", "")
}
