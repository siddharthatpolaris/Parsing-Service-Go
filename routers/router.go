package routers

import (
	"errors"
	"net/http"
	"parsing-service/pkg/logger"
	"parsing-service/routers/middleware"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func customHandleRecovery(c *gin.Context, err interface{}) {
	e := errors.New("unexpected Error Occurred")
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
}

func NewHandler() *Routes {
	router := gin.New()
	allowedHosts := viper.GetString("ALLOWED_HOSTS")

	err := router.SetTrustedProxies([]string{allowedHosts})
	if err != nil {
		logger.GetLogger().Fatalf("error in <NewHandler>: %v", err)
	}

	router.Use(gin.Logger())
	router.Use(gin.CustomRecovery(customHandleRecovery))
	router.Use(middleware.CORSMiddleware())

	return &Routes{
		Router: router,
	}
}
