package middleware

import (
	"context"
	"log"
	"runtime/debug"
	"strconv"
	"time"
	"kube/pkg/errors"

	"github.com/cloudwego/hertz/pkg/app"
)

func ErrorHandlerMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered: %v\nStack trace: %s", r, debug.Stack())

				panicErr := errors.New(
					errors.ErrCodeInternalError,
					"Internal server error",
					"An unexpected error occurred",
				)

				errors.SendError(c, panicErr)
			}
		}()

		c.Next(ctx)
	}
}

func RequestIDMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		requestID := string(c.GetHeader("X-Request-ID"))
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Set("request_id", requestID)

		c.Header("X-Request-ID", requestID)

		c.Next(ctx)
	}
}

func generateRequestID() string {
	return "req-" + time.Now().Format("20060102150405") + "-" + strconv.FormatInt(time.Now().UnixNano()%1000000, 10)
}
