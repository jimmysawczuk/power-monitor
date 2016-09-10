package web

import (
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jimmysawczuk/power-monitor/monitor"

	"strconv"
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
	r.LoadHTMLGlob("web/templates/*")

	r.GET("/", getIndex)
	r.GET("/api/snapshots", getSnapshots)

	r.Use(static.Serve("/", static.LocalFile("web/static/", false)))

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
		res := isTimestampInLast(s.Timestamp, now, 60*time.Second) ||
			(isTimestampInLast(s.Timestamp, now, 5*time.Minute) && isSignificantTimestamp(s.Timestamp, 10*time.Second)) ||
			(isTimestampInLast(s.Timestamp, now, 2*time.Hour) && isSignificantTimestamp(s.Timestamp, 5*time.Minute)) ||
			(isTimestampInLast(s.Timestamp, now, 48*time.Hour) && isSignificantTimestamp(s.Timestamp, 30*time.Minute))

		return res
	})

	limit_str := c.DefaultQuery("limit", "0")
	limit, _ := strconv.ParseInt(limit_str, 10, 64)
	if limit > 0 && limit < int64(len(recent)) {
		recent = recent[0:limit]
	}

	c.JSON(200, recent)
}

func isTimestampInLast(s, now time.Time, dur time.Duration) bool {
	return now.Sub(s) < dur
}

func isSignificantTimestamp(s time.Time, frequency time.Duration) bool {
	return s.UnixNano()%int64(frequency) < int64(active_monitor.Interval)
}
