package ginblog

import (
	"awesomeProject1/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func registerHandler(r *gin.Engine) {
	//Swagger todo
	docs.SwaggerInfo.BasePath = "/api"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}

func registerBaseHandler(r *gin.Engine) {
	base := r.Group("/api")

	base.POST("/login")
}
