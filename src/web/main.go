package web

import (
	"github.com/gin-gonic/gin"
	"monitor"
)

var active_monitor *monitor.Monitor

func New(m *monitor.Monitor) *gin.Engine {
	active_monitor = m

	r := gin.Default()
	r.LoadHTMLGlob("src/web/templates/*")

	r.GET("/", getIndex)
	r.GET("/api/snapshots", getSnapshots)

	return r
}

func getIndex(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{})
}

func getSnapshots(c *gin.Context) {
	c.JSON(200, active_monitor.GetRecentSnapshots())
}
