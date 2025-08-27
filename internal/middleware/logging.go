package middleware

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func LoggingMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()

		hlog.Infof("Request: %s %s from %s",
			string(c.Method()),
			string(c.Request.URI().Path()),
			string(c.ClientIP()),
		)

		c.Next(ctx)

		duration := time.Since(start)
		hlog.Infof("Response: %s %s - %d - %v",
			string(c.Method()),
			string(c.Request.URI().Path()),
			c.Response.StatusCode(),
			duration,
		)
	}
}
