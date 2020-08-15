package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/jimmysawczuk/power-monitor/internal/monitor"
	"github.com/jimmysawczuk/power-monitor/internal/respond"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// version should hold the Git tag of the currently built API, injected via ldflags.
var version = "development"

// revision should hold the Git SHA of the currently built API, injected via ldflags.
var revision = ""

// date should hold the build date of the currently built API, injected via ldflags.
var date = ""

var cfg struct {
	Port int `envconfig:"PORT" default:"3000" required:"true"`

	MonitorInterval time.Duration `envconfig:"MONITOR_INTERVAL" default:"5s" required:"true"`
}

var log *logrus.Logger
var mon *monitor.Monitor
var startTime time.Time

func main() {
	log = logrus.New()

	if err := godotenv.Load(); err != nil {
		log.Println(err)
	}

	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	log.WithField("revision", revision).WithField("version", version).Info("starting up")

	startTime = time.Now()

	mon = monitor.New(cfg.MonitorInterval)
	go mon.Start()

	srv := &http.Server{
		Handler:      getRouter(),
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.WithField("port", cfg.Port).Info("listening")
	err := srv.ListenAndServe()
	log.Println("got here", err)
}

func getRouter() *chi.Mux {
	ch := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
	}).Handler

	r := chi.NewRouter()
	r.Use(logRequest(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(9))
	r.Use(ch)

	r.Get("/api/meta", getMeta)
	r.Get("/api/snapshots", getSnapshots)
	r.Get("/*", getStatic)

	return r
}

func logRequest(log *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			rctx := chi.RouteContext(r.Context())

			log.WithFields(logrus.Fields{
				"pattern": rctx.RoutePattern,
				"method":  r.Method,
				"url":     r.URL.String(),
				"dur":     time.Now().Sub(start),
			}).Infof("%s %s", r.Method, r.URL.String())
		})
	}
}

func getMeta(w http.ResponseWriter, r *http.Request) {
	t, _ := time.Parse(time.RFC3339, date)

	respond.WithSuccess(log, w, r, http.StatusOK, struct {
		StartTime   time.Time `json:"startTime"`
		Version     string    `json:"version"`
		Revision    string    `json:"revision"`
		RevisedDate time.Time `json:"revisedDate"`
	}{
		StartTime:   startTime.Truncate(time.Second),
		Version:     version,
		Revision:    revision,
		RevisedDate: t.Truncate(time.Second),
	})

}

func getSnapshots(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	snapshots := mon.GetRecentSnapshots()

	var recent monitor.SnapshotAverageSlice

	switch {
	case isTimestampInLast(startTime, now, 3*time.Minute):
		recent = snapshots.Rollup(5 * time.Second)

	case isTimestampInLast(startTime, now, 15*time.Minute):
		recent = snapshots.Rollup(10 * time.Second)

	case isTimestampInLast(startTime, now, 1*time.Hour):
		recent = snapshots.Rollup(30 * time.Second)

	case isTimestampInLast(startTime, now, 6*time.Hour):
		recent = snapshots.Rollup(5 * time.Minute)

	case isTimestampInLast(startTime, now, 2*24*time.Hour):
		recent = snapshots.Rollup(30 * time.Minute)

	default:
		recent = snapshots.Rollup(1 * time.Hour)
	}

	var latest monitor.Snapshot
	if snapshots := mon.GetRecentSnapshots(); len(snapshots) > 0 {
		latest = snapshots[0]
	}

	respond.WithSuccess(log, w, r, http.StatusOK, struct {
		Latest monitor.Snapshot             `json:"latest"`
		Recent monitor.SnapshotAverageSlice `json:"recent"`
	}{
		Latest: latest,
		Recent: recent,
	})
}

func getStatic(w http.ResponseWriter, r *http.Request) {
	rctx := chi.RouteContext(r.Context())

	path := rctx.URLParam("*")
	if path == "" {
		path = "index.html"
	}

	v, err := Asset(path)
	if err != nil {
		respond.WithError(log, w, r, http.StatusNotFound, nil)
		return
	}

	http.ServeContent(w, r, path, startTime, bytes.NewReader(v))
}

func isTimestampInLast(s, now time.Time, dur time.Duration) bool {
	return now.Sub(s) < dur
}

func isSignificantTimestamp(s, now time.Time, frequency time.Duration) bool {
	return (now.UnixNano()-s.UnixNano())%int64(frequency) < int64(mon.GetInterval())
}
