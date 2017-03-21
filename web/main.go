package web

import (
	"github.com/gorilla/mux"
	"github.com/jimmysawczuk/power-monitor/monitor"

	"encoding/json"
	"html/template"
	"net/http"
	"path"
	"strconv"
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
	r.HandleFunc("/", getIndex)
	r.HandleFunc("/api/snapshots", getSnapshots)
	r.PathPrefix("/").HandlerFunc(getStaticFile)
	return r
}
func getStaticFile(w http.ResponseWriter, r *http.Request) {
	by, err := Asset("web/static" + r.URL.Path)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	switch path.Ext(r.URL.Path) {
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	}

	w.WriteHeader(200)
	w.Write(by)
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)

	tmpl := indexTmpl
	if tmpl == nil {
		tmpl = template.Must(template.New("name").Parse(string(MustAsset("web/templates/index.html"))))
	}

	tmpl.Execute(w, map[string]interface{}{
		"StartTime": startTime,
		"Interval":  int64(activeMonitor.Interval / 1e6),
		"Mode":      releaseMode,
	})
}

func getSnapshots(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	now := time.Now()

	recent := activeMonitor.GetRecentSnapshots().Filter(func(s monitor.Snapshot) bool {

		switch {
		case isTimestampInLast(startTime, now, time.Minute):
			return true

		case isTimestampInLast(startTime, now, 5*time.Minute):
			return isSignificantTimestamp(s.Timestamp, 10*time.Second)

		case isTimestampInLast(startTime, now, 30*time.Minute):
			return isSignificantTimestamp(s.Timestamp, 30*time.Second)

		case isTimestampInLast(startTime, now, 2*time.Hour):
			return isSignificantTimestamp(s.Timestamp, 5*time.Minute)

		case isTimestampInLast(startTime, now, 2*24*time.Hour):
			return isSignificantTimestamp(s.Timestamp, 30*time.Minute)

		case isTimestampInLast(startTime, now, 4*24*time.Hour):
			return isSignificantTimestamp(s.Timestamp, 1*time.Hour)

		case isTimestampInLast(s.Timestamp, now, 7*24*time.Hour):
			return isSignificantTimestamp(s.Timestamp, 3*time.Hour)

		default:
			return false
		}
	})

	limit_str := r.FormValue("limit")
	limit, _ := strconv.ParseInt(limit_str, 10, 64)
	if limit > 0 && limit < int64(len(recent)) {
		recent = recent[0:limit]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	by, _ := json.Marshal(recent)
	w.Write(by)
}

func isTimestampInLast(s, now time.Time, dur time.Duration) bool {
	return now.Sub(s) < dur
}

func isSignificantTimestamp(s time.Time, frequency time.Duration) bool {
	return s.UnixNano()%int64(frequency) < int64(activeMonitor.Interval)
}
