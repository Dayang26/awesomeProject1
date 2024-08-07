package ginblog

import (
	"awesomeProject1/docs"
	"awesomeProject1/internal/handle"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	// 后台管理系统接口

	userAuthAPI handle.UserAuth // 用户账号

)

func registerHandler(r *gin.Engine) {
	//Swagger todo
	docs.SwaggerInfo.BasePath = "/api"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	RegisterBaseHandler(r)
}

func RegisterBaseHandler(r *gin.Engine) {
	base := r.Group("/api")
	base.POST("/login", userAuthAPI.Login)
	base.POST("/register", userAuthAPI.Register)
	base.GET("/logout", userAuthAPI.Logout)
	base.GET("/re", userAuthAPI.SendCode)
}
