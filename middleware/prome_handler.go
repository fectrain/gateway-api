package middleware

import (
	"github.com/gin-gonic/gin"
	"time"
)

func MwPrometheusHttp(c *gin.Context) {
	start := time.Now()
	path := c.Request.RequestURI
	ApiCount.WithLabelValues(path).Inc()

	c.Next()
	// after request
	end := time.Now()
	d := end.Sub(start) / time.Millisecond
	ApiSum.WithLabelValues(path).Add(float64(d))
}