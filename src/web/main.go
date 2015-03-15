package web

import (
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"monitor"
)

var active_monitor *monitor.Monitor

func New(m *monitor.Monitor) *gin.Engine {
	_ = gzip.Gzip

	active_monitor = m

	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.LoadHTMLGlob("src/web/templates/*")

	r.GET("/", getIndex)
	r.GET("/api/snapshots", getSnapshots) //, gzip.Gzip(gzip.DefaultCompression))

	r.Use(static.Serve("/", static.LocalFile("src/web/static/", false)))

	return r
}

func getIndex(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{})
}

func getSnapshots(c *gin.Context) {
	c.JSON(200, active_monitor.GetRecentSnapshots())
}
