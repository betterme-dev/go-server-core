package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewPingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong", "version": viper.GetString("APP_VERSION")})
	}
}

func NewGoKitPingHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	data := map[string]interface{}{
		"message": "pong",
		"version": viper.GetString("APP_VERSION"),
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Error(err)
	}
}