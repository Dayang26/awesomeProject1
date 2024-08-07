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
	blogInfoAPI handle.BlogInfo // 博客设置

)

func registerHandler(r *gin.Engine) {
	//Swagger todo
	docs.SwaggerInfo.BasePath = "/api"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	RegisterBaseHandler(r)
}

// 通用接口
func RegisterBaseHandler(r *gin.Engine) {
	base := r.Group("/api")
	base.POST("/login", userAuthAPI.Login)
	base.POST("/register", userAuthAPI.Register)
	base.GET("/logout", userAuthAPI.Logout)
	base.GET("/code", userAuthAPI.SendCode)
	base.GET("/config", blogInfoAPI.GetConfigMap)
	base.PATCH("/config", blogInfoAPI.UpdateConfig)
}

// 后台管理接口 全部需要 登录 + 鉴权
func registerAdminHandler(r *gin.Engine) {
	auth := r.Group("/api")
	auth.Use()

}
