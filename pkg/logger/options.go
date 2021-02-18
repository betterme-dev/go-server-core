package logger

type Options struct {
	Level  string
	Sentry struct {
		DSN        string
		AppEnv     string
		AppVersion string
	}
}
