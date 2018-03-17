package web

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jimmysawczuk/power-monitor/monitor"

	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
)

var activeMonitor *monitor.Monitor
var startTime time.Time
var indexTmpl *template.Template
var releaseMode = releaseModeDebug

const (
	releaseModeRelease = "release"
	releaseModeDebug   = "debug"
)

func init() {
	if releaseMode == releaseModeRelease {
		indexTmpl = template.Must(template.New("name").Parse(string(MustAsset("web/templates/index.html"))))
	}
}

func GetRouter(m *monitor.Monitor) *mux.Router {
	activeMonitor = m
	startTime = time.Now()

	r := mux.NewRouter()
	r.Use(handlers.CompressHandler)
	r.Use(logRequest)

	r.Methods(http.MethodGet).Path("/").Handler(http.HandlerFunc(getIndex))
	r.Methods(http.MethodGet).Path("/api/snapshots").Handler(http.HandlerFunc(getSnapshots))
	r.PathPrefix("/").Handler(http.HandlerFunc(getStaticFile))
	return r
}

func getStaticFile(w http.ResponseWriter, r *http.Request) {
	by, err := Asset("web/static" + r.URL.Path)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	types := map[string]string{
		".css":   "text/css",
		".js":    "application/javascript",
		".eot":   "application/vnd.ms-fontobject",
		".otf":   "application/font-sfnt",
		".svg":   "image/svg+xml",
		".ttf":   "application/font-sfnt",
		".woff":  "application/font-woff",
		".woff2": "font/woff2",
	}

	if mime, exists := types[path.Ext(r.URL.Path)]; exists {
		w.Header().Set("Content-Type", mime)
	}

	w.WriteHeader(200)
	w.Write(by)
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := indexTmpl
	if tmpl == nil {
		tmpl = template.Must(template.New("name").Parse(string(MustAsset("web/templates/index.tmpl"))))
	}

	var revision string
	if rev, err := Asset("web/static/REVISION.json"); err != nil {
		revision = "{}"
	} else {
		buf := &bytes.Buffer{}
		json.Compact(buf, rev)
		revision = buf.String()
	}

	w.Header().Set("Content-Type", "text/html; charset=utf8")
	w.WriteHeader(200)
	tmpl.Execute(w, map[string]interface{}{
		"StartTime": startTime,
		"Interval":  30000,
		"Mode":      releaseMode,
		"Revision":  revision,
	})
}

func getSnapshots(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	snapshots := activeMonitor.GetRecentSnapshots()

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
	if snapshots := activeMonitor.GetRecentSnapshots(); len(snapshots) > 0 {
		latest = snapshots[0]
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	by, _ := json.Marshal(struct {
		Latest monitor.Snapshot             `json:"latest"`
		Recent monitor.SnapshotAverageSlice `json:"recent"`
	}{
		Latest: latest,
		Recent: recent,
	})

	w.WriteHeader(200)
	w.Write(by)
}

func isTimestampInLast(s, now time.Time, dur time.Duration) bool {
	return now.Sub(s) < dur
}

func isSignificantTimestamp(s, now time.Time, frequency time.Duration) bool {
	return (now.UnixNano()-s.UnixNano())%int64(frequency) < int64(activeMonitor.GetInterval())
}

type loggableResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (l *loggableResponseWriter) WriteHeader(statusCode int) {
	l.ResponseWriter.WriteHeader(statusCode)
	l.statusCode = statusCode
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw := &loggableResponseWriter{ResponseWriter: w}
		next.ServeHTTP(lw, r)

		if lw.statusCode == 0 {
			lw.statusCode = 200
		}

		log.Printf("[%d] %s %s", lw.statusCode, http.StatusText(lw.statusCode), r.RequestURI)
	})
}
