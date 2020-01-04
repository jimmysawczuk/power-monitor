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
var interval = 5 * time.Second

func main() {
	m := monitor.New(interval)
	go m.Start()

	listen := fmt.Sprintf(":%d", port)

	cert, _ := tls.X509KeyPair(
		MustAsset("tls/certificate.pem"),
		MustAsset("tls/key.pem"),
	)

	srv := &http.Server{
		Addr:    listen,
		Handler: web.GetRouter(m),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	log.Printf("Starting web server in %s mode on %s:", releaseMode, listen)
	srv.ListenAndServeTLS("", "")
}

func MustAsset(name string) []byte {
	asset, err := Asset(name)
	if err != nil {
		panic(err)
	}

	return asset
}
