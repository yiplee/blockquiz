package hc

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yiplee/blockquiz/handler/api/render"
)

type HealthCheckView struct {
	Duration string `json:"duration,omitempty"`
}

func Handle() gin.HandlerFunc {
	start := time.Now()
	return func(c *gin.Context) {
		view := HealthCheckView{Duration: time.Since(start).String()}
		render.OK(c, view)
	}
}
