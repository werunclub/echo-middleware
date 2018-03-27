package middleware

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/Sirupsen/logrus"
)

const x_request_id = "X-Request-Id"

type Logger struct {
	*logrus.Logger
}

func (l Logger) Level() log.Lvl {
	switch l.Logger.Level {
	case logrus.DebugLevel:
		return log.DEBUG
	case logrus.WarnLevel:
		return log.WARN
	case logrus.ErrorLevel:
		return log.ERROR
	case logrus.InfoLevel:
		return log.INFO
	default:
		l.Panic("Invalid level")
	}

	return log.OFF
}

func (l Logger) SetPrefix(s string) {
	// TODO
}

func (l Logger) Prefix() string {
	// TODO.  Is this even valid?  I'm not sure it can be translated since
	// logrus uses a Formatter interface.  Which seems to me to probably be
	// a better way to do it.
	return ""
}

func (l Logger) SetLevel(lvl log.Lvl) {
	switch lvl {
	case log.DEBUG:
		logrus.SetLevel(logrus.DebugLevel)
	case log.WARN:
		logrus.SetLevel(logrus.WarnLevel)
	case log.ERROR:
		logrus.SetLevel(logrus.ErrorLevel)
	case log.INFO:
		logrus.SetLevel(logrus.InfoLevel)
	default:
		l.Panic("Invalid level")
	}
}

func (l Logger) Output() io.Writer {
	return l.Out
}

func (l Logger) SetOutput(w io.Writer) {
	logrus.SetOutput(w)
}

func (l Logger) Printj(j log.JSON) {
	logrus.WithFields(logrus.Fields(j)).Print()
}

func (l Logger) Debugj(j log.JSON) {
	logrus.WithFields(logrus.Fields(j)).Debug()
}

func (l Logger) Infoj(j log.JSON) {
	logrus.WithFields(logrus.Fields(j)).Info()
}

func (l Logger) Warnj(j log.JSON) {
	logrus.WithFields(logrus.Fields(j)).Warn()
}

func (l Logger) Errorj(j log.JSON) {
	logrus.WithFields(logrus.Fields(j)).Error()
}

func (l Logger) Fatalj(j log.JSON) {
	logrus.WithFields(logrus.Fields(j)).Fatal()
}

func (l Logger) Panicj(j log.JSON) {
	logrus.WithFields(logrus.Fields(j)).Panic()
}

func Logrus(log *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			start := time.Now()

			entry := log.WithFields(logrus.Fields{
				"request": ctx.Request().RequestURI,
				"method":  ctx.Request().Method,
				"remote":  ctx.Request().RemoteAddr,
			})

			if requestID := ctx.Request().Header.Get(x_request_id); requestID != "" {
				entry = entry.WithField("request_id", requestID)
			}

			entry.Debug("start handling request")
			ctx.Set("logger", *entry)

			if err := next(ctx); err != nil {
				ctx.Error(err)
			}

			latency := time.Since(start)

			entry.WithFields(logrus.Fields{
				"status":      ctx.Response().Status,
				"text_status": http.StatusText(ctx.Response().Status),
				"took":        latency.Seconds(),
				fmt.Sprintf("measure#%s.latency", "web"): latency.Nanoseconds(),
			}).Info("completed handling request")

			return nil
		}
	}
}
