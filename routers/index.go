package routers

import (
	"net/http"
	"parsing-service/docs"
	"parsing-service/pkg/database"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	
)

type Routes struct {
	Router *gin.Engine
}

func RegisterRoutes(r *Routes) {

	r.Router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Route Not Found"})
	})

	r.Router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"live": "ok"})
	})

	r.Router.GET("/liveness", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"live": "ok"})
	})

	r.Router.GET("/readiness", func(ctx *gin.Context) {
		CheckReadiness(ctx)
	})

	docs.SwaggerInfo.BasePath = "/v1"
	r.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}


func CheckReadiness(ctx *gin.Context) {
	var num int 
	database.GetDB().Raw("Select 1").Scan(&num)
	if num == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "unable to query database"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}