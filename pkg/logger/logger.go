package logger

import (
	"time"

	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Setup initializes all logger settings and log handlers
func Setup() (func(), func(err error)) {
	// detect log level
	logLevel := viper.GetString("LOG_LEVEL")
	switch logLevel {
	case "debug":
		viper.SetDefault("DEBUG", true)
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}

	// set default handlers
	deferHandler := func() {
		if r := recover(); r != nil {
			log.Fatalf("panic: %s", r)
		}
	}
	errHandler := func(err error) {
		log.Error(err)
	}

	// prepare Sentry if needed
	sentryDSN := viper.GetString("SENTRY_DSN")
	if sentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              sentryDSN,
			Environment:      viper.GetString("APPLICATION_ENV"),
			Release:          viper.GetString("APP_VERSION"),
			AttachStacktrace: true,
		})
		if err == nil {
			log.AddHook(new(SentryHook))
			deferHandler = func() {
				sentry.Flush(time.Second * 5)
				sentry.Recover()
			}
			errHandler = func(err error) {
				sentry.CaptureException(err)
			}
		} else {
			log.Errorf("Failed to init sentry logger: %s", err)
		}
	}

	return deferHandler, errHandler
}
