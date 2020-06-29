package web

import (
	"context"
	"github.com/betterme-dev/go-server-core/pkg/web/handlers"
	"github.com/betterme-dev/go-server-core/pkg/web/middleware"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// App is the main structure of a web application. It is recommended that
// an app is created with the web.NewApp() function
type App struct {
	// Web engine of the app
	Engine *gin.Engine

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func (a App) Context() context.Context {
	return a.ctx
}

// NewApp creates a new web Application with some reasonable default settings
func NewApp() *App {
	viper.AutomaticEnv()
	ctx, cancelFunc := context.WithCancel(context.Background())

	return &App{
		ctx:        ctx,
		cancelFunc: cancelFunc,
		Engine:     gin.New(),
	}
}

func (a *App) Run() error {
	if viper.GetBool("DEBUG") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	httpServer := http.Server{
		Addr:           ":" + viper.GetString("PORT"),
		Handler:        a.getHandler(),
		ReadTimeout:    viper.GetDuration("READ_TIMEOUT") * time.Second,
		WriteTimeout:   viper.GetDuration("WRITE_TIMEOUT") * time.Second,
		IdleTimeout:    viper.GetDuration("IDLE_TIMEOUT") * time.Second,
		MaxHeaderBytes: 1 << 20, // prevent headers overflow
	}

	go func() {
		log.Infof("Starting web server on port %s", viper.GetString("PORT"))
		// service connections
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(a.ctx, 10*time.Second)
	defer shutdownCancel()

	return httpServer.Shutdown(shutdownCtx)
}

func (a *App) getHandler() http.Handler {
	ginHandler := a.Engine
	ginHandler.Use(middleware.Logger)

	// Heartbeat route
	ginHandler.GET("/ping", handlers.NewPingHandler())

	return http.TimeoutHandler(ginHandler, viper.GetDuration("REQUEST_TIMEOUT")*time.Second, "request timeout")
}
