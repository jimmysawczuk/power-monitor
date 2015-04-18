package web

import (
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"monitor"
	"time"
)

var active_monitor *monitor.Monitor
var start_time time.Time

func New(m *monitor.Monitor) *gin.Engine {
	_ = gzip.Gzip

	active_monitor = m

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.LoadHTMLGlob("src/web/templates/*")

	r.GET("/", getIndex)
	r.GET("/api/snapshots", getSnapshots) //, gzip.Gzip(gzip.DefaultCompression))

	r.Use(static.Serve("/", static.LocalFile("src/web/static/", false)))

	start_time = time.Now()

	return r
}

func getIndex(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"StartTime": start_time,
		"Interval":  int64(active_monitor.Interval / 1e6),
		"Mode":      gin.Mode(),
	})
}

func getSnapshots(c *gin.Context) {
	now := time.Now()

	recent := active_monitor.GetRecentSnapshots().Filter(func(s monitor.Snapshot) bool {
		dur := now.Sub(s.Timestamp)

		return isTimestampInLast(s.Timestamp, now, 60*time.Second) ||
			isSignificantTimestamp(s.Timestamp, dur, 5*time.Minute, 10*time.Second) ||
			isSignificantTimestamp(s.Timestamp, dur, 2*time.Hour, 5*time.Minute) ||
			isSignificantTimestamp(s.Timestamp, dur, 24*time.Hour, 15*time.Minute) ||
			isSignificantTimestamp(s.Timestamp, dur, 48*time.Hour, 30*time.Minute)
	})

	c.JSON(200, recent)
}

func isTimestampInLast(s, now time.Time, dur time.Duration) bool {
	return s.Sub(now) < dur
}

func isSignificantTimestamp(s time.Time, s_dur time.Duration, cutoff time.Duration, frequency time.Duration) bool {
	return (cutoff == 0 || s_dur < cutoff) && s.UnixNano()%int64(frequency) < int64(active_monitor.Interval)
}
