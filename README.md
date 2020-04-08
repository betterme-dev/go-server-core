# What is this?

Go server core package (all common functionality shared between BetterMe services written on go)

# Whats included?

- DB connection with config (pkg/db)
- Send logs to sentry (pkg/logger)
- AMQP connection support (pkg/mq)
- Base web app with ping route (pkg/web)
- Filesystems support (pkg/filesystem)
- Env to store all open connections (pkg/env)
- Base console app (pkg/console)

# How to use?

- Download as a package `GO111MODULE=on go get github.com/betterme-dev/go-server-core`
- Import statement example `github.com/betterme-dev/go-server-core/pkg/env`
