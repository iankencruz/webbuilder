package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

const (
	cReset  = "\033[0m"
	cGray   = "\033[90m"
	cYellow = "\033[1;33m"
	cOrange = "\033[38;5;208m"
	cBlue   = "\033[94m"
	cGreen  = "\033[32m"
	cWhite  = "\033[37m"
	cRed    = "\033[31m"
	cCyan   = "\033[36m"
	cBold   = "\033[1m"
)

func statusColor(status int) string {
	switch {
	case status >= 500:
		return cRed
	case status >= 400:
		return cOrange
	case status >= 300:
		return cCyan
	default:
		return cGreen
	}
}

func methodColor(method string) string {
	switch method {
	case http.MethodGet:
		return cBlue
	case http.MethodPost:
		return cGreen
	case http.MethodPut:
		return cYellow
	case http.MethodDelete:
		return cRed
	case http.MethodPatch:
		return cOrange
	case http.MethodHead:
		return cCyan
	case http.MethodOptions:
		return cGray
	default:
		return cWhite
	}
}

func formatSize(size int64) string {
	switch {
	case size >= 1024*1024:
		return fmt.Sprintf("%.1fMB", float64(size)/1024/1024)
	case size >= 1024:
		return fmt.Sprintf("%.1fKB", float64(size)/1024)
	default:
		return fmt.Sprintf("%dB", size)
	}
}

func RequestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:       true,
		LogURI:          true,
		LogMethod:       true,
		LogLatency:      true,
		LogRemoteIP:     true,
		LogProtocol:     true,
		LogResponseSize: true,
		HandleError:     true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			ts := v.StartTime.Format("15:04:05")
			method := v.Method
			uri := v.URI
			protocol := v.Protocol
			ip := v.RemoteIP
			status := v.Status
			duration := v.Latency.Round(time.Microsecond)
			size := formatSize(v.ResponseSize)

			fmt.Printf(
				"%s%s%s  %s%s%s  from %s%s%s  Size: %s%s%s  Duration: %s%s%s  Status: %s%d%s  Method: %s%s%s  Path: %s%s%s\n",
				cWhite, ts, cReset,
				cBold+cYellow, protocol, cReset,
				cBold+cYellow, ip, cReset,
				cBold+cCyan, size, cReset,
				cBold+cYellow, duration, cReset,
				cBold+statusColor(status), status, cReset,
				cBold+methodColor(method), method, cReset,
				cBold+cYellow, uri, cReset,
			)

			return nil
		},
	})
}
