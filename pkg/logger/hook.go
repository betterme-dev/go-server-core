package logger

import (
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"strings"
)

type SentryHook struct {
}

func (h *SentryHook) Fire(entry *log.Entry) error {
	sentry.WithScope(func(scope *sentry.Scope) {
		msg := entry.Message
		// check if it is a nested message, and handle it
		if idx := strings.Index(msg, ":"); idx >= 0 {
			scope.SetExtra("message", msg)
			msg = msg[:idx]
		}
		scope.SetTags(map[string]string{
			"logger": "go",
			"level":  entry.Level.String(),
		})
		scope.SetContext("context", entry.Context)
		scope.SetExtra("data", entry.Data)
		scope.SetExtra("caller", entry.Caller)

		sentry.CaptureMessage(msg)
	})

	return nil
}

func (h *SentryHook) Levels() []log.Level {
	return []log.Level{
		log.WarnLevel,
		log.ErrorLevel,
		log.FatalLevel,
		log.PanicLevel,
	}
}
