package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

const RawBodyKey = "raw_body"

// CaptureRawBody reads the request body once, stores the bytes on the context
// and refills the body so later readers still see it.
//
// Echo's binder consumes the body without restoring it, so any handler that
// needs the exact bytes it received - webhook signature verification, for
// instance - has to capture them before binding runs.
func CaptureRawBody() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			if req.Body == nil {
				return next(c)
			}

			raw, err := io.ReadAll(req.Body)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "failed to read body")
			}
			req.Body.Close()

			c.Set(RawBodyKey, raw)
			req.Body = io.NopCloser(bytes.NewReader(raw))

			return next(c)
		}
	}
}

// GetRawBody returns the bytes captured by CaptureRawBody, or nil when the
// middleware is not registered on the route.
func GetRawBody(c echo.Context) []byte {
	if raw, ok := c.Get(RawBodyKey).([]byte); ok {
		return raw
	}
	return nil
}
