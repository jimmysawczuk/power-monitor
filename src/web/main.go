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

		return ((dur < time.Duration(5*time.Minute) && s.Timestamp.UnixNano()%int64(10*time.Second) < int64(active_monitor.Interval)) ||
			(dur < time.Duration(2*time.Hour) && s.Timestamp.UnixNano()%int64(5*time.Minute) < int64(active_monitor.Interval)) ||
			s.Timestamp.UnixNano()%int64(15*time.Minute) < int64(active_monitor.Interval))
	})

	c.JSON(200, recent)
}
