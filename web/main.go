package web

import (
	"github.com/gorilla/mux"
	"github.com/jimmysawczuk/power-monitor/monitor"

	"bytes"
	"compress/gzip"
	"context"
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
	r.HandleFunc("/", http.HandlerFunc(getIndex).ServeHTTP)
	r.HandleFunc("/api/snapshots", http.HandlerFunc(getSnapshots).ServeHTTP)
	r.PathPrefix("/").HandlerFunc(http.HandlerFunc(getStaticFile).ServeHTTP)
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

	if compression := r.Context().Value("compression"); compression != nil {
		switch compression {
		case "gzip":
			w.Header().Set("Content-Encoding", "gzip")
			w.WriteHeader(200)
			gz := gzip.NewWriter(w)
			gz.Write(by)
			gz.Close()
			return
		}
	}

	w.WriteHeader(200)
	w.Write(by)
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := indexTmpl
	if tmpl == nil {
		tmpl = template.Must(template.New("name").Parse(string(MustAsset("web/templates/index.html"))))
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
		"Interval":  int64(activeMonitor.Interval / 1e6),
		"Mode":      releaseMode,
		"Revision":  revision,
	})
}

func getSnapshots(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	now := time.Now()
	var recent monitor.SnapshotSlice

	switch {
	case isTimestampInLast(startTime, now, 3*time.Minute):
		recent = activeMonitor.GetRecentSnapshots().Rollup(1 * time.Second)

	case isTimestampInLast(startTime, now, 10*time.Minute):
		recent = activeMonitor.GetRecentSnapshots().Rollup(10 * time.Second)

	case isTimestampInLast(startTime, now, 1*time.Hour):
		recent = activeMonitor.GetRecentSnapshots().Rollup(30 * time.Second)

	case isTimestampInLast(startTime, now, 6*time.Hour):
		recent = activeMonitor.GetRecentSnapshots().Rollup(5 * time.Minute)

	case isTimestampInLast(startTime, now, 2*24*time.Hour):
		recent = activeMonitor.GetRecentSnapshots().Rollup(30 * time.Minute)

	case isTimestampInLast(startTime, now, 4*24*time.Hour):
		recent = activeMonitor.GetRecentSnapshots().Rollup(1 * time.Hour)

	default:
		recent = activeMonitor.GetRecentSnapshots().Rollup(3 * time.Hour)
	}

	if limitStr := r.FormValue("limit"); limitStr != "" {
		limit, _ := strconv.ParseInt(limitStr, 10, 64)
		if limit > 0 && limit < int64(len(recent)) {
			recent = recent[0:limit]
		}
	}

	w.Header().Set("Content-Type", "application/json")
	by, _ := json.Marshal(recent)

	if compression := r.Context().Value("compression"); compression != nil {
		switch compression {
		case "gzip":
			w.Header().Set("Content-Encoding", "gzip")
			w.WriteHeader(200)
			gz := gzip.NewWriter(w)
			gz.Write(by)
			gz.Close()
			return
		}
	}

	w.WriteHeader(200)
	w.Write(by)
}

func isTimestampInLast(s, now time.Time, dur time.Duration) bool {
	return now.Sub(s) < dur
}

func isSignificantTimestamp(s, now time.Time, frequency time.Duration) bool {
	return (now.UnixNano()-s.UnixNano())%int64(frequency) < int64(activeMonitor.Interval)
}

func tryGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// to do: actually check for compression
		compression := "gzip"

		r = r.WithContext(context.WithValue(r.Context(), "compression", compression))

		next.ServeHTTP(w, r)

	})
}
