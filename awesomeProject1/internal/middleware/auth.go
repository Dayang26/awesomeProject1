package middleware

import (
	g "awesomeProject1/internal/global"
	"awesomeProject1/internal/handle"
	"awesomeProject1/internal/model"
	"awesomeProject1/internal/utils/jwt"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"strings"
	"time"
)

// base on JWT Authorization
// get user information from session if session is existed
// get user information from token if session is not existed

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Debug("[middleware-JWTAuth] user auth not exist ,do jwt auth")

		db := c.MustGet(g.CTX_DB).(*gorm.DB)

		// 系统管理的资源需要做验证，没有加进来的不需要
		url, method := c.FullPath()[4:], c.Request.Method
		resource, err := model.GetResource(db, url, method)
		if err != nil {
			//没找到资源
			if errors.Is(err, gorm.ErrRecordNotFound) {
				slog.Debug("[middleware-JWTAuth] resource not exist,skip jwt auth")
				c.Set("skip_check", true)
				c.Next()
				c.Set("skip_check", false)
				return
			}
			handle.ReturnError(c, g.ErrDbOp, err)
			return
		}

		//匿名资源,直接跳过后续验证
		if resource.Anonymous {
			slog.Debug(fmt.Sprintf("[middleware-JWTAuth] resource: %s %s is anonymous, skip jwt auth!", url, method))
			c.Set("skip_check", true)
			c.Next()
			c.Set("skip_check", false)
			return
		}

		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			handle.ReturnError(c, g.ErrTokenNotExist, nil)
			return
		}

		// token 的正确格式: `Bearer [tokenString]`
		parts := strings.Split(authorization, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			handle.ReturnError(c, g.ErrTokenType, nil)
			return
		}

		claims, err := jwt.ParseToken(g.Conf.JWT.Secret, parts[1])
		if err != nil {
			handle.ReturnError(c, g.ErrTokenWrong, err)
			return
		}

		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			handle.ReturnError(c, g.ErrTokenRuntime, nil)
			return
		}

		user, err := model.GetUserAuthInfoById(db, claims.UserId)
		if err != nil {
			handle.ReturnError(c, g.ErrUserNotExist, err)
			return
		}

		session := sessions.Default(c)
		session.Set(g.CTX_USER_AUTH, claims.UserId)
		err = session.Save()
		if err != nil {
			slog.Debug(fmt.Sprintf("[middleware-JWTAuth] To save user data in session is fail! %s", err))
		}

		c.Set(g.CTX_USER_AUTH, user)
	}
}

// 资源访问权限验证
func PermissionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool("skip_check") {
			c.Next()
			return
		}

		db := c.MustGet(g.CTX_DB).(*gorm.DB)
		auth, err := handle.CurrentUserAuth(c)

		if err != nil {
			handle.ReturnError(c, g.ErrUserNotExist, err)
			return
		}

		if auth.IsSuper {
			slog.Debug("[middleware-PermissionCheck]: super admin no need to check, pass!")
			c.Next()
			return
		}

		fmt.Printf("current request path is %s \n", c.FullPath())
		url := c.FullPath()[4:]

		slog.Debug(fmt.Sprintf("[middlerware-PermissionCheck]"))
	}
}
