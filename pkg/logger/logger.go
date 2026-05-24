package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/rs/zerolog"
)

func RequestLoggerConfig() middleware.RequestLoggerConfig {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
	}
	output.FormatLevel = func(i interface{}) string {
		var color string
		level := fmt.Sprintf("%s", i)
		switch level {
		case "info":
			color = "\033[32m"
		case "warn":
			color = "\033[33m"
		case "error":
			color = "\033[31m"
		case "debug":
			color = "\033[36m"
		default:
			color = "\033[0m"
		}
		return fmt.Sprintf("%s| %-6s|\033[0m", color, strings.ToUpper(level))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("\033[1m%s\033[0m", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("\033[90m%s=\033[0m", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	log := zerolog.New(output).With().Timestamp().Logger()

	return middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogMethod:   true,
		LogLatency:  true,
		HandleError: true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			e := log.Info()
			if v.Status >= 500 {
				e = log.Error()
			} else if v.Status >= 400 {
				e = log.Warn()
			}
			e.Str("method", v.Method).
				Int("status", v.Status).
				Str("uri", v.URI).
				Dur("latency", v.Latency.Round(time.Millisecond)).
				Msg("request")
			if v.Error != nil {
				log.Error().Err(v.Error).Msg("request error")
			}
			return nil
		},
	}
}
